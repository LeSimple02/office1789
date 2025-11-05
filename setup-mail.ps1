# Script pour créer des comptes mail dans Office1789
# Usage: .\setup-mail.ps1

Write-Host "==================================" -ForegroundColor Cyan
Write-Host "Office1789 - Configuration Mail" -ForegroundColor Cyan
Write-Host "==================================" -ForegroundColor Cyan
Write-Host ""

# Vérifier que Docker est en cours d'exécution
$dockerRunning = docker info 2>$null
if (-not $?) {
    Write-Host "❌ Docker n'est pas en cours d'exécution !" -ForegroundColor Red
    Write-Host "Démarrez Docker Desktop et réessayez." -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Docker est en cours d'exécution" -ForegroundColor Green
Write-Host ""

# Vérifier que le container mailserver existe
$mailserverExists = docker ps -a --filter "name=mailserver" --format "{{.Names}}"
if (-not $mailserverExists) {
    Write-Host "❌ Le container 'mailserver' n'existe pas !" -ForegroundColor Red
    Write-Host "Lancez d'abord: docker-compose up -d" -ForegroundColor Yellow
    exit 1
}

# Vérifier que le container mailserver est démarré
$mailserverRunning = docker ps --filter "name=mailserver" --format "{{.Names}}"
if (-not $mailserverRunning) {
    Write-Host "⚠️  Le container mailserver n'est pas démarré. Démarrage..." -ForegroundColor Yellow
    docker start mailserver
    Start-Sleep -Seconds 5
}

Write-Host "✅ Serveur mail prêt" -ForegroundColor Green
Write-Host ""

# Menu principal
Write-Host "Que voulez-vous faire ?" -ForegroundColor Cyan
Write-Host "1. Créer un nouveau compte mail"
Write-Host "2. Lister les comptes mail existants"
Write-Host "3. Supprimer un compte mail"
Write-Host "4. Créer les comptes par défaut (jean, matthis)"
Write-Host "5. Quitter"
Write-Host ""

$choice = Read-Host "Votre choix (1-5)"

switch ($choice) {
    "1" {
        Write-Host ""
        $email = Read-Host "Adresse email (ex: utilisateur@office1789.local)"
        $password = Read-Host "Mot de passe" -AsSecureString
        $plainPassword = [System.Runtime.InteropServices.Marshal]::PtrToStringAuto([System.Runtime.InteropServices.Marshal]::SecureStringToBSTR($password))
        
        Write-Host "Création du compte $email..." -ForegroundColor Yellow
        docker exec -it mailserver setup email add $email $plainPassword
        
        if ($?) {
            Write-Host "✅ Compte créé avec succès !" -ForegroundColor Green
            Write-Host ""
            Write-Host "Vous pouvez maintenant vous connecter à Roundcube avec :" -ForegroundColor Cyan
            Write-Host "  Email: $email" -ForegroundColor White
            Write-Host "  URL: http://localhost:8081" -ForegroundColor White
        } else {
            Write-Host "❌ Erreur lors de la création du compte" -ForegroundColor Red
        }
    }
    
    "2" {
        Write-Host ""
        Write-Host "📋 Liste des comptes mail :" -ForegroundColor Cyan
        docker exec mailserver setup email list
    }
    
    "3" {
        Write-Host ""
        $email = Read-Host "Adresse email à supprimer"
        Write-Host "Suppression du compte $email..." -ForegroundColor Yellow
        docker exec -it mailserver setup email del $email
        
        if ($?) {
            Write-Host "✅ Compte supprimé avec succès !" -ForegroundColor Green
        } else {
            Write-Host "❌ Erreur lors de la suppression" -ForegroundColor Red
        }
    }
    
    "4" {
        Write-Host ""
        Write-Host "Création des comptes par défaut..." -ForegroundColor Yellow
        Write-Host ""
        
        # Compte Jean
        Write-Host "Création: jean@office1789.local" -ForegroundColor Cyan
        docker exec mailserver setup email add jean@office1789.local "jean1789"
        
        # Compte Matthis
        Write-Host "Création: matthis@office1789.local" -ForegroundColor Cyan
        docker exec mailserver setup email add matthis@office1789.local "matthis1789"
        
        Write-Host ""
        Write-Host "✅ Comptes créés avec succès !" -ForegroundColor Green
        Write-Host ""
        Write-Host "📋 Identifiants créés :" -ForegroundColor Cyan
        Write-Host "  jean@office1789.local / jean1789" -ForegroundColor White
        Write-Host "  matthis@office1789.local / matthis1789" -ForegroundColor White
        Write-Host ""
        Write-Host "🌐 Connectez-vous sur: http://localhost:8081" -ForegroundColor Cyan
    }
    
    "5" {
        Write-Host "Au revoir !" -ForegroundColor Cyan
        exit 0
    }
    
    default {
        Write-Host "❌ Choix invalide" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "Appuyez sur une touche pour continuer..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
