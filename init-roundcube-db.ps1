# Script pour initialiser la base de donnees Roundcube

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Initialisation de la base Roundcube" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Verifier que PostgreSQL est en cours d'execution
$pgRunning = docker ps --filter "name=postgres_db" --format "{{.Names}}"
if (-not $pgRunning) {
    Write-Host "Le container PostgreSQL n'est pas demarre !" -ForegroundColor Red
    Write-Host "Lancez: docker-compose up -d" -ForegroundColor Yellow
    exit 1
}

Write-Host "PostgreSQL est en cours d'execution" -ForegroundColor Green
Write-Host ""

# Copier le script SQL dans le container
Write-Host "Copie du script d'initialisation..." -ForegroundColor Cyan
docker cp docker/DockerfileDB/init-roundcube.sql postgres_db:/tmp/init-roundcube.sql

if (-not $?) {
    Write-Host "Erreur lors de la copie du script" -ForegroundColor Red
    exit 1
}

Write-Host "Script copie" -ForegroundColor Green
Write-Host ""

# Executer le script SQL
Write-Host "Creation des tables Roundcube..." -ForegroundColor Cyan
docker exec -i postgres_db psql -U robespierre -d office1789 -f /tmp/init-roundcube.sql

if ($?) {
    Write-Host ""
    Write-Host "Base de donnees Roundcube initialisee avec succes !" -ForegroundColor Green
    Write-Host ""
    Write-Host "Vous pouvez maintenant acceder a Roundcube sur http://localhost:8081" -ForegroundColor Cyan
} else {
    Write-Host ""
    Write-Host "Des erreurs peuvent apparaitre si les tables existent deja" -ForegroundColor Yellow
    Write-Host "Si c'est le cas, ignorez les erreurs." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "Appuyez sur une touche pour continuer..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
