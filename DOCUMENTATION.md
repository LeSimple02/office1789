# 🏛 Office1789 - Documentation complète

[Voir README.md pour le démarrage rapide]

## 📋 Table des matières

1. [Stack technique](#stack-technique)
2. [Prérequis](#prérequis)
3. [Installation](#installation)
4. [Configuration](#configuration)
5. [URLs d'accès](#urls-daccès)
6. [Commandes utiles](#commandes-utiles)
7. [Fonctionnalités](#fonctionnalités)
8. [Dépannage](#dépannage)

## 🚀 Stack technique

- **Frontend** : Vue.js 3 + Vite
- **Backend** : Gin (Go)  
- **Bases de données** : PostgreSQL 16 (x3)
- **Webmail** : Roundcube + Docker Mailserver
- **Chat** : Matrix Synapse + Element
- **Éditeur** : OnlyOffice DocumentServer
- **Conteneurisation** : Docker & Docker Compose

## ⚙️ Prérequis

- Python 3.8+
- Node.js 18+ et npm
- Go 1.21+
- Docker Desktop
- Git

## 🎯 Installation

Voir README.md

## 🌐 Configuration personnalisée

### Script de configuration

```bash
python configure.py
```

Permet de personnaliser :
- Domaines (office1789.com, mail.office1789.com, etc.)
- Ports d'écoute
- Options de démarrage automatique

La configuration est sauvegardée dans `config.json`.

### Structure du config.json

```json
{
  "domains": {
    "main": "office1789.com",
    "mail": "mail.office1789.com",
    "matrix": "matrix.office1789.com",
    "element": "chat.office1789.com",
    "onlyoffice": "docs.office1789.com"
  },
  "ports": {
    "frontend": 5173,
    "backend": 8080,
    "roundcube": 8081,
    "onlyoffice": 8082,
    "element": 8083
  },
  "autostart": {
    "docker_containers": true,
    "backend": true,
    "frontend": true,
    "open_browser": true
  }
}
```

## 📱 URLs d'accès

| Service | Localhost | Domaine personnalisé |
|---------|-----------|---------------------|
| Application | http://localhost:5173 | http://office1789.com:5173 |
| Webmail | http://localhost:8081 | http://mail.office1789.com:8081 |
| Chat | http://localhost:8083 | http://chat.office1789.com:8083 |
| OnlyOffice | http://localhost:8082 | http://docs.office1789.com:8082 |
| API Backend | http://localhost:8080 | http://office1789.com:8080 |

## 🔧 Commandes utiles

### Gestion des conteneurs

```bash
# Démarrer tous les conteneurs
cd docker && docker-compose up -d

# Arrêter tous les conteneurs
docker-compose down

# Voir les logs
docker-compose logs -f

# Redémarrer un conteneur spécifique
docker restart roundcube
```

### Développement

```bash
# Backend avec hot-reload
cd backend && air

# Frontend avec hot-reload
cd webfront2 && npm run dev

# Lancer les tests
npm run test
```

## 🌟 Fonctionnalités

### Implémenté ✅

- Authentification et sessions
- Webmail avec SSO auto-login
- Chat temps réel (Matrix)
- Drive (stockage/partage)
- Calendrier personnel
- Édition de documents
- Synchronisation des mots de passe
- Mode sombre
- Multi-langue (FR/EN)

### En développement 🚧

- Appels audio/vidéo
- Calendrier partagé
- Contacts synchronisés
- Notifications push

## 🛠 Dépannage

### Docker

```bash
# Vérifier Docker
docker --version
docker ps -a

# Reconstruire
docker-compose up -d --build --force-recreate
```

### Backend Go

```bash
# Vérifier Go
go version

# Réinstaller dépendances
cd backend && go mod download
```

### Frontend Vue

```bash
# Vérifier Node
node --version

# Réinstaller
cd webfront2 && rm -rf node_modules && npm install
```

## 📞 Support

- Issues GitHub
- Documentation wiki
- Chat Matrix

---

**Office1789 - Une alternative open-source et française à Office 365** 🇫🇷
