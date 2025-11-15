# 🐧 Office1789 - Guide de déploiement Linux/Debian

## 📋 Prérequis Linux

### Debian/Ubuntu
```bash
# Mettre à jour le système
sudo apt update && sudo apt upgrade -y

# Installer Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# Installer Docker Compose
sudo apt install docker-compose -y

# Installer Go (backend)
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Installer Node.js (frontend)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# Installer PostgreSQL client (optionnel)
sudo apt install postgresql-client -y
```

### CentOS/RHEL/Fedora
```bash
# Installer Docker
sudo dnf install docker docker-compose -y
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker $USER

# Installer Go
sudo dnf install golang -y

# Installer Node.js
sudo dnf install nodejs npm -y
```

---

## 🚀 Installation Office1789 sur Linux

### 1. Cloner le projet
```bash
git clone https://github.com/LeSimple02/office1789.git
cd office1789
```

### 2. Configuration des permissions
```bash
# Rendre les scripts exécutables
chmod +x install_linux.sh
chmod +x test-matrix-sso.sh
chmod +x startOffice1789.py

# Permissions Docker
sudo chmod 666 /var/run/docker.sock  # Si erreurs de permission
```

### 3. Installation automatique
```bash
# Exécuter le script d'installation Linux
./install_linux.sh
```

### 4. Démarrage des services Docker
```bash
cd docker

# Générer les certificats SSL
./generate-mailserver-certs.ps1  # Adapter pour bash si nécessaire

# Démarrer tous les conteneurs
docker-compose up -d

# Vérifier l'état
docker-compose ps
```

### 5. Démarrage du backend Go
```bash
cd backend

# Installer les dépendances
go mod download

# Compiler
go build -o office1789-backend .

# Lancer (ou utiliser Air pour hot-reload)
./office1789-backend
# OU
go run main.go
```

### 6. Démarrage du frontend Vue
```bash
cd webfront2

# Installer les dépendances
npm install

# Mode développement
npm run dev

# OU mode production
npm run build
npm run preview
```

---

## 🔧 Configuration spécifique Linux

### 1. PostgreSQL
```bash
# Accéder à la base Office1789
docker exec -it postgres_db psql -U robespierre -d office1789

# Vérifier les tables
\dt

# Vérifier les utilisateurs
SELECT username, email FROM users;
```

### 2. Mail Server (docker-mailserver)
```bash
# Ajouter un compte mail
docker exec -it mailserver setup email add user@office1789.local password123

# Lister les comptes
docker exec -it mailserver setup email list

# Logs du serveur mail
docker logs mailserver -f
```

### 3. Matrix/Synapse
```bash
# Créer un utilisateur Matrix
docker exec -it synapse register_new_matrix_user \
  -u username \
  -p password \
  -a \
  -c /data/homeserver.yaml \
  http://localhost:8008

# Logs Synapse
docker logs synapse -f
```

### 4. Element avec SSO
```bash
# Rebuild Element avec le plugin SSO
cd docker
docker-compose stop element
docker-compose rm -f element
docker-compose build element
docker-compose up -d element

# Vérifier que le plugin est présent
docker exec element ls -la /app/office1789-sso.js
docker exec element grep "office1789-sso.js" /app/index.html
```

---

## 🧪 Tests sur Linux

### Test complet SSO
```bash
# Rendre le script exécutable
chmod +x test-matrix-sso.sh

# Exécuter les tests
./test-matrix-sso.sh
```

### Tests manuels
```bash
# Test Backend
curl -X POST http://localhost:8080/api/session/check \
  -H "Content-Type: application/json" \
  -d '{"token":"test"}'

# Test Matrix API
curl http://localhost:8008/_matrix/client/versions

# Test Element
curl http://localhost:8083 | grep "office1789-sso.js"

# Test Roundcube
curl http://localhost:8081
```

---

## 🔒 Sécurité Linux

### Firewall (UFW - Ubuntu/Debian)
```bash
# Activer UFW
sudo ufw enable

# Autoriser SSH
sudo ufw allow 22/tcp

# Autoriser les ports Office1789
sudo ufw allow 8080/tcp  # Backend
sudo ufw allow 8081/tcp  # Roundcube
sudo ufw allow 8082/tcp  # OnlyOffice
sudo ufw allow 8083/tcp  # Element
sudo ufw allow 8008/tcp  # Matrix

# Autoriser Mail
sudo ufw allow 25/tcp    # SMTP
sudo ufw allow 587/tcp   # SMTP TLS
sudo ufw allow 993/tcp   # IMAPS

# Statut
sudo ufw status
```

### Firewall (firewalld - CentOS/RHEL/Fedora)
```bash
# Démarrer firewalld
sudo systemctl start firewalld
sudo systemctl enable firewalld

# Autoriser les ports
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --permanent --add-port=8081/tcp
sudo firewall-cmd --permanent --add-port=8082/tcp
sudo firewall-cmd --permanent --add-port=8083/tcp
sudo firewall-cmd --permanent --add-port=8008/tcp
sudo firewall-cmd --permanent --add-port=25/tcp
sudo firewall-cmd --permanent --add-port=587/tcp
sudo firewall-cmd --permanent --add-port=993/tcp

# Recharger
sudo firewall-cmd --reload
```

### SELinux (CentOS/RHEL/Fedora)
```bash
# Vérifier le statut
getenforce

# Si problèmes avec Docker, passer en permissive (temporaire)
sudo setenforce 0

# OU désactiver définitivement (non recommandé en production)
sudo sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config
```

---

## 📊 Systemd Services (démarrage automatique)

### Backend Service
```bash
# Créer le service
sudo nano /etc/systemd/system/office1789-backend.service
```

```ini
[Unit]
Description=Office1789 Backend
After=network.target docker.service postgresql.service

[Service]
Type=simple
User=matthis
WorkingDirectory=/home/matthis/office1789/backend
ExecStart=/home/matthis/office1789/backend/office1789-backend
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

```bash
# Activer et démarrer
sudo systemctl daemon-reload
sudo systemctl enable office1789-backend
sudo systemctl start office1789-backend
sudo systemctl status office1789-backend
```

### Frontend Service (avec Nginx)
```bash
# Installer Nginx
sudo apt install nginx -y  # Debian/Ubuntu
sudo dnf install nginx -y  # CentOS/RHEL/Fedora

# Build frontend
cd /home/matthis/office1789/webfront2
npm run build

# Configuration Nginx
sudo nano /etc/nginx/sites-available/office1789
```

```nginx
server {
    listen 80;
    server_name office1789.com;

    root /home/matthis/office1789/webfront2/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
# Activer le site
sudo ln -s /etc/nginx/sites-available/office1789 /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Docker Compose Service
```bash
# Créer le service
sudo nano /etc/systemd/system/office1789-docker.service
```

```ini
[Unit]
Description=Office1789 Docker Services
After=docker.service
Requires=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/home/matthis/office1789/docker
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down

[Install]
WantedBy=multi-user.target
```

```bash
# Activer
sudo systemctl daemon-reload
sudo systemctl enable office1789-docker
sudo systemctl start office1789-docker
```

---

## 🐛 Troubleshooting Linux

### Permissions Docker
```bash
# Si erreur "permission denied" sur Docker socket
sudo chmod 666 /var/run/docker.sock

# OU ajouter l'utilisateur au groupe docker
sudo usermod -aG docker $USER
newgrp docker  # Appliquer sans reboot
```

### Ports déjà utilisés
```bash
# Trouver le processus utilisant un port
sudo lsof -i :8080
sudo netstat -tulpn | grep :8080

# Tuer le processus
sudo kill -9 <PID>
```

### Logs et debugging
```bash
# Logs backend
journalctl -u office1789-backend -f

# Logs Docker
docker-compose logs -f

# Logs spécifiques
docker logs synapse -f
docker logs element -f
docker logs roundcube -f
docker logs mailserver -f

# Espace disque
df -h
docker system df
docker system prune -a  # Nettoyer
```

### Problèmes réseau
```bash
# Vérifier la connectivité entre conteneurs
docker network inspect office_network

# Tester la connexion
docker exec backend ping postgres_db
docker exec element ping synapse
docker exec roundcube ping mailserver
```

---

## 🔄 Mise à jour Linux

### Backend
```bash
cd /home/matthis/office1789/backend
git pull origin main
go build -o office1789-backend .
sudo systemctl restart office1789-backend
```

### Frontend
```bash
cd /home/matthis/office1789/webfront2
git pull origin main
npm install
npm run build
sudo systemctl restart nginx
```

### Docker Services
```bash
cd /home/matthis/office1789/docker
git pull origin main
docker-compose pull
docker-compose up -d --build
```

---

## 📈 Monitoring Linux

### Ressources système
```bash
# CPU et RAM
htop
# OU
top

# Espace disque
df -h

# Docker stats
docker stats

# Logs système
journalctl -xe
dmesg | tail
```

### Monitoring automatique (Prometheus + Grafana)
```bash
# À ajouter dans docker-compose.yml
# prometheus:
#   image: prom/prometheus
#   ports:
#     - "9090:9090"
#
# grafana:
#   image: grafana/grafana
#   ports:
#     - "3000:3000"
```

---

## ✅ Compatibilité vérifiée

- ✅ **Debian 11/12** (Bullseye/Bookworm)
- ✅ **Ubuntu 20.04/22.04/24.04** LTS
- ✅ **CentOS Stream 8/9**
- ✅ **RHEL 8/9**
- ✅ **Fedora 38/39/40**
- ✅ **Rocky Linux 8/9**
- ✅ **AlmaLinux 8/9**

---

## 🎯 Checklist déploiement Linux

- [ ] Docker et Docker Compose installés
- [ ] Go 1.21+ installé
- [ ] Node.js 20+ installé
- [ ] Ports firewall ouverts (8080-8083, 8008, 25, 587, 993)
- [ ] Permissions Docker configurées
- [ ] Certificats SSL générés
- [ ] Services Docker démarrés (`docker-compose up -d`)
- [ ] Backend compilé et lancé
- [ ] Frontend build et servi (npm ou Nginx)
- [ ] Tests SSO passés (`./test-matrix-sso.sh`)
- [ ] Services systemd configurés (optionnel)
- [ ] Monitoring configuré (optionnel)

---

## 🏛️ Office1789 - Révolution numérique sur Linux !

**Système 100% compatible Linux avec SSO Matrix/Element et Roundcube/Mail en un clic !**
