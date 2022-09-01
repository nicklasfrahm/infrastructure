# Networking

This document contains documentation for my network configuration.

## Cisco NX-OS

This section contains configuration snippets for Cisco NX-OS.

### Add port to LACP group

```sh
interface ethernet 1/<if>
channel-group 1 force mode passive
```

### Add LACP group to VLAN

```sh
interface port-channel <group>
lacp min-links 1
switchport
switchport mode access
switchport access vlan <vlan>
```

### Configure trunk port

```sh
interface ethernet 1/<if>
switchport mode trunk
```

## SSH

```sh
ssh -o PubkeyAcceptedKeyTypes=+ssh-rsa -o HostKeyAlgorithms=+ssh-rsa admin@<ip>
```
