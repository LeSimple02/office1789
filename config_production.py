#!/usr/bin/env python3
"""
Configuration Production Office1789
Simplifié pour déploiement serveur Debian

Usage:
    python3 config_production.py
"""
import json
import secrets
import string
from pathlib import Path

ROOT = Path(__file__).parent
CONFIG_JSON = ROOT / "config.json"
BACKEND_ENV = ROOT / "backend" / ".env"
WEBFRONT_ENV = ROOT / "webfront2" / ".env"
DOCKER_ENV = ROOT / "docker" / ".env"
COTURN_CONF = ROOT / "docker" / "coturn" / "turnserver.conf"

def generate_password(length=32, alphanumeric_only=False):
    """Génère un mot de passe sécurisé."""
    if alphanumeric_only:
        chars = string.ascii_letters + string.digits
    else:
        chars = string.ascii_letters + string.digits + "!@#$%^&*()-_=+[]{}|;:,.<>?"
    return ''.join(secrets.choice(chars) for _ in range(length))

def generate_aes_key():
    """Génère une clé AES-256."""
    return secrets.token_hex(32)

print("""
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║     🏛  OFFICE1789 - Configuration Production             ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
""")

# Demander le domaine principal
print("\n🌐 Configuration du domaine")
main_domain = input("Nom de domaine principal (ex: office1789.com): ").strip()

if not main_domain:
    print("❌ Domaine requis!")
    exit(1)

print(f"\n✅ Domaine: {main_domain}")
print(f"   - Frontend:   https://{main_domain}")
print(f"   - Backend:    https://backend.{main_domain}")
print(f"   - Mail:       https://mail.{main_domain}")
print(f"   - Chat:       https://chat.{main_domain}")
print(f"   - Docs:       https://docs.{main_domain}")
print(f"   - Matrix:     https://matrix.{main_domain}")

confirm = input("\n✓ Confirmer? (o/n) [o]: ").strip().lower() or 'o'
if confirm != 'o':
    print("❌ Configuration annulée")
    exit(0)

# Génération des secrets
print("\n🔐 Génération des secrets sécurisés...")

secrets_dict = {
    # PostgreSQL
    "postgres_user": "office1789_user",
    "postgres_password": generate_password(32, alphanumeric_only=True),
    "postgres_db": "office1789",
    
    # Roundcube
    "roundcube_db_user": "roundcube_user",
    "roundcube_db_password": generate_password(32, alphanumeric_only=True),
    "roundcube_db_name": "roundcube",
    
    # Synapse
    "synapse_db_user": "synapse_user",
    "synapse_db_password": generate_password(32, alphanumeric_only=True),
    "synapse_db_name": "synapse",
    
    # Services
    "coturn_secret": generate_password(48, alphanumeric_only=True),
    "onlyoffice_jwt_secret": generate_password(48),
    "aes_key": generate_aes_key(),
    "matrix_admin_token": generate_password(64, alphanumeric_only=True),
    "mail_admin_user": "admin",
    "mail_admin_password": generate_password(32),
}

print("✅ Secrets générés!")

# Sauvegarder config.json
config_data = {
    "environment": "production",
    "domains": {
        "main": main_domain,
        "mail": f"mail.{main_domain}",
        "matrix": f"matrix.{main_domain}",
        "element": f"chat.{main_domain}",
        "onlyoffice": f"docs.{main_domain}"
    },
    "ports": {
        "backend": 8080,
        "frontend": 80,
        "postgres": 5432,
        "roundcube": 8081,
        "element": 8083,
        "onlyoffice": 8082,
        "synapse": 8008,
        "coturn": 3478
    },
    "secrets": secrets_dict
}

with CONFIG_JSON.open('w', encoding='utf-8') as f:
    json.dump(config_data, f, indent=2)

print(f"✅ Config sauvegardée: {CONFIG_JSON}")

# Créer backend/.env
print("\n📝 Création backend/.env...")
backend_env = f"""# Production - Auto-generated
DB_HOST=localhost
DB_PORT=5432
DB_USER={secrets_dict['postgres_user']}
DB_PASSWORD={secrets_dict['postgres_password']}
DB_NAME={secrets_dict['postgres_db']}

# URLs Production
FRONTEND_URL=https://{main_domain}
ROUNDCUBE_URL=https://mail.{main_domain}
ELEMENT_URL=https://chat.{main_domain}
ONLYOFFICE_URL=https://docs.{main_domain}

# OVH SMS (à configurer)
OVH_SMS_APPLICATION_KEY=your_application_key
OVH_SMS_APPLICATION_SECRET=your_application_secret
OVH_SMS_CONSUMER_KEY=your_consumer_key
OVH_SMS_SERVICE_NAME=sms-xxXXXX-1
OVH_SMS_SENDER=Office1789

# Stripe (à configurer)
STRIPE_SECRET_KEY=sk_live_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
STRIPE_PRICE_STANDARD=price_standard_monthly
STRIPE_PRICE_PROFESSIONAL=price_professional_monthly
STRIPE_PRICE_ENTERPRISE=price_enterprise_monthly

# Security
AES_KEY={secrets_dict['aes_key']}

# Matrix
MATRIX_HOMESERVER=https://matrix.{main_domain}
MATRIX_ADMIN_TOKEN={secrets_dict['matrix_admin_token']}

# Mail
MAIL_SERVER=mail.{main_domain}
MAIL_ADMIN_USER={secrets_dict['mail_admin_user']}
MAIL_ADMIN_PASSWORD={secrets_dict['mail_admin_password']}

# OnlyOffice
ONLYOFFICE_JWT_SECRET={secrets_dict['onlyoffice_jwt_secret']}
"""

BACKEND_ENV.parent.mkdir(parents=True, exist_ok=True)
BACKEND_ENV.write_text(backend_env, encoding='utf-8')
print(f"✅ {BACKEND_ENV}")

# Créer webfront2/.env
print("\n📝 Création webfront2/.env...")
frontend_env = f"""# Production - Auto-generated
VITE_APP_API=https://backend.{main_domain}
VITE_APP_API_LOGIN=https://backend.{main_domain}/api/connect
VITE_APP_API_CREATE_ACCOUNT=https://backend.{main_domain}/api/subscribe
VITE_APP_API_GETINFOP=https://backend.{main_domain}/api/getinfop
VITE_APP_API_CHANGEINFO=https://backend.{main_domain}/api/changeinfo
VITE_APP_API_DELETE_ACCOUNT=https://backend.{main_domain}/api/account/delete
VITE_APP_API_CHANGE_PASSWORD=https://backend.{main_domain}/api/account/change-password
VITE_APP_API_PASSWORD_RESET_REQUEST=https://backend.{main_domain}/api/password/reset/request
VITE_APP_API_PASSWORD_RESET_CONFIRM=https://backend.{main_domain}/api/password/reset/confirm
VITE_APP_API_VERIFICATION_SEND=https://backend.{main_domain}/api/verification/send
VITE_APP_API_VERIFICATION_VERIFY=https://backend.{main_domain}/api/verification/verify

VITE_ROUNDCUBE_URL=https://mail.{main_domain}
VITE_ELEMENT_URL=https://chat.{main_domain}
VITE_ONLYOFFICE_URL=https://docs.{main_domain}
VITE_MATRIX_URL=https://matrix.{main_domain}
VITE_FRONTEND_URL=https://{main_domain}

VITE_ENV=production
"""

WEBFRONT_ENV.parent.mkdir(parents=True, exist_ok=True)
WEBFRONT_ENV.write_text(frontend_env, encoding='utf-8')
print(f"✅ {WEBFRONT_ENV}")

# Créer docker/.env
print("\n📝 Création docker/.env...")
docker_env = f"""# Production - Auto-generated
POSTGRES_USER={secrets_dict['postgres_user']}
POSTGRES_PASSWORD={secrets_dict['postgres_password']}
POSTGRES_DB={secrets_dict['postgres_db']}

ROUNDCUBE_DB_USER={secrets_dict['roundcube_db_user']}
ROUNDCUBE_DB_PASSWORD={secrets_dict['roundcube_db_password']}
ROUNDCUBE_DB_NAME={secrets_dict['roundcube_db_name']}

SYNAPSE_DB_USER={secrets_dict['synapse_db_user']}
SYNAPSE_DB_PASSWORD={secrets_dict['synapse_db_password']}
SYNAPSE_DB_NAME={secrets_dict['synapse_db_name']}

COTURN_SECRET={secrets_dict['coturn_secret']}
ONLYOFFICE_JWT_SECRET={secrets_dict['onlyoffice_jwt_secret']}
MATRIX_ADMIN_TOKEN={secrets_dict['matrix_admin_token']}

MAIL_ADMIN_USER={secrets_dict['mail_admin_user']}
MAIL_ADMIN_PASSWORD={secrets_dict['mail_admin_password']}
"""

DOCKER_ENV.parent.mkdir(parents=True, exist_ok=True)
DOCKER_ENV.write_text(docker_env, encoding='utf-8')
print(f"✅ {DOCKER_ENV}")

# Créer coturn config
print("\n📝 Mise à jour coturn/turnserver.conf...")
coturn_config = f"""use-auth-secret
static-auth-secret={secrets_dict['coturn_secret']}
realm={main_domain}
fingerprint
no-tls
no-dtls
listening-port=3478
min-port=49160
max-port=49200
user-quota=12
total-quota=1200
log-file=stdout
simple-log
"""

COTURN_CONF.parent.mkdir(parents=True, exist_ok=True)
COTURN_CONF.write_text(coturn_config, encoding='utf-8')
print(f"✅ {COTURN_CONF}")

# Résumé
print("\n" + "="*60)
print("✅ CONFIGURATION PRODUCTION TERMINÉE")
print("="*60)
print(f"\n🌐 Domaine principal: {main_domain}")
print(f"\n📁 Fichiers créés:")
print(f"   • {CONFIG_JSON}")
print(f"   • {BACKEND_ENV}")
print(f"   • {WEBFRONT_ENV}")
print(f"   • {DOCKER_ENV}")
print(f"   • {COTURN_CONF}")
print(f"\n🔐 {len(secrets_dict)} secrets générés")
print("\n⚠️  IMPORTANT:")
print("   1. Ne JAMAIS commiter les fichiers .env")
print("   2. Configurer Stripe dans backend/.env")
print("   3. Lancer: python3 start_production.py")
print("   4. Configurer DNS chez Gandi")
print("   5. Générer certificats SSL avec certbot")
print("")
