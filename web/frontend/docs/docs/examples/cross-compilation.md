# Cross-Compilation

Examples of cross-compiling for multiple platforms.

## Build for All Targets

```bash
# Build for all targets in cpx.ci
cpx ci
```

## Build Specific Target

```bash
# Build only for linux-amd64
cpx ci --target linux-amd64
```

## Rebuild Docker Images

```bash
# Force rebuild of Docker images
cpx ci --rebuild
```

## Generate GitHub Actions Workflow

```bash
# Create .github/workflows/ci.yml
cpx ci init --github-actions
```

## Generate GitLab CI Configuration

```bash
# Create .gitlab-ci.yml
cpx ci init --gitlab
```

