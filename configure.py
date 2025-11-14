import json
import os

def update_config():
    """Configure les domaines et ports pour Office1789"""
    config_path = os.path.join(os.path.dirname(__file__), 'config.json')
    
    # Charger la config existante ou créer une nouvelle
    try:
        with open(config_path, 'r', encoding='utf-8') as f:
            config = json.load(f)
    except FileNotFoundError:
        config = {
            "domains": {},
            "ports": {},
            "autostart": {}
        }
    
    print("""
    ╔═══════════════════════════════════════════════════════════╗
    ║                                                           ║
    ║      🏛  OFFICE1789 - Configuration des domaines          ║
    ║                                                           ║
    ╚═══════════════════════════════════════════════════════════╝
    """)
    
    print("\n📝 Configuration des domaines")
    print("Appuyez sur Entrée pour conserver la valeur par défaut entre []")
    
    # Configuration des domaines
    domains = config.get('domains', {})
    
    default_main = domains.get('main', 'office1789.com')
    main_domain = input(f"\n🏠 Domaine principal [{default_main}]: ").strip() or default_main
    
    default_mail = domains.get('mail', f'mail.{main_domain}')
    mail_domain = input(f"📧 Domaine mail [{default_mail}]: ").strip() or default_mail
    
    default_matrix = domains.get('matrix', f'matrix.{main_domain}')
    matrix_domain = input(f"💬 Domaine Matrix [{default_matrix}]: ").strip() or default_matrix
    
    default_element = domains.get('element', f'chat.{main_domain}')
    element_domain = input(f"💭 Domaine Element/Chat [{default_element}]: ").strip() or default_element
    
    default_onlyoffice = domains.get('onlyoffice', f'docs.{main_domain}')
    onlyoffice_domain = input(f"📄 Domaine OnlyOffice [{default_onlyoffice}]: ").strip() or default_onlyoffice
    
    config['domains'] = {
        'main': main_domain,
        'mail': mail_domain,
        'matrix': matrix_domain,
        'element': element_domain,
        'onlyoffice': onlyoffice_domain
    }
    
    # Configuration des ports
    print("\n\n🔌 Configuration des ports")
    
    ports = config.get('ports', {})
    
    default_frontend = ports.get('frontend', 5173)
    frontend_port = input(f"Frontend Vue.js [{default_frontend}]: ").strip() or default_frontend
    
    default_backend = ports.get('backend', 8080)
    backend_port = input(f"Backend Go [{default_backend}]: ").strip() or default_backend
    
    default_roundcube = ports.get('roundcube', 8081)
    roundcube_port = input(f"Roundcube [{default_roundcube}]: ").strip() or default_roundcube
    
    default_element_port = ports.get('element', 8083)
    element_port = input(f"Element [{default_element_port}]: ").strip() or default_element_port
    
    config['ports'] = {
        'frontend': int(frontend_port),
        'backend': int(backend_port),
        'roundcube': int(roundcube_port),
        'onlyoffice': ports.get('onlyoffice', 8082),
        'element': int(element_port),
        'postgres': 5432,
        'mailserver_smtp': 25,
        'mailserver_imap': 143,
        'mailserver_submission': 587,
        'mailserver_imaps': 993,
        'synapse': 8008,
        'coturn': 3478
    }
    
    # Configuration du démarrage automatique
    print("\n\n⚙️ Configuration du démarrage automatique")
    
    autostart = config.get('autostart', {})
    
    docker_default = 'o' if autostart.get('docker_containers', True) else 'n'
    docker_auto = input(f"Démarrer les conteneurs Docker automatiquement? (o/n) [{docker_default}]: ").strip().lower() or docker_default
    
    backend_default = 'o' if autostart.get('backend', True) else 'n'
    backend_auto = input(f"Démarrer le backend Go automatiquement? (o/n) [{backend_default}]: ").strip().lower() or backend_default
    
    frontend_default = 'o' if autostart.get('frontend', True) else 'n'
    frontend_auto = input(f"Démarrer le frontend Vue.js automatiquement? (o/n) [{frontend_default}]: ").strip().lower() or frontend_default
    
    browser_default = 'o' if autostart.get('open_browser', True) else 'n'
    browser_auto = input(f"Ouvrir le navigateur automatiquement? (o/n) [{browser_default}]: ").strip().lower() or browser_default
    
    config['autostart'] = {
        'docker_containers': docker_auto == 'o',
        'backend': backend_auto == 'o',
        'frontend': frontend_auto == 'o',
        'open_browser': browser_auto == 'o'
    }
    
    # Sauvegarder la configuration
    with open(config_path, 'w', encoding='utf-8') as f:
        json.dump(config, f, indent=2, ensure_ascii=False)
    
    print("\n✅ Configuration sauvegardée dans config.json")
    
    # Afficher un récapitulatif
    print("\n" + "="*60)
    print("📋 RÉCAPITULATIF DE LA CONFIGURATION")
    print("="*60)
    
    print("\n🌐 Domaines:")
    for key, value in config['domains'].items():
        print(f"   {key.ljust(12)}: {value}")
    
    print("\n🔌 Ports:")
    for key, value in config['ports'].items():
        print(f"   {key.ljust(20)}: {value}")
    
    print("\n⚙️ Démarrage automatique:")
    for key, value in config['autostart'].items():
        status = "✅ Oui" if value else "❌ Non"
        print(f"   {key.ljust(20)}: {status}")
    
    print("\n" + "="*60)
    
    print("\n💡 Pour configurer vos domaines locaux, ajoutez ces lignes au fichier hosts:")
    print("   C:\\Windows\\System32\\drivers\\etc\\hosts (Windows)")
    print("   /etc/hosts (Linux/Mac)\n")
    for key, domain in config['domains'].items():
        print(f"127.0.0.1    {domain}")
    
    print("\n🚀 Lancez maintenant: python startOffice1789.py")

if __name__ == "__main__":
    try:
        update_config()
    except KeyboardInterrupt:
        print("\n\n❌ Configuration annulée.")
    except Exception as e:
        print(f"\n❌ Erreur: {e}")
