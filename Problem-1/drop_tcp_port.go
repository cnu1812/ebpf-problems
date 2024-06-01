package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
)

const (
	mapName = "port_map"
)

//go:generate go run github.com/cilium/ebpf/cmd/bpf2go -cc $BPF_CLANG -cgo drop_tcp_port.c -- -I$DIR/../headers
func main() {
	// Increase the resource limits for locked memory
	if err := rlimit.RemoveMemlock(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to remove memlock: %v\n", err)
		os.Exit(1)
	}

	// Load the eBPF program
	spec, err := loadDropTCPPort()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load eBPF program: %v\n", err)
		os.Exit(1)
	}
	defer spec.Program.Close()

	// Get the map file descriptor
	mapFD, err := spec.Maps[mapName].FD()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get map FD: %v\n", err)
		os.Exit(1)
	}

	// Check if the port number is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <port_number>\n", os.Args[0])
		os.Exit(1)
	}

	// Parse the port number from the command-line argument
	portNum, err := strconv.ParseUint(os.Args[1], 10, 16)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid port number: %v\n", err)
		os.Exit(1)
	}

	// Update the port number in the eBPF map
	if err := ebpf.UpdateElement(mapFD, 0, uint16(portNum), 0); err != nil {
		fmt.Fprintf(os.Stderr, "failed to update map: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Port number updated to %d\n", portNum)

	// Keep the program running
	var stopCh chan struct{}
	<-stopCh
}
