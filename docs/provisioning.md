# Provisioning 🪄

This document describes how new bare metal servers are deployed. Irrespective of the type of server this process is manual. Currently there is not much value in automating this as it is not done often. Depending on the server hardware there are different steps to provision a server.

## Generic installation 💽

This is the default process and applies to most servers.

### Preparation 📝

You will need the following things to install the server:

- [Ubuntu server ISO][website-ubuntu-server]
- 4 GB USB stick

If you are using Windows, you may use a tool, such as [balenaEtcher][website-balena-etcher] to create to create an installation medium. On Linux you may use `dd` via the following command.

```bash
sudo dd if=ubuntu.iso of=/dev/sdb bs=1M status=progress
```

Now, start the server and navigate to the BIOS. If the system supports UEFI, make sure to disable CSM. Also make sure to set a UEFI administrator password.

If the server has a hardware RAID controller create a RAID1 array. Ideally a seperate RAID1 is set up for the OS and data drives. Save the changes and reboot the server. You are now ready to install the OS.

### OS installation 📦

Boot the server and select the appropriate language and keyboard layout. Once you are asked to set up the network, enable DHCP on the primary network interface and wait for a connection to be established.

When configuring the partitions, make sure to select the **Custom** option. Select a boot drive. If you have multiple drives of the same size and you did not set up a hardware RAID, select a secondary boot drive. Now use the following partitioning scheme.

| Partition | Mount point | Size  | Type    |
| --------- | ----------- | ----- | ------- |
| `sdX1`    | `/boot/efi` | `1G`  | `fat32` |
| `sdX2`    | `/boot`     | `1G`  | `ext4`  |
| `sdX3`    | `n/a`       | `max` | `n/a`   |

If you have multiple drives of the same size for the **operating system**, create two software RAID arrays, `md0` and `md1` and use `md0` with an `ext4` filesystem for `/boot`. If you did not create `md0` use `sdX2` directly.

If you previously created `md1`, use it to create an **encrypted volume group** `vg0`, otherwise use `/dev/sdX3` directly. Within the volume group `vg0` create a logical volume `lv0`, create a `btrfs` filesystem and mount it at `/`.

Make sure to leave all data disks unconfigured as we will configure them later via [Rook][website-rook]. Finalize the installation by importing the SSH public keys and booting into the freshly installed operating system. Remember to type your disk encryption password when prompted.

### Finalization 🧹

Start out by setting up passwordless `sudo`.

```bash
echo "$(whoami) ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/$(whoami)
```

Uninstall `snapd` by running `sudo apt-get purge --autoremove snapd`.

Disable any swap by running `sudo swapoff -a`. Further make this permanent by commenting out all lines in `/etc/fstab` that define a swap.

[website-ubuntu-server]: https://ubuntu.com/download/server
[website-balena-etcher]: https://balena.io/etcher/
[website-rook]: https://rook.io
