#!/bin/bash
# Script de diagnostic et correction complète
# Pour Roundcube, Matrix et variables manquantes

echo "=========================================="
echo "🔍 DIAGNOSTIC COMPLET"
echo "=========================================="

cd /home/debian/office1789/docker

# 1. Vérifier le .env
echo ""
echo "1️⃣ Vérification des variables .env..."
echo ""

# Chercher la variable k problématique
if grep -q "^\${k}" .env 2>/dev/null; then
    echo "⚠️  Variable \${k} trouvée dans .env, suppression..."
    sed -i '/^\${k}/d' .env
fi

# Vérifier les variables essentielles
echo "Variables Matrix:"
grep -E "^SYNAPSE_" .env || echo "❌ Variables SYNAPSE manquantes"

echo ""
echo "Variables PostgreSQL:"
grep -E "^POSTGRES_" .env || echo "❌ Variables POSTGRES manquantes"

# 2. Vérifier la base PostgreSQL Synapse
echo ""
echo "2️⃣ Configuration PostgreSQL Synapse..."

SYNAPSE_DB_USER=$(grep "^SYNAPSE_DB_USER=" .env | cut -d= -f2 || echo "synapse_user")
SYNAPSE_DB_PASSWORD=$(grep "^SYNAPSE_DB_PASSWORD=" .env | cut -d= -f2)
SYNAPSE_DB_NAME=$(grep "^SYNAPSE_DB_NAME=" .env | cut -d= -f2 || echo "synapse")

echo "Utilisateur: $SYNAPSE_DB_USER"
echo "Base: $SYNAPSE_DB_NAME"

# Redémarrer postgres_synapse
docker compose restart postgres_synapse
sleep 5

# Créer l'utilisateur et la base si nécessaire
echo "Création de l'utilisateur et de la base..."
docker compose exec -T postgres_synapse psql -U postgres <<EOF 2>/dev/null
-- Créer l'utilisateur s'il n'existe pas
DO \$\$
BEGIN
  IF NOT EXISTS (SELECT FROM pg_user WHERE usename = '$SYNAPSE_DB_USER') THEN
    CREATE USER $SYNAPSE_DB_USER WITH PASSWORD '$SYNAPSE_DB_PASSWORD';
  END IF;
END
\$\$;

-- Créer la base si elle n'existe pas
SELECT 'CREATE DATABASE $SYNAPSE_DB_NAME OWNER $SYNAPSE_DB_USER'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = '$SYNAPSE_DB_NAME')\gexec

-- Donner les droits
GRANT ALL PRIVILEGES ON DATABASE $SYNAPSE_DB_NAME TO $SYNAPSE_DB_USER;
\q
EOF

echo "✅ Base PostgreSQL configurée"

# 3. Vérifier/Créer homeserver.yaml
echo ""
echo "3️⃣ Configuration homeserver.yaml..."

if [ ! -f synapse/homeserver.yaml ]; then
    echo "Création de homeserver.yaml..."
    mkdir -p synapse
    
    cat > synapse/homeserver.yaml <<HOMESERVER_EOF
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
    user: $SYNAPSE_DB_USER
    password: $SYNAPSE_DB_PASSWORD
    database: $SYNAPSE_DB_NAME
    host: postgres_synapse
    port: 5432
    cp_min: 5
    cp_max: 10

log_config: "/data/office1789.com.log.config"
media_store_path: /data/media_store
registration_shared_secret: "$(openssl rand -base64 32)"
report_stats: false
macaroon_secret_key: "$(openssl rand -base64 32)"
form_secret: "$(openssl rand -base64 32)"
signing_key_path: "/data/office1789.com.signing.key"

trusted_key_servers:
  - server_name: "matrix.org"

enable_registration: false
enable_registration_without_verification: false
suppress_key_server_warning: true
HOMESERVER_EOF

    echo "✅ homeserver.yaml créé avec les bonnes variables"
else
    echo "✅ homeserver.yaml existe"
    # Vérifier les bonnes valeurs
    grep -q "server_name.*office1789.com" synapse/homeserver.yaml || sed -i 's/^server_name:.*/server_name: "office1789.com"/' synapse/homeserver.yaml
    grep -q "public_baseurl.*chat.office1789.com" synapse/homeserver.yaml || sed -i 's|^public_baseurl:.*|public_baseurl: "https://chat.office1789.com"|' synapse/homeserver.yaml
fi

# 4. Redémarrer tous les services
echo ""
echo "4️⃣ Redémarrage des services..."

docker compose stop synapse
sleep 2
docker compose up -d postgres_synapse
sleep 10
docker compose up -d synapse
echo "⏳ Attente du démarrage de Synapse (30s)..."
sleep 30

# 5. Vérifications
echo ""
echo "=========================================="
echo "✅ VÉRIFICATIONS FINALES"
echo "=========================================="

echo ""
echo "📊 Status des containers:"
docker compose ps | grep -E "mailserver|roundcube|synapse|postgres"

echo ""
echo "🔍 Logs Synapse (10 dernières lignes):"
docker compose logs synapse --tail=10 | grep -v "WARN.*variable"

echo ""
echo "🌐 Test homeserver Matrix:"
curl -s https://chat.office1789.com/_matrix/client/versions 2>/dev/null | head -3

echo ""
echo "📧 Test Roundcube/Mailserver:"
docker compose logs mailserver --tail=5 | grep -E "ready|error"

echo ""
echo "=========================================="
echo "✅ TERMINÉ!"
echo "=========================================="
echo ""
echo "🧪 Tests à faire:"
echo "1. Roundcube: https://mail.office1789.com"
echo "2. Matrix: https://chat.office1789.com"
echo "3. SSO depuis office1789.com"
echo ""
