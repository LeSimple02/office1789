import re

# Lire le fichier
with open('drive.go', 'r', encoding='utf-8') as f:
    content = f.read()

# Pattern 1: Supprimer les blocs qui récupèrent userID après validateSession
pattern1 = r'(\s+session, valid := validateSession\([^)]+\)\s*\n\s*if !valid \{[^}]+\})\s*\n\s*// Récupération user_id\s*\n\s*var userID int\s*\n\s*if err := db\.QueryRow\("SELECT user_id FROM Users WHERE username=\$1", [^)]+\)\.Scan\(&userID\); err != nil \{\s*c\.JSON\([^}]+\}\s*return\s*\}'
content = re.sub(pattern1, r'\1', content, flags=re.MULTILINE)

# Pattern 2: Remplacer userID par session.UserID dans les requêtes
content = re.sub(r'(?<!session\.)userID(?![a-zA-Z_])', 'session.UserID', content)

# Pattern 3: Corriger la ligne 301 (sessions[token])
content = re.sub(
    r'session, ok := sessions\[token\]',
    'session, err := getSessionFromDB(token)\n\tif err != nil || session == nil',
    content
)

# Écrire le fichier corrigé
with open('drive.go', 'w', encoding='utf-8') as f:
    f.write(content)

print("Fichier drive.go corrigé!")
