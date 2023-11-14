# Firewall

This document describes the firewall rules for the network.

## Inbound

This section describes the inbound firewall rules.

| Port        | L4 protocol | L7 protocol | Description                |
| ----------- | ----------- | ----------- | -------------------------- |
| `22`        | `tcp`       | `ssh`       | SSH                        |
| `80`        | `tcp`       | `http`      | HTTP reverse proxy         |
| `443`       | `tcp`       | `https`     | HTTPS reverse proxy        |
| `5800-5810` | `udp`       | `wireguard` | Wireguard site-to-site VPN |
| `6443`      | `tcp`       | `https`     | Kubernetes reverse proxy   |
| `7443`      | `tcp`       | `https`     | Kubernetes API server      |
|             | `icmp`      |             | ICMP                       |
|             | `igmp`      |             | IGMP                       |
|             | `icmpv6`    |             | IGMP                       |
