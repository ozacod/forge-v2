# Cpx Documentation

This directory contains the Docusaurus documentation for Cpx.

## Setup

1. Install dependencies:
```bash
npm install
```

2. Start development server:
```bash
npm start
```

3. Build for production:
```bash
npm run build
```

## Integration with Frontend

After building, copy the build output to the frontend public directory:

```bash
npm run build
cp -r build/* ../public/docs/
```

Or use the automated script from the frontend directory:

```bash
cd ..
npm run build:docs
```

