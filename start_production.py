#!/usr/bin/env python3
"""
Démarrage Production Office1789
Script automatisé pour serveur Debian

Usage:
    sudo python3 start_production.py
"""
import json
import os
import sys
import subprocess
from pathlib import Path

ROOT = Path(__file__).parent
CONFIG_JSON = ROOT / "config.json"

def run_command(cmd, cwd=None, check=True):
    """Execute une commande shell."""
    print(f"🔧 {cmd}")
    try:
        result = subprocess.run(
            cmd,
            shell=True,
            cwd=cwd or ROOT,
            check=check,
            capture_output=True,
            text=True
        )
        if result.stdout:
            print(result.stdout)
        return result
    except subprocess.CalledProcessError as e:
        print(f"❌ Erreur: {e}")
        if e.stderr:
            print(e.stderr)
        if check:
            raise
        return e

def check_root():
    """Vérifie si le script est exécuté en root."""
    # Ne plus exiger root - on utilisera sudo uniquement quand nécessaire
    if os.geteuid() == 0:
        print("⚠️  Attention: Script exécuté en root")
        print("   Les commandes npm/node utiliseront l'environnement root")
    else:
        print("✓ Script exécuté en utilisateur normal (recommandé)")

def load_config():
    """Charge config.json."""
    if not CONFIG_JSON.exists():
        print(f"❌ {CONFIG_JSON} introuvable!")
        print("   Exécuter d'abord: python3 config_production.py")
        sys.exit(1)
    
    with CONFIG_JSON.open('r') as f:
        return json.load(f)

def check_dependencies():
    """Vérifie les dépendances système."""
    print("\n📦 Vérification des dépendances...")
    
    deps = {
        "docker": "docker --version",
        "docker-compose": "docker compose version",
        "go": "/usr/local/go/bin/go version",
        "nginx": "/usr/sbin/nginx -v"
    }
    
    missing = []
    for name, cmd in deps.items():
        result = run_command(cmd, check=False)
        if result.returncode != 0:
            missing.append(name)
            print(f"   ❌ {name} non installé")
        else:
            print(f"   ✅ {name}")
    
    # Vérification spéciale pour Node.js (version minimum requise: 18)
    node_result = run_command("node --version", check=False)
    if node_result.returncode == 0:
        node_version = node_result.stdout.strip()
        # Extraire le numéro de version majeure (ex: v20.11.0 -> 20)
        version_match = node_version.replace('v', '').split('.')[0]
        try:
            major_version = int(version_match)
            if major_version >= 18:
                print(f"   ✅ node ({node_version})")
            else:
                print(f"   ❌ node ({node_version}) - Version trop ancienne, >=18 requis")
                print("\n🔄 Installation de Node.js 20 LTS...")
                run_command("curl -fsSL https://deb.nodesource.com/setup_20.x | bash -")
                run_command("apt install -y nodejs")
                print("   ✅ Node.js 20 installé")
        except:
            missing.append("node")
    else:
        print(f"   ❌ node non installé")
        print("\n🔄 Installation de Node.js 20 LTS...")
        run_command("curl -fsSL https://deb.nodesource.com/setup_20.x | bash -")
        run_command("apt install -y nodejs")
        print("   ✅ Node.js 20 installé")
    
    if missing:
        print(f"\n❌ Dépendances manquantes: {', '.join(missing)}")
        print("   Installer avec: sudo bash install_linux.sh")
        sys.exit(1)

def setup_permissions():
    """Configure les permissions Docker."""
    print("\n🔒 Configuration des permissions Docker...")
    
    # Synapse
    synapse_path = ROOT / "docker" / "synapse" / "conf"
    if synapse_path.exists():
        run_command(f"sudo chown -R 991:991 {synapse_path}")
        run_command(f"sudo chmod -R 755 {synapse_path}")
        print("   ✅ Synapse")
    
    # Mailserver SSL
    ssl_path = ROOT / "docker" / "config" / "ssl"
    ssl_path.mkdir(parents=True, exist_ok=True)
    
    key_file = ssl_path / "key.pem"
    cert_file = ssl_path / "cert.pem"
    
    if not key_file.exists() or not cert_file.exists():
        print("   🔐 Génération certificats SSL temporaires (mail)...")
        config = load_config()
        domain = config["domains"]["mail"]
        run_command(
            f'openssl req -x509 -newkey rsa:4096 '
            f'-keyout {key_file} -out {cert_file} '
            f'-days 365 -nodes '
            f'-subj "/C=FR/ST=France/L=Paris/O=Office1789/CN={domain}"'
        )
        print(f"   ✅ SSL temporaire: {ssl_path}")

def install_backend_deps():
    """Installe les dépendances Go."""
    print("\n📦 Installation dépendances Backend (Go)...")
    backend_path = ROOT / "backend"
    
    if not (backend_path / "go.mod").exists():
        print("   ❌ go.mod introuvable!")
        return
    
    # Utiliser le chemin absolu de Go
    go_cmd = "/usr/local/go/bin/go"
    if not Path(go_cmd).exists():
        go_cmd = "go"  # Fallback
    
    run_command(f"{go_cmd} mod tidy", cwd=backend_path)
    run_command(f"{go_cmd} mod download", cwd=backend_path)
    print("   ✅ Dépendances Go installées")

def install_frontend_deps():
    """Installe les dépendances npm et build le frontend."""
    print("\n📦 Installation dépendances Frontend (npm)...")
    frontend_path = ROOT / "webfront2"
    
    if not (frontend_path / "package.json").exists():
        print("   ❌ package.json introuvable!")
        return
    
    run_command("npm install", cwd=frontend_path)
    print("   ✅ Dépendances npm installées")
    
    # Build production Vue.js
    print("\n🏗️  Build production Vue.js...")
    run_command("npm run build", cwd=frontend_path)
    
    dist_path = frontend_path / "dist"
    if dist_path.exists():
        file_count = len(list(dist_path.glob("**/*")))
        print(f"   ✅ Build terminé: {file_count} fichiers dans dist/")
    else:
        print("   ⚠️  Dossier dist/ non trouvé après build")

def start_docker_services():
    """Démarre les services Docker."""
    print("\n🐳 Démarrage des services Docker...")
    docker_path = ROOT / "docker"
    
    # Build et démarrage
    run_command("docker compose pull", cwd=docker_path)
    run_command("docker compose up -d --build", cwd=docker_path)
    
    print("\n📊 État des conteneurs:")
    run_command("docker compose ps", cwd=docker_path)

def start_backend():
    """Démarre le backend Go."""
    print("\n🚀 Démarrage Backend Go...")
    backend_path = ROOT / "backend"
    
    # Utiliser le chemin absolu de Go
    go_cmd = "/usr/local/go/bin/go"
    if not Path(go_cmd).exists():
        go_cmd = "go"  # Fallback
    
    # Build
    run_command(f"{go_cmd} build -o office1789-backend", cwd=backend_path)
    
    # Créer service systemd avec sudo
    service_content = f"""[Unit]
Description=Office1789 Backend API
After=network.target docker.service

[Service]
Type=simple
User=debian
WorkingDirectory={backend_path}
EnvironmentFile={backend_path}/.env
ExecStart={backend_path}/office1789-backend
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
"""
    
    # Écrire dans un fichier temporaire puis copier avec sudo
    temp_service = ROOT / "office1789-backend.service"
    temp_service.write_text(service_content)
    
    run_command(f"sudo mv {temp_service} /etc/systemd/system/office1789-backend.service")
    run_command("sudo systemctl daemon-reload")
    run_command("sudo systemctl enable office1789-backend")
    run_command("sudo systemctl restart office1789-backend")
    
    print("   ✅ Backend démarré (systemd)")
    run_command("sudo systemctl status office1789-backend --no-pager", check=False)

def configure_nginx():
    """Configure nginx reverse proxy."""
    print("\n🌐 Configuration nginx...")
    config = load_config()
    domain = config["domains"]["main"]
    
    nginx_config = f"""# Office1789 - Configuration Production
server {{
    listen 80;
    server_name {domain};
    
    location / {{
        proxy_pass http://localhost:5173;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }}
}}

server {{
    listen 80;
    server_name backend.{domain};
    
    location / {{
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }}
}}

server {{
    listen 80;
    server_name mail.{domain};
    
    location / {{
        proxy_pass http://localhost:8081;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }}
}}

server {{
    listen 80;
    server_name chat.{domain};
    
    location / {{
        proxy_pass http://localhost:8083;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }}
}}

server {{
    listen 80;
    server_name docs.{domain};
    
    location / {{
        proxy_pass http://localhost:8082;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Proto $scheme;
    }}
}}

server {{
    listen 80;
    server_name matrix.{domain};
    
    location / {{
        proxy_pass http://localhost:8008;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }}
}}
"""
    
    nginx_path = Path(f"/etc/nginx/sites-available/office1789")
    
    # Écrire dans un fichier temporaire puis copier avec sudo
    temp_nginx = ROOT / "office1789-nginx.conf"
    temp_nginx.write_text(nginx_config)
    
    run_command(f"sudo mv {temp_nginx} {nginx_path}")
    
    # Activer le site
    link = Path("/etc/nginx/sites-enabled/office1789")
    run_command(f"sudo rm -f {link}", check=False)
    run_command(f"sudo ln -s {nginx_path} {link}")
    
    # Tester et recharger
    run_command("sudo /usr/sbin/nginx -t")
    run_command("sudo systemctl reload nginx")
    
    print("   ✅ nginx configuré")

def configure_firewall():
    """Configure le firewall ufw."""
    print("\n🔥 Configuration firewall...")
    
    # Vérifier si ufw est installé
    check_ufw = run_command("which ufw", check=False)
    if check_ufw.returncode != 0:
        print("   ⚠️  ufw non installé, installation...")
        run_command("sudo apt update", check=False)
        run_command("sudo apt install -y ufw")
    
    ports = [
        ("22", "SSH"),
        ("80", "HTTP"),
        ("443", "HTTPS"),
        ("25", "SMTP"),
        ("587", "Submission"),
        ("993", "IMAPS")
    ]
    
    for port, desc in ports:
        run_command(f"sudo ufw allow {port}/tcp", check=False)
        print(f"   ✅ {port}/tcp ({desc})")
    
    run_command("sudo ufw --force enable", check=False)
    print("   ✅ Firewall configuré")

def print_status():
    """Affiche l'état final."""
    config = load_config()
    domain = config["domains"]["main"]
    
    print("\n" + "="*60)
    print("✅ OFFICE1789 DÉMARRÉ EN PRODUCTION")
    print("="*60)
    print(f"\n🌐 URLs (HTTP - avant SSL):")
    print(f"   • Frontend:   http://{domain}")
    print(f"   • Backend:    http://backend.{domain}")
    print(f"   • Mail:       http://mail.{domain}")
    print(f"   • Chat:       http://chat.{domain}")
    print(f"   • Docs:       http://docs.{domain}")
    print(f"   • Matrix:     http://matrix.{domain}")
    
    print(f"\n📊 Services:")
    print(f"   • Backend:    systemctl status office1789-backend")
    print(f"   • Docker:     cd docker && docker compose ps")
    print(f"   • Nginx:      systemctl status nginx")
    
    print(f"\n⚠️  PROCHAINES ÉTAPES:")
    print(f"   1. Configurer DNS chez Gandi:")
    print(f"      - A    @           → IP serveur")
    print(f"      - A    backend     → IP serveur")
    print(f"      - A    mail        → IP serveur")
    print(f"      - A    chat        → IP serveur")
    print(f"      - A    docs        → IP serveur")
    print(f"      - A    matrix      → IP serveur")
    print(f"      - MX   @           → mail.{domain} (priorité 10)")
    print(f"")
    print(f"   2. Générer certificats SSL:")
    print(f"      sudo apt install certbot python3-certbot-nginx")
    print(f"      sudo certbot --nginx -d {domain} -d backend.{domain} \\")
    print(f"        -d mail.{domain} -d chat.{domain} \\")
    print(f"        -d docs.{domain} -d matrix.{domain}")
    print(f"")
    print(f"   3. Configurer Stripe dans backend/.env")
    print("")

def main():
    """Point d'entrée principal."""
    print("""
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║     🏛  OFFICE1789 - Démarrage Production                 ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
""")
    
    # Vérifications
    check_root()
    config = load_config()
    print(f"✅ Configuration chargée: {config['domains']['main']}")
    
    check_dependencies()
    
    # Configuration
    setup_permissions()
    install_backend_deps()
    install_frontend_deps()
    
    # Démarrage
    start_docker_services()
    start_backend()
    configure_nginx()
    configure_firewall()
    
    # Status final
    print_status()

if __name__ == "__main__":
    main()
