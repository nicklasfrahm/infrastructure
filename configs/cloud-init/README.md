# Cloud-init

This directory contains [cloud-init][docs-cloud-init] configuration files for the automated configuration of manually provisioned devices.

To provision a device, simply copy these appropriate files for the device to the boot partition of the desired system. For a pre-installed Ubuntu SD card image, such as the ones provided for the Raspberry Pi, simply copy theses files to the `system-boot` partition.

If you want to run `k3s` make sure to also modify the `cmdline.txt` file by adding the following to the very beginning of the line:

```txt
cgroup_enable=cpuset cgroup_enable=memory cgroup_memory=1
```

[docs-cloud-init]: https://cloudinit.readthedocs.io/en/latest/
