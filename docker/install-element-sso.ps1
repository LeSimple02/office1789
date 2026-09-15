# Script PowerShell pour installer le plugin SSO Office1789 dans Element

Write-Host "🔧 Installation du plugin SSO Office1789 dans Element..." -ForegroundColor Cyan

# Copier le script SSO dans le conteneur Element
docker cp "docker\element\office1789-sso.js" element:/app/office1789-sso.js

# Injecter le script dans index.html d'Element
docker exec element bash -c "sed -i 's|</head>|<script src=\"/office1789-sso.js\"></script></head>|' /app/index.html"

# Redémarrer Element pour appliquer les changements
docker restart element

Write-Host "✅ Plugin SSO installé avec succès !" -ForegroundColor Green
Write-Host "📝 Element va redémarrer (10-15 secondes)..." -ForegroundColor Yellow
