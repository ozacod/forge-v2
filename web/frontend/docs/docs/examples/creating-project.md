# Creating a Project

Examples of creating different types of projects.

## Basic Executable

```bash
# Create a simple executable project
cpx create my_app
cd my_app
cpx build
cpx run
```

## Library Project

```bash
# Create a library project
cpx create my_lib --lib
cd my_lib
cpx build
```

## From Template

```bash
# Create from default template (googletest)
cpx create my_project --template default

# Create with Catch2 template
cpx create my_project --template catch

cd my_project
cpx build
```

