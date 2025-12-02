#!/bin/sh
# cpx installer - C++ Project Generator
# Usage: sh -c "$(curl -fsSL https://raw.githubusercontent.com/ozacod/cpx/master/install.sh)"

set -e

REPO="ozacod/cpx"
BINARY_NAME="cpx"

# Determine install directory (prefer user-local, fallback to system with sudo)
get_install_dir() {
    # Check for user-specified directory
    if [ -n "$CPX_INSTALL_DIR" ]; then
        echo "$CPX_INSTALL_DIR"
        return
    fi
    
    # Prefer ~/.local/bin (no sudo needed)
    LOCAL_BIN="$HOME/.local/bin"
    if [ -d "$LOCAL_BIN" ] && [ -w "$LOCAL_BIN" ]; then
        echo "$LOCAL_BIN"
        return
    fi
    
    # Check if /usr/local/bin is writable
    if [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
        return
    fi
    
    # Default to ~/.local/bin (will be created)
    echo "$LOCAL_BIN"
}

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

print_banner() {
    printf "\n"
    printf "%b   ██████╗██████╗ ██╗  ██╗%b\n" "$CYAN" "$NC"
    printf "%b  ██╔════╝██╔══██╗╚██╗██╔╝%b\n" "$CYAN" "$NC"
    printf "%b  ██║     ██████╔╝ ╚███╔╝ %b\n" "$CYAN" "$NC"
    printf "%b  ██║     ██╔═══╝  ██╔██╗ %b\n" "$CYAN" "$NC"
    printf "%b  ╚██████╗██║     ██╔╝ ██╗%b\n" "$CYAN" "$NC"
    printf "%b   ╚═════╝╚═╝     ╚═╝  ╚═╝%b\n" "$CYAN" "$NC"
    printf "\n"
    printf "  %bC++ Project Generator - Cpx Your Code!%b\n" "$YELLOW" "$NC"
    printf "\n"
}

detect_os() {
    OS="$(uname -s)"
    case "$OS" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *)       echo "unknown" ;;
    esac
}

detect_arch() {
    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)            echo "unknown" ;;
    esac
}

check_dependencies() {
    if ! command -v curl > /dev/null 2>&1; then
        if ! command -v wget > /dev/null 2>&1; then
            printf "%bError: curl or wget is required%b\n" "$RED" "$NC"
            exit 1
        fi
        DOWNLOADER="wget"
    else
        DOWNLOADER="curl"
    fi
}

# Check if vcpkg is already installed
check_vcpkg() {
    # Check VCPKG_ROOT environment variable
    if [ -n "$VCPKG_ROOT" ] && [ -f "$VCPKG_ROOT/vcpkg" ] || [ -f "$VCPKG_ROOT/vcpkg.exe" ]; then
        echo "$VCPKG_ROOT"
        return 0
    fi
    
    # Check common installation locations
    COMMON_LOCATIONS="$HOME/vcpkg $HOME/.local/vcpkg $HOME/.vcpkg /opt/vcpkg /usr/local/vcpkg"
    for loc in $COMMON_LOCATIONS; do
        if [ -f "$loc/vcpkg" ] || [ -f "$loc/vcpkg.exe" ]; then
            echo "$loc"
            return 0
        fi
    done
    
    return 1
}

# Install vcpkg
install_vcpkg() {
    printf "\n%bChecking for vcpkg...%b\n" "$CYAN" "$NC"
    
    # Check if vcpkg is already installed
    VCPKG_PATH=$(check_vcpkg)
    if [ -n "$VCPKG_PATH" ]; then
        printf "%bvcpkg found at: %s%b\n" "$GREEN" "$VCPKG_PATH" "$NC"
        configure_vcpkg "$VCPKG_PATH"
        return 0
    fi
    
    # Check if git is available (needed to clone vcpkg)
    if ! command -v git > /dev/null 2>&1; then
        printf "%bWarning: git is not installed. Skipping vcpkg installation.%b\n" "$YELLOW" "$NC"
        printf "You can install vcpkg manually and configure it with: %bcpx config set-vcpkg-root <path>%b\n" "$CYAN" "$NC"
        return 1
    fi
    
    # Ask user if they want to install vcpkg (skip in non-interactive mode)
    if [ ! -t 0 ]; then
        # Non-interactive mode - skip installation
        printf "%bSkipping vcpkg installation (non-interactive mode).%b\n" "$YELLOW" "$NC"
        printf "You can install it later and configure with: %bcpx config set-vcpkg-root <path>%b\n" "$CYAN" "$NC"
        return 1
    fi
    
    printf "%bvcpkg not found. Would you like to install it? (y/n): %b" "$YELLOW" "$NC"
    read -r response
    case "$response" in
        [yY]|[yY][eE][sS])
            ;;
        *)
            printf "%bSkipping vcpkg installation.%b\n" "$YELLOW" "$NC"
            printf "You can install it later and configure with: %bcpx config set-vcpkg-root <path>%b\n" "$CYAN" "$NC"
            return 1
            ;;
    esac
    
    # Determine vcpkg installation directory
    VCPKG_INSTALL_DIR="$HOME/.local/vcpkg"
    if [ -n "$CPX_VCPKG_DIR" ]; then
        VCPKG_INSTALL_DIR="$CPX_VCPKG_DIR"
    fi
    
    printf "%bInstalling vcpkg to %s...%b\n" "$CYAN" "$VCPKG_INSTALL_DIR" "$NC"
    
    # Remove existing directory if it exists (incomplete installation)
    if [ -d "$VCPKG_INSTALL_DIR" ]; then
        printf "%bRemoving existing directory...%b\n" "$YELLOW" "$NC"
        rm -rf "$VCPKG_INSTALL_DIR"
    fi
    
    # Clone vcpkg
    printf "%bCloning vcpkg from GitHub...%b\n" "$CYAN" "$NC"
    if ! git clone https://github.com/microsoft/vcpkg.git "$VCPKG_INSTALL_DIR"; then
        printf "%bError: Failed to clone vcpkg%b\n" "$RED" "$NC"
        return 1
    fi
    
    # Bootstrap vcpkg
    printf "%bBootstrapping vcpkg...%b\n" "$CYAN" "$NC"
    cd "$VCPKG_INSTALL_DIR" || return 1
    
    OS=$(detect_os)
    if [ "$OS" = "windows" ]; then
        if ! ./bootstrap-vcpkg.bat; then
            printf "%bError: Failed to bootstrap vcpkg%b\n" "$RED" "$NC"
            cd - > /dev/null || true
            return 1
        fi
    else
        if ! ./bootstrap-vcpkg.sh; then
            printf "%bError: Failed to bootstrap vcpkg%b\n" "$RED" "$NC"
            cd - > /dev/null || true
            return 1
        fi
    fi
    
    cd - > /dev/null || true
    
    printf "%bSuccessfully installed vcpkg to %s%b\n" "$GREEN" "$VCPKG_INSTALL_DIR" "$NC"
    
    # Configure cpx to use vcpkg
    configure_vcpkg "$VCPKG_INSTALL_DIR"
}

# Configure cpx to use vcpkg
configure_vcpkg() {
    VCPKG_PATH=$1
    
    # Check if cpx is in PATH
    if ! command -v "$BINARY_NAME" > /dev/null 2>&1; then
        INSTALL_DIR=$(get_install_dir)
        CPX_BINARY="$INSTALL_DIR/$BINARY_NAME"
    else
        CPX_BINARY=$(command -v "$BINARY_NAME")
    fi
    
    # Check if cpx binary exists
    if [ ! -f "$CPX_BINARY" ]; then
        printf "%bWarning: cpx binary not found. Cannot configure vcpkg automatically.%b\n" "$YELLOW" "$NC"
        printf "Run this after cpx is in your PATH: %bcpx config set-vcpkg-root %s%b\n" "$CYAN" "$VCPKG_PATH" "$NC"
        return 1
    fi
    
    # Configure vcpkg root
    printf "%bConfiguring cpx to use vcpkg...%b\n" "$CYAN" "$NC"
    if "$CPX_BINARY" config set-vcpkg-root "$VCPKG_PATH" 2>/dev/null; then
        printf "%bSuccessfully configured cpx to use vcpkg%b\n" "$GREEN" "$NC"
    else
        printf "%bWarning: Failed to configure cpx automatically.%b\n" "$YELLOW" "$NC"
        printf "Run this manually: %bcpx config set-vcpkg-root %s%b\n" "$CYAN" "$VCPKG_PATH" "$NC"
    fi
}

get_latest_version() {
    if [ "$DOWNLOADER" = "curl" ]; then
        VERSION=$(curl -sI "https://github.com/$REPO/releases/latest" | grep -i "location:" | sed 's/.*tag\///' | tr -d '\r\n')
    else
        VERSION=$(wget -qO- --server-response "https://github.com/$REPO/releases/latest" 2>&1 | grep -i "location:" | sed 's/.*tag\///' | tr -d '\r\n')
    fi
    
    if [ -z "$VERSION" ]; then
        VERSION="v1.0.2"
    fi
    echo "$VERSION"
}

download_binary() {
    OS=$1
    ARCH=$2
    VERSION=$3
    
    if [ "$OS" = "unknown" ] || [ "$ARCH" = "unknown" ]; then
        printf "%bError: Unsupported platform: %s/%s%b\n" "$RED" "$OS" "$ARCH" "$NC"
        exit 1
    fi
    
    BINARY_NAME_PLATFORM="$BINARY_NAME-$OS-$ARCH"
    if [ "$OS" = "windows" ]; then
        BINARY_NAME_PLATFORM="${BINARY_NAME_PLATFORM}.exe"
    fi
    
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_NAME_PLATFORM"
    
    printf "%bDownloading %s from %s...%b\n" "$CYAN" "$BINARY_NAME_PLATFORM" "$DOWNLOAD_URL" "$NC"
    
    INSTALL_DIR=$(get_install_dir)
    TARGET_PATH="$INSTALL_DIR/$BINARY_NAME"
    
    # Create install directory if it doesn't exist
    mkdir -p "$INSTALL_DIR"
    
    if [ "$DOWNLOADER" = "curl" ]; then
        if ! curl -fSL "$DOWNLOAD_URL" -o "$TARGET_PATH"; then
            printf "%bError: Failed to download binary%b\n" "$RED" "$NC"
            exit 1
        fi
    else
        if ! wget -q "$DOWNLOAD_URL" -O "$TARGET_PATH"; then
            printf "%bError: Failed to download binary%b\n" "$RED" "$NC"
            exit 1
        fi
    fi
    
    chmod +x "$TARGET_PATH"
    
    printf "%bSuccessfully installed %s to %s%b\n" "$GREEN" "$BINARY_NAME" "$TARGET_PATH" "$NC"
    
    # Check if binary is in PATH
    if ! command -v "$BINARY_NAME" > /dev/null 2>&1; then
        printf "%bWarning: %s is not in your PATH.%b\n" "$YELLOW" "$BINARY_NAME" "$NC"
        printf "Add this to your shell profile (.bashrc, .zshrc, etc.):\n"
        printf "  export PATH=\"\$PATH:%s\"\n" "$INSTALL_DIR"
    else
        printf "%b%s is ready to use!%b\n" "$GREEN" "$BINARY_NAME" "$NC"
        printf "Run '%b%s version%b' to verify installation.\n" "$CYAN" "$BINARY_NAME" "$NC"
    fi
}

main() {
    print_banner
    
    check_dependencies
    
    OS=$(detect_os)
    ARCH=$(detect_arch)
    VERSION=$(get_latest_version)
    
    printf "%bDetected: %s/%s%b\n" "$CYAN" "$OS" "$ARCH" "$NC"
    printf "%bLatest version: %s%b\n" "$CYAN" "$VERSION" "$NC"
    printf "\n"
    
    download_binary "$OS" "$ARCH" "$VERSION"
    
    # Try to install/configure vcpkg (non-fatal if it fails)
    # Skip on Windows unless in Git Bash/MSYS2
    if [ "$OS" != "windows" ] || [ -n "$MSYSTEM" ]; then
        install_vcpkg || true
    else
        printf "\n%bNote: vcpkg installation on Windows requires manual setup.%b\n" "$YELLOW" "$NC"
        printf "After installing vcpkg, run: %bcpx config set-vcpkg-root <path>%b\n" "$CYAN" "$NC"
    fi
}

main
