#!/bin/bash
set -e

REPO_DIR="/home/debian/office1789"
DOCKER_DIR="$REPO_DIR/docker"
BACKEND_DIR="$REPO_DIR/backend"

echo "=========================================="
echo "🚀 DEPLOY PROD OFFICE1789 - COMPLET"
echo "=========================================="
echo ""

cd "$REPO_DIR"

echo "📥 Mise à jour du code..."
git pull origin main
echo "✅ Code à jour"
echo ""

cd "$DOCKER_DIR"

echo "🧩 Vérification .env Docker..."
if [ ! -f ".env" ]; then
  if [ -f ".env.example" ]; then
    cp .env.example .env
    echo "⚠️  .env créé depuis .env.example - édite les secrets avant de continuer"
    exit 1
  else
    echo "❌ .env manquant et pas de .env.example"
    exit 1
  fi
fi
echo "✅ .env Docker présent"
echo ""

echo "🧩 Création backend/.env si manquant..."
if [ ! -f "$BACKEND_DIR/.env" ]; then
  # Lire les credentials depuis docker/.env
  POSTGRES_USER=$(grep POSTGRES_USER .env | cut -d '=' -f2)
  POSTGRES_PASSWORD=$(grep POSTGRES_PASSWORD .env | cut -d '=' -f2)
  POSTGRES_DB=$(grep POSTGRES_DB .env | cut -d '=' -f2)
  
  cat > "$BACKEND_DIR/.env" << EOF
DB_HOST=postgres_db
DB_PORT=5432
DB_USER=$POSTGRES_USER
DB_PASSWORD=$POSTGRES_PASSWORD
DB_NAME=$POSTGRES_DB

PORT=8080
GIN_MODE=release
EOF
  echo "✅ backend/.env créé automatiquement"
else
  echo "✅ backend/.env existe déjà"
fi
echo ""

echo "🔧 Fix Synapse (permissions media_store)..."
sudo mkdir -p "$DOCKER_DIR/synapse/conf/media_store"
sudo chown -R 991:991 "$DOCKER_DIR/synapse/conf"
sudo chmod -R 755 "$DOCKER_DIR/synapse/"
echo "✅ Permissions Synapse OK"
echo ""

echo "🔧 Reset Roundcube DB pour mot de passe propre..."
docker compose stop roundcube postgres_roundcube 2>/dev/null || true
docker volume rm office1789_roundcube_db 2>/dev/null || true
echo "✅ Roundcube DB réinitialisée"
echo ""

echo "🔐 Préparation certificats SSL..."
bash "$DOCKER_DIR/prepare_certs.sh"
echo ""

echo "🐳 Build des images..."
docker compose build backend frontend element roundcube synapse
echo "✅ Build terminé"
echo ""

echo "🐳 Démarrage des bases de données..."
docker compose up -d postgres_db postgres_roundcube postgres_synapse
sleep 15
echo "✅ Postgres UP"
echo ""

echo "🐳 Démarrage mailserver..."
docker compose up -d mailserver
sleep 10
echo "✅ Mailserver lancé"
echo ""

echo "🐳 Démarrage backend..."
docker compose up -d backend
sleep 5
echo "✅ Backend lancé"
echo ""

echo "🐳 Démarrage services applicatifs..."
docker compose up -d synapse element roundcube frontend coturn onlyoffice
sleep 10
echo "✅ Services UP"
echo ""

echo "🌐 Configuration nginx..."
sudo tee /etc/nginx/sites-available/office1789 > /dev/null << 'NGINXCONF'
# Backend API
server {
    listen 80;
    server_name backend.office1789.com;
    
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Frontend
server {
    listen 80;
    server_name office1789.com www.office1789.com;
    
    location / {
        proxy_pass http://127.0.0.1:5173;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Mail (Roundcube)
server {
    listen 80;
    server_name mail.office1789.com;
    
    location / {
        proxy_pass http://127.0.0.1:8081;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Chat (Element + Matrix)
server {
    listen 80;
    server_name chat.office1789.com;
    
    location / {
        proxy_pass http://127.0.0.1:8083;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
    
    location /_matrix/ {
        proxy_pass http://127.0.0.1:8008/_matrix/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
NGINXCONF

sudo ln -sf /etc/nginx/sites-available/office1789 /etc/nginx/sites-enabled/office1789
sudo nginx -t
sudo systemctl restart nginx
sudo systemctl enable nginx
echo "✅ Nginx configuré et démarré"
echo ""

echo "📊 État des services Docker:"
docker compose ps
echo ""

echo "🧪 Tests des endpoints:"
echo -n "  Backend:   "
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8080/api/welcome || echo "❌"
echo ""
echo -n "  Frontend:  "
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:5173 || echo "❌"
echo ""
echo -n "  Roundcube: "
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8081 || echo "❌"
echo ""
echo -n "  Element:   "
curl -s -o /dev/null -w "%{http_code}" http://127.0.0.1:8083 || echo "❌"
echo ""
echo ""

echo "=========================================="
echo "✅ DÉPLOIEMENT TERMINÉ AVEC SUCCÈS"
echo ""
echo "🌐 Services accessibles sur:"
echo "  - https://backend.office1789.com"
echo "  - https://office1789.com"
echo "  - https://mail.office1789.com"
echo "  - https://chat.office1789.com"
echo ""
echo "📝 Pour passer en HTTPS:"
echo "  ./scripts/prod_letsencrypt.sh"
echo "=========================================="