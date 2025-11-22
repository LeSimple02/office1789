# Office1789 Electron Desktop App

Application desktop Office1789 avec Electron.

## Installation

```bash
npm install
```

## Développement

```bash
# Build Vue app + copy to www + launch Electron
npm run electron:dev
```

## Build Production

```bash
# Create distributable packages (Windows/Mac/Linux)
npm run electron:build
```

Les packages seront générés dans le dossier `dist/`.

## Structure

```
electron-app/
├── main.js          # Electron main process
├── preload.js       # Preload script for security
├── copy-dist.js     # Script to copy Vue build
├── www/             # Vue app files (copied from ../webfront2/dist)
└── package.json     # Electron config
```
