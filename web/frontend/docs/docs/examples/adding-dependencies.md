# Adding Dependencies

Examples of adding dependencies to your project.

## Using vcpkg Commands

```bash
# Add dependencies directly
cpx add port spdlog
cpx add port fmt
cpx add port nlohmann-json

# Or use vcpkg commands directly
cpx install spdlog
cpx list
```

## Manual vcpkg.json Edit

You can also edit `vcpkg.json` directly:

```json
{
  "dependencies": [
    "spdlog",
    "fmt",
    "nlohmann-json"
  ]
}
```

