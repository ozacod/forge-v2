# Installation

Install Cpx on your system using one of the methods below.

## Quick Install (Recommended)

The quickest way to install Cpx is using our install script:

```bash
curl -f https://raw.githubusercontent.com/ozacod/cpx/master/install.sh | sh
```

This script will:
- Detect your OS and architecture
- Download the latest Cpx binary
- Set up vcpkg
- Configure Cpx with the vcpkg root directory

## Manual Installation

### 1. Download the Binary

Visit the [GitHub releases page](https://github.com/ozacod/cpx/releases/latest) and download the binary for your platform.

### 2. Make it Executable and Move to PATH

```bash
chmod +x cpx-linux-amd64
sudo mv cpx-linux-amd64 /usr/local/bin/cpx
```

### 3. Configure vcpkg

```bash
cpx config set-vcpkg-root /path/to/vcpkg
```

## Verify Installation

```bash
cpx version
```

You should see the Cpx version number if installation was successful.

