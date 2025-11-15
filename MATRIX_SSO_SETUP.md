# 🏛️ Office1789 - SSO Matrix/Element

## 📋 Vue d'ensemble

Système SSO complet permettant l'accès **en un clic** à Matrix/Element ET Roundcube/Mail, **sans modal de mot de passe**, comme Outlook ou OneDrive.

---

## 🔐 Architecture SSO

### Stockage du mot de passe
- **En RAM uniquement** : `session.Password` (string)
- Jamais persisté en base de données
- Expire après 24h (durée de la session)
- Protégé par `sync.RWMutex` (thread-safe)

### Flux d'authentification

```
1. Login Office1789 → Password stocké dans session (RAM)
2. Click "Ouvrir Element (Chat)" → POST /api/matrix/sso
3. Backend génère token SSO (HMAC-SHA256, 5 min)
4. Element s'ouvre avec ?sso_token=XXX
5. Plugin JS valide le token → Authentification Matrix
6. Stockage des credentials dans localStorage
```

---

## 🛠️ Composants implémentés

### Backend Go (`backend/`)

#### 1. `security.go` (120 lignes)
```go
type session struct {
    Username  string
    Email     string
    Password  string // Stocké en RAM pour SSO
    ExpiresAt int64
}

var (
    sessions      = make(map[string]session)
    sessionsMutex sync.RWMutex
)
```

#### 2. `matrix.go` (95 lignes)
```go
func GenerateMatrixSSOAuto(c *gin.Context) {
    // Récupère le password depuis session.Password
    // Génère token HMAC-SHA256 signé
    // User ID: @username:office1789.com
    // Return: http://localhost:8083/?sso_token=XXX
}
```

#### 3. `main.go`
```go
r.POST("/api/matrix/sso", GenerateMatrixSSOAuto)
r.POST("/api/mail/sso", GenerateMailSSOAuto)
```

### Frontend Vue 3 (`webfront2/src/`)

#### 1. `components/MatrixAccessSSO.vue` (140 lignes)
- Bouton unique "Ouvrir Element (Chat)"
- Appel `POST /api/matrix/sso` avec `{username, token}`
- Ouvre Element dans nouvelle fenêtre
- Style: Gradient bleu (vs purple pour mail)

#### 2. `views/ChatView.vue`
```vue
<MatrixAccessSSO 
  :username="username" 
  :token="sessionToken" 
/>
```

### Element Plugin (`docker/element/`)

#### 1. `office1789-sso.js` (220 lignes)
```javascript
// Web Crypto API - HMAC-SHA256 verification
// Intercepte ?sso_token= dans l'URL
// Décode token JWT (base64url)
// Valide signature HMAC
// Authentifie via Matrix API:
//   POST http://localhost:8008/_matrix/client/r0/login
// Stocke: mx_access_token, mx_user_id, mx_device_id
```

#### 2. `DockerfileElement/Dockerfile`
```dockerfile
FROM vectorim/element-web:latest
USER root
COPY office1789-sso.js /app/
RUN cp /app/index.html /tmp/index.html && \
    sed 's|</head>|<script src="/office1789-sso.js"></script></head>|' /tmp/index.html > /app/index.html.new && \
    mv /app/index.html.new /app/index.html
COPY config.json /app/config.json
USER nginx
```

---

## 🔑 Secrets & Configuration

### Backend Secrets
```go
// mail.go
secretKey := []byte("Office1789-SecretKey-ChangeInProduction")

// matrix.go
secretKey := []byte("Office1789-Matrix-SecretKey-ChangeInProduction")
```

### Matrix Configuration
- **Server Name**: `office1789.com`
- **User ID Format**: `@username:office1789.com`
- **Matrix API**: `http://localhost:8008/_matrix/client/r0/login`
- **Element URL**: `http://localhost:8083`

### Mail Configuration
- **IMAP**: `ssl://mailserver:993` (IMAPS)
- **SMTP**: `tls://mailserver:587`
- **Domain**: `office1789.local`
- **Email Format**: `username@office1789.local`

---

## 📦 Installation

### Prérequis
- Docker & Docker Compose
- Backend Office1789 démarré (`startOffice1789.py`)

### Déploiement Element avec SSO

```powershell
# Windows
cd docker
docker-compose stop element
docker-compose rm -f element
docker-compose build element
docker-compose up -d element
```

```bash
# Linux/Debian
cd docker
docker-compose stop element
docker-compose rm -f element
docker-compose build element
docker-compose up -d element
```

### Vérification
```bash
# Script SSO présent
docker exec element ls -la /app/office1789-sso.js

# Script injecté dans index.html
docker exec element grep "office1789-sso.js" /app/index.html
```

---

## 🧪 Tests SSO

### Test Mail SSO
1. Login Office1789 avec un compte existant
2. Aller dans **Mails** (MailView.vue)
3. Cliquer sur **"Ouvrir Roundcube (Mail)"**
4. ✅ Roundcube s'ouvre authentifié, sans modal de mot de passe

### Test Matrix SSO
1. Login Office1789 avec le même compte
2. Aller dans **Chat** (ChatView.vue)
3. Cliquer sur **"Ouvrir Element (Chat)"**
4. ✅ Element s'ouvre avec écran de chargement Office1789
5. ✅ Authentification automatique vers Matrix
6. ✅ Accès direct au chat

### Vérification des tokens
```bash
# Console navigateur (F12) sur Element
# Devrait afficher: "✅ Token SSO validé avec succès"
# Vérifier localStorage:
localStorage.getItem('mx_access_token')  // Token Matrix
localStorage.getItem('mx_user_id')       // @username:office1789.com
localStorage.getItem('mx_device_id')     // Device ID
```

---

## 🔒 Sécurité

### ✅ Points forts
- Password jamais persisté en base de données
- Token SSO expire après 5 minutes
- Signature HMAC-SHA256 pour validation
- Session expire après 24h
- Thread-safe (sync.RWMutex)

### ⚠️ Considérations
- **Password en RAM en clair** : Acceptable pour SSO, mais vulnérable si RAM compromise
- **Secrets hardcodés** : À changer en production (variables d'environnement)
- **SSL/TLS** : Activer HTTPS en production
- **CORS** : Configurer correctement pour éviter XSS

### 🚀 Production
```go
// Utiliser variables d'environnement
secretKey := []byte(os.Getenv("OFFICE1789_SSO_SECRET"))
matrixSecretKey := []byte(os.Getenv("OFFICE1789_MATRIX_SECRET"))

// Activer HTTPS
r.Use(tls.Middleware())

// Configurer CORS
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"https://office1789.com"},
}))
```

---

## 🐛 Troubleshooting

### Element ne s'authentifie pas
```bash
# Vérifier que Synapse est up
docker logs synapse

# Vérifier le script SSO
docker exec element cat /app/office1789-sso.js | grep "Office1789"

# Vérifier Matrix API
curl http://localhost:8008/_matrix/client/versions
```

### Token invalide
```bash
# Vérifier que les secrets correspondent
# Backend matrix.go vs Element office1789-sso.js
grep "SecretKey" backend/matrix.go
docker exec element grep "SECRET_KEY" /app/office1789-sso.js
```

### Session expirée
```bash
# Vérifier la session côté backend
# POST /api/session/check avec {token: "..."}
curl -X POST http://localhost:8080/api/session/check \
  -H "Content-Type: application/json" \
  -d '{"token": "your-session-token"}'
```

---

## 📊 État du projet

### ✅ Terminé
- [x] Stockage password en RAM (session.Password)
- [x] Thread-safety avec sync.RWMutex
- [x] Mail SSO (GenerateMailSSOAuto)
- [x] Matrix SSO (GenerateMatrixSSOAuto)
- [x] MailAccessSSO.vue (sans modal)
- [x] MatrixAccessSSO.vue (sans modal)
- [x] Element plugin JavaScript (office1789-sso.js)
- [x] Dockerfile Element personnalisé
- [x] Installation cross-platform (Python script)
- [x] IMAPS ssl://mailserver:993 pour Roundcube

### 🧪 À tester
- [ ] Mail SSO end-to-end avec vrais passwords
- [ ] Matrix SSO end-to-end avec authentification
- [ ] Expiration token SSO (5 min)
- [ ] Expiration session (24h)
- [ ] Thread-safety sous charge (concurrent requests)

### 🚀 Améliorations futures
- [ ] Variables d'environnement pour secrets
- [ ] HTTPS/TLS en production
- [ ] Refresh token pour Matrix
- [ ] Audit logs pour SSO
- [ ] Rate limiting sur endpoints SSO
- [ ] 2FA compatible avec SSO

---

## 📝 Notes techniques

### Format Token SSO
```
Header.Payload.Signature (base64url)

Payload (Matrix):
{
  "username": "matthis",
  "matrixUserId": "@matthis:office1789.com",
  "password": "plaintextPassword",
  "exp": 1234567890
}

Signature: HMAC-SHA256(Header + Payload, secretKey)
```

### Synchronisation passwords
- Passwords Office1789 = Passwords Matrix = Passwords Mail
- Synchronisés lors de la création de compte
- Stored hashed in database (bcrypt)
- Stored plaintext in RAM session (for SSO only)

### Compatibilité
- ✅ Windows (PowerShell)
- ✅ Linux (bash)
- ✅ Debian (systemd)
- ✅ Docker Desktop
- ✅ Docker Engine

---

## 🎯 Résumé

**L'utilisateur se connecte UNE FOIS sur Office1789, et accède ensuite à Mail et Matrix EN UN CLIC, sans ressaisir son mot de passe, comme Outlook/OneDrive.**

🏛️ **Office1789 - La révolution numérique !**
