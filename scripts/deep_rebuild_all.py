import os
import subprocess
import sys

# Chemin du dossier docker où se trouve le docker-compose.yml
DOCKER_DIR = os.path.join(os.path.dirname(os.path.abspath(__file__)), '..', 'docker')

def deep_rebuild_all():
    print('🔄 Rebuild complet de toutes les images Docker du projet...')
    cmd_build = ['docker-compose', 'build', '--no-cache']
    cmd_up = ['docker-compose', 'up', '-d']

    try:
        subprocess.run(cmd_build, cwd=DOCKER_DIR, check=True)
        print('✅ Build terminé.')
        subprocess.run(cmd_up, cwd=DOCKER_DIR, check=True)
        print('🚀 Tous les conteneurs sont relancés !')
    except subprocess.CalledProcessError as e:
        print(f'❌ Erreur lors du rebuild ou du redémarrage : {e}')
        sys.exit(1)

if __name__ == '__main__':
    deep_rebuild_all()
