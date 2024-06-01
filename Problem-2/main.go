package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cilium/ebpf"
	"github.com/vishvananda/netlink"
)

const (
	targetPort        = 4040
	maxProcessNameLen = 16
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -Wall -Werror" bpf filter.c -- -I/usr/include

func main() {
	var targetProcess string
	var interfaceName string

	flag.StringVar(&targetProcess, "process", "myprocess", "Target process name")
	flag.StringVar(&interfaceName, "interface", "eth0", "Network interface to attach the filter")
	flag.Parse()

	
	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		log.Fatalf("loading objects: %v", err)
	}
	defer objs.Close()

	// Update the target process name in the map
	key := uint32(0)
	value := make([]byte, maxProcessNameLen)
	copy(value, targetProcess)

	if err := objs.targetProcess.Put(key, value); err != nil {
		log.Fatalf("updating target process: %v", err)
	}

	// Get the network interface
	link, err := netlink.LinkByName(interfaceName)
	if err != nil {
		log.Fatalf("getting interface %q: %v", interfaceName, err)
	}

	// Attach the eBPF program to the interface
	l := netlink.Link(link).Attrs()
	qdisc := &netlink.GenericQdisc{
		QdiscAttrs: netlink.QdiscAttrs{
			LinkIndex: l.Index,
			Handle:   netlink.MakeHandle(0xffff, 0),
			Parent:   netlink.HANDLE_CLSACT,
		},
		QdiscType: "clsact",
	}

	if err := netlink.QdiscAdd(qdisc); err != nil {
		log.Fatalf("adding clsact qdisc: %v", err)
	}

	classId := netlink.MakeHandle(1, 1)
	filter := &netlink.BpfFilter{
		FilterAttrs: netlink.FilterAttrs{
			LinkIndex: l.Index,
			Parent:   netlink.HANDLE_MIN_EGRESS,
			Handle:   classId,
			Protocol: syscall.ETH_P_ALL,
			Priority: 1,
		},
		Fd:           objs.filterTraffic.FD(),
		Name:        "filter_traffic",
		DirectAction: true,
	}

	if err := netlink.FilterAdd(filter); err != nil {
		log.Fatalf("attaching filter: %v", err)
	}

	log.Printf("Attached filter to interface %s. Allowing traffic only on port %d for process %q", interfaceName, targetPort, targetProcess)

	
	stopper := make(chan os.Signal, 1)
	signal.Notify(stopper, os.Interrupt, syscall.SIGTERM)
	<-stopper

	log.Println("Removing filter...")
	netlink.FilterDel(filter)
	netlink.QdiscDel(qdisc)
}