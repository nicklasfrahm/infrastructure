# Appliance lifycycle

An appliance is a piece of hardware that comes with a pre-installed operating system. This document describes how the lifecycle of appliances is managed.

## Provisioning 🪄

For different hardware the appliance provisioning process is different. Please see the different sections below for more information.

### Nano Pi R5S 📡

1. Build a fresh appliance image:

   ```bash
   BOARD=nanopi-r5s make build-appliance
   ```

2. Use [Etcher][etcher] to flash the image to an SD card.

[etcher]: https://etcher.balena.io/
