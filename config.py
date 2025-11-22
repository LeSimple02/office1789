"""Configuration complète Office1789 - Setup interactif et runtime URLs.

Ce module unifie la gestion de configuration:
1. Configuration initiale interactive (domaines, ports, autostart)
2. Génération automatique de mots de passe sécurisés
3. Construction des URLs runtime depuis config.json + variables d'environnement
4. Création des fichiers .env pour backend, frontend et Docker
5. Export de la configuration pour le frontend

Variables d'environnement supportées (prioritaires sur config.json):
  OFFICE1789_ENV                -> local | staging | production (default: local)
  OFFICE1789_BACKEND_URL        -> URL complète du backend (ex: https://api.office1789.com)
  OFFICE1789_FRONTEND_URL       -> URL web front (ex: https://app.office1789.com)
  OFFICE1789_ROUNDCUBE_URL      -> Webmail Roundcube
  OFFICE1789_ONLYOFFICE_URL     -> Document Server OnlyOffice
  OFFICE1789_MATRIX_URL         -> Synapse base URL
  OFFICE1789_ELEMENT_URL        -> Element Web
  OFFICE1789_MAIL_URL           -> Mail domain (si exposé séparément)
  OFFICE1789_STRIPE_PUBLIC_KEY  -> Clé publique Stripe

Usage:
  python config.py                 # configuration interactive (crée/modifie config.json + .env)
  python config.py print           # affiche la configuration courante
  python config.py export-frontend # génère webfront2/src/config/runtime-config.json
  python config.py generate-secrets # regénère tous les secrets (USE WITH CAUTION!)
"""
from __future__ import annotations
import os
import json
import secrets
import string
from dataclasses import dataclass, asdict
from pathlib import Path
from typing import Optional

ROOT = Path(__file__).parent
CONFIG_JSON = ROOT / "config.json"
FRONTEND_RUNTIME_JSON = ROOT / "webfront2" / "src" / "config" / "runtime-config.json"
BACKEND_ENV = ROOT / "backend" / ".env"
FRONTEND_ENV = ROOT / "webfront2" / ".env"
DOCKER_ENV = ROOT / "docker" / ".env"
COTURN_CONF = ROOT / "docker" / "coturn" / "turnserver.conf"

# Valeurs par défaut (mode local)
DEFAULT_PORTS = {
    "backend": 8080,
    "frontend": 5173,
    "roundcube": 8081,
    "onlyoffice": 8082,
    "element": 8083,
}
DEFAULT_DOMAINS = {
    "main": "office1789.com",
    "mail": "mail.office1789.com",
    "matrix": "matrix.office1789.com",
    "element": "chat.office1789.com",
    "onlyoffice": "docs.office1789.com",
}


def generate_password(length: int = 32, alphanumeric_only: bool = False) -> str:
    """Génère un mot de passe sécurisé aléatoire."""
    if alphanumeric_only:
        chars = string.ascii_letters + string.digits
    else:
        chars = string.ascii_letters + string.digits + "!@#$%^&*()-_=+[]{}|;:,.<>?"
    return ''.join(secrets.choice(chars) for _ in range(length))


def generate_aes_key() -> str:
    """Génère une clé AES-256 (32 bytes en hex = 64 chars)."""
    return secrets.token_hex(32)

@dataclass
class ServiceURLs:
    backend: str
    frontend: str
    roundcube: str
    onlyoffice: str
    matrix: str
    element: str
    mail: str

@dataclass
class AppConfig:
    environment: str
    services: ServiceURLs
    stripe_public_key: Optional[str] = None

    def to_dict(self) -> dict:
        d = asdict(self)
        return d

    def write_frontend_runtime(self, path: Path = FRONTEND_RUNTIME_JSON) -> None:
        """Écrit un fichier JSON minimal consommable par le frontend.
        Exemple: { "API_BASE": "https://api.office1789.com" }
        """
        path.parent.mkdir(parents=True, exist_ok=True)
        data = {
            "API_BASE": self.services.backend.rstrip('/'),
            "FRONTEND_BASE": self.services.frontend.rstrip('/'),
            "MATRIX_BASE": self.services.matrix.rstrip('/'),
            "ELEMENT_BASE": self.services.element.rstrip('/'),
            "ONLYOFFICE_BASE": self.services.onlyoffice.rstrip('/'),
            "ROUNDCUBE_BASE": self.services.roundcube.rstrip('/'),
            "MAIL_BASE": self.services.mail.rstrip('/'),
            "ENV": self.environment,
        }
        with path.open('w', encoding='utf-8') as f:
            json.dump(data, f, indent=2, ensure_ascii=False)


def _load_raw_config() -> dict:
    """Charge config.json ou retourne les valeurs par défaut."""
    if CONFIG_JSON.exists():
        try:
            return json.loads(CONFIG_JSON.read_text(encoding='utf-8'))
        except Exception:
            pass
    return {
        "domains": DEFAULT_DOMAINS.copy(),
        "ports": DEFAULT_PORTS.copy(),
        "autostart": {
            "docker_containers": True,
            "backend": True,
            "frontend": True,
            "open_browser": True
        },
        "secrets": {}  # Sera rempli lors du premier configure
    }


def _save_config(config_data: dict) -> None:
    """Sauvegarde la configuration dans config.json."""
    with CONFIG_JSON.open('w', encoding='utf-8') as f:
        json.dump(config_data, f, indent=2, ensure_ascii=False)


def generate_all_secrets() -> dict:
    """Génère tous les secrets nécessaires pour l'application."""
    print("\n🔐 Génération des secrets sécurisés...")
    secrets_dict = {
        # PostgreSQL
        "postgres_user": "office1789_user",
        "postgres_password": generate_password(32, alphanumeric_only=True),
        "postgres_db": "office1789",
        
        # Roundcube PostgreSQL
        "roundcube_db_user": "roundcube_user",
        "roundcube_db_password": generate_password(32, alphanumeric_only=True),
        "roundcube_db_name": "roundcube",
        
        # Synapse PostgreSQL
        "synapse_db_user": "synapse_user",
        "synapse_db_password": generate_password(32, alphanumeric_only=True),
        "synapse_db_name": "synapse",
        
        # Coturn
        "coturn_secret": generate_password(48, alphanumeric_only=True),
        
        # OnlyOffice
        "onlyoffice_jwt_secret": generate_password(48),
        
        # AES Encryption
        "aes_key": generate_aes_key(),
        
        # Matrix Admin Token
        "matrix_admin_token": generate_password(64, alphanumeric_only=True),
        
        # Mail Admin
        "mail_admin_user": "admin",
        "mail_admin_password": generate_password(32),
    }
    print("✅ Secrets générés avec succès!")
    return secrets_dict


def create_backend_env(config: dict, secrets: dict, environment: str = "local") -> None:
    """Crée le fichier .env pour le backend Go."""
    print("\n📝 Création de backend/.env...")
    
    domains = config.get('domains', DEFAULT_DOMAINS)
    ports = config.get('ports', DEFAULT_PORTS)
    
    # Construction des URLs selon l'environnement
    if environment == 'local':
        frontend_url = f"http://localhost:{ports.get('frontend', 5173)}"
        matrix_url = f"http://localhost:8008"
        roundcube_url = f"http://localhost:{ports.get('roundcube', 8081)}"
        element_url = f"http://localhost:{ports.get('element', 8083)}"
        onlyoffice_url = f"http://localhost:{ports.get('onlyoffice', 8082)}"
    else:
        frontend_url = f"https://{domains.get('main', 'office1789.com')}"
        matrix_url = f"https://{domains.get('matrix', 'matrix.office1789.com')}"
        roundcube_url = f"https://{domains.get('mail', 'mail.office1789.com')}"
        element_url = f"https://{domains.get('element', 'chat.office1789.com')}"
        onlyoffice_url = f"https://{domains.get('onlyoffice', 'docs.office1789.com')}"
    
    env_content = f"""# Database Configuration (Auto-generated by config.py)
DB_HOST=localhost
DB_PORT=5432
DB_USER={secrets['postgres_user']}
DB_PASSWORD={secrets['postgres_password']}
DB_NAME={secrets['postgres_db']}

# Frontend URL
FRONTEND_URL={frontend_url}

# Services URLs
ROUNDCUBE_URL={roundcube_url}
ELEMENT_URL={element_url}
ONLYOFFICE_URL={onlyoffice_url}

# OVH SMS Configuration (pour l'envoi de SMS de vérification)
# Obtenir ces valeurs sur https://api.ovh.com/createToken/
# Si non configuré, les codes SMS seront loggés dans la console (mode dev)
OVH_SMS_APPLICATION_KEY=your_application_key
OVH_SMS_APPLICATION_SECRET=your_application_secret
OVH_SMS_CONSUMER_KEY=your_consumer_key
OVH_SMS_SERVICE_NAME=sms-xxXXXX-1
OVH_SMS_SENDER=Office1789

# Stripe Configuration
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret

# Stripe Price IDs (créés dans le dashboard Stripe)
STRIPE_PRICE_STANDARD=price_standard_monthly
STRIPE_PRICE_PROFESSIONAL=price_professional_monthly
STRIPE_PRICE_ENTERPRISE=price_enterprise_monthly

# AES Encryption Key (32 bytes for AES-256)
AES_KEY={secrets['aes_key']}

# Matrix Configuration
MATRIX_HOMESERVER={matrix_url}
MATRIX_ADMIN_TOKEN={secrets['matrix_admin_token']}

# Mail Configuration
MAIL_SERVER={domains.get('mail', 'mail.office1789.com')}
MAIL_ADMIN_USER={secrets['mail_admin_user']}
MAIL_ADMIN_PASSWORD={secrets['mail_admin_password']}

# OnlyOffice JWT Secret (for document server token signing)
ONLYOFFICE_JWT_SECRET={secrets.get('onlyoffice_jwt_secret', '')}
"""
    
    BACKEND_ENV.parent.mkdir(parents=True, exist_ok=True)
    BACKEND_ENV.write_text(env_content, encoding='utf-8')
    print(f"✅ Fichier créé: {BACKEND_ENV}")


def create_frontend_env(config: dict, environment: str = "local") -> None:
    """Crée le fichier .env pour le frontend Vue."""
    print("\n📝 Création de webfront2/.env...")
    
    domains = config.get('domains', DEFAULT_DOMAINS)
    ports = config.get('ports', DEFAULT_PORTS)
    
    # Construction des URLs selon l'environnement
    if environment == 'local':
        backend_url = f"http://localhost:{ports.get('backend', 8080)}"
        frontend_url = f"http://localhost:{ports.get('frontend', 5173)}"
        roundcube_url = f"http://localhost:{ports.get('roundcube', 8081)}"
        element_url = f"http://localhost:{ports.get('element', 8083)}"
        onlyoffice_url = f"http://localhost:{ports.get('onlyoffice', 8082)}"
        matrix_url = f"http://localhost:8008"
    else:
        backend_url = f"https://backend.{domains.get('main', 'office1789.com')}"
        frontend_url = f"https://{domains.get('main', 'office1789.com')}"
        roundcube_url = f"https://{domains.get('mail', 'mail.office1789.com')}"
        element_url = f"https://{domains.get('element', 'chat.office1789.com')}"
        onlyoffice_url = f"https://{domains.get('onlyoffice', 'docs.office1789.com')}"
        matrix_url = f"https://{domains.get('matrix', 'matrix.office1789.com')}"
    
    env_content = f"""# Auto-generated by config.py
# API Backend Base
VITE_APP_API={backend_url}

# Auth Routes
VITE_APP_API_LOGIN={backend_url}/api/connect
VITE_APP_API_CREATE_ACCOUNT={backend_url}/api/subscribe
VITE_APP_API_GETINFOP={backend_url}/api/getinfop
VITE_APP_API_CHANGEINFO={backend_url}/api/changeinfo

# Account Routes
VITE_APP_API_DELETE_ACCOUNT={backend_url}/api/account/delete
VITE_APP_API_CHANGE_PASSWORD={backend_url}/api/account/change-password

# Password Reset Routes
VITE_APP_API_PASSWORD_RESET_REQUEST={backend_url}/api/password/reset/request
VITE_APP_API_PASSWORD_RESET_CONFIRM={backend_url}/api/password/reset/confirm

# Verification Routes
VITE_APP_API_VERIFICATION_SEND={backend_url}/api/verification/send
VITE_APP_API_VERIFICATION_VERIFY={backend_url}/api/verification/verify

# Drive Routes (optionnel - utilisera VITE_APP_API par défaut)
# VITE_API_DRIVE={backend_url}/api/drive/getfiles
# VITE_API_DRIVE_UPLOAD={backend_url}/api/drive/upload
# VITE_API_DRIVE_DOWNLOAD={backend_url}/api/drive/download

# Services URLs
VITE_ROUNDCUBE_URL={roundcube_url}
VITE_ELEMENT_URL={element_url}
VITE_ONLYOFFICE_URL={onlyoffice_url}
VITE_MATRIX_URL={matrix_url}

# Frontend Base URL
VITE_FRONTEND_URL={frontend_url}

# Environment
VITE_ENV={environment}
"""
    
    FRONTEND_ENV.parent.mkdir(parents=True, exist_ok=True)
    FRONTEND_ENV.write_text(env_content, encoding='utf-8')
    print(f"✅ Fichier créé: {FRONTEND_ENV}")


def create_docker_env(secrets: dict) -> None:
    """Crée le fichier .env pour Docker Compose."""
    print("\n📝 Création de docker/.env...")
    
    env_content = f"""# Auto-generated by config.py - DO NOT COMMIT THIS FILE!

# PostgreSQL Main Database
POSTGRES_USER={secrets['postgres_user']}
POSTGRES_PASSWORD={secrets['postgres_password']}
POSTGRES_DB={secrets['postgres_db']}

# Roundcube Database
ROUNDCUBE_DB_USER={secrets['roundcube_db_user']}
ROUNDCUBE_DB_PASSWORD={secrets['roundcube_db_password']}
ROUNDCUBE_DB_NAME={secrets['roundcube_db_name']}

# Synapse Database
SYNAPSE_DB_USER={secrets['synapse_db_user']}
SYNAPSE_DB_PASSWORD={secrets['synapse_db_password']}
SYNAPSE_DB_NAME={secrets['synapse_db_name']}

# Coturn (TURN server)
COTURN_SECRET={secrets['coturn_secret']}

# OnlyOffice
ONLYOFFICE_JWT_SECRET={secrets['onlyoffice_jwt_secret']}

# Matrix
MATRIX_ADMIN_TOKEN={secrets['matrix_admin_token']}

# Mail Server
MAIL_ADMIN_USER={secrets['mail_admin_user']}
MAIL_ADMIN_PASSWORD={secrets['mail_admin_password']}
"""
    
    DOCKER_ENV.parent.mkdir(parents=True, exist_ok=True)
    DOCKER_ENV.write_text(env_content, encoding='utf-8')
    print(f"✅ Fichier créé: {DOCKER_ENV}")


def update_coturn_config(secrets: dict) -> None:
    """Met à jour le fichier de configuration Coturn avec le secret généré."""
    print("\n📝 Mise à jour de coturn/turnserver.conf...")
    
    coturn_content = f"""# Auto-updated by config.py
use-auth-secret
static-auth-secret={secrets['coturn_secret']}
realm=office1789.com
fingerprint
no-tls
no-dtls
listening-port=3478
# Narrow relay port range for Docker Desktop
min-port=49160
max-port=49200
user-quota=12
total-quota=1200
log-file=stdout
simple-log
"""
    
    COTURN_CONF.parent.mkdir(parents=True, exist_ok=True)
    COTURN_CONF.write_text(coturn_content, encoding='utf-8')
    print(f"✅ Fichier mis à jour: {COTURN_CONF}")


def interactive_configure() -> None:
    """Configuration interactive des domaines, ports et autostart (remplace configure.py)."""
    config = _load_raw_config()
    
    print("""
    ╔═══════════════════════════════════════════════════════════╗
    ║                                                           ║
    ║      🏛  OFFICE1789 - Configuration des domaines          ║
    ║                                                           ║
    ╚═══════════════════════════════════════════════════════════╝
    """)
    
    # Choix de l'environnement
    print("\n🌍 Environnement de déploiement")
    print("1. local (développement)")
    print("2. staging (pré-production)")
    print("3. production")
    env_choice = input("\nChoisissez l'environnement [1]: ").strip() or "1"
    env_map = {"1": "local", "2": "staging", "3": "production"}
    environment = env_map.get(env_choice, "local")
    config['environment'] = environment
    
    print(f"\n✅ Environnement sélectionné: {environment}")
    
    print("\n📝 Configuration des domaines")
    print("Appuyez sur Entrée pour conserver la valeur par défaut entre []")
    
    # Configuration des domaines
    domains = config.get('domains', {})
    
    default_main = domains.get('main', DEFAULT_DOMAINS['main'])
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
    
    default_frontend = ports.get('frontend', DEFAULT_PORTS['frontend'])
    frontend_port = input(f"Frontend Vue.js [{default_frontend}]: ").strip() or default_frontend
    
    default_backend = ports.get('backend', DEFAULT_PORTS['backend'])
    backend_port = input(f"Backend Go [{default_backend}]: ").strip() or default_backend
    
    default_roundcube = ports.get('roundcube', DEFAULT_PORTS['roundcube'])
    roundcube_port = input(f"Roundcube [{default_roundcube}]: ").strip() or default_roundcube
    
    default_element_port = ports.get('element', DEFAULT_PORTS['element'])
    element_port = input(f"Element [{default_element_port}]: ").strip() or default_element_port
    
    config['ports'] = {
        'frontend': int(frontend_port),
        'backend': int(backend_port),
        'roundcube': int(roundcube_port),
        'onlyoffice': DEFAULT_PORTS['onlyoffice'],
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
    
    # Génération ou réutilisation des secrets
    existing_secrets = config.get('secrets', {})
    if existing_secrets:
        print("\n\n🔐 Secrets existants détectés")
        regen = input("Voulez-vous REGÉNÉRER tous les secrets? (o/n) [n]: ").strip().lower()
        if regen == 'o':
            print("\n⚠️  ATTENTION: La régénération des secrets nécessitera la recréation des bases de données!")
            confirm = input("Êtes-vous SÛR? (tapez 'OUI' pour confirmer): ").strip()
            if confirm == 'OUI':
                config['secrets'] = generate_all_secrets()
            else:
                print("❌ Régénération annulée, secrets existants conservés")
        else:
            print("✅ Secrets existants conservés")
    else:
        config['secrets'] = generate_all_secrets()
    
    # Sauvegarder la configuration
    _save_config(config)
    print("\n✅ Configuration sauvegardée dans config.json")
    
    # Créer les fichiers .env
    create_backend_env(config, config['secrets'], environment)
    create_frontend_env(config, environment)
    create_docker_env(config['secrets'])
    update_coturn_config(config['secrets'])
    
    # Afficher un récapitulatif
    print("\n" + "="*60)
    print("📋 RÉCAPITULATIF DE LA CONFIGURATION")
    print("="*60)
    
    print(f"\n🌍 Environnement: {environment}")
    
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
    
    print("\n🔐 Sécurité:")
    print(f"   Secrets générés      : ✅ {len(config['secrets'])} éléments")
    print(f"   backend/.env         : ✅ Créé")
    print(f"   webfront2/.env       : ✅ Créé")
    print(f"   docker/.env          : ✅ Créé")
    print(f"   coturn config        : ✅ Mis à jour")
    
    print("\n" + "="*60)
    
    if environment == 'local':
        print("\n💡 Pour configurer vos domaines locaux, ajoutez ces lignes au fichier hosts:")
        print("   C:\\Windows\\System32\\drivers\\etc\\hosts (Windows)")
        print("   /etc/hosts (Linux/Mac)\n")
        for key, domain in config['domains'].items():
            print(f"127.0.0.1    {domain}")
    
    print("\n⚠️  IMPORTANT: Les fichiers .env contiennent des secrets sensibles!")
    print("   Ajoutez-les à .gitignore et ne les commitez JAMAIS!")
    
    print("\n🚀 Prochaines étapes:")
    print("   1. Vérifiez les fichiers .env créés")
    print("   2. Configurez Stripe (backend/.env) si nécessaire")
    print("   3. Lancez: python startOffice1789.py")


def _env_override(name: str) -> Optional[str]:
    val = os.getenv(name)
    if val:
        return val.strip()
    return None


def _build_url(domain: str, port: int | None, env: str) -> str:
    """Construit l'URL selon l'environnement.
    - local: http://domain:port (si port)
    - staging/production: https://domain (port ignoré pour standardisation)
    """
    if env == 'local':
        if port:
            return f"http://{domain}:{port}"
        return f"http://{domain}"
    # staging / production -> https, pas de port explicite (derrière reverse proxy)
    return f"https://{domain}".rstrip('/')


def load_config() -> AppConfig:
    """Charge la configuration runtime (config.json + variables d'environnement)."""
    raw = _load_raw_config()
    env = _env_override('OFFICE1789_ENV') or raw.get('environment', 'local')

    # Lecture des domaines / ports issus de la configuration
    domains = raw.get('domains', {})
    ports = raw.get('ports', {})

    # Appliquer overrides d'URL complètes si fournis
    backend_url = _env_override('OFFICE1789_BACKEND_URL')
    frontend_url = _env_override('OFFICE1789_FRONTEND_URL')
    roundcube_url = _env_override('OFFICE1789_ROUNDCUBE_URL')
    onlyoffice_url = _env_override('OFFICE1789_ONLYOFFICE_URL')
    matrix_url = _env_override('OFFICE1789_MATRIX_URL')
    element_url = _env_override('OFFICE1789_ELEMENT_URL')
    mail_url = _env_override('OFFICE1789_MAIL_URL')

    # Fallback domaines + ports
    backend = backend_url or _build_url(domains.get('main', DEFAULT_DOMAINS['main']), ports.get('backend', DEFAULT_PORTS['backend']), env)
    frontend = frontend_url or _build_url(domains.get('main', DEFAULT_DOMAINS['main']), ports.get('frontend', DEFAULT_PORTS['frontend']), env)
    roundcube = roundcube_url or _build_url(domains.get('mail', DEFAULT_DOMAINS['mail']), ports.get('roundcube', DEFAULT_PORTS['roundcube']), env)
    onlyoffice = onlyoffice_url or _build_url(domains.get('onlyoffice', DEFAULT_DOMAINS['onlyoffice']), ports.get('onlyoffice', DEFAULT_PORTS['onlyoffice']), env)
    matrix = matrix_url or _build_url(domains.get('matrix', DEFAULT_DOMAINS['matrix']), ports.get('synapse', 8008 if env == 'local' else None), env)
    element = element_url or _build_url(domains.get('element', DEFAULT_DOMAINS['element']), ports.get('element', DEFAULT_PORTS['element']), env)
    mail = mail_url or _build_url(domains.get('mail', DEFAULT_DOMAINS['mail']), None, env)

    stripe_key = _env_override('OFFICE1789_STRIPE_PUBLIC_KEY')

    services = ServiceURLs(
        backend=backend,
        frontend=frontend,
        roundcube=roundcube,
        onlyoffice=onlyoffice,
        matrix=matrix,
        element=element,
        mail=mail,
    )
    return AppConfig(environment=env, services=services, stripe_public_key=stripe_key)


def main_cli():
    import sys
    
    # Sans argument ou avec 'configure' -> mode interactif
    if len(sys.argv) == 1 or sys.argv[1] in {"configure", "setup", "init"}:
        try:
            interactive_configure()
        except KeyboardInterrupt:
            print("\n\n❌ Configuration annulée.")
        except Exception as e:
            print(f"\n❌ Erreur: {e}")
            import traceback
            traceback.print_exc()
        return
    
    cmd = sys.argv[1]
    
    # Commande pour regénérer tous les secrets (DANGER!)
    if cmd == "generate-secrets":
        print("\n⚠️  ATTENTION: Cette commande va REGÉNÉRER tous les secrets!")
        print("Cela nécessitera la recréation complète des bases de données.")
        confirm = input("\nÊtes-vous ABSOLUMENT SÛR? (tapez 'YES I AM SURE'): ").strip()
        if confirm == 'YES I AM SURE':
            config = _load_raw_config()
            config['secrets'] = generate_all_secrets()
            _save_config(config)
            env = config.get('environment', 'local')
            create_backend_env(config, config['secrets'], env)
            create_frontend_env(config, env)
            create_docker_env(config['secrets'])
            update_coturn_config(config['secrets'])
            print("\n✅ Tous les secrets ont été régénérés et les fichiers .env mis à jour!")
            print("⚠️  N'oubliez pas de recréer vos bases de données Docker!")
        else:
            print("❌ Opération annulée")
        return
    
    # Autres commandes nécessitent la config chargée
    cfg = load_config()
    
    if cmd in {"print", "show"}:
        config = _load_raw_config()
        print("Environment:", cfg.environment)
        print("\nURLs des services:")
        for k, v in cfg.services.__dict__.items():
            print(f"{k:12} -> {v}")
        if cfg.stripe_public_key:
            print("stripe_public_key ->", cfg.stripe_public_key)
        print(f"\nSecrets configurés: {len(config.get('secrets', {}))} éléments")
    elif cmd == "export-frontend":
        cfg.write_frontend_runtime()
        print(f"✅ Fichier runtime écrit: {FRONTEND_RUNTIME_JSON}")
    elif cmd == "json":
        print(json.dumps(cfg.to_dict(), indent=2, ensure_ascii=False))
    else:
        print("Commandes disponibles:")
        print("  python config.py                 # Configuration interactive")
        print("  python config.py print           # Afficher la configuration")
        print("  python config.py export-frontend # Exporter config frontend")
        print("  python config.py generate-secrets # REGÉNÉRER tous les secrets (DANGER!)")
        print("  python config.py json            # Export JSON complet")

if __name__ == "__main__":
    main_cli()
