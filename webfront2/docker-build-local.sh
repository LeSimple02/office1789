#!/bin/bash
# Build local pour webfront2 : utilise .env.local comme .env.production
set -e
cd "$(dirname "$0")"
cp .env.local .env.production
npm install
npm run build
