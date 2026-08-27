import subprocess
import sys
import os

# Chemin absolu du dossier webfront2
ROOT = os.path.dirname(os.path.abspath(__file__))

# 1. Basculer les variables d'environnement en mode local
subprocess.check_call([sys.executable, os.path.join(ROOT, 'switch_env.py'), 'local'])

# 2. Build local (npm run build ou docker, à adapter selon ton workflow)
# Ici, on suppose build local classique
subprocess.check_call(['npm', 'install'], cwd=ROOT)
subprocess.check_call(['npm', 'run', 'build'], cwd=ROOT)

print('✅ Build local terminé avec variables localhost !')
