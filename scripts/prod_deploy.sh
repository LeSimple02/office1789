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

echo "🐳 Démarrage Synapse + Element + Roundcube + frontend..."
docker compose up -d synapse element roundcube frontend coturn onlyoffice
echo "✅ Services applicatifs lancés"
echo ""

echo "📊 docker compose ps:"
docker compose ps
echo ""

echo "🧪 Logs Synapse (30 dernières lignes):"
docker compose logs synapse --tail=30 || true
echo ""

echo "🧪 Logs Mailserver (30 dernières lignes):"
docker compose logs mailserver --tail=30 || true
echo ""

echo "🧪 Logs Roundcube (30 dernières lignes):"
docker compose logs roundcube --tail=30 || true
echo ""

echo "=========================================="
echo "✅ DEPLOY PROD TERMINÉ"
echo "À tester:"
echo "- Mail:     https://mail.office1789.com"
echo "- Chat:     https://chat.office1789.com"
echo "- Frontend: https://office1789.com (via ton nginx)"
echo "=========================================="