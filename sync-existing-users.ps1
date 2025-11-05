# Script pour synchroniser les comptes mail avec les utilisateurs existants

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Office1789 - Sync comptes mail" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

Write-Host "Ce script va créer des comptes mail pour tous les utilisateurs" -ForegroundColor Yellow
Write-Host "existants dans la base de données PostgreSQL." -ForegroundColor Yellow
Write-Host ""

$continue = Read-Host "Continuer? (O/N)"
if ($continue -ne "O" -and $continue -ne "o") {
    Write-Host "Annulé." -ForegroundColor Red
    exit 0
}

Write-Host ""
Write-Host "Compilation du script..." -ForegroundColor Cyan
cd backend
go run sync-mail-accounts.go

cd ..

Write-Host ""
Write-Host "Terminé !" -ForegroundColor Green
