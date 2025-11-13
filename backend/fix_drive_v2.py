# Script de correction robuste pour drive.go
import re

with open('drive.go', 'r', encoding='utf-8') as f:
    lines = f.readlines()

output = []
i = 0
while i < len(lines):
    line = lines[i]
    
    # Pattern 1: Remplacer if session, ok := sessions[xxx]
    if 'if session, ok := sessions[' in line:
        # Extraire les paramètres
        match = re.search(r'if session, ok := sessions\[([^\]]+)\]; !ok \|\| session\.Username != ([^\s]+) \|\| \2 == ""', line)
        if match:
            token = match.group(1)
            username = match.group(2)
            indent = re.match(r'(\s*)', line).group(1)
            output.append(f'{indent}session, valid := validateSession({token}, {username})\n')
            output.append(f'{indent}if !valid {{\n')
            i += 1
            continue
    
    # Pattern 2: Supprimer les blocs de récupération userID
    if 'var userID int' in line and i + 1 < len(lines):
        # Vérifier si c'est un bloc de récupération userID
        next_line = lines[i + 1]
        if 'db.QueryRow("SELECT user_id FROM Users WHERE username=' in next_line:
            # Sauter ce bloc (var userID int + les 4 lignes suivantes)
            i += 5
            continue
    
    # Pattern 3: Remplacer userID par session.UserID (sauf dans var userID int)
    if 'userID' in line and 'var userID int' not in line and 'var session.UserID' not in line:
        line = re.sub(r'\buserID\b', 'session.UserID', line)
    
    output.append(line)
    i += 1

# Écrire le résultat
with open('drive.go', 'w', encoding='utf-8') as f:
    f.writelines(output)

print("drive.go corrigé!")
