# Build Options

Examples of different build configurations.

## Release Build

```bash
cpx build --release
```

## Optimization Levels

```bash
cpx build -O3        # Maximum optimization
cpx build -O2        # Standard release (default)
cpx build -O1        # Light optimization
cpx build -O0        # No optimization (debug)
```

## Parallel Build

```bash
cpx build -j 8       # Use 8 parallel jobs
cpx build -j 4       # Use 4 parallel jobs
```

