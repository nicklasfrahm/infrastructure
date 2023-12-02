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

  # TODO: Create "/usr/local/bin/kboot" script.
  # #!/bin/bash

  # [[ "$1" != '-' ]] && kernel="$1"
  # shift
  # if [[ "$1" == '-' ]]; then
  #     reuse=--reuse-cmdline
  #     shift
  # fi
  # [[ $# == 0 ]] && reuse=--reuse-cmdline
  # kernel="${kernel:-$(uname -r)}"
  # kargs="/boot/vmlinuz-$kernel --initrd=/boot/initrd.img-$kernel"

  # kexec -l -t bzImage $kargs $reuse --append="$*" && systemctl kexec
  true
}

# Set up cloud-init.
configure_cloud_init() {
  # Install cloud-init.
  apt-get install -y cloud-init

  # TODO: Configure cloud-init data source.
  # Reference: https://forum.armbian.com/topic/14616-cloud-init/

  # TODO: Add cloud-init configuration.
  # TODO: Set up user account
  # TODO: Set up ssh keys
  # TODO: Set up hostname
  true
}

# Configure CPU and memory sets.
configure_cpu_memory_sets() {
  # TODO: Add "cgroup_enable=cpuset cgroup_memory=1 cgroup_enable=memory"
  # to command line. Use "armbianEnv.txt" for this.
  true
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
  # - Banner [ no || none ]
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
