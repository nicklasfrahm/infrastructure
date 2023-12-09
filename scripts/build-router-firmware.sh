#!/usr/bin/env bash

set -eou pipefail

# Global constants.
RED='\033[0;31m'
CLEAR='\033[0m'
BUILD_CUSTOMIZATION_DIR="configs/armbian-build"
PATCH_DIR="$BUILD_CUSTOMIZATION_DIR/userpatches"
BUILD_DIR="third_party/armbian-build"
KERNEL_CONFIG_FILE="config/kernel/linux-rk3568-odroid-edge.config"

cleanup() {
  # Restore the kernel config.
  restore_kernel_config
}

# Restore the kernel config.
restore_kernel_config() {
  # Do surgical reset rather than a coarse reset using
  # "git submodule foreach" and "git reset --hard".
  pushd "$BUILD_DIR" >/dev/null
  git checkout HEAD -- "$KERNEL_CONFIG_FILE"
  git clean -fd
  popd >/dev/null
}

# Ensure that the clean up function is called on SIGINT.
trap cleanup SIGINT EXIT

# Prepare build system.
setup_toolchain() {
  # TODO: Enable this again once the following PR is merged:
  # https://github.com/armbian/build/pull/6021
  # Ensure that the build system is up to date.
  # git submodule update --remote third_party/armbian-build

  # Copy the patch files into the build system.
  cp -r "$PATCH_DIR" third_party/armbian-build
}

# Compare the kernel config and install the patched config that enables wireguard.
update_kernel_config() {
  # Display diff of the kernel config. We expect a diff, so we ignore the exit code.
  diff --color=always -u "$BUILD_DIR/$KERNEL_CONFIG_FILE" "$BUILD_CUSTOMIZATION_DIR/$KERNEL_CONFIG_FILE" || true
  cp "$BUILD_CUSTOMIZATION_DIR/$KERNEL_CONFIG_FILE" "$BUILD_DIR/$KERNEL_CONFIG_FILE"
}

# Build the firmware image.
build_firmware() {
  "./$BUILD_DIR/compile.sh" build \
    BOARD=nanopi-r5s \
    BRANCH=edge \
    BUILD_DESKTOP=no \
    BUILD_MINIMAL=no \
    KERNEL_CONFIGURE=no \
    ROOTFS_TYPE=btrfs \
    CRYPTROOT_ENABLE=yes \
    CRYPTROOT_PARAMETERS="--type luks2 --use-random --cipher aes-xts-plain64 --key-size 512 --hash sha512" \
    CRYPTROOT_PASSPHRASE="nicklasfrahm" \
    CRYPTROOT_SSH_UNLOCK=yes \
    CRYPTROOT_SSH_UNLOCK_PORT=2222 \
    RELEASE=jammy

  # Cryptroot parameters are configured via build parameters above. For more information, see:
  # Reference: https://github.com/armbian/build/commit/681e58b6689acda6a957e325f12e7b748faa8330
  echo
  echo -e "${RED}CRYPTROOT PASSPHRASE MUST BE ROTATED ON FIRST LOGIN:${CLEAR}"
  echo -e "  sudo cryptsetup luksChangeKey /dev/mmcblk0p2"
  echo
}

main() {
  setup_toolchain

  update_kernel_config

  build_firmware
}

main "$@"
