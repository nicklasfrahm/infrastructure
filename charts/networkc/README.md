# 📡 networkc

`networkc` is a software-defined routing stack meant to be deployed to a Kubernetes cluster.

## 🚀 Quickstart

```bash
helm upgrade --install networkc charts/networkc --namespace networkc
```

## 📦 Components

For now the implementation is naive and simple. It consists of the following components:

- `CronJob` for updating the firewall configuration via `nftables`
  <!-- - `CronJob` to update the `isc-dhcp-server` configuration -->
  <!-- - `CronJob` to update the `netplan` configuration -->
  <!-- - `CronJob` to run `ddclient` -->
