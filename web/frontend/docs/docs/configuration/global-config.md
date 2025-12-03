# Global Configuration

Cpx stores its global configuration in a YAML file.

## Configuration Location

- **Linux/macOS**: `~/.config/cpx/config.yaml`
- **Windows**: `%APPDATA%/cpx/config.yaml`

## Configuration File

The `config.yaml` file contains:

```yaml
vcpkg_root: "/path/to/vcpkg"
```

## Setting Configuration

Use the config command to set values:

```bash
cpx config set-vcpkg-root /path/to/vcpkg
```

## Getting Configuration

Retrieve configuration values:

```bash
cpx config get-vcpkg-root
```

