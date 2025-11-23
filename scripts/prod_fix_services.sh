#!/bin/bash
set -e

REPO_DIR="/home/debian/office1789"
DOCKER_DIR="$REPO_DIR/docker"

echo "=========================================="
echo "🔧 FIX SERVICES PROD OFFICE1789"
echo "=========================================="
echo ""

cd "$DOCKER_DIR"

echo "📋 Vérification de l'environnement..."
if [ ! -f ".env" ]; then
  echo "❌ Fichier .env manquant dans $DOCKER_DIR"
  exit 1
fi
echo "✅ .env présent"
echo ""

echo "🔧 Correction Synapse (permissions media_store)..."
if [ ! -d "synapse/conf" ]; then
  echo "⚠️  Dossier synapse/conf n'existe pas, création..."
  sudo mkdir -p synapse/conf
fi

sudo mkdir -p synapse/conf/media_store
sudo chown -R 991:991 synapse/conf
echo "✅ Permissions Synapse corrigées (UID 991)"
echo ""

echo "🔧 Correction Roundcube (réinitialisation base de données)..."
echo "⚠️  Cette opération va supprimer les données Roundcube existantes (pas les emails)"
read -p "Continuer? (OUI pour confirmer): " confirm
if [ "$confirm" != "OUI" ]; then
  echo "❌ Abandon de la réinitialisation Roundcube"
else
  echo "🛑 Arrêt de Roundcube et sa base..."
  docker compose stop roundcube postgres_roundcube
  
  echo "🗑️  Suppression du volume roundcube_db..."
  docker volume rm office1789_roundcube_db 2>/dev/null || echo "   (volume déjà supprimé ou inexistant)"
  
  echo "🔄 Redémarrage de postgres_roundcube..."
  docker compose up -d postgres_roundcube
  sleep 10
  
  echo "🔄 Redémarrage de roundcube..."
  docker compose up -d roundcube
  sleep 5
  
  echo "✅ Roundcube réinitialisé"
fi
echo ""

echo "🔄 Redémarrage de Synapse pour appliquer les corrections..."
docker compose restart synapse
sleep 5
echo ""

echo "📊 État des services:"
docker compose ps
echo ""

echo "🧪 Logs Synapse (dernières lignes):"
docker compose logs synapse --tail=20
echo ""

echo "🧪 Logs Roundcube (dernières lignes):"
docker compose logs roundcube --tail=20
echo ""

echo "=========================================="
echo "✅ FIX SERVICES TERMINÉ"
echo "Si des erreurs persistent, vérifier:"
echo "  - Les mots de passe dans docker/.env"
echo "  - Les logs complets: docker compose logs <service>"
echo "=========================================="
