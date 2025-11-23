#!/bin/bash
set -e

REPO_DIR="/home/debian/office1789"
DOCKER_DIR="$REPO_DIR/docker"

echo "=========================================="
echo "🧨 TEARDOWN COMPLET DEPLOIEMENT OFFICE1789 (DOCKER)"
echo "=========================================="
echo "⚠️  ATTENTION: ceci va arrêter et supprimer tous les containers et certains volumes Docker liés à Office1789."
echo "    Les données mail, Matrix, Roundcube, DB seront PERDUES si tu supprimes les volumes."
echo "=========================================="

read -p "Taper OUI pour continuer: " CONFIRM
if [ "$CONFIRM" != "OUI" ]; then
  echo "❌ Abandon."
  exit 1
fi

echo "✅ Confirmation reçue, on continue."
echo ""

cd "$DOCKER_DIR"

echo "🛑 Arrêt et suppression des containers (docker compose down)..."
docker compose down

echo ""
read -p "Supprimer aussi les VOLUMES de données Docker (OUI/NON) ? " DEL_VOLUMES
if [ "$DEL_VOLUMES" = "OUI" ]; then
  echo "🧹 Suppression des volumes Docker liés à Office1789..."
  docker volume rm office1789_pgdata \
                    office1789_maildata \
                    office1789_mailstate \
                    office1789_maillogs \
                    office1789_roundcube_db \
                    office1789_synapse_db 2>/dev/null || true
  echo "✅ Volumes supprimés (si présents)."
else
  echo "ℹ️  Volumes conservés."
fi

echo ""
read -p "Supprimer aussi les IMAGES Docker construites (frontend, element, roundcube, synapse) ? (OUI/NON) " DEL_IMAGES
if [ "$DEL_IMAGES" = "OUI" ]; then
  echo "🧹 Suppression des images Docker personnalisées..."
  docker rmi docker-frontend docker-element docker-roundcube 2>/dev/null || true
  echo "✅ Images supprimées (si présentes)."
else
  echo "ℹ️  Images conservées."
fi

echo "=========================================="
echo "✅ TEARDOWN TERMINÉ"
echo "Tu peux maintenant relancer un déploiement propre avec :"
echo "  cd $REPO_DIR && ./scripts/prod_deploy.sh"
echo "=========================================="