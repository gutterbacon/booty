#!/usr/bin/env bash

BIN_NAME=booty
EX_CURL="$(command -v 'curl' &>/dev/null && echo true || echo false)"
EX_GPG="$(command -v 'gpg' &>/dev/null && echo true || echo false)"
OS="$(uname -s | tr [:upper:] [:lower:])"
OS_ARCH=""
RELEASE_URL="https://github.com/amplify-edge/booty/releases"
LATEST_RELEASE_URL="${RELEASE_URL}/latest"
INSTALL_LOC=/usr/local/bin

# checks if we're running as root
runAsRoot() {
  if [ $EUID -ne 0 ]; then
    sudo "${@}"
  else
    "${@}"
  fi
}

osArch() {
  ARCH="$(uname -m)"
  case "$ARCH" in
    x86-64) ARCH="amd64" ;;
    aarch64) ARCH="arm64";;
  esac
  OS_ARCH="$(echo ${OS}_${ARCH})"
}

validOsArchCombo() {
  case "$OS_ARCH" in:
    linux_amd64) return 1 ;;
    darwin_amd64) return 1;;
    darwin_arm64) return 1;;
    *) return 0;;
  esac
}

latestVersion() {
  curl -sL "${LATEST_RELEASE_URL}" |grep "Release" | head -n 1 | awk '{n=split($2,a,"·");print a[n]}'
}

fetchAndInstall() {
  local tag_format="${BIN_NAME}-${TAG}-${OS_ARCH}.tar.gz"
  # download
  curl -L -o "/tmp/${tag_format}" "${RELEASE_URL}/download/${TAG}/${tag_format}"
  # extract
  tar -zxvf "/tmp/${tag_format}"
  # install
  install -m755 "/tmp/${BIN_NAME}" "${INSTALL_LOC}/${BIN_NAME}"
}

if [ !$EX_CURL ];
} then
  echo "please install curl to proceed, exiting..."
  exit 1
fi
runAsRoot
osArch
if [!validOsArchCombo]; then
  echo "unsupported os and arch ${OS_ARCH}"
  exit 1
fi
if [$(latestVersion) -eq "" ]; then
  echo "latest version not found, exiting..."
  exit 1
fi
TAG=$(latestVersion)
fetchAndInstall