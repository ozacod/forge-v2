# Configuration Commands

Commands for managing Cpx configuration.

## cpx config set-vcpkg-root

Set the vcpkg installation directory.

```bash
cpx config set-vcpkg-root <path>
```

## cpx config get-vcpkg-root

Get the current vcpkg root directory.

```bash
cpx config get-vcpkg-root
```

## cpx hooks install

Install git hooks based on cpx.yaml configuration.

```bash
cpx hooks install
```

This command:
- Creates hooks for precommit/prepush if configured in cpx.yaml
- Creates .sample files if hooks are not configured
- Uses defaults if no cpx.yaml exists

