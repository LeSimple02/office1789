# 🌍 Office1789 - Compatibilité Multi-Plateforme

## ✅ Systèmes supportés

| Composant | Windows | Linux | macOS | Docker |
|-----------|---------|-------|-------|--------|
| **Backend Go** | ✅ | ✅ | ✅ | ✅ |
| **Frontend Vue** | ✅ | ✅ | ✅ | ✅ |
| **PostgreSQL** | ✅ | ✅ | ✅ | ✅ |
| **Mail Server** | ⚠️ (Docker) | ✅ | ⚠️ (Docker) | ✅ |
| **Roundcube** | ⚠️ (Docker) | ✅ | ⚠️ (Docker) | ✅ |
| **Matrix/Synapse** | ⚠️ (Docker) | ✅ | ⚠️ (Docker) | ✅ |
| **Element** | ⚠️ (Docker) | ✅ | ⚠️ (Docker) | ✅ |
| **OnlyOffice** | ⚠️ (Docker) | ✅ | ⚠️ (Docker) | ✅ |

⚠️ = Nécessite Docker Desktop

---

## 🔧 Architecture Cross-Platform

### Backend Go
```
✅ Go 1.21+ fonctionne identiquement sur Windows/Linux/macOS
✅ Binaires natifs pour chaque plateforme
✅ Pas de dépendances système spécifiques
✅ Compilation: go build -o office1789-backend .
```

**Compilation cross-platform:**
```bash
# Windows depuis Linux
GOOS=windows GOARCH=amd64 go build -o office1789-backend.exe .

# Linux depuis Windows
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o office1789-backend .

# macOS depuis Windows/Linux
GOOS=darwin GOARCH=amd64 go build -o office1789-backend-macos .
```

### Frontend Vue 3
```
✅ Node.js 20+ fonctionne identiquement sur toutes les plateformes
✅ npm run dev / npm run build identique
✅ Build génère du HTML/JS/CSS statique (portable)
✅ Peut être servi par n'importe quel serveur web
```

### Services Docker
```
✅ docker-compose.yml identique sur toutes les plateformes
✅ Images Docker identiques (Linux containers)
✅ Networking Docker fonctionne pareil
✅ Volumes Docker gérés par Docker Engine
```

---

## 📦 Installation par plateforme

### Windows
```powershell
# Installation automatique
.\install_windows.ps1

# Démarrage
python startOffice1789.py
```

**Prérequis:**
- Docker Desktop for Windows
- Go 1.21+
- Node.js 20+
- Python 3.8+

### Linux (Debian/Ubuntu)
```bash
# Installation automatique
chmod +x install_linux.sh
./install_linux.sh

# Démarrage
python3 startOffice1789.py
```

**Prérequis:**
- Docker CE + Docker Compose
- Go 1.21+
- Node.js 20+
- Python 3.8+

### macOS
```bash
# Installation automatique
chmod +x install_linux.sh  # Même script que Linux
./install_linux.sh

# Démarrage
python3 startOffice1789.py
```

**Prérequis:**
- Docker Desktop for Mac
- Go 1.21+ (via Homebrew: brew install go)
- Node.js 20+ (via Homebrew: brew install node)
- Python 3.8+ (préinstallé sur macOS 10.15+)

---

## 🔐 SSO Cross-Platform

### Fonctionnement identique sur toutes les plateformes

**Backend:**
- `backend/matrix.go` - HMAC-SHA256 (crypto/hmac, standard Go)
- `backend/mail.go` - HMAC-SHA256 (crypto/hmac, standard Go)
- `backend/security.go` - sync.RWMutex (standard Go)

**Frontend:**
- `webfront2/src/components/MatrixAccessSSO.vue` - JavaScript pur
- `webfront2/src/components/MailAccessSSO.vue` - JavaScript pur
- Pas de dépendances système

**Element Plugin:**
- `docker/element/office1789-sso.js` - Web Crypto API (standard W3C)
- Fonctionne dans tous les navigateurs modernes
- Pas de dépendances natives

**Roundcube Plugin:**
- `docker/DockerfileRC/plugins/office1789_sso/` - PHP pur
- IMAP via extension PHP standard
- Fonctionne dans tous les conteneurs PHP

---

## 🧪 Tests Cross-Platform

### Windows
```powershell
.\test-matrix-sso.ps1
```

### Linux/macOS
```bash
chmod +x test-matrix-sso.sh
./test-matrix-sso.sh
```

**Tests identiques:**
1. Backend API (curl/Invoke-WebRequest)
2. Matrix API (curl/Invoke-WebRequest)
3. Element (curl/Invoke-WebRequest)
4. Script SSO (docker exec)
5. Endpoint SSO (manuel)

---

## 🌐 URLs et Ports (identiques partout)

| Service | URL | Port | Protocole |
|---------|-----|------|-----------|
| Backend | http://localhost:8080 | 8080 | HTTP |
| Roundcube | http://localhost:8081 | 8081 | HTTP |
| OnlyOffice | http://localhost:8082 | 8082 | HTTP |
| Element | http://localhost:8083 | 8083 | HTTP |
| Matrix | http://localhost:8008 | 8008 | HTTP |
| PostgreSQL | localhost:5432 | 5432 | TCP |
| SMTP | localhost:25 | 25 | SMTP |
| SMTP TLS | localhost:587 | 587 | SMTP+TLS |
| IMAPS | localhost:993 | 993 | IMAPS |

---

## 🔄 Différences par plateforme

### Chemins de fichiers

**Windows:**
```
C:\Users\Matthis\Documents\office1789\
C:\Users\Matthis\Documents\office1789\backend\
C:\Users\Matthis\Documents\office1789\webfront2\
```

**Linux/macOS:**
```
/home/matthis/office1789/
/home/matthis/office1789/backend/
/home/matthis/office1789/webfront2/
```

### Commandes terminal

| Action | Windows (PowerShell) | Linux/macOS (Bash) |
|--------|---------------------|-------------------|
| Lister fichiers | `Get-ChildItem` / `ls` | `ls -la` |
| Changer répertoire | `cd` | `cd` |
| Variables env | `$env:VAR="value"` | `export VAR=value` |
| Exécuter script | `.\script.ps1` | `./script.sh` |
| Joindre commandes | `;` | `;` ou `&&` |
| Copier fichier | `Copy-Item` | `cp` |
| Supprimer fichier | `Remove-Item` | `rm` |

### Docker

**Windows (Docker Desktop):**
- Interface graphique disponible
- WSL2 backend recommandé
- Volumes Windows: `/c/Users/...`

**Linux (Docker CE):**
- Ligne de commande uniquement
- Natif (pas de VM)
- Volumes Linux: `/home/...`

**macOS (Docker Desktop):**
- Interface graphique disponible
- HyperKit/VirtioFS backend
- Volumes macOS: `/Users/...`

---

## 🏗️ Build Production Cross-Platform

### Backend
```bash
# Build pour toutes les plateformes
./build-all.sh

# OU manuellement
GOOS=windows GOARCH=amd64 go build -o bin/office1789-backend-windows.exe .
GOOS=linux GOARCH=amd64 go build -o bin/office1789-backend-linux .
GOOS=darwin GOARCH=amd64 go build -o bin/office1789-backend-macos .
```

### Frontend
```bash
# Build production (identique partout)
cd webfront2
npm run build

# Dossier dist/ contient les fichiers statiques
# Peut être servi par:
# - Nginx (Linux)
# - Apache (Linux/Windows)
# - IIS (Windows)
# - Caddy (tous)
# - Python http.server (tous)
```

### Docker Images
```bash
# Build images (identiques partout)
cd docker
docker-compose build

# Export images pour transfert
docker save -o office1789-images.tar \
  docker-element:latest \
  docker-roundcube:latest

# Import sur autre machine
docker load -i office1789-images.tar
```

---

## 🔒 Sécurité Cross-Platform

### Secrets
```bash
# Variables d'environnement (tous)
# backend/.env
OFFICE1789_SSO_SECRET=your-secret-key
OFFICE1789_MATRIX_SECRET=your-matrix-secret
POSTGRES_PASSWORD=your-db-password
```

### Certificats SSL/TLS
```bash
# Génération identique (OpenSSL disponible partout)
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout key.pem -out cert.pem
```

### Firewall

**Windows (Windows Defender Firewall):**
```powershell
New-NetFirewallRule -DisplayName "Office1789 Backend" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow
```

**Linux (UFW):**
```bash
sudo ufw allow 8080/tcp
```

**macOS (PF):**
```bash
# Généralement pas nécessaire en développement
# Docker Desktop gère les ports automatiquement
```

---

## 📊 Performance par plateforme

### Backend Go
- **Windows**: Performance native, identique à Linux
- **Linux**: Performance optimale (natif)
- **macOS**: Performance native, légèrement inférieure sur ARM (M1/M2)

### Frontend Vue
- **Tous**: Identique (JavaScript dans le navigateur)

### Docker
- **Windows**: Légère overhead (WSL2 ou Hyper-V)
- **Linux**: Performance maximale (natif)
- **macOS**: Légère overhead (HyperKit/VirtioFS)

### Base de données
- **Tous**: Performance identique (PostgreSQL dans Docker)

---

## 🚀 Recommandations par plateforme

### Développement
- **Windows**: ✅ Excellent (WSL2 + Docker Desktop)
- **Linux**: ✅ Excellent (natif)
- **macOS**: ✅ Excellent (M1/M2 supportés)

### Production
- **Linux**: ⭐ **RECOMMANDÉ** (Debian/Ubuntu/CentOS)
  - Performance maximale
  - Coût le plus bas
  - Stabilité éprouvée
  - Outils de monitoring natifs

- **Windows Server**: ✅ OK (si infrastructure Windows existante)
  - Support Microsoft
  - Active Directory integration
  - Coût plus élevé

- **macOS Server**: ❌ Non recommandé
  - Pas conçu pour production
  - Coût très élevé
  - Support limité

---

## 📝 Checklist Compatibilité

### ✅ Identique sur toutes les plateformes
- [x] Code backend Go (100% portable)
- [x] Code frontend Vue (100% portable)
- [x] Services Docker (docker-compose.yml)
- [x] Images Docker (Linux containers)
- [x] Algorithmes SSO (HMAC-SHA256)
- [x] Tokens JWT (base64url)
- [x] API endpoints (/api/*)
- [x] Ports et URLs
- [x] Base de données (PostgreSQL)
- [x] Tests automatisés

### ⚠️ Adapté par plateforme
- [x] Scripts d'installation (.ps1 vs .sh)
- [x] Scripts de test (.ps1 vs .sh)
- [x] Chemins de fichiers (Windows vs Unix)
- [x] Variables d'environnement (syntaxe)
- [x] Services système (systemd vs Windows Services)
- [x] Firewall (Windows Defender vs UFW vs firewalld)

---

## 🎯 Résumé

**Office1789 est 100% cross-platform !**

- ✅ **Backend Go**: Fonctionne identiquement sur Windows/Linux/macOS
- ✅ **Frontend Vue**: JavaScript pur, portable partout
- ✅ **Docker**: Même configuration sur toutes les plateformes
- ✅ **SSO Matrix/Element**: Web Crypto API standard
- ✅ **SSO Roundcube**: PHP standard
- ✅ **Tests**: Scripts adaptés (.ps1/.sh) mais logique identique

**Déploiement recommandé:**
- 🧪 **Développement**: Windows/Linux/macOS (au choix)
- 🚀 **Production**: Linux (Debian 12 / Ubuntu 24.04 LTS)

---

## 🏛️ Office1789 - La révolution numérique, partout !
