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

if [ $# -ne 4 ]; then
  echo "error: this script is meant to be called by the armbian build system"
  exit 1
fi
