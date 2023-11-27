#!/usr/bin/env bash

set -eou pipefail

# Global constants.
BUILD_CUSTOMIZATION_DIR="configs/armbian-build"
PATCH_DIR="$BUILD_CUSTOMIZATION_DIR/userpatches"
BUILD_DIR="third_party/armbian-build"

# Prepare build system.
setup_toolchain() {
  # Ensure that the build system is up to date.
  git submodule update third_party/armbian-build

  # Copy the patch files into the build system.
  cp -r "$PATCH_DIR" third_party/armbian-build/userpatches
}

# Compare the kernel config.
diff_kernel_config() {
  KERNEL_CONFIG_FILE="config/kernel/linux-rk3568-odroid-edge.config"
  diff -u "$BUILD_CUSTOMIZATION_DIR/$KERNEL_CONFIG_FILE" "$BUILD_DIR/$KERNEL_CONFIG_FILE"
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
  diff_kernel_config
  build_firmware
}

main "$@"
