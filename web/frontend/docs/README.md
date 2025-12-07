# Cpx Documentation

This directory contains the Docusaurus documentation for Cpx.

## Setup

1. Install dependencies:
```bash
bun install
```

2. Start the docs dev server:
```bash
bun run start
```

3. Build for production:
```bash
bun run build
```

## Integration with Frontend

After building, copy the output to the frontend public directory:
```bash
bun run build
cp -r build/* ../public/docs/
```

Or use the automated script from the frontend directory:
```bash
cd ..
bun run build:docs
```

## Notes
- Only `bun.lockb` is kept; use Bun for reproducible installs.
- Use `bunx docusaurus <command>` (or `npx`) for one-off Docusaurus CLI tasks outside of package scripts.

