#  Office1789

**Office1789** est une suite collaborative complète et open-source, alternative française à Office 365.

##  Fonctionnalités

-  **Webmail** (Roundcube) avec auto-login SSO
-  **Chat temps réel** (Matrix/Element)  
-  **Édition de documents** (OnlyOffice)
-  **Stockage de fichiers** (Drive)
-  **Calendrier**
-  **Gestion de comptes** avec synchronisation automatique

##  Démarrage rapide

### 1. Prérequis

- Python 3.8+
- Node.js 18+ et npm
- Go 1.21+
- Docker Desktop

### 2. Installation

```bash
git clone https://github.com/LeSimple02/office1789.git
cd office1789
```

### 3. Configuration (optionnel)

Personnalisez vos domaines et ports :

```bash
python configure.py
```

### 4. Premier lancement

Créez les conteneurs Docker :

```bash
cd docker
docker-compose up -d
cd ..
```

### 5. Démarrage automatique

```bash
python startOffice1789.py
```

Le script va automatiquement :
-  Démarrer tous les conteneurs Docker
-  Lancer le backend Go
-  Lancer le frontend Vue.js
-  Ouvrir votre navigateur

##  Accès aux services

| Service | URL |
|---------|-----|
| Application | http://localhost:5173 |
| Webmail | http://localhost:8081 |
| Chat | http://localhost:8083 |
| OnlyOffice | http://localhost:8082 |

##  Domaines personnalisés

Configurez votre fichier hosts :

**Windows** : `C:\Windows\System32\drivers\etc\hosts`

```
127.0.0.1    office1789.com
127.0.0.1    mail.office1789.com
127.0.0.1    chat.office1789.com
```

##  Documentation

- [DOCUMENTATION.md](DOCUMENTATION.md) - Documentation complète
- [config.json](config.json) - Configuration personnalisée

##  Contribution

Les contributions sont bienvenues ! Voir [CONTRIBUTING.md](CONTRIBUTING.md)

##  Licence

MIT License - voir [LICENSE](LICENSE)

---

**Développé avec  en France** 
