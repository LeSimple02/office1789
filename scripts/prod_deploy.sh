#!/bin/bash
set -e

REPO_DIR="/home/debian/office1789"
DOCKER_DIR="$REPO_DIR/docker"

echo "=========================================="
echo "🚀 DEPLOY PROD OFFICE1789"
echo "=========================================="
echo ""

cd "$REPO_DIR"

echo "📥 Mise à jour du code..."
git pull origin main
echo "✅ Code à jour"
echo ""

cd "$DOCKER_DIR"

echo "🧩 Vérification .env..."
if [ ! -f ".env" ]; then
  if [ -f ".env.example" ]; then
    cp .env.example .env
    echo "⚠️  .env créé depuis .env.example - pense à éditer les secrets (CHANGE_ME_*)"
  else
    echo "❌ .env manquant et pas de .env.example"
    exit 1
  fi
fi
echo "✅ .env présent"
echo ""

echo "🔐 Préparation certificats SSL..."
bash "$DOCKER_DIR/prepare_certs.sh"
echo ""

echo "🐳 Build des images clés..."
docker compose build frontend element roundcube synapse
echo "✅ Build terminé"
echo ""

echo "🐳 Démarrage des bases de données..."
docker compose up -d postgres_db postgres_roundcube postgres_synapse
sleep 10
echo "✅ Postgres UP"
echo ""

echo "🐳 Démarrage mailserver..."
docker compose up -d mailserver
sleep 10
echo "✅ Mailserver lancé (vérifier logs si souci)"
echo ""

echo "🐳 Démarrage backend..."
docker compose up -d backend
sleep 5
echo "✅ Backend lancé"
echo ""

echo "🐳 Démarrage Synapse + Element + Roundcube + frontend..."
docker compose up -d synapse element roundcube frontend coturn onlyoffice
echo "✅ Services applicatifs lancés"
echo ""

echo "📊 docker compose ps:"
docker compose ps
echo ""

echo "🧪 Test backend (interne):"
curl -s http://127.0.0.1:8080/api/welcome | head -n 2 || echo "⚠️  Backend non accessible"
echo ""

echo "🧪 Logs Backend (20 dernières lignes):"
docker compose logs backend --tail=20 || true
echo ""

echo "🧪 Logs Synapse (20 dernières lignes):"
docker compose logs synapse --tail=20 || true
echo ""

echo "🧪 Logs Mailserver (20 dernières lignes):"
docker compose logs mailserver --tail=20 || true
echo ""

echo "🧪 Logs Roundcube (20 dernières lignes):"
docker compose logs roundcube --tail=20 || true
echo ""

echo "=========================================="
echo "✅ DEPLOY PROD TERMINÉ"
echo ""
echo "📝 Services déployés:"
echo "  - Backend:   http://127.0.0.1:8080 (via backend.office1789.com)"
echo "  - Frontend:  http://127.0.0.1:5173 (via office1789.com)"
echo "  - Mail:      http://127.0.0.1:8081 (via mail.office1789.com)"
echo "  - Chat:      http://127.0.0.1:8083 (via chat.office1789.com)"
echo "  - Matrix:    http://127.0.0.1:8008 (API)"
echo "  - OnlyOffice: http://127.0.0.1:8082"
echo "  - Coturn:    3478 (TURN/STUN)"
echo ""
echo "⚠️  IMPORTANT: Configure nginx pour exposer ces services:"
echo "  sudo systemctl start nginx"
echo ""
echo "Si erreurs Synapse/Roundcube, exécuter:"
echo "  cd $REPO_DIR && ./scripts/prod_fix_services.sh"
echo "=========================================="