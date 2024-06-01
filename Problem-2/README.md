## Install the required dependencies:

``` go get github.com/cilium/ebpf```
```go get github.com/vishvananda/netlink```

## Compiling the eBPF code

``` go run github.com/cilium/ebpf/cmd/bpf2go -cc clang -cflags "-O2 -g -Wall -Werror" bpf filter.c -- -I/usr/include ```

OR

``` go generate ```

## Compiling the go code

``` go build -o filter ```

## Run the Program

```sudo ./filter --process "firefox" --interface "eth0"```

Replace "firefox" with the name of the process you want to restrict, and "eth0" with your actual network interface. You can find your network interface names using ```ip link show```.

Example: $ ip link show
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN mode DEFAULT group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP mode DEFAULT group default qlen 1000
    link/ether 00:22:48:6e:c3:60 brd ff:ff:ff:ff:ff:ff
3: docker0: <NO-CARRIER,BROADCAST,MULTICAST,UP> mtu 1500 qdisc noqueue state DOWN mode DEFAULT group default 
    link/ether 02:42:c8:26:1b:bf brd ff:ff:ff:ff:ff:ff

lo is not we want


Output: ```2024/05/30 12:17:11 Attached filter to interface eth0. Allowing traffic only on port 4040 for process "firefox"```

