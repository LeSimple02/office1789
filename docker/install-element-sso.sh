#!/bin/bash
# Script pour installer le plugin SSO Office1789 dans Element

echo "🔧 Installation du plugin SSO Office1789 dans Element..."

# Copier le script SSO dans le conteneur Element
docker cp docker/element/office1789-sso.js element:/app/office1789-sso.js

# Injecter le script dans index.html d'Element
docker exec element bash -c "sed -i 's|</head>|<script src=\"/office1789-sso.js\"></script></head>|' /app/index.html"

# Redémarrer Element pour appliquer les changements
docker restart element

echo "✅ Plugin SSO installé avec succès !"
echo "📝 Element va redémarrer (10-15 secondes)..."
