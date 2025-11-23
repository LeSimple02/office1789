# Script de developpement local Windows
# Usage: .\scripts\local_dev.ps1

Write-Host "Demarrage environnement local Office1789..." -ForegroundColor Green
Write-Host ""

$ROOT_DIR = Split-Path -Parent (Split-Path -Parent $PSCommandPath)
$DOCKER_DIR = Join-Path $ROOT_DIR "docker"
$BACKEND_DIR = Join-Path $ROOT_DIR "backend"

Set-Location $ROOT_DIR

# Verifier Docker
Write-Host "Verification Docker..." -ForegroundColor Cyan
$dockerRunning = docker info 2>$null
if (-not $dockerRunning) {
    Write-Host "Docker n est pas demarre" -ForegroundColor Red
    exit 1
}
Write-Host "Docker actif" -ForegroundColor Green
Write-Host ""

# Creer .env
Set-Location $DOCKER_DIR
if (-not (Test-Path ".env")) {
    if (Test-Path ".env.example") {
        Copy-Item ".env.example" ".env"
        $content = Get-Content ".env" -Raw
        $content = $content -replace "CHANGE_ME_MAIN_DB_PASSWORD", "devpass123"
        $content = $content -replace "CHANGE_ME_ROUNDCUBE_DB_PASSWORD", "devpass123"
        $content = $content -replace "CHANGE_ME_SYNAPSE_DB_PASSWORD", "devpass123"
        $content = $content -replace "CHANGE_ME_COTURN_SECRET", "devpass123"
        $content = $content -replace "CHANGE_ME_ONLYOFFICE_JWT_SECRET", "devpass123"
        $content = $content -replace "CHANGE_ME_MATRIX_ADMIN_TOKEN", "devpass123"
        $content = $content -replace "CHANGE_ME_MAIL_ADMIN_PASSWORD", "devpass123"
        Set-Content ".env" $content -Encoding UTF8
        Write-Host ".env cree" -ForegroundColor Green
    }
}

# Creer backend/.env
Set-Location $BACKEND_DIR
@"
DB_HOST=localhost
DB_USER=office1789
DB_PASSWORD=devpass123
DB_NAME=office1789db
DB_PORT=5432
JWT_SECRET=dev_secret_key_local_only
"@ | Set-Content ".env" -Encoding UTF8

# Demarrer Docker
Set-Location $DOCKER_DIR
docker compose up -d postgres_db mailserver postgres_roundcube roundcube postgres_synapse synapse element onlyoffice coturn

Write-Host ""
Write-Host "Services demarres !" -ForegroundColor Green
Write-Host "Backend: cd backend ; go run ." -ForegroundColor Cyan
Write-Host "Frontend: cd webfront2 ; npm run dev" -ForegroundColor Cyan
