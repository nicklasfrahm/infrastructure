#!/bin/bash

# This script will be invoked by the armbian build system
# with the following arguments:
# ./userpatches/customize_image.sh $RELEASE $LINUXFAMILY $BOARD $BUILD_DESKTOP

# NOTE: It is copied to /tmp directory inside the image
# and executed there inside chroot environment, so don't
# reference any files that are not already installed.

# NOTE: If you want to transfer files between chroot and
# host, then the userpatches/overlay directory on the host
# is bind-mounted to /tmp/overlay in chroot. The SD card's
# root path is accessible via $SDCARD variable.

export DEBIAN_FRONTEND=noninteractive
export APT_LISTCHANGES_FRONTEND=none

# Configure the corresponding option in the openssh-server.
set_openssh_server_option() {
  option="$1"
  value="$2"

  # Check if a comment with the option exists and uncomment it.
  if grep -q "^#$option" /etc/ssh/sshd_config; then
    sed -i "s|^#$option|$option|" /etc/ssh/sshd_config
  fi

  if grep -q "^$option" /etc/ssh/sshd_config; then
    sed -i "s|^$option.*|$option $value|" /etc/ssh/sshd_config
  else
    echo "$option $value" >>/etc/ssh/sshd_config
  fi
}

# Set up users and initial passwords that will
# be rotated once the image was flashed.
configure_users() {
  # Avoid user config on first boot.
  rm /root/.not_logged_in_yet

  # Set random root password.
  ROOT_PASSWORD=$(openssl rand -hex 32)
  echo "root:${ROOT_PASSWORD}" | chpasswd

  # Disable autologin.
  rm -f /etc/systemd/system/getty@.service.d/override.conf
  rm -f /etc/systemd/system/serial-getty@.service.d/override.conf
  systemctl daemon-reload
}

configure_extra_packages() {
  apt-get install -y vim bind9-dnsutils curl
}

# Disable Armbian services that may cause additional resource
# consumption, such as ramlog and avoid services that conflict
# with cloud-init.
configure_armbian_services() {
  sed -i "s|^ENABLED=.*|ENABLED=false|" /etc/default/armbian-ramlog

  # Disable RAM logging to free up memory
  # for compute workloads.
  systemctl stop armbian-ramlog
  systemctl disable armbian-ramlog
  systemctl mask armbian-ramlog

  systemctl stop armbian-zram-config
  systemctl disable armbian-zram-config
  systemctl mask armbian-zram-config

  sed -i "s|^ENABLED=.*|ENABLED=false|" /etc/default/armbian-zram-config
  sed -i "s|^# SWAP=.*|SWAP=false|" /etc/default/armbian-zram-config

  rm /etc/cron.d/armbian-truncate-logs
  rm /etc/cron.daily/armbian-ram-logging

  # Disable armbian firstrun script because
  # it does not play well with the hardened
  # OpenSSH server configuration.
  systemctl stop armbian-firstrun
  systemctl disable armbian-firstrun
  systemctl mask armbian-firstrun

  systemctl daemon-reload
}

# Uninstall NetworkManager and install netplan.
configure_netplan() {
  apt-get purge -y network-manager
  apt-get install -y netplan.io

  rm /etc/netplan/armbian-default.yaml
  chmod 600 /etc/netplan/*.yaml

  systemctl unmask systemd-networkd
  systemctl enable systemd-networkd

  systemctl daemon-reload
}

# Ensure that I can logon to dropbear with my SSH key.
configure_cryptroot() {
  dropbear_initramfs="/etc/dropbear/initramfs"

  mkdir -p "$dropbear_initramfs"
  grep -oE "ssh-ed25519(.[^\"])*" /tmp/overlay/cloud-init/user-data | head -n 1 >"$dropbear_initramfs/authorized_keys"
  chmod 700 "$dropbear_initramfs"
  chmod 600 "$dropbear_initramfs/authorized_keys"
}

# Set up kboot for faster kernel updates.
configure_kboot() {
  apt-get install -y kexec-tools

  # Install kboot for faster kernel updates.
  cp /tmp/overlay/kboot /usr/local/bin/kboot
  chown root:root /usr/local/bin/kboot
  chmod +x /usr/local/bin/kboot
}

# Set up cloud-init.
# Reference: https://forum.armbian.com/topic/14616-cloud-init/
configure_cloud_init() {
  apt-get install -y cloud-init

  mkdir -p /boot/cloud-init
  cp /tmp/overlay/cloud-init/user-data /boot/cloud-init/user-data
  cp /tmp/overlay/cloud-init/meta-data /boot/cloud-init/meta-data
  INSTANCE_ID=$(uuidgen -r) envsubst </tmp/overlay/meta-data >/boot/cloud-init/meta-data

  # Configure cloud-init data source via kernel command line.
  # TODO: Investigate if hosting a metadata service makes sense.
  echo "extraargs=ds=nocloud;s=file://boot/cloud-init/" >>/boot/armbianEnv.txt

  # Configure cloud-init to use the network configuration from netplan.
  cp /tmp/overlay/cloud-init/network-config /etc/cloud/cloud.cfg.d/90-network-config.cfg
}

# Harden OpenSSH server.
configure_openssh_server() {
  set_openssh_server_option PasswordAuthentication no
  set_openssh_server_option PermitRootLogin no
  set_openssh_server_option PubkeyAuthentication yes
  set_openssh_server_option UseDNS no
  set_openssh_server_option PrintMotd no
  set_openssh_server_option UsePAM no
  set_openssh_server_option Banner no
  set_openssh_server_option X11Forwarding no
  set_openssh_server_option KbdInteractiveAuthentication no
}

main() {
  if [ $# -lt 1 ]; then
    echo "error: this script is meant to be called by the armbian build system"
    exit 1
  fi

  RELEASE="$1"

  case "$RELEASE" in
  jammy)
    # Ensure package index is up to date.
    apt-get update

    configure_users
    configure_extra_packages
    configure_armbian_services
    configure_netplan
    configure_cryptroot
    configure_kboot
    configure_cloud_init
    configure_openssh_server
    ;;
  *)
    echo "error: unknown release $RELEASE"
    exit 1
    ;;
  esac
}

main "$@"
