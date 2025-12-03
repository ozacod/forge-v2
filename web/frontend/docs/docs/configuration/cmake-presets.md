# CMake Presets

Cpx generates `CMakePresets.json` for seamless IDE integration.

## CMakePresets.json

The generated presets file includes:

- Uses environment variables (`$env{VCPKG_ROOT}`)
- Safe to commit to version control
- `VCPKG_ROOT` is automatically set by cpx build commands
- Works seamlessly with IDEs like VS Code, CLion, and Qt Creator

## Example

```json
{
  "version": 3,
  "cmakeMinimumRequired": {
    "major": 3,
    "minor": 20,
    "patch": 0
  },
  "configurePresets": [
    {
      "name": "default",
      "displayName": "Default Config",
      "generator": "Unix Makefiles",
      "binaryDir": "${sourceDir}/build",
      "cacheVariables": {
        "CMAKE_TOOLCHAIN_FILE": "$env{VCPKG_ROOT}/scripts/buildsystems/vcpkg.cmake"
      }
    }
  ]
}
```

## IDE Integration

Once generated, CMake presets are automatically detected by:
- **VS Code** - CMake Tools extension
- **CLion** - Native CMake support
- **Qt Creator** - CMake project support

