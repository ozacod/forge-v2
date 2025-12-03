# CI Commands

Commands for CI/CD integration and cross-compilation.

## cpx ci

Build for all targets defined in `cpx.ci`.

```bash
cpx ci
```

### Options

- `--target <name>` - Build only specific target
- `--rebuild` - Rebuild Docker images even if they exist

## cpx ci init --github-actions

Generate GitHub Actions workflow file (`.github/workflows/ci.yml`).

```bash
cpx ci init --github-actions
```

## cpx ci init --gitlab

Generate GitLab CI configuration file (`.gitlab-ci.yml`).

```bash
cpx ci init --gitlab
```

These commands create CI workflow files that call `cpx ci` automatically.

