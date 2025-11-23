#!/bin/bash
# Script de réparation Roundcube et Matrix/Synapse
# Usage: ./fix_roundcube_matrix.sh

set -e

echo "=========================================="
echo "🔧 RÉPARATION ROUNDCUBE ET MATRIX"
echo "=========================================="

cd /home/debian/office1789/docker

# Étape 1: Vérification des variables d'environnement
echo ""
echo "📋 Étape 1: Vérification des variables Matrix..."
if [ ! -f .env ]; then
    echo "❌ Fichier .env introuvable!"
    exit 1
fi

# Afficher les variables Matrix (sans les mots de passe complets)
echo "Variables Matrix actuelles:"
grep -E "SYNAPSE_|POSTGRES_" .env | sed 's/\(PASSWORD=\).*/\1***/'

# Étape 2: Fix PostgreSQL pour Synapse
echo ""
echo "🔧 Étape 2: Réparation PostgreSQL Synapse..."

# Arrêter Synapse
echo "Arrêt de Synapse..."
docker compose stop synapse

# Récupérer les variables de la base
POSTGRES_PASSWORD=$(grep "^POSTGRES_PASSWORD=" .env | cut -d= -f2)
SYNAPSE_DB_USER=$(grep "^SYNAPSE_DB_USER=" .env | cut -d= -f2 || echo "synapse")
SYNAPSE_DB_PASSWORD=$(grep "^SYNAPSE_DB_PASSWORD=" .env | cut -d= -f2)
SYNAPSE_DB_NAME=$(grep "^SYNAPSE_DB_NAME=" .env | cut -d= -f2 || echo "synapse")

echo "POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:0:3}***"
echo "SYNAPSE_DB_USER: $SYNAPSE_DB_USER"
echo "SYNAPSE_DB_PASSWORD: ${SYNAPSE_DB_PASSWORD:0:3}***"
echo "SYNAPSE_DB_NAME: $SYNAPSE_DB_NAME"

# Redémarrer PostgreSQL Synapse
echo "Redémarrage de postgres_synapse..."
docker compose restart postgres_synapse
sleep 10

# Vérifier la connexion PostgreSQL avec le bon utilisateur
echo "Test de connexion PostgreSQL..."
docker compose exec -T postgres_synapse psql -U "$SYNAPSE_DB_USER" -d "$SYNAPSE_DB_NAME" -c "SELECT version();" 2>&1 | head -5 || {
    echo "⚠️  Échec de connexion avec $SYNAPSE_DB_USER"
    echo "Tentative de création de l'utilisateur et de la base..."
    docker compose exec -T postgres_synapse psql -U postgres -c "CREATE USER $SYNAPSE_DB_USER WITH PASSWORD '$SYNAPSE_DB_PASSWORD';" 2>/dev/null || echo "Utilisateur existe déjà"
    docker compose exec -T postgres_synapse psql -U postgres -c "CREATE DATABASE $SYNAPSE_DB_NAME OWNER $SYNAPSE_DB_USER;" 2>/dev/null || echo "Base existe déjà"
    docker compose exec -T postgres_synapse psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE $SYNAPSE_DB_NAME TO $SYNAPSE_DB_USER;" 2>/dev/null
    echo "✅ Utilisateur et base configurés"
}

# Étape 3: Configuration du homeserver Synapse
echo ""
echo "🔧 Étape 3: Configuration homeserver.yaml..."

# Vérifier si le fichier homeserver.yaml existe
if [ ! -f synapse/homeserver.yaml ]; then
    echo "❌ homeserver.yaml introuvable!"
    echo "Création d'un homeserver.yaml minimal..."
    
    cat > synapse/homeserver.yaml << 'HOMESERVER_EOF'
server_name: "office1789.com"
public_baseurl: "https://chat.office1789.com"
pid_file: /data/homeserver.pid
listeners:
  - port: 8008
    tls: false
    type: http
    x_forwarded: true
    bind_addresses: ['0.0.0.0']
    resources:
      - names: [client, federation]
        compress: false

database:
  name: psycopg2
  args:
    user: ${SYNAPSE_DB_USER}
    password: ${SYNAPSE_DB_PASSWORD}
    database: ${SYNAPSE_DB_NAME}
    host: postgres_synapse
    cp_min: 5
    cp_max: 10

log_config: "/data/office1789.com.log.config"
media_store_path: /data/media_store
registration_shared_secret: "votre_secret_registration"
report_stats: false
macaroon_secret_key: "votre_secret_macaroon"
form_secret: "votre_secret_form"
signing_key_path: "/data/office1789.com.signing.key"
trusted_key_servers:
  - server_name: "matrix.org"

enable_registration: false
enable_registration_without_verification: false

suppress_key_server_warning: true
HOMESERVER_EOF

    echo "✅ homeserver.yaml créé"
else
    echo "✅ homeserver.yaml existe"
    
    # Vérifier la configuration du homeserver
    echo "Vérification server_name..."
    grep -q "server_name.*office1789.com" synapse/homeserver.yaml || {
        echo "⚠️  server_name incorrect, correction..."
        sed -i 's/^server_name:.*/server_name: "office1789.com"/' synapse/homeserver.yaml
    }
    
    echo "Vérification public_baseurl..."
    grep -q "public_baseurl.*chat.office1789.com" synapse/homeserver.yaml || {
        echo "⚠️  public_baseurl incorrect, correction..."
        sed -i 's|^public_baseurl:.*|public_baseurl: "https://chat.office1789.com"|' synapse/homeserver.yaml
    }
fi

# Étape 4: Redémarrer Synapse
echo ""
echo "🔧 Étape 4: Redémarrage de Synapse..."
docker compose up -d synapse
echo "⏳ Attente du démarrage de Synapse (30s)..."
sleep 30

# Étape 5: Vérifier Mailserver
echo ""
echo "🔧 Étape 5: Vérification Mailserver..."
docker compose ps mailserver
MAILSERVER_STATUS=$(docker compose ps mailserver --format "{{.Status}}" || echo "error")

if [[ "$MAILSERVER_STATUS" =~ "Up" ]]; then
    echo "✅ Mailserver est UP"
else
    echo "❌ Mailserver n'est pas UP: $MAILSERVER_STATUS"
    echo "Le mailserver doit être réparé avec fix_mailserver.sh"
fi

# Étape 6: Test Roundcube
echo ""
echo "🔧 Étape 6: Test de connexion Roundcube..."
docker compose ps roundcube
docker compose logs roundcube --tail=20 | grep -i "error\|warning\|imap" || echo "Pas d'erreurs IMAP récentes"

# Étape 7: Vérifications finales
echo ""
echo "=========================================="
echo "📊 VÉRIFICATIONS FINALES"
echo "=========================================="

echo ""
echo "1. Status des containers:"
docker compose ps | grep -E "mailserver|roundcube|synapse|postgres_synapse|element"

echo ""
echo "2. Logs Synapse (dernières lignes):"
docker compose logs synapse --tail=10

echo ""
echo "3. Test de connexion Matrix homeserver:"
curl -s https://chat.office1789.com/_matrix/client/versions | jq '.' 2>/dev/null || echo "❌ Homeserver non accessible"

echo ""
echo "4. Variables backend pour SSO:"
ssh -o StrictHostKeyChecking=no debian@localhost "grep -E 'ROUNDCUBE_URL|ELEMENT_URL' /home/debian/office1789/backend/.env" 2>/dev/null || echo "⚠️  Impossible de vérifier les variables backend"

echo ""
echo "=========================================="
echo "✅ RÉPARATION TERMINÉE"
echo "=========================================="
echo ""
echo "📝 Prochaines étapes:"
echo "1. Tester Roundcube: https://mail.office1789.com"
echo "2. Tester Matrix: https://chat.office1789.com"
echo "3. Tester SSO depuis Office1789:"
echo "   - Connexion -> Ouvrir Email"
echo "   - Connexion -> Chat"
echo ""
echo "Si Roundcube ne fonctionne toujours pas:"
echo "   → Exécuter ./fix_mailserver.sh pour réparer les certificats SSL"
echo ""
