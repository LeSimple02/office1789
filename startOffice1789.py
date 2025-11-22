import subprocess
import webbrowser
import docker
import json
import os
import time
import sys
import urllib.request
try:
    from tqdm import tqdm
except ImportError:
    tqdm = None

def check_env_files():
    """Vérifie si les fichiers .env existent et propose de les générer"""
    backend_env = os.path.join(os.path.dirname(__file__), 'backend', '.env')
    frontend_env = os.path.join(os.path.dirname(__file__), 'webfront2', '.env')
    docker_env = os.path.join(os.path.dirname(__file__), 'docker', '.env')
    
    missing = []
    if not os.path.exists(backend_env):
        missing.append('backend/.env')
    if not os.path.exists(frontend_env):
        missing.append('webfront2/.env')
    if not os.path.exists(docker_env):
        missing.append('docker/.env')
    
    if missing:
        print("\n⚠️  Fichiers .env manquants:")
        for f in missing:
            print(f"   - {f}")
        print("\n🔑 Génération des secrets et fichiers .env...")
        
        try:
            result = subprocess.run(
                [sys.executable, 'config.py', 'generate-secrets'],
                input='YES I AM SURE\n',
                text=True,
                capture_output=True
            )
            if result.returncode == 0:
                print("✅ Fichiers .env générés avec succès!")
                return True
            else:
                print(f"❌ Erreur lors de la génération: {result.stderr}")
                print("\n💡 Exécutez manuellement: python config.py generate-secrets")
                return False
        except Exception as e:
            print(f"❌ Erreur: {e}")
            print("\n💡 Exécutez manuellement: python config.py generate-secrets")
            return False
    
    return True

def load_config():
    """Charge la configuration depuis config.json"""
    config_path = os.path.join(os.path.dirname(__file__), 'config.json')
    try:
        with open(config_path, 'r', encoding='utf-8') as f:
            return json.load(f)
    except FileNotFoundError:
        print("❌ Fichier config.json introuvable.")
        print("💡 Exécutez: python config.py configure")
        print("\nUtilisation de la configuration par défaut...")
        return {
            "domains": {
                "main": "office1789.com",
                "mail": "mail.office1789.com",
                "matrix": "matrix.office1789.com",
                "element": "chat.office1789.com"
            },
            "ports": {
                "frontend": 5173,
                "backend": 8080,
                "roundcube": 8081,
                "element": 8083,
                "onlyoffice": 8082
            },
            "autostart": {
                "docker_containers": True,
                "backend": True,
                "frontend": True,
                "open_browser": True
            }
        }

def update_hosts_file(domains):
    """Propose d'ajouter les domaines au fichier hosts (Windows)"""
    hosts_path = r"C:\Windows\System32\drivers\etc\hosts"
    
    print("\n📝 Configuration des domaines locaux...")
    print("Pour utiliser des noms de domaine personnalisés, ajoutez ces lignes à votre fichier hosts:")
    print(f"   {hosts_path}\n")
    
    for key, domain in domains.items():
        print(f"127.0.0.1    {domain}")
    
    print("\n💡 Vous pouvez le faire manuellement ou exécuter ce script en tant qu'administrateur.")

def start_docker_containers():
    """Démarre tous les conteneurs Docker (les crée avec docker-compose si nécessaire)"""
    print("\n🐳 Démarrage des conteneurs Docker...")
    
    try:
        client = docker.from_env()
        docker_dir = os.path.join(os.path.dirname(__file__), 'docker')
        
        # Liste des conteneurs à démarrer dans l'ordre (sans backend/frontend qui sont en local)
        containers = [
            "postgres_db",
            "postgres_roundcube", 
            "postgres_synapse",
            "mailserver",
            "roundcube",
            "synapse",
            "element",
            "coturn",
            "onlyoffice"
        ]
        
        # Lancer docker compose up pour créer/démarrer tous les conteneurs
        print(f"   📦 Lancement de docker compose (cela peut prendre quelques minutes)...")
        try:
            process = subprocess.Popen(
                ['docker', 'compose', 'up', '-d', '--build'],
                cwd=docker_dir,
                stdout=subprocess.PIPE,
                stderr=subprocess.STDOUT,
                text=True
            )
            
            # Afficher une barre de progression basée sur le temps estimé
            if tqdm:
                total_steps = len(containers) * 10
                with tqdm(total=total_steps, desc="   Conteneurs", unit="step") as pbar:
                    last_update = time.time()
                    for line in process.stdout:
                        # Mettre à jour la barre toutes les 0.5 secondes
                        current = time.time()
                        if current - last_update > 0.5 and pbar.n < total_steps - 5:
                            pbar.update(1)
                            last_update = current
                        # Afficher les lignes importantes
                        line = line.strip()
                        if any(keyword in line for keyword in ['Creating', 'Starting', 'Created', 'Started', 'Recreating', 'Building']):
                            pbar.set_postfix_str(line[-50:] if len(line) > 50 else line)
                    # Compléter la barre
                    pbar.n = total_steps
                    pbar.refresh()
            else:
                # Fallback sans tqdm - afficher les lignes importantes
                for line in process.stdout:
                    line = line.strip()
                    if any(keyword in line for keyword in ['Pulling', 'Creating', 'Starting', 'Created', 'Started', 'Recreating', 'Building']):
                        print(f"      {line}")
            
            process.wait()
            
            if process.returncode == 0:
                print("   ✅ Conteneurs créés avec succès!")
            else:
                print(f"   ⚠️  Quelques warnings (c'est normal lors de la première création)")
        except FileNotFoundError:
            print("   ❌ 'docker compose' introuvable. Installez Docker Desktop.")
            return False
        except Exception as e:
            print(f"   ⚠️  Erreur: {e}")
        
        # Attendre un peu que les conteneurs démarrent
        print("\n   ⏳ Attente du démarrage complet (5 secondes)...")
        time.sleep(5)
        
        # Vérifier que les conteneurs sont bien démarrés
        print("   🔍 Vérification des conteneurs...")
        running_count = 0
        for container_name in containers:
            try:
                container = client.containers.get(container_name)
                status = container.status
                if status == "running":
                    running_count += 1
                    print(f"   ✅ {container_name} est actif")
                elif status in ["created", "restarting"]:
                    print(f"   ⏳ {container_name} démarre...")
                    container.start()
                    running_count += 1
                else:
                    print(f"   ⚠️  {container_name} : {status}")
            except docker.errors.NotFound:
                print(f"   ⚠️  {container_name} introuvable")
            except Exception as e:
                print(f"   ⚠️  {container_name} : {e}")
        
        print(f"\n   ✅ {running_count}/{len(containers)} conteneurs actifs!")
        
    except docker.errors.DockerException as e:
        print(f"❌ Erreur Docker: {e}")
        print("⚠️  Assurez-vous que Docker Desktop est en cours d'exécution.")
        return False
    
    return True

def kill_process_on_port(port):
    """Tue le processus utilisant un port spécifique (Windows)"""
    try:
        if sys.platform == 'win32':
            # Trouver le PID utilisant le port
            result = subprocess.run(
                f'netstat -ano | findstr :{port}',
                shell=True,
                capture_output=True,
                text=True
            )
            if result.stdout:
                lines = result.stdout.strip().split('\n')
                for line in lines:
                    if 'LISTENING' in line:
                        parts = line.split()
                        pid = parts[-1]
                        print(f"   🔧 Arrêt du processus {pid} utilisant le port {port}...")
                        subprocess.run(f'taskkill /F /PID {pid}', shell=True, capture_output=True)
                        time.sleep(1)
                        return True
        return False
    except Exception as e:
        print(f"   ⚠️ Erreur lors de l'arrêt du processus: {e}")
        return False

def start_backend(config):
    """Démarre le backend Go avec logs visibles"""
    print("\n🚀 Démarrage du backend Go...")
    backend_dir = os.path.join(os.path.dirname(__file__), 'backend')
    
    # Vérifier si le port 8080 est déjà utilisé
    port = config['ports']['backend']
    print(f"   🔍 Vérification du port {port}...")
    kill_process_on_port(port)
    
    try:
        # Windows peut nécessiter shell=True pour air
        if sys.platform == 'win32':
            process = subprocess.Popen('air', cwd=backend_dir, shell=True)
        else:
            process = subprocess.Popen(['air'], cwd=backend_dir)
        print(f"✅ Backend démarré sur http://localhost:{port}")
        return process
    except FileNotFoundError:
        print("⚠️  'air' introuvable, tentative avec 'go run'...")
        try:
            if sys.platform == 'win32':
                process = subprocess.Popen('go run .', cwd=backend_dir, shell=True)
            else:
                process = subprocess.Popen(['go', 'run', '.'], cwd=backend_dir)
            print(f"✅ Backend démarré sur http://localhost:{port}")
            return process
        except Exception as e:
            print(f"❌ Erreur lors du démarrage du backend: {e}")
            return None
    except Exception as e:
        print(f"❌ Erreur lors du démarrage du backend: {e}")
        return None

def start_frontend(config):
    """Démarre le frontend Vue.js"""
    print("\n🎨 Démarrage du frontend Vue.js...")
    frontend_dir = os.path.join(os.path.dirname(__file__), 'webfront2')
    
    try:
        # Windows nécessite shell=True pour npm
        if sys.platform == 'win32':
            process = subprocess.Popen('npm run dev', cwd=frontend_dir, shell=True)
        else:
            process = subprocess.Popen(['npm', 'run', 'dev'], cwd=frontend_dir)
        print(f"✅ Frontend démarré sur http://localhost:{config['ports']['frontend']}\n")
        return process
    except Exception as e:
        print(f"❌ Erreur lors du démarrage du frontend: {e}")
        return None

def display_urls(config):
    """Affiche les URLs d'accès"""
    print("\n" + "="*60)
    print("🏛  OFFICE1789 - URLs d'accès")
    print("="*60)
    
    domains = config['domains']
    ports = config['ports']
    
    print(f"\n📱 Application principale:")
    print(f"   http://localhost:{ports['frontend']}")
    print(f"   http://{domains['main']}:{ports['frontend']} (si configuré dans hosts)")
    
    print(f"\n📧 Webmail (Roundcube):")
    print(f"   http://localhost:{ports['roundcube']}")
    print(f"   http://{domains['mail']}:{ports['roundcube']} (si configuré dans hosts)")
    
    print(f"\n💬 Chat (Element/Matrix):")
    print(f"   http://localhost:{ports['element']}")
    print(f"   http://{domains['element']}:{ports['element']} (si configuré dans hosts)")
    
    print(f"\n📄 OnlyOffice:")
    print(f"   http://localhost:{ports['onlyoffice']}")
    
    print(f"\n🔧 Backend API:")
    print(f"   http://localhost:{ports['backend']}")
    
    print("\n" + "="*60 + "\n")

if __name__ == "__main__":
    print("""
    ╔═══════════════════════════════════════════════════════════╗
    ║                                                           ║
    ║           🏛  OFFICE1789 - Démarrage automatique          ║
    ║                                                           ║
    ╚═══════════════════════════════════════════════════════════╝
    """)
    
    # Charger la configuration
    config = load_config()
    
    # Afficher les informations de configuration
    print(f"📋 Configuration chargée:")
    print(f"   Domaine principal: {config['domains']['main']}")
    print(f"   Port frontend: {config['ports']['frontend']}")
    print(f"   Port backend: {config['ports']['backend']}")
    
    # Proposer la configuration du fichier hosts
    update_hosts_file(config['domains'])
    
    # Démarrer les conteneurs Docker si configuré
    if config['autostart']['docker_containers']:
        if not start_docker_containers():
            print("\n⚠️  Impossible de démarrer les conteneurs Docker.")
            print("💡 Vérifiez que Docker Desktop est lancé.")
            sys.exit(1)
        else:
            # Attendre que les conteneurs soient prêts
            print("\n⏳ Attente du démarrage des services Docker (10 secondes)...")
            time.sleep(10)
    
    # Démarrer le backend si configuré (en local, pas dans Docker)
    backend_process = None
    if config['autostart']['backend']:
        backend_process = start_backend(config)
        time.sleep(2)  # Attendre que le backend démarre
    
    # Démarrer le frontend si configuré (en local, pas dans Docker)
    frontend_process = None
    if config['autostart']['frontend']:
        frontend_process = start_frontend(config)
        time.sleep(3)  # Attendre que le frontend démarre
    
    # Afficher les URLs
    display_urls(config)
    
    # Ouvrir le navigateur si configuré
    if config['autostart']['open_browser']:
        print("🌐 Ouverture du navigateur...")
        time.sleep(2)
        webbrowser.open(f"http://localhost:{config['ports']['frontend']}")
    
    print("\n" + "="*60)
    print("✅ Office1789 est prêt à l'emploi!")
    print("   Backend et frontend en local, services Docker actifs.")
    print("   Appuyez sur Ctrl+C pour arrêter tous les services.")
    print("="*60 + "\n")
    
    # Attendre que les processus se terminent (ou Ctrl+C)
    try:
        while True:
            time.sleep(1)
            # Vérifier si les processus sont toujours vivants
            if backend_process and backend_process.poll() is not None:
                print("\n❌ Le backend s'est arrêté de manière inattendue!")
                break
            if frontend_process and frontend_process.poll() is not None:
                print("\n❌ Le frontend s'est arrêté de manière inattendue!")
                break
    except KeyboardInterrupt:
        print("\n\n🛑 Arrêt d'Office1789...")
        if frontend_process:
            frontend_process.terminate()
        if backend_process:
            backend_process.terminate()
        print("💡 Les conteneurs Docker continuent de tourner.")
        print("   Pour les arrêter: cd docker && docker compose down")
        print("👋 Au revoir!")
