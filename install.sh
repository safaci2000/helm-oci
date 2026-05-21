#!/usr/bin/env sh

PROJECT_NAME="helm-oci"
PROJECT_GH="esnet/$PROJECT_NAME"

if command -v cygpath >/dev/null 2>&1; then
  HELM_BIN="$(cygpath -u "${HELM_BIN}")"
  HELM_PLUGIN_DIR="$(cygpath -u "${HELM_PLUGIN_DIR}")"
fi

[ -z "$HELM_BIN" ] && HELM_BIN=$(command -v helm)
[ -z "$HELM_HOME" ] && HELM_HOME=$(helm env | grep 'HELM_DATA_HOME' | cut -d '=' -f2 | tr -d '"')

mkdir -p "$HELM_HOME"
: "${HELM_PLUGIN_DIR:="$HELM_HOME/plugins/helm-oci"}"

if [ "$SKIP_BIN_INSTALL" = "1" ]; then
  echo "Skipping binary install"
  exit
fi

SCRIPT_MODE="install"
if [ "$1" = "-u" ]; then
  SCRIPT_MODE="update"
fi

initArch() {
  ARCH=$(uname -m)
  case $ARCH in
  aarch64) ARCH="arm64" ;;
  x86_64) ARCH="amd64" ;;
  i686) ARCH="amd64" ;;
  i386) ARCH="amd64" ;;
  esac
}

initOS() {
  OS=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$OS" in
  msys*|mingw*|cygwin*) OS='windows' ;;
  esac
}

getDownloadURL() {
  version=$(git -C "$HELM_PLUGIN_DIR" describe --tags --exact-match 2>/dev/null || :)
  if [ "$SCRIPT_MODE" = "install" ] && [ -n "$version" ]; then
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/download/$version/helm-oci-$OS-$ARCH.tgz"
  else
    DOWNLOAD_URL="https://github.com/$PROJECT_GH/releases/latest/download/helm-oci-$OS-$ARCH.tgz"
  fi
}

mkTempDir() {
  HELM_TMP="$(mktemp -d -t "${PROJECT_NAME}-XXXXXX")"
}

rmTempDir() {
  if [ -d "${HELM_TMP:-/dev/null}" ]; then
    rm -rf "${HELM_TMP}"
  fi
}

downloadFile() {
  PLUGIN_TMP_FILE="${HELM_TMP}/${PROJECT_NAME}.tgz"
  echo "Downloading $DOWNLOAD_URL"
  if command -v curl >/dev/null 2>&1; then
    curl -sSf -L "$DOWNLOAD_URL" >"$PLUGIN_TMP_FILE"
  elif command -v wget >/dev/null 2>&1; then
    wget -q -O - "$DOWNLOAD_URL" >"$PLUGIN_TMP_FILE"
  else
    echo "Either curl or wget is required"
    exit 1
  fi
}

installFile() {
  tar xzf "$PLUGIN_TMP_FILE" -C "$HELM_TMP"
  echo "Preparing to install into ${HELM_PLUGIN_DIR}"
  mkdir -p "$HELM_PLUGIN_DIR/bin"
  cp "$HELM_TMP/helm-oci/bin/helm-oci" "$HELM_PLUGIN_DIR/bin/"
}

exit_trap() {
  result=$?
  rmTempDir
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    printf '\tFor support, go to https://github.com/%s.\n' "$PROJECT_GH"
  fi
  exit $result
}

trap "exit_trap" EXIT
set -e
initArch
initOS
getDownloadURL
mkTempDir
downloadFile
installFile
