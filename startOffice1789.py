import subprocess
import webbrowser
import docker
import json
import os
import time
import sys

def load_config():
    """Charge la configuration depuis config.json"""
    config_path = os.path.join(os.path.dirname(__file__), 'config.json')
    try:
        with open(config_path, 'r', encoding='utf-8') as f:
            return json.load(f)
    except FileNotFoundError:
        print("❌ Fichier config.json introuvable. Utilisation de la configuration par défaut.")
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
                "roundcube": 8081
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
    """Démarre tous les conteneurs Docker"""
    print("\n🐳 Démarrage des conteneurs Docker...")
    
    try:
        client = docker.from_env()
        
        # Liste des conteneurs à démarrer dans l'ordre
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
        
        for container_name in containers:
            try:
                container = client.containers.get(container_name)
                if container.status != "running":
                    print(f"   ▶️  Démarrage de {container_name}...")
                    container.start()
                    time.sleep(1)
                else:
                    print(f"   ✅ {container_name} est déjà en cours d'exécution")
            except docker.errors.NotFound:
                print(f"   ⚠️  Conteneur {container_name} introuvable (peut-être pas encore créé)")
            except Exception as e:
                print(f"   ❌ Erreur lors du démarrage de {container_name}: {e}")
        
        print("\n✅ Conteneurs Docker démarrés avec succès!")
        
    except docker.errors.DockerException as e:
        print(f"❌ Erreur Docker: {e}")
        print("⚠️  Assurez-vous que Docker Desktop est en cours d'exécution.")
        return False
    
    return True

def start_backend(config):
    """Démarre le backend Go"""
    print("\n🚀 Démarrage du backend Go...")
    backend_dir = os.path.join(os.path.dirname(__file__), 'backend')
    
    try:
        # Lancer air directement (comme dans l'ancien script)
        if sys.platform == 'win32':
            # Windows : lancer dans une nouvelle fenêtre de console
            process = subprocess.Popen('air', cwd=backend_dir, shell=True,
                                     creationflags=subprocess.CREATE_NEW_CONSOLE)
        else:
            # Linux/Mac
            process = subprocess.Popen(['air'], cwd=backend_dir)
        
        print(f"✅ Backend démarré sur http://localhost:{config['ports']['backend']}")
        print("   (Une nouvelle fenêtre de terminal s'est ouverte pour le backend)")
        return process
    except Exception as e:
        print(f"⚠️  Erreur avec 'air': {e}")
        print("⚠️  Tentative avec 'go run'...")
        try:
            if sys.platform == 'win32':
                process = subprocess.Popen('go run .', cwd=backend_dir, shell=True,
                                         creationflags=subprocess.CREATE_NEW_CONSOLE)
            else:
                process = subprocess.Popen(['go', 'run', '.'], cwd=backend_dir)
            print(f"✅ Backend démarré sur http://localhost:{config['ports']['backend']}")
            return process
        except Exception as e:
            print(f"❌ Erreur lors du démarrage du backend: {e}")
            return None

def start_frontend(config):
    """Démarre le frontend Vue.js"""
    print("\n🎨 Démarrage du frontend Vue.js...")
    frontend_dir = os.path.join(os.path.dirname(__file__), 'webfront2')
    
    try:
        process = subprocess.Popen(['npm', 'run', 'dev'], cwd=frontend_dir, shell=True)
        print(f"✅ Frontend démarré sur http://localhost:{config['ports']['frontend']}")
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
            print("\n⚠️  Erreur lors du démarrage des conteneurs Docker.")
            print("Voulez-vous continuer sans Docker? (o/n)")
            response = input().lower()
            if response != 'o':
                sys.exit(1)
        
        # Attendre que les conteneurs soient prêts
        print("\n⏳ Attente du démarrage des services (10 secondes)...")
        time.sleep(10)
    
    # Démarrer le backend si configuré
    backend_process = None
    if config['autostart']['backend']:
        backend_process = start_backend(config)
        time.sleep(3)  # Attendre que le backend démarre
    
    # Démarrer le frontend si configuré
    frontend_process = None
    if config['autostart']['frontend']:
        frontend_process = start_frontend(config)
        time.sleep(5)  # Attendre que le frontend démarre
    
    # Afficher les URLs
    display_urls(config)
    
    # Ouvrir le navigateur si configuré
    if config['autostart']['open_browser']:
        print("🌐 Ouverture du navigateur...")
        time.sleep(2)
        webbrowser.open(f"http://localhost:{config['ports']['frontend']}")
    
    print("\n✅ Office1789 est prêt à l'emploi!")
    print("   Appuyez sur Ctrl+C pour arrêter tous les services.\n")
    
    # Attendre que le frontend se termine (ou Ctrl+C)
    try:
        if frontend_process:
            frontend_process.wait()
    except KeyboardInterrupt:
        print("\n\n🛑 Arrêt d'Office1789...")
        if frontend_process:
            frontend_process.terminate()
        if backend_process:
            backend_process.terminate()
        print("👋 Au revoir!")
