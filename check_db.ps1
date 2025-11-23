# Script PowerShell pour vérifier les tables dans les 3 bases PostgreSQL

Write-Host "======================================================" -ForegroundColor Cyan
Write-Host "🔍 VÉRIFICATION DES BASES DE DONNÉES PostgreSQL" -ForegroundColor Cyan
Write-Host "======================================================" -ForegroundColor Cyan
Write-Host ""

# Fonction pour vérifier une base de données
function Check-Database {
    param(
        [string]$Container,
        [string]$DbName,
        [string]$DbUser,
        [string]$Title
    )
    
    Write-Host "📦 $Title" -ForegroundColor Yellow
    Write-Host "   Container: $Container"
    Write-Host "   Database: $DbName"
    Write-Host ""
    
    # Vérifier si le conteneur existe et est en cours d'exécution
    $containerExists = docker ps --format "{{.Names}}" | Select-String -Pattern "^$Container$"
    
    if (-not $containerExists) {
        Write-Host "   ❌ Conteneur non trouvé ou arrêté" -ForegroundColor Red
        Write-Host ""
        return
    }
    
    # Lister les tables
    Write-Host "   📋 Tables trouvées:"
    $tables = docker exec $Container psql -U $DbUser -d $DbName -c "\dt" 2>$null
    
    if ($tables) {
        $tables | Select-String -Pattern "^ " | ForEach-Object { Write-Host "   $_" }
    } else {
        Write-Host "   (aucune table ou erreur de connexion)"
    }
    
    # Compter les tables
    $tableCount = docker exec $Container psql -U $DbUser -d $DbName -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>$null
    
    if ($tableCount) {
        $tableCount = $tableCount.Trim()
        if ([int]$tableCount -gt 0) {
            Write-Host "   ✅ $tableCount table(s) trouvée(s)" -ForegroundColor Green
        } else {
            Write-Host "   ❌ Aucune table trouvée" -ForegroundColor Red
        }
    } else {
        Write-Host "   ❌ Erreur de connexion" -ForegroundColor Red
    }
    
    Write-Host ""
    Write-Host "------------------------------------------------------"
    Write-Host ""
}

# Vérifier postgres_db (base principale Office1789)
Check-Database -Container "postgres_db" -DbName "office1789" -DbUser "office1789user" -Title "BASE PRINCIPALE (office1789)"

# Vérifier postgres_roundcube
Check-Database -Container "postgres_roundcube" -DbName "roundcube" -DbUser "roundcube" -Title "BASE ROUNDCUBE (Webmail)"

# Vérifier postgres_synapse
Check-Database -Container "postgres_synapse" -DbName "synapse" -DbUser "synapse_user" -Title "BASE SYNAPSE (Matrix)"

Write-Host "======================================================" -ForegroundColor Cyan
Write-Host "🏁 VÉRIFICATION TERMINÉE" -ForegroundColor Cyan
Write-Host "======================================================" -ForegroundColor Cyan
