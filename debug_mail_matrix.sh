#!/bin/bash
# Script de debug pour Roundcube et Matrix SSO

echo "========================================="
echo "🔍 DEBUG MAIL & MATRIX SSO"
echo "========================================="
echo ""

echo "📦 1. État des conteneurs"
echo "========================================="
cd /home/debian/office1789/docker
docker compose ps | grep -E "mailserver|roundcube|synapse|element"
echo ""

echo "📧 2. Logs Mailserver (dernières 20 lignes)"
echo "========================================="
docker compose logs mailserver --tail=20
echo ""

echo "📬 3. Logs Roundcube (dernières 20 lignes - erreurs uniquement)"
echo "========================================="
docker compose logs roundcube --tail=50 | grep -i "error\|failed\|IMAP Error"
echo ""

echo "💬 4. Logs Matrix/Synapse (dernières 20 lignes)"
echo "========================================="
docker compose logs synapse --tail=20
echo ""

echo "🔐 5. Vérification certificats SSL mailserver"
echo "========================================="
ls -lh /home/debian/office1789/docker/config/ssl/
echo ""

echo "🌐 6. Test connexion Roundcube -> Mailserver"
echo "========================================="
docker compose exec roundcube ping -c 2 mailserver 2>/dev/null || echo "❌ Ping failed"
docker compose exec roundcube nslookup mailserver 2>/dev/null || echo "❌ DNS lookup failed"
echo ""

echo "📝 7. Configuration Roundcube IMAP"
echo "========================================="
docker compose exec roundcube cat /var/www/html/config/config.inc.php | grep -E "imap_host|smtp_server" || echo "Config file not found"
echo ""

echo "🔑 8. Test backend SSO"
echo "========================================="
echo "Backend URL: https://backend.office1789.com"
echo "Roundcube URL from backend .env:"
grep ROUNDCUBE_URL /home/debian/office1789/backend/.env
echo "Element URL from backend .env:"
grep ELEMENT_URL /home/debian/office1789/backend/.env
echo ""

echo "✅ Debug terminé!"
echo ""
echo "Actions recommandées:"
echo "- Si mailserver est 'Restarting': vérifier certificats SSL"
echo "- Si Roundcube ne ping pas mailserver: problème réseau Docker"
echo "- Si IMAP Error: vérifier que mailserver écoute sur port 143"
echo "- Si Matrix SSO échoue: vérifier logs synapse et configuration SSO"
