# Script pour corriger toutes les validations de session

$files = @("drive.go", "chat.go")

foreach ($file in $files) {
    $content = Get-Content $file -Raw
    
    # Pattern 1: session, valid := validateSession(...); if !valid
    $content = $content -replace 'session, valid := validateSession\(([^;]+)\); if !valid', 'session, valid := validateSession($1)`nif !valid'
    
    # Pattern 2: Supprimer les var userID int inutiles après validateSession
    $content = $content -replace '(?m)^\s*var userID int\s*$\s*if err := db\.QueryRow\("SELECT user_id FROM Users WHERE username=\$1", [^)]+\)\.Scan\(&userID\); err != nil \{\s*c\.JSON\([^}]+\}\s*return\s*\}', ''
    
    # Pattern 3: Remplacer userID par session.UserID
    # (sera fait manuellement car trop complexe)
    
    Set-Content -Path $file -Value $content
}

Write-Host "Fichiers corrigés!"
