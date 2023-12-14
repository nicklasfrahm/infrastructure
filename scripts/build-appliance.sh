#!/usr/bin/env bash

set -eou pipefail

# Global constants.
RED='\033[0;31m'
CLEAR='\033[0m'
BUILD_CUSTOMIZATION_DIR="configs/armbian-build"
PATCH_DIR="$BUILD_CUSTOMIZATION_DIR/userpatches"
BUILD_DIR="third_party/armbian-build"
BOARD_NANOPI_R5S_KERNEL_CONFIG_FILE="config/kernel/linux-rockchip64-edge.config"
BOARD_NANOPI_R5S="nanopi-r5s"
USERNAME="nicklasfrahm"

# Global variables.
board=""

cleanup() {
  # Restore the kernel config.
  restore_kernel_config
}

# Restore the kernel config.
restore_kernel_config() {
  if [[ "$board" == "$BOARD_NANOPI_R5S" ]]; then
    # Do surgical reset rather than a coarse reset using
    # "git submodule foreach" and "git reset --hard".
    pushd "$BUILD_DIR" >/dev/null
    git checkout HEAD -- "$BOARD_NANOPI_R5S_KERNEL_CONFIG_FILE"
    git clean -fd
    popd >/dev/null
  fi
}

# Ensure that the clean up function is called on SIGINT.
trap cleanup SIGINT EXIT

parse_args() {
  if [[ $# -ne 1 ]]; then
    echo "usage: $0 <board>"
    exit 1
  fi
  board="$1"

  supported_boards=("$BOARD_NANOPI_R5S")
  is_supported_board=false
  for supported_board in "${supported_boards[@]}"; do
    if [[ "$supported_board" == "$board" ]]; then
      is_supported_board=true
    fi
  done

  if [[ $is_supported_board == false ]]; then
    echo "error: unsupported board: $board"
    exit 1
  fi
}

# Add customizations to the armbian build system.
apply_customizations() {
  # TODO: Enable this again once the following PR is merged:
  # https://github.com/armbian/build/pull/6021
  # Ensure that the build system is up to date.
  # git submodule update --remote third_party/armbian-build

  # Copy the patch files into the build system.
  cp -r "$PATCH_DIR" third_party/armbian-build
}

# Compare the kernel config and install the patched config that enables wireguard.
patch_kernel_config() {
  if [[ "$board" == "$BOARD_NANOPI_R5S" ]]; then
    # Display diff of the kernel config. We expect a diff, so we ignore the exit code.
    kernel_config_file="$BOARD_NANOPI_R5S_KERNEL_CONFIG_FILE"
    diff --color=always -u "$BUILD_DIR/$kernel_config_file" "$BUILD_CUSTOMIZATION_DIR/$kernel_config_file" || true
    cp -f "$BUILD_CUSTOMIZATION_DIR/$kernel_config_file" "$BUILD_DIR/$kernel_config_file"
  fi
}

# Build the firmware image.
build_firmware() {
  cryptroot_enabled=${CRYPTROOT_ENABLE:-false}
  if [[ "$cryptroot_enabled" == true ]]; then
    "./$BUILD_DIR/compile.sh" build \
      BOARD="$board" \
      BRANCH=edge \
      BUILD_DESKTOP=no \
      BUILD_MINIMAL=yes \
      KERNEL_CONFIGURE=no \
      ROOTFS_TYPE=btrfs \
      CRYPTROOT_ENABLE=yes \
      CRYPTROOT_PARAMETERS="--type luks2 --use-random --cipher aes-xts-plain64 --key-size 512 --hash sha512" \
      CRYPTROOT_PASSPHRASE="$USERNAME" \
      CRYPTROOT_SSH_UNLOCK=yes \
      CRYPTROOT_SSH_UNLOCK_PORT=2222 \
      RELEASE=jammy

    show_notes
  fi

  "./$BUILD_DIR/compile.sh" build \
    BOARD="$board" \
    BRANCH=edge \
    BUILD_DESKTOP=no \
    BUILD_MINIMAL=yes \
    KERNEL_CONFIGURE=no \
    ROOTFS_TYPE=btrfs \
    RELEASE=jammy
}

# Move the firmware image to the output directory in the root repo.
move_firmware() {
  version=$(git describe --always --tags --dirty)

  mkdir -p output
  image_file=$(find "$BUILD_DIR/output/images/" -iname "*$board*.img" | sort -rV | head -n1)
  mv "$image_file" "output/${board}-${version}.img"
}

show_notes() {
  # Cryptroot parameters are configured via build parameters above. For more information, see:
  # Reference: https://github.com/armbian/build/commit/681e58b6689acda6a957e325f12e7b748faa8330
  echo
  echo -e "${RED}CRYPTROOT PASSPHRASE MUST BE ROTATED ON FIRST LOGIN:${CLEAR}"
  echo -e "  sudo cryptsetup luksChangeKey /dev/mmcblk0p2"
  echo
}

main() {
  parse_args "$@"

  apply_customizations
  patch_kernel_config

  build_firmware
  move_firmware
}

main "$@"
