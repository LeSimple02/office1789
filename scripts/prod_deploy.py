import os
import re
import subprocess
import sys
from pathlib import Path
import secrets
import json
import crypt
from typing import Optional

print("==========================================")
print("DEPLOY OFFICE1789 (prod)")
print("==========================================")

REPO_DIR = Path("/home/debian/office1789")
DOCKER_DIR = REPO_DIR / "docker"
ENV_FILE = DOCKER_DIR / ".env"
BACKEND_ENV = REPO_DIR / "backend/.env"
USE_PREBUILT_IMAGES = os.environ.get("USE_PREBUILT_IMAGES", "0") == "1"
BUILD_LOCAL_IMAGES = os.environ.get("BUILD_LOCAL_IMAGES", "0") == "1"
CONFIG_DIR = DOCKER_DIR / "config"


def get_domain_base() -> str:
    """
    Determine DOMAIN_BASE for prod, preferring docker/.env if it exists,
    otherwise falling back to the DOMAIN_BASE environment variable or
    'office1789.com'.
    """
    if ENV_FILE.exists():
        try:
            for line in ENV_FILE.read_text(encoding="utf-8").splitlines():
                line = line.strip()
                if line.startswith("DOMAIN_BASE="):
                    return line.split("=", 1)[1].strip()
        except Exception:
            pass
    return os.environ.get("DOMAIN_BASE", "office1789.com")


def update_element_config_prod(domain_base: str) -> None:
    """
    Ensure docker/DockerfileElement/config.json points Element to the
    correct Matrix homeserver in production.
    """
    config_path = DOCKER_DIR / "DockerfileElement" / "config.json"
    base_url = f"https://matrix.{domain_base}"
    data = {
        "default_server_config": {
            "m.homeserver": {
                "base_url": base_url,
                "server_name": domain_base,
            }
        },
        "brand": "Office1789",
        "disable_custom_urls": True,
        "disable_guests": True,
        "default_theme": "light",
        "features": {"voip": True},
    }
    config_path.write_text(json.dumps(data, indent=2), encoding="utf-8")
    print(f"[ELEMENT] Updated {config_path} with hs={base_url}")


def get_env_value(env_path: Path, key: str) -> Optional[str]:
    if not env_path.exists():
        return None
    for line in env_path.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line or line.startswith("#"):
            continue
        if "=" not in line:
            continue
        k, v = line.split("=", 1)
        if k.strip() == key:
            return v.strip()
    return None


def ensure_mail_admin_account(config_dir: Path, username: str, password: str, domain_base: str) -> None:
    """
    Ensure postfix-accounts.cf contains a mailbox for the SMTP admin user so
    SMTP auth works out of the box in prod. Always rewrites the line to keep
    the password in sync with docker/.env.
    """
    if not username or not password:
        print("[MAIL][WARN] MAIL_ADMIN_USER/PASSWORD not set; cannot seed admin mailbox.")
        return
    email = f"{username}@{domain_base}"
    accounts_path = config_dir / "postfix-accounts.cf"
    config_dir.mkdir(parents=True, exist_ok=True)

    try:
        hashed = crypt.crypt(password, crypt.mksalt(crypt.METHOD_SHA512))
    except Exception as e:
        print(f"[MAIL][WARN] Could not hash password for admin mailbox: {e}")
        return

    # Rewrite file without existing admin line, then append current hash.
    lines: list[str] = []
    if accounts_path.exists():
        for line in accounts_path.read_text(encoding="utf-8").splitlines():
            if line.split("|", 1)[0].strip() == email:
                continue
            if line.strip():
                lines.append(line.strip())
    lines.append(f"{email}|{hashed}")
    accounts_path.write_text("\n".join(lines) + "\n", encoding="utf-8")
    print(f"[MAIL] Synced admin mailbox {email} in postfix-accounts.cf for SMTP auth")


def ensure_noreply_alias(config_dir: Path, target_email: str, domain_base: str) -> None:
    """
    Ensure postfix-virtual.cf contains noreply@... -> target_email alias.
    Rewrites the noreply line to ensure it always points to target_email.
    """
    if not target_email:
        return
    alias_path = config_dir / "postfix-virtual.cf"
    config_dir.mkdir(parents=True, exist_ok=True)
    line = f"noreply@{domain_base} {target_email}"
    lines: list[str] = []
    if alias_path.exists():
        for existing in alias_path.read_text(encoding="utf-8").splitlines():
            if existing.startswith(f"noreply@{domain_base} "):
                continue
            if existing.strip():
                lines.append(existing.strip())
    lines.append(line)
    alias_path.write_text("\n".join(lines) + "\n", encoding="utf-8")
    print(f"[MAIL] Synced noreply alias: {line}")


def upsert_env_var(text: str, key: str, value: str) -> tuple[str, bool]:
    """
    Ensure the env string contains key=value (overwrite if present).
    Returns (new_text, changed_flag).
    """
    changed = False
    pattern = rf"^{key}=.*$"
    replacement = f"{key}={value}"
    if re.search(pattern, text, flags=re.MULTILINE):
        new_text = re.sub(pattern, replacement, text, flags=re.MULTILINE)
        if new_text != text:
            changed = True
            text = new_text
    else:
        text = text.rstrip() + f"\n{replacement}\n"
        changed = True
    return text, changed

# Switch frontend env to prod
webfront_dir = REPO_DIR / "webfront2"
switch_env = webfront_dir / "switch_env.py"
if switch_env.exists():
    print("[ENV] switch_env.py prod...")
    subprocess.run([sys.executable, str(switch_env), "prod"], check=False)
else:
    print("[WARN] switch_env.py not found in webfront2, skipping env switch.")


# 1. Git update
def git_pull():
    print("[GIT] git pull...")
    subprocess.run(["git", "pull"], check=False)


git_pull()

# 2. Init env files if missing
print("[ENV] Checking environment files...")

# docker/.env
if not ENV_FILE.exists():
    print("[ENV] Generating docker/.env with secrets...")
    POSTGRES_USER = "office1789"
    POSTGRES_PASSWORD = secrets.token_hex(32)
    POSTGRES_DB = "office1789"
    DB_HOST = "postgres_db"
    DB_PORT = "5432"
    DB_USER = POSTGRES_USER
    DB_PASSWORD = POSTGRES_PASSWORD
    DB_NAME = POSTGRES_DB
    ROUNDCUBE_DB_NAME = "roundcube"
    ROUNDCUBE_DB_USER = "roundcube"
    ROUNDCUBE_DB_PASSWORD = secrets.token_hex(32)
    SYNAPSE_DB_NAME = "synapse"
    SYNAPSE_DB_USER = "synapse"
    # Keep Synapse DB password simple and in sync with homeserver.yaml
    SYNAPSE_DB_PASSWORD = "synapse"

    COTURN_SECRET = secrets.token_hex(32)
    ONLYOFFICE_JWT_SECRET = os.environ.get("ONLYOFFICE_JWT_SECRET") or secrets.token_hex(32)
    MATRIX_ADMIN_TOKEN = secrets.token_hex(32)
    MAIL_ADMIN_PASSWORD = secrets.token_hex(32)
    MAIL_ADMIN_USER = os.environ.get("MAIL_ADMIN_USER", "admin")
    DOMAIN_BASE = os.environ.get("DOMAIN_BASE", "office1789.com")
    if not DOMAIN_BASE:
        print("Erreur : DOMAIN_BASE n'est pas défini. Modifie prod_deploy.py ou exporte DOMAIN_BASE avant d'exécuter le script.")
        sys.exit(1)
    BACKEND_URL = f"https://backend.{DOMAIN_BASE}"
    FRONTEND_URL = f"https://{DOMAIN_BASE}"
    MAIL_URL = f"https://mail.{DOMAIN_BASE}"
    CHAT_URL = f"https://chat.{DOMAIN_BASE}"
    DOCS_URL = f"https://docs.{DOMAIN_BASE}"
    ROUNDCUBE_URL = f"https://mail.{DOMAIN_BASE}"
    # En prod, Element est servi sur chat.<domaine>
    ELEMENT_URL = f"https://chat.{DOMAIN_BASE}"
    env_text = f"""# === MAIN DB ===
POSTGRES_USER={POSTGRES_USER}
POSTGRES_PASSWORD={POSTGRES_PASSWORD}
POSTGRES_DB={POSTGRES_DB}

# === BACKEND DB VARS ===
DB_HOST={DB_HOST}
DB_PORT={DB_PORT}
DB_USER={DB_USER}
DB_PASSWORD={DB_PASSWORD}
DB_NAME={DB_NAME}

# === ROUNDCUBE DB ===
ROUNDCUBE_DB_NAME={ROUNDCUBE_DB_NAME}
ROUNDCUBE_DB_USER={ROUNDCUBE_DB_USER}
ROUNDCUBE_DB_PASSWORD={ROUNDCUBE_DB_PASSWORD}

# === SYNAPSE DB ===
SYNAPSE_DB_NAME={SYNAPSE_DB_NAME}
SYNAPSE_DB_USER={SYNAPSE_DB_USER}
SYNAPSE_DB_PASSWORD={SYNAPSE_DB_PASSWORD}

# === SECRETS ===
COTURN_SECRET={COTURN_SECRET}
ONLYOFFICE_JWT_SECRET={ONLYOFFICE_JWT_SECRET}
MATRIX_ADMIN_TOKEN={MATRIX_ADMIN_TOKEN}
MAIL_ADMIN_USER={MAIL_ADMIN_USER}
MAIL_ADMIN_PASSWORD={MAIL_ADMIN_PASSWORD}
STRIPE_SECRET_KEY={os.environ.get("STRIPE_SECRET_KEY", "set-me")}
STRIPE_WEBHOOK_SECRET={os.environ.get("STRIPE_WEBHOOK_SECRET", "set-me")}
STRIPE_PRICE_STANDARD={os.environ.get("STRIPE_PRICE_STANDARD", "set-me")}
STRIPE_PRICE_PROFESSIONAL={os.environ.get("STRIPE_PRICE_PROFESSIONAL", "set-me")}
STRIPE_PRICE_ENTERPRISE={os.environ.get("STRIPE_PRICE_ENTERPRISE", "set-me")}

# === FRONTEND URLS ===
DOMAIN_BASE={DOMAIN_BASE}
BACKEND_URL={BACKEND_URL}
FRONTEND_URL={FRONTEND_URL}
MAIL_URL={MAIL_URL}
CHAT_URL={CHAT_URL}
DOCS_URL={DOCS_URL}
ROUNDCUBE_URL={ROUNDCUBE_URL}
ELEMENT_URL={ELEMENT_URL}
# Explicit Roundcube IMAP/SMTP hosts (useful for docker-compose runtime)
ROUNDCUBE_IMAP_HOST=mailserver:143
ROUNDCUBE_SMTP_HOST=mailserver:587
# === SMTP (backend mail/verification)
SMTP_HOST=mailserver
SMTP_PORT=587
SMTP_FROM=noreply@{DOMAIN_BASE}
SMTP_USERNAME={MAIL_ADMIN_USER}@{DOMAIN_BASE}
SMTP_PASSWORD={MAIL_ADMIN_PASSWORD}
# Mailserver relay policy (allow docker compose network) + enable submission/auth
PERMIT_DOCKER=connected-networks
ENABLE_SUBMISSION=1
ENABLE_SMTP=1
ENABLE_SMTP_AUTH=1
ENABLE_AMAVIS=0

# === IMAGES (prebuilt) ===
BACKEND_IMAGE=office1789/backend:latest
FRONTEND_IMAGE=office1789/frontend:latest
"""
    DOCKER_DIR.mkdir(parents=True, exist_ok=True)
    ENV_FILE.write_text(env_text, encoding="utf-8")
    print("[ENV] docker/.env created (secrets generated).")
else:
    print("[ENV] docker/.env found, keeping it.")
    # Migration douce: si ELEMENT_URL pointe encore vers element.<domaine>, le basculer sur chat.<domaine>
    try:
        text = ENV_FILE.read_text(encoding="utf-8")
        if "ELEMENT_URL=" in text:
            domain_base_effective = get_domain_base()
            desired_line = f"ELEMENT_URL=https://chat.{domain_base_effective}"
            new_text = re.sub(r"^ELEMENT_URL=.*$", desired_line, text, flags=re.MULTILINE)
            if new_text != text:
                ENV_FILE.write_text(new_text, encoding="utf-8")
                print(f"[ENV] Updated ELEMENT_URL to {desired_line}")
                text = new_text

        # Ensure SMTP defaults are present and aligned with MAIL_ADMIN_USER/PASSWORD
        domain_base_effective = get_domain_base()
        mail_admin_user = get_env_value(ENV_FILE, "MAIL_ADMIN_USER") or os.environ.get("MAIL_ADMIN_USER", "admin")
        mail_admin_password = get_env_value(ENV_FILE, "MAIL_ADMIN_PASSWORD") or os.environ.get("MAIL_ADMIN_PASSWORD", "admin")
        desired_values = {
            "SMTP_HOST": "mailserver",
            "SMTP_PORT": "587",
            "SMTP_FROM": f"noreply@{domain_base_effective}",
            "SMTP_USERNAME": f"{mail_admin_user}@{domain_base_effective}",
            "SMTP_PASSWORD": mail_admin_password,
            "FRONTEND_URL": f"https://{domain_base_effective}",
            "MAIL_ADMIN_USER": mail_admin_user,
            "MAIL_ADMIN_PASSWORD": mail_admin_password,
            "PERMIT_DOCKER": "connected-networks",
            "ENABLE_SUBMISSION": "1",
            "ENABLE_SMTP": "1",
            "ENABLE_SMTP_AUTH": "1",
            "ENABLE_AMAVIS": "0",
            "STRIPE_SECRET_KEY": os.environ.get("STRIPE_SECRET_KEY", "set-me"),
            "STRIPE_WEBHOOK_SECRET": os.environ.get("STRIPE_WEBHOOK_SECRET", "set-me"),
            "STRIPE_PRICE_STANDARD": os.environ.get("STRIPE_PRICE_STANDARD", "set-me"),
            "STRIPE_PRICE_PROFESSIONAL": os.environ.get("STRIPE_PRICE_PROFESSIONAL", "set-me"),
            "STRIPE_PRICE_ENTERPRISE": os.environ.get("STRIPE_PRICE_ENTERPRISE", "set-me"),
        }
        changed = False
        for key, value in desired_values.items():
            text, updated = upsert_env_var(text, key, value)
            changed = changed or updated
        if changed:
            ENV_FILE.write_text(text, encoding="utf-8")
            print("[ENV] Updated SMTP_*/MAIL_ADMIN_* values in docker/.env")
    except Exception as e:
        print(f"[WARN] Could not migrate ELEMENT_URL in docker/.env: {e}")

# Ensure Element config points to the right Matrix homeserver for prod
domain_base_effective = get_domain_base()
update_element_config_prod(domain_base_effective)

# Ensure SMTP admin mailbox exists for docker-mailserver
mail_admin_user = get_env_value(ENV_FILE, "MAIL_ADMIN_USER") or os.environ.get("MAIL_ADMIN_USER", "admin")
mail_admin_password = get_env_value(ENV_FILE, "MAIL_ADMIN_PASSWORD") or os.environ.get("MAIL_ADMIN_PASSWORD", "admin")
ensure_mail_admin_account(CONFIG_DIR, mail_admin_user, mail_admin_password, domain_base_effective)
ensure_noreply_alias(CONFIG_DIR, f"{mail_admin_user}@{domain_base_effective}", domain_base_effective)

# POST: regenerate postfix accounts/aliases in the running mailserver
print("[POST] Updating mailserver accounts/aliases")
subprocess.run(["docker", "compose", "-f", "docker-compose.yml", "exec", "-T", "mailserver", "postmap", "/tmp/docker-mailserver/postfix-virtual.cf"], cwd=DOCKER_DIR, check=False)
subprocess.run(["docker", "compose", "-f", "docker-compose.yml", "exec", "-T", "mailserver", "postfix", "reload"], cwd=DOCKER_DIR, check=False)

# backend/.env
if not BACKEND_ENV.exists():
    print("[ENV] Generating backend/.env...")
    BACKEND_ENV.write_text(
        "DB_HOST=postgres_db\nDB_PORT=5432\nDB_USER=office1789\nDB_PASSWORD="
        f"{os.environ.get('POSTGRES_PASSWORD', 'office1789')}\nDB_NAME=office1789\n",
        encoding="utf-8",
    )
    print("[ENV] backend/.env created.")
else:
    print("[ENV] backend/.env already exists, keeping it.")

# backend/.env.production for Matrix (prod)
backend_env_production = REPO_DIR / "backend/.env.production"
backend_env_content = f"""# Backend prod env auto-generated by prod_deploy.py
DB_HOST=postgres_db
DB_PORT=5432
DB_USER=office1789
DB_PASSWORD={os.environ.get('POSTGRES_PASSWORD', 'office1789')}
DB_NAME=office1789
FRONTEND_URL=https://office1789.com
MAIL_SERVER=mail.office1789.com
MAIL_ADMIN_USER={os.environ.get('MAIL_ADMIN_USER', 'admin')}
MAIL_ADMIN_PASSWORD={os.environ.get('MAIL_ADMIN_PASSWORD', 'admin')}
MATRIX_HOMESERVER=https://matrix.office1789.com
MATRIX_ADMIN_TOKEN={os.environ.get('MATRIX_ADMIN_TOKEN', 'changeme')}
SMTP_HOST=mailserver
SMTP_PORT=587
SMTP_FROM=noreply@{get_domain_base()}
SMTP_USERNAME={os.environ.get('MAIL_ADMIN_USER', 'admin')}@{get_domain_base()}
SMTP_PASSWORD={os.environ.get('MAIL_ADMIN_PASSWORD', 'admin')}
PERMIT_DOCKER=connected-networks
STRIPE_SECRET_KEY={os.environ.get('STRIPE_SECRET_KEY', 'set-me')}
STRIPE_WEBHOOK_SECRET={os.environ.get('STRIPE_WEBHOOK_SECRET', 'set-me')}
STRIPE_PRICE_STANDARD={os.environ.get('STRIPE_PRICE_STANDARD', 'set-me')}
STRIPE_PRICE_PROFESSIONAL={os.environ.get('STRIPE_PRICE_PROFESSIONAL', 'set-me')}
STRIPE_PRICE_ENTERPRISE={os.environ.get('STRIPE_PRICE_ENTERPRISE', 'set-me')}
"""
backend_env_production.write_text(backend_env_content, encoding="utf-8")
print(f"Generated {backend_env_production} (Matrix prod config)")

# frontend .env.production
webfront_env_production = REPO_DIR / "webfront2/.env.production"
webfront_env_content = """# URLs de production (aucun localhost)
VITE_APP_API=https://backend.office1789.com
VITE_APP_API_LOGIN=https://backend.office1789.com/api/connect
VITE_APP_API_CREATE_ACCOUNT=https://backend.office1789.com/api/subscribe
VITE_APP_API_GETINFOP=https://backend.office1789.com/api/getinfop
VITE_APP_API_CHANGEINFO=https://backend.office1789.com/api/changeinfo

# Account routes
VITE_APP_API_DELETE_ACCOUNT=https://backend.office1789.com/api/account/delete
VITE_APP_API_CHANGE_PASSWORD=https://backend.office1789.com/api/account/change-password

# Password reset
VITE_APP_API_PASSWORD_RESET_REQUEST=https://backend.office1789.com/api/password/reset/request
VITE_APP_API_PASSWORD_RESET_CONFIRM=https://backend.office1789.com/api/password/reset/confirm

# Verification
VITE_APP_API_VERIFICATION_SEND=https://backend.office1789.com/api/verification/send
VITE_APP_API_VERIFICATION_VERIFY=https://backend.office1789.com/api/verification/verify
VITE_MATRIX_URL=https://matrix.office1789.com
VITE_ROUNDCUBE_URL=https://mail.office1789.com
VITE_CHAT_URL=https://chat.office1789.com
VITE_ELEMENT_URL=https://chat.office1789.com
VITE_ONLYOFFICE_URL=https://docs.office1789.com
VITE_FRONTEND_URL=https://office1789.com
VITE_ENV=production
"""
webfront_env_production.write_text(webfront_env_content, encoding="utf-8")
print(f"Generated {webfront_env_production} (frontend prod env)")

# 3. Docker compose
print("[DOCKER] Updating containers...")
os.chdir(DOCKER_DIR)

compose_base = ["docker", "compose", "-f", "docker-compose.yml"]

if USE_PREBUILT_IMAGES:
    print("[DOCKER] Pulling prebuilt images (backend, frontend)...")
    r1 = subprocess.run(compose_base + ["pull", "backend"], check=False)
    r2 = subprocess.run(compose_base + ["pull", "frontend"], check=False)
    if r1.returncode != 0 or r2.returncode != 0:
        print("[WARN] Pull failed (missing image or auth). Building locally instead.")
        subprocess.run(compose_base + ["up", "-d", "--build"], check=False)
    else:
        subprocess.run(compose_base + ["up", "-d"], check=False)
else:
    # Default: build locally using compose build definitions
    subprocess.run(compose_base + ["up", "-d", "--build"], check=False)

# 3.c Align Synapse DB password with docker/.env (best effort)
fix_synapse_script = REPO_DIR / "scripts/fix_synapse_password.py"
if fix_synapse_script.exists():
    print("[SYNAPSE] Aligning DB password with docker/.env...")
    subprocess.run([sys.executable, str(fix_synapse_script)], check=False)
    print("[SYNAPSE] You can now ensure Synapse starts with: cd docker && docker compose up -d synapse")

# 3.b Optional: fix Synapse host directory ownership
SYNAPSE_HOST_DIR = DOCKER_DIR / "synapse/conf"
if SYNAPSE_HOST_DIR.exists():
    print(f"[DOCKER] Checking Synapse ownership for {SYNAPSE_HOST_DIR}")
    synapse_uid = None
    try:
        result = subprocess.run(
            ["docker", "compose", "run", "--rm", "--entrypoint", "sh", "synapse", "-c", "id -u"],
            capture_output=True,
            text=True,
            check=True,
        )
        synapse_uid = result.stdout.strip()
    except Exception:
        try:
            result = subprocess.run(
                ["docker", "run", "--rm", "--entrypoint", "sh", "matrixdotorg/synapse:latest", "-c", "id -u"],
                capture_output=True,
                text=True,
                check=True,
            )
            synapse_uid = result.stdout.strip()
        except Exception:
            pass
    if synapse_uid:
        print(f"[DOCKER] Detected synapse UID={synapse_uid} -> attempting chown on host dir (may require sudo)")
        try:
            subprocess.run(["chown", "-R", f"{synapse_uid}:{synapse_uid}", str(SYNAPSE_HOST_DIR)], check=True)
            print(f"[DOCKER] Ownership updated for {SYNAPSE_HOST_DIR}")
            subprocess.run(["docker", "compose", "restart", "synapse"], check=False)
        except Exception:
            print(f"[WARN] Could not chown {SYNAPSE_HOST_DIR}. Run with elevated rights if Synapse fails:")
            print(f"  sudo chown -R {synapse_uid}:{synapse_uid} {SYNAPSE_HOST_DIR}")
    else:
        print("[WARN] Unable to determine Synapse UID automatically.")
else:
    print(f"[WARN] Synapse host directory {SYNAPSE_HOST_DIR} not found, skipping ownership check")

# 4. Post-deploy: mailserver & provisioning helpers
print("[POST] Ensuring mailserver has virtual mailboxes/maps")
postfix_map = DOCKER_DIR / "config/postfix-virtual-mailboxes.cf"
if postfix_map.exists():
    print("[POST] Found postfix-virtual-mailboxes.cf -> copying and postmap")
    subprocess.run(["docker", "cp", str(postfix_map), "mailserver:/tmp/docker-mailserver/postfix-virtual-mailboxes.cf"])
    subprocess.run(["docker", "exec", "mailserver", "postmap", "/tmp/docker-mailserver/postfix-virtual-mailboxes.cf"])
    subprocess.run(["docker", "exec", "mailserver", "postfix", "reload"])
    print("[POST] postfix maps updated")
else:
    print("[POST] No postfix-virtual-mailboxes.cf found, skipping map update")

virtual_aliases = DOCKER_DIR / "config/postfix-virtual.cf"
if virtual_aliases.exists():
    print("[POST] Found postfix-virtual.cf -> copying and postmap")
    subprocess.run(["docker", "cp", str(virtual_aliases), "mailserver:/tmp/docker-mailserver/postfix-virtual.cf"])
    subprocess.run(["docker", "exec", "mailserver", "postmap", "/tmp/docker-mailserver/postfix-virtual.cf"])
    subprocess.run(["docker", "exec", "mailserver", "postfix", "reload"])
    print("[POST] postfix virtual aliases updated")
else:
    print("[POST] No postfix-virtual.cf found, skipping virtual alias update")


print("==========================================")
print("DEPLOYMENT COMPLETE")
print("==========================================")
