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

# Check if a file exists and create it if not.
ensure_file_exists() {
  if [ ! -f "$1" ]; then
    touch "$1"
  fi
}

# Add kernel command line arguments via armbianEnv.txt.
append_armbian_extraargs() {
  armbian_env_file="/boot/armbianEnv.txt"
  new_kargs="$1"

  ensure_file_exists "$armbian_env_file"

  if grep -q "^extraargs=" "$armbian_env_file"; then
    # Check if value is quoted.
    if grep -q "^extraargs=\".*\"" "$armbian_env_file"; then
      sed -i "s/^extraargs=\"\(.*\)\"/extraargs=\"\1 $new_kargs\"/" "$armbian_env_file"
    else
      sed -i "s/^extraargs=\(.*\)/extraargs=\"\1 $new_kargs\"/" "$armbian_env_file"
    fi
  else
    echo "extraargs=\"$new_kargs\"" >>"$armbian_env_file"
  fi
}

# Set up users and initial passwords that will
# be rotated once the image was flashed.
configure_users() {
  # Avoid user config on first boot.
  ROOT_PASSWORD=$(openssl rand -hex 32)
  echo "root:${ROOT_PASSWORD}" | chpasswd
  rm /root/.not_logged_in_yet
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

  cp -r /tmp/overlay/cloud-init /boot/cloud-init

  # Configure cloud-init data source via kernel command line.
  append_armbian_extraargs "ds=nocloud;s=/boot/cloud-init/"
}

# Configure CPU and memory sets.
configure_cpu_memory_sets() {
  append_armbian_extraargs "cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory"
}

# Set up OpenSSH server.
configure_openssh_server() {
  # TODO: Harden OpenSSH server.
  # - HostKey /etc/ssh/ssh_host_ed25519_key
  # - PasswordAuthentication no
  # - PermitRootLogin no
  # - PubkeyAuthentication yes
  # - PermitEmptyPasswords no
  # - UseDns no
  # - PrintMotd no
  # - UsePAM no
  # - Banner [ no ]
  # - X11Forwarding no
  # - KbdInteractiveAuthentication no
  true
}

main() {
  if [ $# -ne 4 ]; then
    echo "error: this script is meant to be called by the armbian build system"
    exit 1
  fi

  RELEASE="$1"

  case "$RELEASE" in
  jammy)
    # Ensure package index is up to date.
    apt-get update

    configure_users
    configure_cryptroot
    configure_kboot
    configure_cloud_init
    configure_cpu_memory_sets
    configure_openssh_server
    ;;
  *)
    echo "error: unknown release $RELEASE"
    exit 1
    ;;
  esac
}

main "$@"
