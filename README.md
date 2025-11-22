docker-compose up -d
<div align="center">
	<h1>Office1789</h1>
	<p><strong>Suite collaborative open‑source (mail, chat, documents, stockage, calendrier) — alternative souveraine à Office 365.</strong></p>
	<p>
		<img src="webfront2/public/logo.png" alt="Office1789" height="90" />
	</p>
</div>

## Sommaire
1. Vue d’ensemble
2. Architecture & composants
3. Fonctionnalités clés
4. Démarrage rapide (local)
5. Configuration & environnements
6. Mobile / Desktop (Cordova & Electron)
7. Déploiement production (aperçu)
8. Contribution & Licence

## 1. Vue d’ensemble
Office1789 réunit sous une interface intégrée : Webmail (Roundcube), Chat temps réel (Matrix / Element), édition collaborative (OnlyOffice), stockage fichiers (Drive), calendrier et gestion des comptes (SSO, TOTP, Stripe). Le backend est écrit en Go (Gin), le frontend en Vue 3 (Vite), empaqueté pour mobile (Cordova Android) et desktop (Electron).

## 2. Architecture & composants
| Couche | Tech | Rôle |
|--------|------|------|
| Frontend | Vue 3 + Vite | Interface utilisateur web + logique API unifiée |
| Mobile | Cordova Android | Wrapper léger (splash + stockage local) |
| Desktop | Electron | Application bureau (session persistante) |
| Backend | Go (Gin) | API REST, auth, routing, CORS, intégrations externes |
| Chat | Matrix Synapse | Messagerie temps réel |
| Client Chat | Element Web | Interface Matrix utilisateur |
| Docs | OnlyOffice DocServer | Édition collaborative (documents, tableurs) |
| Webmail | Roundcube | Courriel IMAP/SMTP via conteneur mail |
| SSO / Sécurité | TOTP, reset password | Sécurisation des accès |
| Paiements | Stripe | Abonnements / portails billing |

## 3. Fonctionnalités clés
✅ SSO entre services (Roundcube, Element, OnlyOffice)  
✅ Gestion d’organisation, utilisateurs, abonnements  
✅ Stockage & édition collaborative de documents  
✅ Chat (Matrix) + notifications temps réel  
✅ Mobile & Desktop packaging  
✅ Configuration centralisée des domaines / ports (`configure.py`, `config.py`)  

## 4. Démarrage rapide (local)
### Prérequis
| Outil | Version recommandée |
|-------|---------------------|
| Python | 3.10+ |
| Node.js | 18+ |
| Go | 1.21+ |
| Docker | Desktop / Engine |
| Android SDK (optionnel) | Cordova build |

### Installation
```bash
git clone https://github.com/LeSimple02/office1789.git
cd office1789
```

### Configuration interactive (facultative)
```bash
python configure.py    # génère / met à jour config.json
```

### Configuration programmatique / environnement
`config.py` charge `config.json` + variables d’environnement (ex: `OFFICE1789_BACKEND_URL`).
```bash
python config.py print
python config.py export-frontend    # écrit webfront2/src/config/runtime-config.json
```

### Conteneurs de base
```bash
cd docker
docker compose up -d
cd ..
```

### Backend Go
```bash
cd backend
go run .
```

### Frontend Vue (dev)
```bash
cd webfront2
npm install
npm run dev
```

### Build production + intégration Cordova
```bash
cd webfront2
npm run build
cd ..\cordova-app
cordova build android --debug
```

### Script de démarrage global
```bash
python startOffice1789.py
```

## 5. Configuration & environnements
`configure.py` : assistant interactif domaines / ports (écrit `config.json`).  
`config.py` : assemble URLs (local vs production) + overrides env (ex: `OFFICE1789_ENV=production`).  
Frontend consomme `runtime-config.json` si présent pour déterminer `API_BASE`.

Variables d’environnement courantes:
| Nom | Description |
|-----|-------------|
| OFFICE1789_ENV | local / staging / production |
| OFFICE1789_BACKEND_URL | URL explicite API backend |
| OFFICE1789_FRONTEND_URL | URL application web |
| OFFICE1789_MATRIX_URL | Base Synapse publique |
| OFFICE1789_ELEMENT_URL | Interface Element |
| OFFICE1789_ONLYOFFICE_URL | Serveur de documents |
| OFFICE1789_ROUNDCUBE_URL | Webmail |
| OFFICE1789_STRIPE_PUBLIC_KEY | Clé publique Stripe |

## 6. Mobile / Desktop
### Cordova Android
Sources dans `cordova-app/` (splash Android 12+, icône adaptative, stockage local). Build debug :
```bash
cd cordova-app
cordova build android --debug
adb install -r platforms/android/app/build/outputs/apk/debug/app-debug.apk
```
### Electron
Wrapper bureau dans `electron-app/` pour usage desktop hors navigateur. Build script (exemple à compléter selon packaging souhaité).

## 7. Déploiement production (aperçu)
Résumé rapide (document unique):
1. DNS: `app`, `api`, `mail`, `chat`, `matrix`, `docs` pointent vers votre reverse proxy.
2. Nginx (exemple API):
	 ```nginx
	 server {
		 server_name api.example.com;
		 listen 80; listen 443 ssl http2;
		 ssl_certificate /etc/letsencrypt/live/api.example.com/fullchain.pem;
		 ssl_certificate_key /etc/letsencrypt/live/api.example.com/privkey.pem;
		 location / {
			 proxy_pass http://127.0.0.1:8080/;
			 proxy_set_header Host $host;
			 proxy_set_header X-Real-IP $remote_addr;
			 proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
			 proxy_set_header X-Forwarded-Proto $scheme;
		 }
	 }
	 ```
3. Variables d'environnement (exemple):
	 ```bash
	 export OFFICE1789_ENV=production
	 export OFFICE1789_BACKEND_URL=https://api.example.com
	 export OFFICE1789_FRONTEND_URL=https://app.example.com
	 export OFFICE1789_MATRIX_URL=https://matrix.example.com
	 export OFFICE1789_ELEMENT_URL=https://chat.example.com
	 export OFFICE1789_ONLYOFFICE_URL=https://docs.example.com
	 export OFFICE1789_ROUNDCUBE_URL=https://mail.example.com
	 export OFFICE1789_STRIPE_PUBLIC_KEY=pk_live_xxx
	 python config.py export-frontend
	 ```
4. CORS: limiter aux domaines front officiels (`app`, `chat`, `mail`).
5. Headers sécurité: HSTS, X-Frame-Options=SAMEORIGIN, X-Content-Type-Options=nosniff, CSP durcie.
6. SSO: jetons générés par backend pour Element / Roundcube / OnlyOffice (aligner domaines exacts).
7. Stripe: définir clé publique (env), secret côté backend hors dépôt; configurer webhook.
8. Checklist finale: HTTPS ok, CORS restreint, logs activés, sauvegardes DB + volumes, dépendances à jour.

## 8. Contribution & Licence
Les contributions sont bienvenues (issues, PR).  
Licence: MIT – voir `LICENSE`.

---
Made in France 🇫🇷

