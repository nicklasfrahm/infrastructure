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
3. Log in to the appliance using the default credentials:

   - Username: `nicklasfrahm`
   - Password: `nicklasfrahm`

4. Fetch the IP address of the appliance:

   ```bash
   ip addr show wan | grep 'inet ' | awk '{print $2}'
   ```

5. On a different machine, run the bootstrap command:

   ```bash
   ic metal bootstrap <ip>
   ```

[etcher]: https://etcher.balena.io/
