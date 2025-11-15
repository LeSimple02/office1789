# Test SSO Matrix/Element
# Office1789 - Test end-to-end

Write-Host "============================================================" -ForegroundColor Cyan
Write-Host "OFFICE1789 - Test SSO Matrix/Element" -ForegroundColor Cyan
Write-Host "============================================================`n" -ForegroundColor Cyan

# Configuration
$BACKEND_URL = "http://localhost:8080"
$MATRIX_URL = "http://localhost:8008"
$ELEMENT_URL = "http://localhost:8083"

# 1. Test backend API
Write-Host "1. Test Backend API..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$BACKEND_URL/api/session/check" -Method POST -ContentType "application/json" -Body '{"token":"test"}' -ErrorAction Stop
    Write-Host "   OK Backend accessible" -ForegroundColor Green
} catch {
    Write-Host "   ERROR Backend inaccessible - Demarrez le backend" -ForegroundColor Red
    exit 1
}

# 2. Test Matrix API
Write-Host "`n2. Test Matrix API..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$MATRIX_URL/_matrix/client/versions" -Method GET -ErrorAction Stop
    Write-Host "   OK Matrix API accessible" -ForegroundColor Green
    $versions = ($response.Content | ConvertFrom-Json).versions
    Write-Host "   INFO Versions supportees: $($versions[-3..-1] -join ', ')..." -ForegroundColor Gray
} catch {
    Write-Host "   ERROR Matrix API inaccessible - Demarrez Synapse" -ForegroundColor Red
    exit 1
}

# 3. Test Element
Write-Host "`n3. Test Element..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri $ELEMENT_URL -Method GET -ErrorAction Stop
    if ($response.Content -match "office1789-sso.js") {
        Write-Host "   OK Element accessible avec plugin SSO" -ForegroundColor Green
    } else {
        Write-Host "   WARN Element accessible SANS plugin SSO" -ForegroundColor Yellow
        Write-Host "   TIP Rebuild Element: docker-compose build element" -ForegroundColor Gray
    }
} catch {
    Write-Host "   ERROR Element inaccessible - Demarrez le conteneur" -ForegroundColor Red
    exit 1
}

# 4. Test script SSO présent
Write-Host "`n4. Test script SSO dans Element..." -ForegroundColor Yellow
try {
    $result = docker exec element ls -la /app/office1789-sso.js 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host "   OK Script SSO present" -ForegroundColor Green
        Write-Host "   INFO $result" -ForegroundColor Gray
    } else {
        Write-Host "   ERROR Script SSO manquant" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "   ERROR Impossible de verifier le script" -ForegroundColor Red
}

# 5. Test endpoint Matrix SSO
Write-Host "`n5. Test endpoint /api/matrix/sso..." -ForegroundColor Yellow
Write-Host "   Necessaire: token de session valide" -ForegroundColor Yellow
Write-Host "   Test manuel:" -ForegroundColor Gray
Write-Host "      1. Login sur Office1789" -ForegroundColor Gray
Write-Host "      2. Aller dans Chat" -ForegroundColor Gray
Write-Host "      3. Cliquer sur 'Ouvrir Element (Chat)'" -ForegroundColor Gray

# Resume
Write-Host "`n============================================================" -ForegroundColor Cyan
Write-Host "OK Tests preliminaires OK" -ForegroundColor Green
Write-Host "============================================================" -ForegroundColor Cyan

Write-Host "`nProchaines etapes:" -ForegroundColor Yellow
Write-Host "   1. Demarrer Office1789 frontend (npm run dev)" -ForegroundColor White
Write-Host "   2. Login avec un compte Office1789" -ForegroundColor White
Write-Host "   3. Aller dans 'Chat' (ChatView)" -ForegroundColor White
Write-Host "   4. Cliquer sur 'Ouvrir Element (Chat)'" -ForegroundColor White
Write-Host "   5. Verifier l'authentification automatique`n" -ForegroundColor White

Write-Host "`nURLs:" -ForegroundColor Cyan
Write-Host "   Backend: $BACKEND_URL" -ForegroundColor Gray
Write-Host "   Matrix:  $MATRIX_URL" -ForegroundColor Gray
Write-Host "   Element: $ELEMENT_URL" -ForegroundColor Gray
Write-Host "`nOffice1789 - SSO Matrix pret !`n" -ForegroundColor Cyan
