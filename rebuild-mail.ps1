# Script de rebuild complet pour Office1789 avec SSO et design custom

Write-Host "🔨 Reconstruction de l'environnement Office1789..." -ForegroundColor Cyan

# Arrêter les containers
Write-Host "`n📦 Arrêt des containers..." -ForegroundColor Yellow
Set-Location docker
docker-compose down

# Reconstruire Roundcube avec le plugin SSO et le design
Write-Host "`n🔧 Reconstruction de Roundcube avec SSO et design Office1789..." -ForegroundColor Yellow
docker-compose build roundcube

# Redémarrer tous les services
Write-Host "`n🚀 Démarrage des services..." -ForegroundColor Yellow
docker-compose up -d

# Attendre que les services soient prêts
Write-Host "`n⏳ Attente du démarrage des services..." -ForegroundColor Yellow
Start-Sleep -Seconds 15

# Vérifier le statut
Write-Host "`n✅ Statut des containers:" -ForegroundColor Green
docker-compose ps

Write-Host "`n📧 Vérification des comptes mail..." -ForegroundColor Yellow
docker exec mailserver setup email list

Write-Host "`n✨ Reconstruction terminée !" -ForegroundColor Green
Write-Host "`n📝 Prochaines étapes:" -ForegroundColor Cyan
Write-Host "  1. Démarrer le backend Go: cd .. && cd backend && go run ." -ForegroundColor White
Write-Host "  2. Démarrer le frontend Vue: cd .. && cd webfront2 && npm run dev" -ForegroundColor White
Write-Host "  3. Se connecter sur http://localhost:5173" -ForegroundColor White
Write-Host "  4. Aller dans Mail - la connexion SSO devrait être automatique!" -ForegroundColor White
Write-Host "`n🎨 Le design Office1789 avec dégradé bleu-rouge et mode sombre est appliqué!" -ForegroundColor Magenta
