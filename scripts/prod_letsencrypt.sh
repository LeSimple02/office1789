#!/bin/bash
set -e

DOMAIN_BASE="office1789.com"
DOMAINS="-d office1789.com -d mail.office1789.com -d chat.office1789.com"

echo "=========================================="
echo "🔐 Génération / renouvellement certificats Let's Encrypt"
echo "=========================================="
echo ""

# 1. Installer certbot si nécessaire
if ! command -v certbot >/dev/null 2>&1; then
  echo "📦 Installation de certbot..."
  sudo apt update
  sudo apt install -y certbot
  echo "✅ certbot installé"
else
  echo "✅ certbot déjà installé"
fi

# 2. Arrêter nginx ou autres services qui écoutent sur le port 80 (si présents)
echo "🛑 Arrêt éventuel de nginx sur le port 80..."
sudo systemctl stop nginx 2>/dev/null || true

echo "🚀 Lancement certbot (standalone) pour $DOMAINS"
sudo certbot certonly --standalone $DOMAINS --agree-tos -m "admin@$DOMAIN_BASE" --non-interactive || {
  echo "❌ Échec de la génération des certificats Let's Encrypt"
  exit 1
}

echo ""
echo "✅ Certificats Let's Encrypt générés dans /etc/letsencrypt/live/$DOMAIN_BASE/"
ls -lh "/etc/letsencrypt/live/$DOMAIN_BASE" || true

echo ""
echo "Tu peux maintenant lancer :"
echo "  cd /home/debian/office1789 && ./scripts/prod_deploy.sh"