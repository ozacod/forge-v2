# Project Configuration

Project-specific configuration files and their purposes.

## vcpkg.json

The `vcpkg.json` file is the vcpkg manifest that lists all dependencies. This file is auto-generated and managed via `cpx add port` commands.

```json
{
  "dependencies": [
    "spdlog",
    "fmt",
    "nlohmann-json"
  ]
}
```

## cpx.yaml

The `cpx.yaml` file is used as a template for project creation. It defines the project structure, build configuration, testing framework, and git hooks.

```yaml
package:
  name: my_project
  version: "0.1.0"
  cpp_standard: 17
  project_type: exe

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

**Note**: Dependencies are managed in `vcpkg.json`, not `cpx.yaml`. The `cpx.yaml` file is only used as a template for project creation.

