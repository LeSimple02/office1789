# Script pour changer le mot de passe d'un utilisateur Office1789
# Synchronise automatiquement : Office1789 + Mail + Matrix

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "Office1789 - Password Change Tool" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Ce script change le mot de passe sur tous les services :" -ForegroundColor Yellow
Write-Host "  - Office1789 (PostgreSQL)" -ForegroundColor Yellow
Write-Host "  - Mail (Docker Mailserver)" -ForegroundColor Yellow
Write-Host "  - Matrix (Synapse)" -ForegroundColor Yellow
Write-Host ""

# Vérifier si le binaire existe, sinon compiler
if (-not (Test-Path ".\change-password-tool.exe")) {
    Write-Host "Compilation du script..." -ForegroundColor Cyan
    go build -o change-password-tool.exe change-password-tool.go
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Erreur de compilation !" -ForegroundColor Red
        exit 1
    }
}

# Exécuter le script
.\change-password-tool.exe
