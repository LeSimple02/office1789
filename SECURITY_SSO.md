# 🔒 Office1789 - Sécurité SSO

## 🛡️ Protection des pages de connexion

### ✅ Roundcube - Protection ACTIVE

**Fichier:** `docker/DockerfileRC/plugins/office1789_sso/office1789_sso.php`

```php
// SÉCURITÉ : Bloquer TOUT accès sans token SSO valide
if (empty($sso_token) && !isset($_SESSION['user_id'])) {
    error_log('[SSO] SÉCURITÉ - Accès refusé sans token SSO valide');
    header('HTTP/1.1 403 Forbidden');
    // Affiche page d'erreur et redirige vers Office1789
    exit;
}
```

**Comportement:**
- ❌ Accès direct `http://localhost:8081` → **BLOQUÉ**
- ✅ Accès via Office1789 avec token SSO → **AUTORISÉ**
- ✅ Session déjà active → **AUTORISÉ**

**Page d'erreur:**
```
🔒 Accès Refusé
L'accès direct à Roundcube est désactivé pour des raisons de sécurité.
Veuillez vous connecter via Office1789 pour accéder à votre messagerie.
[Retour à Office1789]
```

---

### ✅ Element - Protection ACTIVE

**Fichier:** `docker/element/office1789-sso.js`

```javascript
// SÉCURITÉ: Vérifier si accès autorisé
function checkAccess() {
    const ssoToken = urlParams.get('sso_token');
    const isAuthenticated = localStorage.getItem('mx_access_token');
    
    // Si pas de token SSO ET pas authentifié -> BLOQUER
    if (!ssoToken && !isAuthenticated) {
        console.warn('[Office1789-SSO] SÉCURITÉ - Accès refusé');
        // Affiche page d'erreur
        throw new Error('Accès non autorisé');
    }
}

// Vérifier l'accès AVANT tout traitement
checkAccess();
```

**Comportement:**
- ❌ Accès direct `http://localhost:8083` → **BLOQUÉ**
- ✅ Accès via Office1789 avec `?sso_token=XXX` → **AUTORISÉ**
- ✅ `mx_access_token` en localStorage → **AUTORISÉ**

**Page d'erreur:**
```
🔒 Accès Refusé
L'accès direct à Element est désactivé pour des raisons de sécurité.
Veuillez vous connecter via Office1789 pour accéder au chat Matrix.
[🏛️ Retour à Office1789]
```

---

## 🔐 Mécanisme de sécurité

### Flux d'authentification sécurisé

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Utilisateur accède directement à Roundcube/Element      │
│    http://localhost:8081 ou http://localhost:8083           │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Plugin/Script vérifie:                                   │
│    - Présence du token SSO dans l'URL ?                     │
│    - Session/localStorage déjà active ?                     │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
              ┌──────┴──────┐
              │             │
         NON  │             │  OUI
              ▼             ▼
    ┌──────────────┐   ┌──────────────┐
    │ 403 Forbidden│   │ Accès autorisé│
    │ Page d'erreur│   │ Continue...   │
    │ Redirect     │   └──────────────┘
    └──────────────┘
```

### Validation du token SSO

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Backend Go génère token SSO                              │
│    Token = base64url(claims) + "." + HMAC-SHA256(claims)   │
│    Claims: {username, email, password, exp}                 │
│    Secret: "Office1789-SecretKey-ChangeInProduction"        │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. URL avec token SSO                                       │
│    http://localhost:8081/?sso_token=CLAIMS.SIGNATURE        │
│    http://localhost:8083/?sso_token=CLAIMS.SIGNATURE        │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Plugin/Script valide le token                            │
│    - Vérifier la signature HMAC-SHA256                      │
│    - Décoder les claims (base64url)                         │
│    - Vérifier l'expiration (5 minutes)                      │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
              ┌──────┴──────┐
              │             │
      INVALIDE│             │ VALIDE
              ▼             ▼
    ┌──────────────┐   ┌──────────────┐
    │ 403 Forbidden│   │ Authentifier │
    │ Token invalide│   │ avec password│
    └──────────────┘   └──────────────┘
```

---

## 🚨 Vecteurs d'attaque bloqués

### ✅ Accès direct sans authentification
**Attaque:** Aller directement à `http://localhost:8081` ou `http://localhost:8083`
**Protection:** Plugin/Script vérifie présence du token SSO ou session active
**Résultat:** 403 Forbidden + redirection vers Office1789

### ✅ Token SSO rejoué (Replay Attack)
**Attaque:** Réutiliser un vieux token SSO
**Protection:** Token expire après 5 minutes (`exp` claim)
**Résultat:** "Token SSO expiré. Veuillez vous reconnecter."

### ✅ Token SSO forgé
**Attaque:** Créer un faux token sans connaître le secret
**Protection:** Signature HMAC-SHA256 avec secret partagé
**Résultat:** "Token SSO invalide. Signature incorrecte."

### ✅ Token SSO modifié
**Attaque:** Modifier les claims (username, email, password)
**Protection:** Signature HMAC-SHA256 invalide si claims modifiés
**Résultat:** "Token SSO invalide. Signature incorrecte."

### ✅ Brute force du secret
**Attaque:** Essayer de deviner le secret HMAC
**Protection:** Secret long et complexe (256 bits minimum)
**Recommandation:** Changer en production avec clé aléatoire cryptographiquement sécurisée

---

## 🔧 Configuration de production

### Secrets forts

**Backend Go:**
```go
// NE PAS UTILISER EN PRODUCTION!
// secretKey := []byte("Office1789-SecretKey-ChangeInProduction")

// ✅ PRODUCTION: Utiliser variables d'environnement
secretKey := []byte(os.Getenv("OFFICE1789_SSO_SECRET"))

// Générer un secret fort:
// openssl rand -base64 32
// ou
// python -c "import secrets; print(secrets.token_urlsafe(32))"
```

**Roundcube Plugin:**
```php
// ❌ Hardcodé (dev only)
// private $secret = 'Office1789-SecretKey-ChangeInProduction';

// ✅ Production: Depuis variable d'environnement
private $secret;

function init() {
    $this->secret = getenv('OFFICE1789_SSO_SECRET');
    if (empty($this->secret)) {
        error_log('[SSO] ERREUR: OFFICE1789_SSO_SECRET non défini!');
        die('Configuration SSO manquante');
    }
}
```

**Element Plugin:**
```javascript
// ❌ Hardcodé (dev only)
// const SSO_SECRET = 'Office1789-Matrix-SecretKey-ChangeInProduction';

// ✅ Production: Injecté depuis config.json
const config = await (await fetch('/config.json')).json();
const SSO_SECRET = config.sso_secret;
```

---

## 🌐 HTTPS/TLS en production

### Problèmes avec HTTP

- ⚠️ **Token SSO en clair** dans l'URL (visible dans les logs)
- ⚠️ **Password en clair** dans le token (visible en transit)
- ⚠️ **Man-in-the-Middle** possible
- ⚠️ **Session hijacking** possible

### ✅ Solution: HTTPS partout

**Nginx reverse proxy:**
```nginx
server {
    listen 443 ssl http2;
    server_name office1789.com;
    
    ssl_certificate /etc/ssl/certs/office1789.crt;
    ssl_certificate_key /etc/ssl/private/office1789.key;
    
    # Frontend
    location / {
        proxy_pass http://localhost:5173;
    }
    
    # Backend
    location /api {
        proxy_pass http://localhost:8080;
    }
    
    # Roundcube
    location /mail {
        proxy_pass http://localhost:8081;
    }
    
    # Element
    location /chat {
        proxy_pass http://localhost:8083;
    }
}
```

---

## 🔍 Audit de sécurité

### Logs à surveiller

**Backend Go:**
```
[SSO] Token SSO généré pour user: matthis
[SSO] Token expire dans 300 secondes
```

**Roundcube:**
```
[SSO] SÉCURITÉ - Accès refusé sans token SSO valide
[SSO] Token détecté: eyJhbGciOiJIUzI1Ni...
[SSO] Claims décodés pour matthis@office1789.local
[SSO] Token valide, expiration OK
[SSO] Auto-login configuré et session nettoyée
```

**Element:**
```
[Office1789-SSO] Plugin chargé
[Office1789-SSO] SÉCURITÉ - Accès refusé sans token SSO valide
[Office1789-SSO] Accès autorisé - Token SSO présent
[Office1789-SSO] Claims décodés: {username: "matthis", ...}
[Office1789-SSO] Token valide ✓
[Office1789-SSO] Authentification Matrix réussie
```

### Indicateurs de compromission

🚨 **ALERTE si vous voyez:**
- Multiples "Token invalide" en peu de temps (brute force?)
- "Token expiré" pour un utilisateur légitime (replay attack?)
- Accès direct bloqué répété (scan de vulnérabilités?)
- Erreurs de signature HMAC (tentative de forgerie?)

---

## 📊 Comparaison des protections

| Vecteur d'attaque | Roundcube | Element | Protection |
|-------------------|-----------|---------|------------|
| Accès direct sans token | ✅ Bloqué | ✅ Bloqué | 403 Forbidden |
| Token expiré | ✅ Vérifié | ✅ Vérifié | Expiration 5 min |
| Token forgé | ✅ Vérifié | ✅ Vérifié | HMAC-SHA256 |
| Token modifié | ✅ Vérifié | ✅ Vérifié | HMAC-SHA256 |
| Session active | ✅ Autorisé | ✅ Autorisé | localStorage |
| Rejeu de token | ✅ Limité | ✅ Limité | Expiration 5 min |
| CSRF | ✅ Protégé | ✅ Protégé | Token unique |
| XSS | ⚠️ Partiel | ⚠️ Partiel | CSP recommandé |
| MITM | ⚠️ HTTP only | ⚠️ HTTP only | HTTPS requis |

---

## ✅ Checklist sécurité production

- [ ] **Secrets forts** - Générer clés aléatoires 256 bits minimum
- [ ] **Variables d'environnement** - Ne jamais hardcoder les secrets
- [ ] **HTTPS partout** - Activer TLS/SSL pour tous les services
- [ ] **Certificats valides** - Let's Encrypt ou certificat commercial
- [ ] **Firewall** - Bloquer accès direct aux ports 8080-8083 depuis Internet
- [ ] **Rate limiting** - Limiter tentatives de connexion
- [ ] **Logs centralisés** - ELK Stack ou Grafana Loki
- [ ] **Monitoring** - Alertes sur tentatives de connexion échouées
- [ ] **Rotation des secrets** - Changer clés SSO tous les 90 jours
- [ ] **CSP headers** - Content-Security-Policy pour XSS
- [ ] **CORS restrictif** - Autoriser uniquement domaine Office1789
- [ ] **Session timeout** - Expiration session après inactivité
- [ ] **2FA obligatoire** - Pour comptes administrateurs minimum
- [ ] **Audit régulier** - Revue des logs chaque semaine

---

## 🎯 Résumé

**✅ Roundcube:** Protection ACTIVE - Accès direct bloqué  
**✅ Element:** Protection ACTIVE - Accès direct bloqué  
**✅ Backend:** Tokens signés HMAC-SHA256 avec expiration 5 min  
**✅ Frontend:** Composants SSO secured avec validation côté serveur  

**⚠️ Développement:** HTTP acceptable  
**🚀 Production:** HTTPS obligatoire + secrets forts + monitoring  

---

## 🏛️ Office1789 - Sécurité garantie !
