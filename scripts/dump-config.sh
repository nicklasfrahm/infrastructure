#!/usr/bin/env bash
# Usage: ./scripts/dump-config.sh <host>

config_dir="configs/sites"

if [ $# -lt 1 ] || [ "$1" == "-h" ] || [ "$1" == "--help" ]; then
  echo "Usage: $0 <host>"
  echo
  echo "Dump configuration of routers and baremetal"
  echo "hosts into the folder \"$config_dir\"."
  exit 1
fi
host="$1"

config_files=(
  "/etc/nftables.conf"
  "/etc/default/isc-dhcp-server"
  "/etc/dhcp/dhcpd.conf"
  "/etc/netplan/00-installer-config.yaml"
  "/etc/netplan/50-cloud-init.yaml"
)

echo "Probing host availability: $host"
if ! ssh "$host" "sudo hostname" >/dev/null 2>&1; then
  echo "error: failed to probe host: $host"
  exit 1
fi

mkdir -p "$config_dir/$host"

for config_file in "${config_files[@]}"; do
  file=$(basename "$config_file")
  dump_file="$config_dir/$host/$file"
  # shellcheck disable=SC2029
  if ! ssh "$host" "sudo cat $config_file" >"$dump_file" 2>/dev/null; then
    rm "$dump_file"
  fi
done
