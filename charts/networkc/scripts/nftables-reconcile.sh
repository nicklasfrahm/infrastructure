#!/bin/bash

apk add --no-cache nftables

if [[ $# -ne 2 ]]; then
  echo "usage: $0 <old-config> <new-config>"
  exit 1
fi
old_config="$1"
new_config="$2"

if ! diff -q "$old_config" "$new_config"; then
  echo "Detected changes in nftables configuration. Reconciling..."
  cp "$new_config" "$old_config"
  nft -f "$old_config"
fi

# TODO: Add connectivity check to ensure that
# the nftables rules are working as expected.
