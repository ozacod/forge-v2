# Generated Project Structure

When you create a new project with Cpx, the following structure is generated:

```
my_project/
├── CMakeLists.txt           # Main CMake file with vcpkg integration
├── CMakePresets.json        # CMake presets (safe to commit, uses env vars)
├── vcpkg.json               # vcpkg manifest (dependencies)
├── vcpkg-configuration.json # vcpkg configuration (auto-generated)
├── cpx.ci                   # Cross-compilation configuration
├── include/
│   └── my_project/
│       ├── my_project.hpp
│       └── version.hpp
├── src/
│   ├── main.cpp             # Main executable (if exe)
│   └── my_project.cpp
├── tests/
│   ├── CMakeLists.txt
│   └── test_main.cpp
├── .gitignore
├── .clang-format
└── README.md
```

## Directory Structure

- **include/** - Header files organized by project name
- **src/** - Source files including main.cpp for executables
- **tests/** - Test files with CMakeLists.txt for test framework
- **build/** - Build directory (gitignored)

## Key Files

- **CMakeLists.txt** - Main build configuration with vcpkg integration
- **CMakePresets.json** - IDE integration presets
- **vcpkg.json** - Dependency manifest
- **cpx.ci** - Cross-compilation configuration
- **.clang-format** - Code formatting configuration

