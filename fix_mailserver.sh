#!/bin/bash
# Script de réparation rapide pour Mailserver

echo "🔧 Réparation Mailserver..."
echo ""

# 1. Copier certificats SSL
echo "📋 1. Copie des certificats SSL..."
sudo mkdir -p /home/debian/office1789/docker/config/ssl
sudo cp /etc/letsencrypt/live/office1789.com-0001/fullchain.pem /home/debian/office1789/docker/config/ssl/cert.pem
sudo cp /etc/letsencrypt/live/office1789.com-0001/privkey.pem /home/debian/office1789/docker/config/ssl/key.pem
sudo chown -R debian:debian /home/debian/office1789/docker/config/ssl
echo "✅ Certificats copiés"
echo ""

# 2. Redémarrer mailserver
echo "🔄 2. Redémarrage du mailserver..."
cd /home/debian/office1789/docker
docker compose restart mailserver
echo "⏳ Attente 30 secondes..."
sleep 30
echo ""

# 3. Vérifier l'état
echo "📊 3. État du mailserver:"
docker compose ps mailserver
echo ""

# 4. Vérifier les logs
echo "📝 4. Logs récents:"
docker compose logs mailserver --tail=15
echo ""

# 5. Test de connexion IMAP
echo "🔌 5. Test connexion IMAP depuis Roundcube:"
docker compose exec roundcube nc -zv mailserver 143 2>&1
echo ""

echo "✅ Réparation terminée!"
