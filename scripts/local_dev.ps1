# Script de développement local Windows
# Usage: .\scripts\local_dev.ps1

Write-Host "🚀 Démarrage environnement local Office1789..." -ForegroundColor Green
Write-Host ""

$ROOT_DIR = Split-Path -Parent (Split-Path -Parent $PSCommandPath)
$DOCKER_DIR = Join-Path $ROOT_DIR "docker"
$BACKEND_DIR = Join-Path $ROOT_DIR "backend"
$FRONTEND_DIR = Join-Path $ROOT_DIR "webfront2"

Set-Location $ROOT_DIR

# Vérifier Docker
Write-Host "🐳 Vérification Docker..." -ForegroundColor Cyan
$dockerRunning = docker info 2>$null
if (-not $dockerRunning) {
    Write-Host "❌ Docker n'est pas démarré. Lancez Docker Desktop." -ForegroundColor Red
    exit 1
}
Write-Host "✅ Docker actif" -ForegroundColor Green
Write-Host ""

# Creer .env si necessaire
Set-Location $DOCKER_DIR
Write-Host "🧩 Verification .env Docker..." -ForegroundColor Cyan
if (-not (Test-Path ".env")) {
    if (Test-Path ".env.example") {
        Write-Host "📝 Creation .env depuis .env.example..." -ForegroundColor Yellow
        Copy-Item ".env.example" ".env"
        
        # Generer des mots de passe simples pour le dev local
        $content = Get-Content ".env" -Raw
        $content = $content -replace "CHANGE_ME_MAIN_DB_PASSWORD", "devpass123"
        $content = $content -replace "CHANGE_ME_ROUNDCUBE_DB_PASSWORD", "devpass123"
        $content = $content -replace "CHANGE_ME_SYNAPSE_DB_PASSWORD", "devpass123"
        $content = $content -replace "CHANGE_ME_COTURN_SECRET", "devpass123"
        $content = $content -replace "CHANGE_ME_ONLYOFFICE_JWT_SECRET", "devpass123"
        $content = $content -replace "CHANGE_ME_MATRIX_ADMIN_TOKEN", "devpass123"
        $content = $content -replace "CHANGE_ME_MAIL_ADMIN_PASSWORD", "devpass123"
        
        Set-Content ".env" $content
        Write-Host "✅ .env cree avec mots de passe dev" -ForegroundColor Green
    } else {
        Write-Host "❌ .env manquant et pas de .env.example" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "✅ .env Docker existe deja" -ForegroundColor Green
}
Write-Host ""

# Creer backend/.env
Set-Location $BACKEND_DIR
Write-Host "🔧 Configuration backend/.env..." -ForegroundColor Cyan
$envContent = @"
DB_HOST=localhost
DB_USER=office1789
DB_PASSWORD=devpass123
DB_NAME=office1789db
DB_PORT=5432
JWT_SECRET=dev_secret_key_local_only
"@
Set-Content ".env" $envContent -Encoding UTF8
Write-Host "✅ backend/.env cree" -ForegroundColor Green
Write-Host ""

# Demarrer les services Docker
Set-Location $DOCKER_DIR
Write-Host "🐳 Demarrage services Docker..." -ForegroundColor Cyan
Write-Host "   (postgres, mailserver, roundcube, synapse, element, onlyoffice, coturn)" -ForegroundColor Gray
docker compose up -d postgres_db mailserver postgres_roundcube roundcube postgres_synapse synapse element onlyoffice coturn
Write-Host ""

# Attendre que Postgres soit pret
Write-Host "⏳ Attente PostgreSQL..." -ForegroundColor Cyan
Start-Sleep -Seconds 5
$retries = 0
while ($retries -lt 30) {
    $pgReady = docker compose exec -T postgres_db pg_isready -U office1789 2>$null
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ PostgreSQL pret" -ForegroundColor Green
        break
    }
    Start-Sleep -Seconds 1
    $retries++
}
Write-Host ""

# Instructions pour backend et frontend
Write-Host "=" -NoNewline -ForegroundColor Yellow
Write-Host ("=" * 70) -ForegroundColor Yellow
Write-Host "✅ Services Docker demarres !" -ForegroundColor Green
Write-Host ""
Write-Host "Pour lancer le BACKEND (dans un nouveau terminal):" -ForegroundColor Cyan
Write-Host "  cd backend" -ForegroundColor White
Write-Host "  go run ." -ForegroundColor White
Write-Host "  → Backend sur http://localhost:8080" -ForegroundColor Gray
Write-Host ""
Write-Host "Pour lancer le FRONTEND (dans un autre terminal):" -ForegroundColor Cyan
Write-Host "  cd webfront2" -ForegroundColor White
Write-Host "  npm install" -ForegroundColor White
Write-Host "  npm run dev" -ForegroundColor White
Write-Host "  → Frontend sur http://localhost:5173" -ForegroundColor Gray
Write-Host ""
Write-Host "Services disponibles:" -ForegroundColor Cyan
Write-Host "  • Roundcube:   http://localhost:8081" -ForegroundColor White
Write-Host "  • Element:     http://localhost:8083" -ForegroundColor White
Write-Host "  • OnlyOffice:  http://localhost:8082" -ForegroundColor White
Write-Host ""
Write-Host "Pour arreter:" -ForegroundColor Cyan
Write-Host "  cd docker" -ForegroundColor White
Write-Host "  docker compose down" -ForegroundColor White
Write-Host "=" -NoNewline -ForegroundColor Yellow
Write-Host ("=" * 70) -ForegroundColor Yellow
