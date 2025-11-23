#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import re

# Lire le fichier avec encodage UTF-8
with open('webfront2/src/traduction.json', 'r', encoding='utf-8') as f:
    content = f.read()

# Remplacer uniquement les @ dans les lignes yourEmail
content = re.sub(r'("yourEmail":\s*"[^"]*?)@([^"]*?")', r"\1{'@'}\2", content)

# Écrire le fichier avec encodage UTF-8
with open('webfront2/src/traduction.json', 'w', encoding='utf-8', newline='') as f:
    f.write(content)

print("✅ Fichier traduction.json corrigé avec encodage UTF-8 préservé")
