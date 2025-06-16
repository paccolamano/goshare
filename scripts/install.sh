#!/bin/bash
set -eu -o pipefail

ARCHITECTURE=""
case $(uname -m) in
    x86_64)                     ARCHITECTURE="amd64" ;;
    arm64)                      ARCHITECTURE="arm64" ;;
    ppc64le)                    ARCHITECTURE="ppc64le" ;;
    s390x)                      ARCHITECTURE="s390x" ;;
    arm|armv7l|armv8l|aarch64)  dpkg --print-architecture | grep -q "arm64" && ARCHITECTURE="arm64" || ARCHITECTURE="arm" ;;
esac

INSTALL_OS=""
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     INSTALL_OS=linux;;
    Darwin*)    INSTALL_OS=darwin;;
esac

for product in $*; do
  if ! command -v $product > /dev/null; then
    read -p "${product} is not installed on your machine. Do you want to install it? [Y/n] " choice;
    if [ "$choice" != "n" ] && [ "$choice" != "N" ]; then
      ARCHITECTURE=$ARCHITECTURE INSTALL_OS=$INSTALL_OS "$(dirname $0)/installers/install-${product}.sh"
    else
      echo "You chose not to install ${product}. Exiting...";
      exit 1;
    fi;
  fi
done
