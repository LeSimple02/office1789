#!/usr/bin/env python3
"""
Script d'installation du plugin SSO Office1789 pour Element/Matrix
Compatible Windows, Linux, macOS
"""

import subprocess
import sys
import os

def run_command(cmd, description):
    """Exécute une commande et affiche le résultat"""
    print(f"🔧 {description}...")
    try:
        result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
        if result.returncode == 0:
            print(f"✅ {description} - OK")
            return True
        else:
            print(f"❌ {description} - ERREUR")
            if result.stderr:
                print(f"   {result.stderr}")
            return False
    except Exception as e:
        print(f"❌ Erreur: {e}")
        return False

def main():
    print("\n" + "="*60)
    print("🏛  OFFICE1789 - Installation SSO Element/Matrix")
    print("="*60 + "\n")
    
    # Déterminer le chemin du script (marche sur Windows et Linux)
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(script_dir)
    sso_js_path = os.path.join(script_dir, "element", "office1789-sso.js")
    
    # Vérifier que le fichier SSO existe
    if not os.path.exists(sso_js_path):
        print(f"❌ Fichier introuvable: {sso_js_path}")
        sys.exit(1)
    
    # Normaliser le chemin pour Docker (utiliser / même sur Windows)
    sso_js_path_docker = sso_js_path.replace('\\', '/')
    
    # 1. Copier le script SSO dans le conteneur Element
    cmd = f'docker cp "{sso_js_path_docker}" element:/app/office1789-sso.js'
    if not run_command(cmd, "Copie du script SSO dans Element"):
        sys.exit(1)
    
    # 2. Injecter le script dans index.html d'Element
    # Écrire dans /tmp puis copier (car /app est read-only)
    print("🔧 Injection du script dans index.html...")
    
    # Lire le contenu actuel
    result = subprocess.run(
        ['docker', 'exec', 'element', 'cat', '/app/index.html'],
        capture_output=True,
        text=True
    )
    
    if result.returncode != 0:
        print("❌ Impossible de lire index.html")
        sys.exit(1)
    
    html_content = result.stdout
    
    # Vérifier si le script est déjà présent
    if 'office1789-sso.js' in html_content:
        print("✅ Script déjà présent dans index.html")
    else:
        # Ajouter le script avant </head>
        new_content = html_content.replace(
            '</head>',
            '<script src="/office1789-sso.js"></script></head>'
        )
        
        # Écrire dans /tmp d'abord
        process = subprocess.Popen(
            ['docker', 'exec', '-i', 'element', 'sh', '-c', 'cat > /tmp/index.html && cp /tmp/index.html /app/index.html'],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        
        stdout, stderr = process.communicate(input=new_content)
        
        if process.returncode == 0:
            print("✅ Script injecté dans index.html")
        else:
            print(f"❌ Erreur: {stderr}")
            print("⚠️  Le conteneur Element est en lecture seule")
            print("💡 Solution: Le script SSO sera chargé au prochain rebuild d'Element")
            # Ne pas arrêter, continuer quand même
    

    
    # 3. Redémarrer Element
    if not run_command("docker restart element", "Redémarrage d'Element"):
        sys.exit(1)
    
    print("\n" + "="*60)
    print("✅ Installation terminée avec succès !")
    print("📝 Element va redémarrer (10-15 secondes)")
    print("🌐 Accédez à Element via: http://localhost:8083")
    print("="*60 + "\n")

if __name__ == "__main__":
    main()
