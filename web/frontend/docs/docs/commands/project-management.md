# Project Management Commands

Commands for creating and managing C++ projects.

## cpx create

Create a new C++ project.

```bash
cpx create <name>
```

### Options

- `--template <name>` - Create project from template (default, catch, or path to .yaml file)
- `--lib` - Create library project instead of executable

### Examples

```bash
# Create executable project
cpx create my_app

# Create library project
cpx create my_lib --lib

# Create from template
cpx create my_project --template default
cpx create my_project --template catch
```

