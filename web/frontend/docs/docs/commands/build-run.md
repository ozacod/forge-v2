# Build & Run Commands

Commands for building and running your C++ projects.

## cpx build

Compile the project using CMake presets if available.

```bash
cpx build
```

### Options

- `--release` - Build in release mode
- `-O<level>` - Optimization level: 0, 1, 2, 3, s, fast
- `--clean` - Clean and rebuild
- `-j <n>` - Use n parallel jobs

### Examples

```bash
# Debug build
cpx build

# Release build
cpx build --release

# Maximum optimization
cpx build -O3

# Parallel build
cpx build -j 8
```

## cpx run

Build and run the executable.

```bash
cpx run
```

### Options

- `--release` - Run in release mode

## cpx test

Build and run tests.

```bash
cpx test
```

### Options

- `-v, --verbose` - Verbose test output
- `--filter <name>` - Filter tests by name

## cpx clean

Remove build artifacts.

```bash
cpx clean
```

### Options

- `--all` - Also remove generated files

