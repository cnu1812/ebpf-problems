Write an eBPF code to allow traffic only at a specific TCP port (default 4040) for a
given process name (for e.g, &quot;myprocess&quot;). All the traffic to all other ports for only
that process should be dropped.