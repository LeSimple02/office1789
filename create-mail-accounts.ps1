# Script pour creer les comptes mail dans le mailserver

Write-Host "======================================" -ForegroundColor Cyan
Write-Host "Creation des comptes mail" -ForegroundColor Cyan
Write-Host "======================================" -ForegroundColor Cyan
Write-Host ""

# Verifier que mailserver est en cours d'execution
$mailRunning = docker ps --filter "name=mailserver" --filter "status=running" --format "{{.Names}}"
if (-not $mailRunning) {
    Write-Host "Le container mailserver n'est pas demarre !" -ForegroundColor Red
    Write-Host "Lancez: cd docker; docker-compose up -d mailserver" -ForegroundColor Yellow
    exit 1
}

Write-Host "Mailserver est en cours d'execution" -ForegroundColor Green
Write-Host ""

# Creer les comptes mail
Write-Host "Creation du compte jean@office1789.local..." -ForegroundColor Cyan
docker exec mailserver setup email add jean@office1789.local password123

Write-Host ""
Write-Host "Creation du compte matthis@office1789.local..." -ForegroundColor Cyan
docker exec mailserver setup email add matthis@office1789.local password123

Write-Host ""
Write-Host "======================================" -ForegroundColor Green
Write-Host "Comptes mail crees avec succes !" -ForegroundColor Green
Write-Host "======================================" -ForegroundColor Green
Write-Host ""
Write-Host "Utilisateurs:" -ForegroundColor Cyan
Write-Host "  - jean@office1789.local (password: password123)" -ForegroundColor Yellow
Write-Host "  - matthis@office1789.local (password: password123)" -ForegroundColor Yellow
Write-Host ""
Write-Host "Vous pouvez maintenant vous connecter a Roundcube sur http://localhost:8081" -ForegroundColor Green
Write-Host ""
