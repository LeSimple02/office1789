#!/bin/bash
set -e

REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

echo "=========================================="
echo "🛠  Environnement local Office1789"
echo "=========================================="
echo ""

cd "$REPO_DIR"

echo "📦 Vérification dépendances (Go, Node, npm)..."
command -v go >/dev/null 2>&1 || { echo "❌ go manquant"; exit 1; }
command -v npm >/dev/null 2>&1 || { echo "❌ npm manquant"; exit 1; }
echo "✅ Go et npm trouvés"
echo ""

echo "📦 Backend Go: vérification modules..."
cd "$REPO_DIR/backend"
go mod tidy
echo "✅ Modules Go OK"
echo ""

echo "📦 Frontend webfront2: installation deps..."
cd "$REPO_DIR/webfront2"
npm install
echo "✅ Dépendances npm OK"
echo ""

echo "🚀 Lancer backend (Go)..."
cd "$REPO_DIR/backend"
go run main.go &
BACKEND_PID=$!
echo "Backend lancé (PID=$BACKEND_PID)"
echo ""

echo "🚀 Lancer frontend (Vite)..."
cd "$REPO_DIR/webfront2"
npm run dev &
FRONT_PID=$!
echo "Frontend lancé (PID=$FRONT_PID)"

echo ""
echo "=========================================="
echo "✅ Env local lancé"
echo "Backend: http://localhost:8080 (selon config)"
echo "Frontend: http://localhost:5173" 
echo "=========================================="