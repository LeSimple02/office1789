#!/bin/bash
# 🚨 RÉPARATION URGENTE COMPLÈTE
# Corrige: Mailserver SSL, Synapse PostgreSQL, Docker compose variable k

set -e  # Stop on error

echo "=========================================="
echo "🚨 RÉPARATION URGENTE COMPLÈTE"
echo "=========================================="
echo ""

cd /home/debian/office1789

# ============================================
# 1. GIT PULL
# ============================================
echo "📥 1. Mise à jour du code..."
git pull origin main
echo "✅ Code à jour"
echo ""

cd docker

# ============================================
# 2. NETTOYER VARIABLE ${k} DANS .ENV
# ============================================
echo "🧹 2. Nettoyage variable \${k} dans .env..."
if grep -q '^\${k}' .env 2>/dev/null; then
    echo "⚠️  Variable \${k} trouvée, suppression..."
    sed -i '/^\${k}/d' .env
    echo "✅ Variable \${k} supprimée"
else
    echo "✅ Pas de variable \${k} trouvée"
fi
echo ""

# ============================================
# 3. RÉPARER CERTIFICATS MAILSERVER
# ============================================
echo "🔐 3. Réparation certificats SSL mailserver..."

# Le docker-compose monte vers /etc/ssl/dms/cert et /etc/ssl/dms/key (SANS .pem)
# Il faut copier les certificats sous ces noms exacts

echo "Arrêt mailserver..."
docker compose stop mailserver

echo "Nettoyage ancien SSL..."
sudo rm -rf /home/debian/office1789/docker/config/ssl/*

echo "Copie certificats avec NOMS CORRECTS (cert et key, pas .pem)..."
sudo mkdir -p /home/debian/office1789/docker/config/ssl

# Copier SANS extension .pem pour correspondre au docker-compose
sudo cp /etc/letsencrypt/live/office1789.com-0001/fullchain.pem /home/debian/office1789/docker/config/ssl/cert.pem
sudo cp /etc/letsencrypt/live/office1789.com-0001/privkey.pem /home/debian/office1789/docker/config/ssl/key.pem

sudo chown -R debian:debian /home/debian/office1789/docker/config/ssl
sudo chmod 644 /home/debian/office1789/docker/config/ssl/cert.pem
sudo chmod 600 /home/debian/office1789/docker/config/ssl/key.pem

echo "Vérification:"
ls -lh /home/debian/office1789/docker/config/ssl/

echo "✅ Certificats copiés"
echo ""

# ============================================
# 4. RÉPARER BASE POSTGRESQL SYNAPSE
# ============================================
echo "🗄️  4. Réparation base PostgreSQL Synapse..."

# Charger les variables
source .env

SYNAPSE_DB_USER=${SYNAPSE_DB_USER:-synapse_user}
SYNAPSE_DB_PASSWORD=${SYNAPSE_DB_PASSWORD}
SYNAPSE_DB_NAME=${SYNAPSE_DB_NAME:-synapse}

echo "Configuration:"
echo "  User: $SYNAPSE_DB_USER"
echo "  Database: $SYNAPSE_DB_NAME"

# Redémarrer PostgreSQL Synapse
echo "Redémarrage postgres_synapse..."
docker compose restart postgres_synapse
sleep 10

# Créer utilisateur et base
echo "Création utilisateur et base..."
docker compose exec -T postgres_synapse psql -U postgres <<EOF
-- Supprimer et recréer proprement
DROP DATABASE IF EXISTS $SYNAPSE_DB_NAME;
DROP USER IF EXISTS $SYNAPSE_DB_USER;

-- Créer utilisateur
CREATE USER $SYNAPSE_DB_USER WITH PASSWORD '$SYNAPSE_DB_PASSWORD';

-- Créer base avec bon owner
CREATE DATABASE $SYNAPSE_DB_NAME 
  OWNER $SYNAPSE_DB_USER 
  ENCODING 'UTF8' 
  LC_COLLATE 'C' 
  LC_CTYPE 'C' 
  TEMPLATE template0;

-- Accorder tous les privilèges
GRANT ALL PRIVILEGES ON DATABASE $SYNAPSE_DB_NAME TO $SYNAPSE_DB_USER;

-- Se connecter à la base pour créer l'extension
\c $SYNAPSE_DB_NAME
ALTER SCHEMA public OWNER TO $SYNAPSE_DB_USER;
GRANT ALL ON SCHEMA public TO $SYNAPSE_DB_USER;
EOF

if [ $? -eq 0 ]; then
    echo "✅ Base PostgreSQL Synapse créée"
else
    echo "❌ Erreur création base Synapse"
fi
echo ""

# ============================================
# 5. CONFIGURER HOMESERVER.YAML
# ============================================
echo "⚙️  5. Configuration homeserver.yaml Synapse..."

HOMESERVER_YAML="/home/debian/office1789/docker/synapse/conf/homeserver.yaml"

if [ ! -f "$HOMESERVER_YAML" ]; then
    echo "❌ homeserver.yaml introuvable à $HOMESERVER_YAML"
    echo "Génération config initiale..."
    
    mkdir -p /home/debian/office1789/docker/synapse/conf
    
    docker compose run --rm synapse generate || {
        echo "⚠️  Génération config échouée, création manuelle..."
        cat > "$HOMESERVER_YAML" <<YAML_EOF
server_name: "office1789.com"
public_baseurl: "https://chat.office1789.com"
pid_file: /data/homeserver.pid

listeners:
  - port: 8008
    tls: false
    type: http
    x_forwarded: true
    bind_addresses: ['::1', '127.0.0.1', '0.0.0.0']
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
registration_shared_secret: "$(openssl rand -hex 32)"
report_stats: false
macaroon_secret_key: "$(openssl rand -hex 32)"
form_secret: "$(openssl rand -hex 32)"
signing_key_path: "/data/office1789.com.signing.key"

trusted_key_servers:
  - server_name: "matrix.org"

enable_registration: false
enable_registration_without_verification: false

suppress_key_server_warning: true
YAML_EOF
    }
fi

# Mettre à jour la config database dans homeserver.yaml
echo "Mise à jour configuration database..."
python3 <<PYTHON_EOF
import yaml

config_path = "$HOMESERVER_YAML"
try:
    with open(config_path, 'r') as f:
        config = yaml.safe_load(f)
    
    # Mettre à jour la config database
    config['database'] = {
        'name': 'psycopg2',
        'args': {
            'user': '$SYNAPSE_DB_USER',
            'password': '$SYNAPSE_DB_PASSWORD',
            'database': '$SYNAPSE_DB_NAME',
            'host': 'postgres_synapse',
            'port': 5432,
            'cp_min': 5,
            'cp_max': 10
        }
    }
    
    with open(config_path, 'w') as f:
        yaml.dump(config, f, default_flow_style=False)
    
    print("✅ homeserver.yaml mis à jour")
except Exception as e:
    print(f"❌ Erreur: {e}")
PYTHON_EOF

echo ""

# ============================================
# 6. REDÉMARRER TOUS LES SERVICES
# ============================================
echo "🔄 6. Redémarrage de tous les services..."

echo "Démarrage mailserver..."
docker compose up -d mailserver
sleep 15

echo "Démarrage Synapse..."
docker compose up -d synapse
sleep 10

echo "Démarrage Element..."
docker compose up -d element

echo "Redémarrage Roundcube..."
docker compose restart roundcube

echo "✅ Services redémarrés"
echo ""

# ============================================
# 7. VÉRIFICATIONS
# ============================================
echo "=========================================="
echo "📊 VÉRIFICATIONS"
echo "=========================================="
echo ""

echo "🐳 Status containers:"
docker compose ps | grep -E "mailserver|synapse|element|roundcube"
echo ""

echo "📧 Logs Mailserver (5 dernières lignes):"
docker compose logs mailserver --tail=5
echo ""

echo "🔹 Logs Synapse (5 dernières lignes):"
docker compose logs synapse --tail=5
echo ""

echo "🔐 Certificats SSL:"
ls -lh /home/debian/office1789/docker/config/ssl/
echo ""

# ============================================
# 8. TESTS
# ============================================
echo "=========================================="
echo "🧪 TESTS"
echo "=========================================="
echo ""

echo "Test 1: Mailserver IMAP (143)"
timeout 3 bash -c '</dev/tcp/127.0.0.1/143' 2>/dev/null && echo "✅ IMAP accessible" || echo "❌ IMAP non accessible"

echo ""
echo "Test 2: Mailserver SMTP (587)"
timeout 3 bash -c '</dev/tcp/127.0.0.1/587' 2>/dev/null && echo "✅ SMTP accessible" || echo "❌ SMTP non accessible"

echo ""
echo "Test 3: Synapse Homeserver (8008)"
SYNAPSE_TEST=$(curl -s http://localhost:8008/_matrix/client/versions 2>/dev/null | grep -o "versions" | head -1)
if [ "$SYNAPSE_TEST" = "versions" ]; then
    echo "✅ Synapse homeserver répond"
else
    echo "❌ Synapse homeserver ne répond pas"
fi

echo ""
echo "Test 4: Element UI (8083)"
ELEMENT_TEST=$(curl -s http://localhost:8083 2>/dev/null | grep -o "element" | head -1)
if [ -n "$ELEMENT_TEST" ]; then
    echo "✅ Element accessible"
else
    echo "❌ Element non accessible"
fi

echo ""
echo "=========================================="
echo "✅ RÉPARATION TERMINÉE"
echo "=========================================="
echo ""
echo "🎯 PROCHAINES ÉTAPES:"
echo ""
echo "1. Vérifier que mailserver est 'Up' (pas 'Restarting')"
echo "   docker compose ps mailserver"
echo ""
echo "2. Vérifier que Synapse est 'Up (healthy)'"
echo "   docker compose ps synapse"
echo ""
echo "3. Tester Roundcube:"
echo "   https://mail.office1789.com"
echo ""
echo "4. Tester Matrix:"
echo "   https://chat.office1789.com"
echo ""
echo "5. Si problème persiste, exécuter:"
echo "   ./configure_backend_sso.sh"
echo ""
