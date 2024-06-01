#include <linux/bpf.h>
#include <bpf/bpf_helpers.h>

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __uint(max_entries, 1);
    __type(key, __u32);
    __type(value, __u16);
} port_map SEC(".maps");

SEC("tc")
int drop_tcp_port(struct __sk_buff *skb)
{
    __u16 dport;
    __u16 configured_port;

    dport = bpf_ntohs(load_half(skb, ETH_HLEN + sizeof(struct iphdr) + offsetof(struct tcphdr, dest)));
    configured_port = bpf_map_lookup_elem(&port_map, 0);

    if (dport == configured_port) {
        return TC_ACT_SHOT;
    }

    return TC_ACT_OK;
}

char __license[] SEC("license") = "GPL";