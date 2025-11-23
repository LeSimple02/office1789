#!/bin/bash
# Configuration complète du backend pour SSO Roundcube et Matrix

echo "=========================================="
echo "⚙️  CONFIGURATION BACKEND SSO"
echo "=========================================="

BACKEND_ENV="/home/debian/office1789/backend/.env"

if [ ! -f "$BACKEND_ENV" ]; then
    echo "❌ Fichier $BACKEND_ENV introuvable!"
    exit 1
fi

echo ""
echo "📋 Configuration actuelle:"
echo "=========================================="
grep -E "ROUNDCUBE_URL|ELEMENT_URL|MATRIX_" "$BACKEND_ENV" || echo "⚠️ Variables SSO non trouvées"

echo ""
echo "🔧 Mise à jour des variables SSO..."
echo "=========================================="

# Fonction pour ajouter ou mettre à jour une variable
update_env_var() {
    local key=$1
    local value=$2
    
    if grep -q "^${key}=" "$BACKEND_ENV"; then
        # Variable existe, la mettre à jour
        sudo sed -i "s|^${key}=.*|${key}=${value}|" "$BACKEND_ENV"
        echo "✅ ${key} mis à jour"
    else
        # Variable n'existe pas, l'ajouter
        echo "${key}=${value}" | sudo tee -a "$BACKEND_ENV" > /dev/null
        echo "✅ ${key} ajouté"
    fi
}

# Configurer les URLs pour SSO
update_env_var "ROUNDCUBE_URL" "https://mail.office1789.com"
update_env_var "ELEMENT_URL" "https://chat.office1789.com"

# Configurer Matrix homeserver pour l'authentification
update_env_var "MATRIX_HOMESERVER" "http://synapse:8008"
update_env_var "MATRIX_SERVER_NAME" "office1789.com"

# Secret partagé Matrix/Element (doit correspondre au plugin Element)
if ! grep -q "^MATRIX_SSO_SECRET=" "$BACKEND_ENV"; then
    SSO_SECRET="Office1789-Matrix-SecretKey-$(openssl rand -hex 16)"
    update_env_var "MATRIX_SSO_SECRET" "$SSO_SECRET"
fi

echo ""
echo "📊 Nouvelle configuration:"
echo "=========================================="
sudo cat "$BACKEND_ENV" | grep -E "ROUNDCUBE_URL|ELEMENT_URL|MATRIX_"

echo ""
echo "🔄 Redémarrage du backend..."
echo "=========================================="
sudo systemctl restart office1789-backend
sleep 3

echo ""
echo "✅ Status du backend:"
sudo systemctl status office1789-backend --no-pager | head -20

echo ""
echo "🧪 Test de l'API backend:"
echo "=========================================="

echo "1. Test /api/session/check:"
RESPONSE=$(curl -s -X POST https://backend.office1789.com/api/session/check \
  -H "Content-Type: application/json" \
  -d '{"username":"test","token":"test"}')
echo "$RESPONSE" | head -3

echo ""
echo "2. Logs backend (dernières lignes):"
sudo journalctl -u office1789-backend --since "1 minute ago" --no-pager | tail -10

echo ""
echo "=========================================="
echo "✅ CONFIGURATION TERMINÉE"
echo "=========================================="
echo ""
echo "📝 Prochaines étapes:"
echo "1. Vérifier que mailserver est UP: docker compose ps mailserver"
echo "2. Vérifier que synapse est UP: docker compose ps synapse"
echo "3. Tester SSO Roundcube depuis office1789.com (connecté)"
echo "4. Tester SSO Matrix depuis office1789.com (connecté)"
echo ""
echo "🔍 Debug si ça ne fonctionne pas:"
echo "   sudo journalctl -u office1789-backend -f"
echo ""
