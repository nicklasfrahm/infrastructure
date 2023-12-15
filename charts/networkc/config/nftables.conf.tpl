#!/usr/sbin/nft -f

# This setup is based on the Red Hat Documentation of nftables:
# https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/7/html/security_guide/sec-configuring_nat_using_nftables

# Reset firewall flushing the ruleset.
flush ruleset

# Declare variables.
{{- if not (.Values.interfaces.wan) }}
{{- fail "failed to find WAN interface: .interfaces.wan" }}
{{- end }}
define iface_wan = "{{ .Values.interfaces.wan }}"
{{- if not (.Values.interfaces.lan) }}
{{- fail "failed to find LAN interface: .interfaces.lan" }}
{{- end }}
define iface_lan = "{{ .Values.interfaces.lan }}"
{{- if not (.Values.subnets.kubernetes.pods) }}
{{- fail "failed to find Kubernetes pod subnet: .subnets.kubernetes.pods" }}
{{- end }}
{{- if not (.Values.subnets.kubernetes.services) }}
{{- fail "failed to find Kubernetes service subnet: .subnets.kubernetes.services" }}
{{- end }}
define cidr_kubernetes = "{ {{ .Values.subnets.kubernetes.pods }}, {{ .Values.subnets.kubernetes.services }} }"

# Create a table for firewalling.
table inet firewall {
  # Handle incoming traffic for host processes.
  chain input {
    # Drop all traffic by default.
    type filter hook input priority filter; policy drop;

    # Let invalid connections time out, while accepting established
    # or related connections.
    ct state invalid counter drop comment "invalid input packets"
    ct state established,related counter accept comment "accepted input packets"

    # Accept connections from the loopback interface to itself while
    # protecting it from external traffic.
    iif lo accept
    iif != lo ip daddr 127.0.0.1/8 counter drop comment "dropped lo packets"
    iif != lo ip6 daddr ::1/128 counter drop comment "dropped lo packets"

    # Allow ICMP and IGMP traffic for correct network operation.
    # Reference: http://shouldiblockicmp.com/
    ip protocol icmp counter accept comment "icmp input packets"
    ip protocol igmp counter accept comment "igmp input packets"
    ip6 nexthdr icmpv6 counter accept comment "icmpv6 input packets"

    # Allow SSH, HTTP and HTTPS.
    tcp dport { ssh, 80, 443, 6443, 7443 } counter accept comment "accepted input packets"

    # Allow wireguard traffic.
    udp dport { 5800, 5801, 5802, 5803 } counter accept comment "accepted wireguard packets"

    # Ensure that Kubernetes node metrics can be collected and BGP routes can be advertised.
    iif != $iface_wan tcp dport { 179, 9100, 10250 } counter accept comment "accepted input packets"

    # Ensure that VXLAN works.
    iif != $iface_wan udp dport { 8472 } counter accept comment "accepted input packets"

    # Count dropped packets.
    counter comment "dropped input packets"
  }

  # Handle forwarded traffic between network interfaces.
  chain forward {
    # Drop all traffic by default.
    type filter hook forward priority 0; policy drop;

    # Clamp MSS to MTU based on routing cache via PMTUD. Without this
    # TCP traffic might not work if your ISP uses PPPoE or FTTH.
    tcp flags syn tcp option maxseg size set rt mtu

    # Allow forwarded traffic only from internal network interfaces
    # or if packets are established or related.
    ct state established,related counter accept comment "accepted related packets"

    # Allow forwarded traffic to Kubernetes services and pods.
    ip daddr $cidr_kubernetes counter accept comment "accepted kubernetes packets"
    ip saddr $cidr_kubernetes counter accept comment "accepted kubernetes packets"

    # Allow forwarded traffic towards the internet.
    iif != $iface_wan counter accept comment "accepted internet packets"

    # Count dropped packets.
    counter comment "dropped alien packets"
  }

  # Handle outgoing traffic from host processes.
  chain output {
    # Allows all outgoing traffic by default.
    type filter hook output priority 0; policy accept;
    counter comment "accepted output packets"
  }

  # Handle outgoing packets before leaving the system.
  chain postrouting {
    # Register a hook to in the postrouting chain.
    type nat hook postrouting priority srcnat; policy accept;

    # Rewrite the source address with the address of the outgoing interface.
    oifname $iface_wan masquerade
  }
}
