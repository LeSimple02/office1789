#!/bin/bash
# Diagnostic et réparation complète - Mailserver + Roundcube + Matrix SSO

echo "=========================================="
echo "🔍 DIAGNOSTIC COMPLET SYSTÈME"
echo "=========================================="

cd /home/debian/office1789/docker

# 1. Vérifier le statut de tous les containers
echo ""
echo "1️⃣ Status des containers Docker:"
echo "=========================================="
docker compose ps

# 2. Vérifier le MAILSERVER en priorité
echo ""
echo "2️⃣ DIAGNOSTIC MAILSERVER (cause Roundcube):"
echo "=========================================="

MAILSERVER_STATUS=$(docker compose ps mailserver --format "{{.Status}}")
echo "Status: $MAILSERVER_STATUS"

if [[ ! "$MAILSERVER_STATUS" =~ "Up" ]]; then
    echo "❌ MAILSERVER EST DOWN - C'est pour ça que Roundcube ne fonctionne pas!"
    echo ""
    echo "🔍 Logs mailserver (20 dernières lignes):"
    docker compose logs mailserver --tail=20
    
    echo ""
    echo "📋 Vérification des certificats SSL:"
    ls -lh /home/debian/office1789/docker/config/ssl/
    
    # Vérifier si ce sont des fichiers ou des dossiers
    if [ -d "/home/debian/office1789/docker/config/ssl/cert.pem" ]; then
        echo "❌ cert.pem est un DOSSIER (devrait être un fichier)!"
        echo "🔧 Correction en cours..."
        
        cd /home/debian/office1789/docker
        sudo rm -rf config/ssl/*
        sudo mkdir -p config/ssl
        sudo cp /etc/letsencrypt/live/office1789.com-0001/fullchain.pem config/ssl/cert.pem
        sudo cp /etc/letsencrypt/live/office1789.com-0001/privkey.pem config/ssl/key.pem
        sudo chown -R debian:debian config/ssl/
        
        echo "✅ Certificats corrigés (fichiers)"
        ls -lh /home/debian/office1789/docker/config/ssl/
        
        echo ""
        echo "🔄 Redémarrage du mailserver..."
        docker compose up -d mailserver
        echo "⏳ Attente 30s..."
        sleep 30
        
        MAILSERVER_STATUS=$(docker compose ps mailserver --format "{{.Status}}")
        echo "Nouveau status: $MAILSERVER_STATUS"
    else
        echo "✅ Certificats sont des fichiers (bon)"
        echo "Problème ailleurs, vérification des logs..."
    fi
else
    echo "✅ MAILSERVER est UP"
fi

# 3. Test connexion Roundcube -> Mailserver
echo ""
echo "3️⃣ TEST CONNEXION ROUNDCUBE:"
echo "=========================================="

echo "Test depuis le container Roundcube vers mailserver..."
docker compose exec -T roundcube nc -zv mailserver 143 2>&1 || echo "❌ Roundcube ne peut pas joindre mailserver:143"
docker compose exec -T roundcube nc -zv mailserver 587 2>&1 || echo "❌ Roundcube ne peut pas joindre mailserver:587"

echo ""
echo "Logs Roundcube (erreurs IMAP):"
docker compose logs roundcube --tail=30 | grep -i "imap\|error\|connection" | tail -10

# 4. Vérifier la configuration Roundcube
echo ""
echo "4️⃣ CONFIGURATION ROUNDCUBE:"
echo "=========================================="

echo "Vérification config IMAP dans Roundcube..."
docker compose exec -T roundcube cat /var/www/html/config/config.inc.php 2>/dev/null | grep -A2 "default_host\|smtp_host" || echo "⚠️ Impossible de lire la config Roundcube"

# 5. MATRIX SSO - Vérifier le backend
echo ""
echo "5️⃣ DIAGNOSTIC MATRIX SSO:"
echo "=========================================="

echo "Variables backend pour SSO:"
sudo cat /home/debian/office1789/backend/.env | grep -E "ELEMENT_URL|ROUNDCUBE_URL|MATRIX_"

echo ""
echo "Test du backend (session check):"
curl -s -X POST https://backend.office1789.com/api/session/check \
  -H "Content-Type: application/json" \
  -d '{"username":"test","token":"test"}' | head -5

echo ""
echo "Status du service backend:"
sudo systemctl status office1789-backend --no-pager | head -15

echo ""
echo "Logs backend (SSO Matrix):"
sudo journalctl -u office1789-backend --since "5 minutes ago" | grep -i "matrix\|element\|sso" | tail -10

# 6. Vérifier Matrix homeserver
echo ""
echo "6️⃣ MATRIX HOMESERVER:"
echo "=========================================="

echo "Status Synapse:"
docker compose ps synapse

echo ""
echo "Test homeserver accessible:"
curl -s https://chat.office1789.com/_matrix/client/versions 2>&1 | head -5

echo ""
echo "Variables Synapse dans homeserver.yaml:"
if [ -f synapse/homeserver.yaml ]; then
    grep -E "server_name|public_baseurl" synapse/homeserver.yaml
else
    echo "❌ homeserver.yaml introuvable!"
fi

# 7. Test complet du flux SSO
echo ""
echo "7️⃣ TEST FLUX SSO MATRIX:"
echo "=========================================="

echo "Le flux SSO Matrix devrait être:"
echo "1. Frontend -> Backend /api/matrix/sso"
echo "2. Backend génère token Matrix avec homeserver"
echo "3. Backend redirige vers Element avec token"
echo ""
echo "Vérifier que MATRIX_HOMESERVER est défini dans backend/.env:"
sudo grep "MATRIX_HOMESERVER" /home/debian/office1789/backend/.env || echo "❌ MATRIX_HOMESERVER non défini!"

# RÉSUMÉ
echo ""
echo "=========================================="
echo "📊 RÉSUMÉ DES PROBLÈMES"
echo "=========================================="

ISSUES=()

if [[ ! "$MAILSERVER_STATUS" =~ "Up" ]]; then
    ISSUES+=("❌ Mailserver DOWN (bloque Roundcube)")
fi

ROUNDCUBE_STATUS=$(docker compose ps roundcube --format "{{.Status}}")
if [[ ! "$ROUNDCUBE_STATUS" =~ "Up" ]]; then
    ISSUES+=("❌ Roundcube DOWN")
fi

SYNAPSE_STATUS=$(docker compose ps synapse --format "{{.Status}}")
if [[ ! "$SYNAPSE_STATUS" =~ "Up" ]]; then
    ISSUES+=("❌ Synapse DOWN")
fi

if ! sudo grep -q "MATRIX_HOMESERVER" /home/debian/office1789/backend/.env; then
    ISSUES+=("❌ Variable MATRIX_HOMESERVER manquante dans backend")
fi

if [ ${#ISSUES[@]} -eq 0 ]; then
    echo "✅ Tous les services sont UP"
    echo ""
    echo "Si SSO ne fonctionne toujours pas:"
    echo "1. Vérifier les logs backend: sudo journalctl -u office1789-backend -f"
    echo "2. Tester manuellement: curl -X POST https://backend.office1789.com/api/matrix/sso"
else
    echo "Problèmes détectés:"
    for issue in "${ISSUES[@]}"; do
        echo "$issue"
    done
fi

echo ""
echo "=========================================="
echo "🔧 ACTIONS RECOMMANDÉES"
echo "=========================================="

if [[ ! "$MAILSERVER_STATUS" =~ "Up" ]]; then
    echo "1. Exécuter: cd /home/debian/office1789 && ./fix_mailserver.sh"
fi

if ! sudo grep -q "MATRIX_HOMESERVER" /home/debian/office1789/backend/.env; then
    echo "2. Ajouter dans backend/.env:"
    echo "   MATRIX_HOMESERVER=http://synapse:8008"
    echo "   Puis: sudo systemctl restart office1789-backend"
fi

echo ""
echo "3. Une fois tout UP, tester:"
echo "   - Roundcube: https://mail.office1789.com"
echo "   - Matrix: https://chat.office1789.com"
echo "   - SSO depuis: https://office1789.com (connecté)"
echo ""
