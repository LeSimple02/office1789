#!/bin/bash
# Script de réparation rapide pour Mailserver et Frontend

echo "🔧 Réparation Mailserver et mise à jour Frontend..."
echo ""

# 0. Git pull pour mettre à jour le code
echo "📥 0. Mise à jour du code depuis GitHub..."
cd /home/debian/office1789
git pull origin main
echo "✅ Code mis à jour"
echo ""

# 1. Supprimer les dossiers SSL incorrects et copier correctement
echo "📋 1. Nettoyage et copie des certificats SSL..."
sudo rm -rf /home/debian/office1789/docker/config/ssl/cert.pem
sudo rm -rf /home/debian/office1789/docker/config/ssl/key.pem
sudo rm -rf /home/debian/office1789/docker/config/ssl/cert
sudo rm -rf /home/debian/office1789/docker/config/ssl/key
sudo mkdir -p /home/debian/office1789/docker/config/ssl
sudo cp /etc/letsencrypt/live/office1789.com-0001/fullchain.pem /home/debian/office1789/docker/config/ssl/cert.pem
sudo cp /etc/letsencrypt/live/office1789.com-0001/privkey.pem /home/debian/office1789/docker/config/ssl/key.pem
sudo chown -R debian:debian /home/debian/office1789/docker/config/ssl
echo "✅ Certificats copiés correctement"
echo ""

# 2. Rebuild et redémarrer le frontend
echo "🎨 2. Rebuild du frontend avec nouveau code..."
cd /home/debian/office1789/docker
docker compose stop frontend
docker compose rm -f frontend
docker rmi docker-frontend 2>/dev/null || true
docker compose build --no-cache frontend
docker compose up -d frontend
echo "✅ Frontend redémarré"
echo ""

# 3. Redémarrer mailserver
# 3. Redémarrer mailserver
echo "🔄 3. Redémarrage du mailserver..."
docker compose restart mailserver
echo "⏳ Attente 30 secondes..."
sleep 30
echo ""

# 4. Vérifier l'état
echo "📊 4. État des services:"
docker compose ps | grep -E "frontend|mailserver"
echo ""

# 5. Vérifier les logs mailserver
echo "📝 5. Logs mailserver récents:"
docker compose logs mailserver --tail=15
echo ""

# 6. Vérifier les certificats
echo "🔐 6. Vérification des certificats (doivent être des FICHIERS, pas des dossiers):"
ls -lh /home/debian/office1789/docker/config/ssl/
echo ""

# 7. Test de connexion IMAP
echo "🔌 7. Test connexion IMAP depuis Roundcube:"
docker compose exec -T roundcube sh -c "command -v nc >/dev/null 2>&1 && nc -zv mailserver 143 2>&1 || echo 'nc not available, trying telnet...'; command -v telnet >/dev/null 2>&1 && timeout 2 telnet mailserver 143 2>&1 || echo 'Connection tools not available in container'"
echo ""

echo "✅ Réparation terminée!"
echo ""
echo "🎯 Prochaines étapes:"
echo "- Vérifier que mailserver est 'Up' (pas 'Restarting')"
echo "- Vérifier que cert.pem et key.pem sont des fichiers (-) pas des dossiers (d)"
echo "- Tester Roundcube: https://mail.office1789.com"
