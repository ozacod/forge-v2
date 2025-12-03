# Dependency Management

Manage project dependencies using vcpkg.

## cpx add port

Add a dependency to your project.

```bash
cpx add port <package>
```

This command calls `vcpkg add port <package>` to add the dependency to your `vcpkg.json` manifest.

### Examples

```bash
cpx add port spdlog
cpx add port fmt
cpx add port nlohmann-json
```

## cpx remove

Remove a dependency from your project.

```bash
cpx remove <package>
```

## cpx list

List installed packages.

```bash
cpx list
```

## cpx search

Search for available packages.

```bash
cpx search <query>
```

## Note

All vcpkg commands pass through automatically. You can use:
- `cpx install <package>`
- `cpx list`
- `cpx search <query>`
- And any other vcpkg command

