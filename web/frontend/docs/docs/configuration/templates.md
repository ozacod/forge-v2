# Project Templates

Cpx provides project templates that define the project structure, build configuration, testing framework, and git hooks.

## Using Templates

Templates are automatically downloaded from the GitHub repository when needed.

```bash
# Use default template
cpx create my_project --template default

# Use Catch2 template
cpx create my_project --template catch
```

## Available Templates

### default

The default template uses Google Test framework and includes standard git hooks configuration.

**Configuration:**
```yaml
package:
  version: 0.1.0
  cpp_standard: 17

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: googletest

hooks:
  precommit:
    - fmt
    - lint
  prepush:
    - test
```

### catch

The catch template uses Catch2 test framework. Catch2 is automatically downloaded via FetchContent.

**Configuration:**
```yaml
package:
  version: 0.1.0
  cpp_standard: 17

build:
  shared_libs: false
  clang_format: Google

testing:
  framework: catch2

hooks:
  precommit:
    - fmt
    - lint
  prepush:
    - test
```

## Template Features

- **Automatic Download**: Templates are downloaded from GitHub when needed
- **No Local Storage**: Templates are not stored locally, always fetched from the repository
- **Testing Framework**: Choose between googletest (default) or catch2
- **Git Hooks**: Templates can include pre-configured git hooks
- **Build Configuration**: C++ standard, clang-format style, and library settings

