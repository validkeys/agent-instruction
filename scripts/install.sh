#!/usr/bin/env bash
#
# Install script for agent-instruction CLI
#
# Usage:
#   ./scripts/install.sh [INSTALL_DIR]
#
# Default INSTALL_DIR is /usr/local/bin

set -euo pipefail

BINARY_NAME="agent-instruction"
DEFAULT_INSTALL_DIR="/usr/local/bin"
INSTALL_DIR="${1:-$DEFAULT_INSTALL_DIR}"
BUILD_DIR="build"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

error() {
    echo -e "${RED}Error: $1${NC}" >&2
    exit 1
}

info() {
    echo -e "${GREEN}$1${NC}"
}

warn() {
    echo -e "${YELLOW}$1${NC}"
}

# Check if binary exists
if [ ! -f "${BUILD_DIR}/${BINARY_NAME}" ]; then
    error "Binary not found at ${BUILD_DIR}/${BINARY_NAME}. Run 'make build' first."
fi

# Check if install directory exists
if [ ! -d "$INSTALL_DIR" ]; then
    error "Install directory does not exist: $INSTALL_DIR"
fi

# Check write permissions
if [ ! -w "$INSTALL_DIR" ]; then
    warn "No write permission to $INSTALL_DIR. You may need to run with sudo."
    warn "Retrying with sudo..."
    sudo install -m 755 "${BUILD_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    install -m 755 "${BUILD_DIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

info "✓ Installed ${BINARY_NAME} to ${INSTALL_DIR}/${BINARY_NAME}"
info "Run '${BINARY_NAME}' to use."

# Verify installation
if command -v "${BINARY_NAME}" &> /dev/null; then
    VERSION=$("${BINARY_NAME}" --version)
    info "Installation verified: $VERSION"
else
    warn "Installation complete, but '${BINARY_NAME}' is not in your PATH."
    warn "Add ${INSTALL_DIR} to your PATH or use the full path: ${INSTALL_DIR}/${BINARY_NAME}"
fi
