#!/usr/bin/env bash

set -eou pipefail

# Global constants.
BUILD_CUSTOMIZATION_DIR="configs/armbian-build"
PATCH_DIR="$BUILD_CUSTOMIZATION_DIR/userpatches"
BUILD_DIR="third_party/armbian-build"
KERNEL_CONFIG_FILE="config/kernel/linux-rk3568-odroid-edge.config"

# Prepare build system.
setup_toolchain() {
  # Ensure that the build system is up to date.
  git submodule update third_party/armbian-build

  # Copy the patch files into the build system.
  cp -r "$PATCH_DIR" third_party/armbian-build/userpatches
}

# Compare the kernel config and install the patched config that enables wireguard.
update_kernel_config() {
  # Display diff of the kernel config. We expect a diff, so we ignore the exit code.
  diff --color=always -u "$BUILD_CUSTOMIZATION_DIR/$KERNEL_CONFIG_FILE" "$BUILD_DIR/$KERNEL_CONFIG_FILE" || true
  cp "$BUILD_CUSTOMIZATION_DIR/$KERNEL_CONFIG_FILE" "$BUILD_DIR/$KERNEL_CONFIG_FILE"
}

# Restore the kernel config.
restore_kernel_config() {
  # Do surgical reset rather than a coarse reset using
  # "git submodule foreach" and "git reset --hard".
  cd "$BUILD_DIR"
  git checkout HEAD -- "$KERNEL_CONFIG_FILE"
  cd - >/dev/null
}

# Build the firmware image.
build_firmware() {
  # ./compile.sh build \
  #   BOARD=nanopi-r5s \
  #   BRANCH=edge \
  #   BUILD_DESKTOP=no \
  #   BUILD_MINIMAL=no \
  #   KERNEL_CONFIGURE=no \
  #   RELEASE=jammy
  true
}

main() {
  setup_toolchain

  update_kernel_config

  build_firmware

  restore_kernel_config
}

main "$@"
