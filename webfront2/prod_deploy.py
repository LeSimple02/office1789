import subprocess
import sys
import os

# Chemin absolu du dossier webfront2
ROOT = os.path.dirname(os.path.abspath(__file__))

# 1. Remettre les variables d'environnement prod
subprocess.check_call([sys.executable, os.path.join(ROOT, 'switch_env.py'), 'prod'])

# 2. Build prod (npm run build ou docker, à adapter selon ton workflow)
subprocess.check_call(['npm', 'install'], cwd=ROOT)
subprocess.check_call(['npm', 'run', 'build'], cwd=ROOT)

print('✅ Build prod terminé avec variables de production !')
