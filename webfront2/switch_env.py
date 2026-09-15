import shutil
import sys
import os

"""
Usage:
  python switch_env.py local   # copie .env.local sur .env.production
  python switch_env.py prod    # restaure .env.production.bak si backup existe
"""

ROOT = os.path.dirname(os.path.abspath(__file__))
PROD = os.path.join(ROOT, '.env.production')
LOCAL = os.path.join(ROOT, '.env.local')
BAK = os.path.join(ROOT, '.env.production.bak')

if len(sys.argv) != 2 or sys.argv[1] not in ('local', 'prod'):
    print('Usage: python switch_env.py [local|prod]')
    sys.exit(1)

mode = sys.argv[1]

if mode == 'local':
    # Backup prod if not already
    if not os.path.exists(BAK):
        shutil.copy2(PROD, BAK)
    shutil.copy2(LOCAL, PROD)
    print('✅ .env.production remplacé par .env.local (localhost)')
elif mode == 'prod':
    if os.path.exists(BAK):
        shutil.copy2(BAK, PROD)
        print('✅ .env.production restauré depuis le backup (prod)')
    else:
        print('Aucun backup .env.production.bak trouvé, rien à faire.')
