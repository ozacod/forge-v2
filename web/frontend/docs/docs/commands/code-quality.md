# Code Quality Commands

Tools for maintaining code quality and security.

## cpx fmt

Format code with clang-format.

```bash
cpx fmt
```

### Options

- `--check` - Check formatting without modifying files

## cpx lint

Run clang-tidy static analysis.

```bash
cpx lint
```

### Options

- `--fix` - Auto-fix lint issues

## cpx flawfinder

Run Flawfinder security analysis for C/C++.

```bash
cpx flawfinder
```

### Options

- `--minlevel <0-5>` - Minimum risk level to report (default: 1)
- `--html` - Output results in HTML format
- `--csv` - Output results in CSV format
- `--output <file>` - Output file path (required for HTML/CSV)
- `--dataflow` - Enable dataflow analysis
- `--quiet` - Quiet mode (minimal output)
- `--context <n>` - Number of lines of context to show (default: 2)

## cpx cppcheck

Run Cppcheck static analysis for C/C++.

```bash
cpx cppcheck
```

### Options

- `--enable <checks>` - Enable checks (all, style, performance, portability, etc.)
- `--xml` - Output results in XML format
- `--csv` - Output results in CSV format
- `--output <file>` - Output file path (for XML/CSV output)
- `--quiet` - Quiet mode (suppress progress messages)
- `--force` - Force checking of all configurations
- `--platform <name>` - Target platform (unix32, unix64, win32A, win64, etc.)
- `--std <standard>` - C/C++ standard (c++17, c++20, etc.)

## cpx check

Check code compiles with sanitizers.

```bash
cpx check
```

### Options

- `--asan` - Build with AddressSanitizer (detects memory errors)
- `--tsan` - Build with ThreadSanitizer (detects data races)
- `--msan` - Build with MemorySanitizer (detects uninitialized memory)
- `--ubsan` - Build with UndefinedBehaviorSanitizer (detects undefined behavior)

