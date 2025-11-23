#!/usr/bin/env python3
# -*- coding: utf-8 -*-

import json
import re

# Read the file
with open('webfront2/src/traduction.json', 'r', encoding='utf-8') as f:
    content = f.read()

# Fix all occurrences of {'@'} syntax
# Pattern 1: ({'@'}office1789.com) -> (@office1789.com)
content = content.replace("({'@'}office1789.com)", "(@office1789.com)")

# Pattern 2: votre{'{'@'}'}email.com -> votre@email.com
content = re.sub(r"(\w+)\{'{'@'\}'\}email\.com", r"\1@email.com", content)

# Pattern 3: any variations of {'@'} -> @
content = content.replace("{'@'}", "@")
content = content.replace("{'{'@'}'}", "@")

# Write back with UTF-8 encoding
with open('webfront2/src/traduction.json', 'w', encoding='utf-8') as f:
    f.write(content)

print("✅ All translations fixed! @ symbols are now hardcoded.")
print("Fixed patterns:")
print("  - ({'@'}office1789.com) → (@office1789.com)")
print("  - text{'{'@'}'}email.com → text@email.com")
print("  - {'@'} → @")
