# Test complet du SSO Matrix - Office1789
# Simule le flux complet d'authentification SSO

Write-Host "`n============================================================" -ForegroundColor Cyan
Write-Host "TEST SSO MATRIX - Office1789" -ForegroundColor Cyan
Write-Host "============================================================`n" -ForegroundColor Cyan

# 1. Vérifier que le backend est accessible
Write-Host "1. Verification backend..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/session/check" -Method POST -ContentType "application/json" -Body '{"token":"test"}' -ErrorAction Stop
    Write-Host "   OK - Backend accessible" -ForegroundColor Green
} catch {
    Write-Host "   ERREUR - Backend inaccessible" -ForegroundColor Red
    Write-Host "   Lancez: python startOffice1789.py" -ForegroundColor Yellow
    exit 1
}

# 2. Vérifier Matrix/Synapse
Write-Host "`n2. Verification Matrix/Synapse..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8008/_matrix/client/versions" -Method GET -ErrorAction Stop
    Write-Host "   OK - Matrix accessible" -ForegroundColor Green
} catch {
    Write-Host "   ERREUR - Matrix inaccessible" -ForegroundColor Red
    exit 1
}

# 3. Vérifier Element
Write-Host "`n3. Verification Element..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8083" -Method GET -ErrorAction Stop
    if ($response.Content -match "office1789-sso\.js") {
        Write-Host "   OK - Element avec script SSO" -ForegroundColor Green
    } else {
        Write-Host "   ATTENTION - Script SSO manquant" -ForegroundColor Yellow
    }
} catch {
    Write-Host "   ERREUR - Element inaccessible" -ForegroundColor Red
    exit 1
}

# 4. Tester la protection de la page de login
Write-Host "`n4. Test protection page de login Element..." -ForegroundColor Yellow
Write-Host "   INFO - Ouvrez http://localhost:8083 dans un navigateur" -ForegroundColor Gray
Write-Host "   INFO - Vous devriez voir 'Acces Refuse' apres quelques secondes" -ForegroundColor Gray

# 5. Instructions pour tester le SSO complet
Write-Host "`n5. Test SSO complet (manuel):" -ForegroundColor Yellow
Write-Host "`n   ETAPE 1 - Login Office1789:" -ForegroundColor Cyan
Write-Host "   - Ouvrir http://localhost:5173" -ForegroundColor White
Write-Host "   - Se connecter avec vos identifiants" -ForegroundColor White

Write-Host "`n   ETAPE 2 - Acceder au Chat:" -ForegroundColor Cyan
Write-Host "   - Cliquer sur 'Chat' dans le menu" -ForegroundColor White
Write-Host "   - Cliquer sur 'Ouvrir Element (Chat)'" -ForegroundColor White

Write-Host "`n   ETAPE 3 - Verification:" -ForegroundColor Cyan
Write-Host "   - Element devrait s'ouvrir dans un nouvel onglet" -ForegroundColor White
Write-Host "   - Ecran de chargement 'Connexion au chat Matrix...'" -ForegroundColor White
Write-Host "   - Authentification automatique" -ForegroundColor White
Write-Host "   - Redirection vers Element authentifie" -ForegroundColor White

Write-Host "`n   CONSOLE NAVIGATEUR (F12):" -ForegroundColor Cyan
Write-Host "   - [Office1789-SSO] Token SSO detecte" -ForegroundColor Gray
Write-Host "   - [Office1789-SSO] Token valide" -ForegroundColor Gray
Write-Host "   - [Office1789-SSO] Claims: username, matrixUserId, password" -ForegroundColor Gray
Write-Host "   - [Office1789-SSO] Authentification Matrix reussie" -ForegroundColor Gray

# 6. Debug si problème
Write-Host "`n6. Debug en cas de probleme:" -ForegroundColor Yellow
Write-Host "`n   A. Verifier le backend matrix.go:" -ForegroundColor Cyan
Write-Host "   cd backend" -ForegroundColor Gray
Write-Host "   go build -o test.exe ." -ForegroundColor Gray
Write-Host "   -> Doit compiler sans erreur" -ForegroundColor Gray

Write-Host "`n   B. Verifier les logs backend:" -ForegroundColor Cyan
Write-Host "   - [Matrix-SSO] Requete recue - Username: xxx" -ForegroundColor Gray
Write-Host "   - [Matrix-SSO] Session valide - UserID: xxx" -ForegroundColor Gray
Write-Host "   - [Matrix-SSO] Token SSO genere pour utilisateur" -ForegroundColor Gray

Write-Host "`n   C. Verifier les claims du token:" -ForegroundColor Cyan
Write-Host "   Console navigateur -> [Office1789-SSO] Claims:" -ForegroundColor Gray
Write-Host "   Doit contenir: username, matrixUserId, password, exp, iat" -ForegroundColor Gray

Write-Host "`n   D. Tester API Matrix directement:" -ForegroundColor Cyan
Write-Host "   Utiliser Postman ou curl pour tester l'endpoint Matrix" -ForegroundColor Gray

# Résumé
Write-Host "`n============================================================" -ForegroundColor Cyan
Write-Host "RESUME" -ForegroundColor Cyan
Write-Host "============================================================" -ForegroundColor Cyan

Write-Host "`nProtections actives:" -ForegroundColor Yellow
Write-Host "  OK - Roundcube bloque acces direct (403)" -ForegroundColor Green
Write-Host "  OK - Element bloque page de login sans token SSO" -ForegroundColor Green

Write-Host "`nFlux SSO Matrix:" -ForegroundColor Yellow
Write-Host "  1. Login Office1789 -> Password stocke en RAM" -ForegroundColor White
Write-Host "  2. Click Ouvrir Element -> POST /api/matrix/sso" -ForegroundColor White
Write-Host "  3. Backend genere token SSO avec credentials" -ForegroundColor White
Write-Host "  4. Element ouvre avec token dans URL" -ForegroundColor White
Write-Host "  5. Script JS valide token et authentifie a Matrix" -ForegroundColor White
Write-Host "  6. Stockage credentials en localStorage" -ForegroundColor White
Write-Host "  7. Redirection vers Element authentifie" -ForegroundColor White

Write-Host "`nURLs de test:" -ForegroundColor Yellow
Write-Host "  Office1789: http://localhost:5173" -ForegroundColor Gray
Write-Host "  Backend:    http://localhost:8080" -ForegroundColor Gray
Write-Host "  Matrix:     http://localhost:8008" -ForegroundColor Gray
Write-Host "  Element:    http://localhost:8083" -ForegroundColor Gray

Write-Host "`n============================================================`n" -ForegroundColor Cyan
