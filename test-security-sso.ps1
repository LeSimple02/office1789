# Test de sécurité SSO - Office1789
# Vérifie que Roundcube ET Element bloquent l'accès direct

Write-Host "`n============================================================" -ForegroundColor Cyan
Write-Host "TEST DE SECURITE SSO - Office1789" -ForegroundColor Cyan
Write-Host "============================================================`n" -ForegroundColor Cyan

# Test 1: Roundcube
Write-Host "1. Test Roundcube (http://localhost:8081)..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8081" -Method GET -ErrorAction Stop
    if ($response.Content -match "Accès Refusé|403") {
        Write-Host "   ERREUR - Roundcube devrait bloquer mais retourne 200" -ForegroundColor Red
    } else {
        Write-Host "   ERREUR - Roundcube accessible sans token SSO!" -ForegroundColor Red
    }
} catch {
    if ($_.Exception.Message -match "403|Forbidden|Accès Refusé") {
        Write-Host "   OK - Roundcube bloque l'accès (403 Forbidden)" -ForegroundColor Green
    } else {
        Write-Host "   ERREUR - Erreur inattendue: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# Test 2: Element
Write-Host "`n2. Test Element (http://localhost:8083)..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8083" -Method GET -ErrorAction Stop
    
    # Element retourne 200 mais le JavaScript bloque côté client
    if ($response.Content -match "office1789-sso\.js") {
        Write-Host "   OK - Script SSO présent dans Element" -ForegroundColor Green
        
        # Vérifier que le script contient la protection
        $scriptContent = Invoke-WebRequest -Uri "http://localhost:8083/office1789-sso.js" -Method GET
        if ($scriptContent.Content -match "Accès refusé|checkAccess|isAuthenticated") {
            Write-Host "   OK - Script SSO contient la protection" -ForegroundColor Green
        } else {
            Write-Host "   ATTENTION - Script SSO peut ne pas bloquer correctement" -ForegroundColor Yellow
        }
    } else {
        Write-Host "   ERREUR - Script SSO manquant dans Element" -ForegroundColor Red
    }
} catch {
    Write-Host "   ERREUR - Element inaccessible: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 3: Vérifier le script SSO dans le conteneur
Write-Host "`n3. Vérification du script SSO dans Element..." -ForegroundColor Yellow
try {
    $scriptCheck = docker exec element cat /app/office1789-sso.js 2>&1
    if ($LASTEXITCODE -eq 0) {
        if ($scriptCheck -match "Accès refusé|ssoToken|isAuthenticated") {
            Write-Host "   OK - Script SSO avec protection présent" -ForegroundColor Green
        } else {
            Write-Host "   ERREUR - Script SSO sans protection" -ForegroundColor Red
        }
    } else {
        Write-Host "   ERREUR - Script SSO introuvable" -ForegroundColor Red
    }
} catch {
    Write-Host "   ERREUR - Impossible de vérifier: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: Roundcube - vérifier le plugin PHP
Write-Host "`n4. Vérification du plugin Roundcube..." -ForegroundColor Yellow
try {
    $pluginCheck = docker exec roundcube cat /var/www/html/plugins/office1789_sso/office1789_sso.php 2>&1
    if ($LASTEXITCODE -eq 0) {
        if ($pluginCheck -match "SÉCURITÉ.*Accès refusé|empty.*sso_token.*user_id") {
            Write-Host "   OK - Plugin Roundcube avec protection présent" -ForegroundColor Green
        } else {
            Write-Host "   ATTENTION - Plugin Roundcube peut manquer de protection" -ForegroundColor Yellow
        }
    } else {
        Write-Host "   ERREUR - Plugin Roundcube introuvable" -ForegroundColor Red
    }
} catch {
    Write-Host "   ERREUR - Impossible de vérifier: $($_.Exception.Message)" -ForegroundColor Red
}

# Résumé
Write-Host "`n============================================================" -ForegroundColor Cyan
Write-Host "RESUME DES TESTS" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan

Write-Host "`nProtection attendue:" -ForegroundColor Yellow
Write-Host "  - Roundcube: Bloque avec 403 Forbidden (protection PHP)" -ForegroundColor White
Write-Host "  - Element: Retourne 200 mais JavaScript bloque immédiatement" -ForegroundColor White

Write-Host "`nTest manuel recommandé:" -ForegroundColor Yellow
Write-Host "  1. Ouvrir http://localhost:8081 dans un navigateur" -ForegroundColor Gray
Write-Host "     -> Devrait afficher 'Accès Refusé'" -ForegroundColor Gray
Write-Host "  2. Ouvrir http://localhost:8083 dans un navigateur" -ForegroundColor Gray
Write-Host "     -> Devrait afficher 'Accès Refusé' (page blanche)" -ForegroundColor Gray
Write-Host "  3. Vérifier la console navigateur (F12)" -ForegroundColor Gray
Write-Host "     -> '[Office1789-SSO] Accès refusé sans token SSO'" -ForegroundColor Gray

Write-Host "`n============================================================`n" -ForegroundColor Cyan
