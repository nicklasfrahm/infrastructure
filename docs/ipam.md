# IP Address Management

This document describes the allocation of IP addresses.

## Routers

This section covers the autonomous system (AS) numbers and IP addresses of the router loopback interfaces.

| Hostname  | IP address (`lo`) | AS number | Description                     |
| --------- | ----------------- | --------- | ------------------------------- |
| `alfa`    | `172.31.255.0/32` | `65000`   | Edge router for site `dktil01`. |
| `bravo`   | `172.31.255.1/32` | `65001`   | Edge router for site `deflf01`. |
| `charlie` | `172.31.255.2/32` | `65002`   | Edge router for site `deflf02`. |
| `delta`   | `172.31.255.3/32` | `65003`   | Edge router for site `dksjb01`. |

## Transit networks

This section contains the IP addresses of transit networks.

| Network          | Description                                             |
| ---------------- | ------------------------------------------------------- |
| `172.31.0.0/31`  | Wireguard tunnel between `alfa wg1` and `bravo wg0`.    |
| `172.31.0.2/31`  | Wireguard tunnel between `alfa wg2` and `charlie wg0`.  |
| `172.31.0.4/31`  | Wireguard tunnel between `alfa wg3` and `delta wg0`.    |
| `172.31.0.6/31`  | Wireguard tunnel between `bravo wg2` and `charlie wg1`. |
| `172.31.0.8/31`  | Wireguard tunnel between `bravo wg3` and `delta wg1`.   |
| `172.31.0.10/31` | Wireguard tunnel between `charlie wg3` and `delta wg2`. |

## VPN networks

This section contains the IP addresses of VPN networks.

| Network         | Gateway        | Description                                   |
| --------------- | -------------- | --------------------------------------------- |
| `172.30.0.0/24` | `172.30.0.254` | Wireguard peer-to-site VPN network on `alfa`. |
