#include <linux/bpf.h>
#include <linux/if_ether.h>
#include <linux/ip.h>
#include <linux/tcp.h>
#include <linux/in.h>
#include <bpf/bpf_helpers.h>
#include <bpf/bpf_endian.h>

#define TARGET_PORT 4040
#define MAX_PROCESS_NAME_LEN 16

struct {
    __uint(type, BPF_MAP_TYPE_ARRAY);
    __type(key, __u32);
    __type(value, char[MAX_PROCESS_NAME_LEN]);
    __uint(max_entries, 1);
} target_process SEC(".maps");

SEC("socket")
int filter_traffic(struct __sk_buff *skb) {
    void *data = (void *)(long)skb->data;
    void *data_end = (void *)(long)skb->data_end;
    
    struct ethhdr *eth = data;
    if ((void *)(eth + 1) > data_end)
        return 0;
    
    if (eth->h_proto != bpf_htons(ETH_P_IP))
        return 0;
    
    struct iphdr *ip = (struct iphdr *)(eth + 1);
    if ((void *)(ip + 1) > data_end)
        return 0;
    
    if (ip->protocol != IPPROTO_TCP)
        return 0;
    
    struct tcphdr *tcp = (struct tcphdr *)(ip + 1);
    if ((void *)(tcp + 1) > data_end)
        return 0;
    
    __u16 dst_port = bpf_ntohs(tcp->dest);
    
    char process_name[MAX_PROCESS_NAME_LEN] = {};
    bpf_get_current_comm(process_name, sizeof(process_name));
    
    __u32 key = 0;
    char *target_name = bpf_map_lookup_elem(&target_process, &key);
    if (!target_name)
        return 0;
    
    for (int i = 0; i < MAX_PROCESS_NAME_LEN; i++) {
        if (process_name[i] != target_name[i])
            return 0;
        if (process_name[i] == '\0' && target_name[i] == '\0')
            break;
    }
    
    if (dst_port == TARGET_PORT)
        return 1;
    
    return 0;
}
char _license[] SEC("license") = "GPL";