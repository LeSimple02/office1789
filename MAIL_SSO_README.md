# Intégration Mail avec SSO - Office1789

## Vue d'ensemble

Ce système permet aux utilisateurs d'Office1789 d'accéder à leur boîte mail Roundcube sans avoir à se reconnecter, grâce à un système de Single Sign-On (SSO).

## Architecture

```
┌─────────────┐     1. Demande SSO      ┌──────────────┐
│             │ ───────────────────────> │              │
│  Frontend   │                          │   Backend    │
│  (Vue.js)   │ <─────────────────────── │   (Go/Gin)   │
│             │  2. URL + Token SSO      │              │
└─────────────┘                          └──────────────┘
      │                                         │
      │ 3. Ouvre URL avec token                │
      │                                         │
      ▼                                         │
┌─────────────┐     4. Valide token     ┌──────▼───────┐
│             │ ───────────────────────> │              │
│  Roundcube  │                          │   Backend    │
│  (PHP)      │ <─────────────────────── │   (Go/Gin)   │
│             │  5. Token valide         │              │
└─────────────┘                          └──────────────┘
      │
      │ 6. Connexion automatique
      ▼
   [Session créée]
```

## Flux d'authentification

### 1. Vérification de session
- Le composant Vue vérifie régulièrement si l'utilisateur est toujours connecté
- Endpoint: `POST /api/session/check`
- Paramètres: `{username, token}`

### 2. Demande d'accès mail
- L'utilisateur clique sur "Ouvrir ma boîte mail"
- Appel à `POST /api/mail/auth` avec les credentials de session
- Le backend génère un token SSO temporaire (valide 5 minutes)

### 3. Redirection vers Roundcube
- L'URL générée contient: `?_autologin=1&_user=USERNAME&_token=SSO_TOKEN`
- Ouvre dans un nouvel onglet

### 4. Validation du token SSO
- Roundcube appelle `GET /api/mail/validate-sso?token=...&user=...`
- Le backend vérifie la validité du token
- Si valide, le token est supprimé (usage unique)

### 5. Création de session Roundcube
- Si le token est valide, Roundcube crée une session automatiquement
- L'utilisateur est connecté sans saisir de mot de passe

## Fichiers modifiés/créés

### Backend (Go)
- **backend/mail.go** (nouveau) : Gestion SSO et tokens
- **backend/security.go** : Ajout des tokens SSO
- **backend/main.go** : Nouvelles routes mail

### Frontend (Vue.js)
- **webfront2/src/components/MailAccess.vue** (nouveau) : Composant d'accès mail

### Docker
- **docker/docker-compose.yml** : Configuration Roundcube mise à jour
- **docker/DockerfileRC/office1789_sso.php** (nouveau) : Plugin Roundcube

## Configuration requise

### 1. Variables d'environnement Backend

Dans `backend/mail.go`, changez la clé secrète :
```go
secret := "VOTRE_CLE_SECRETE_FORTE"
```

### 2. Installation du plugin Roundcube

Copiez `office1789_sso.php` dans le container Roundcube :
```dockerfile
# Dans DockerfileRC/Dockerfile
COPY office1789_sso.php /var/www/html/plugins/office1789_sso/office1789_sso.php
```

Activez le plugin dans la config Roundcube :
```php
$config['plugins'] = array('archive', 'zipdownload', 'office1789_sso');
```

### 3. Configuration serveur mail

Vous devez configurer un serveur mail (SMTP/IMAP) :

**Option A : Serveur mail local (Docker)**
```yaml
mailserver:
  image: mailserver/docker-mailserver:latest
  ports:
    - "25:25"    # SMTP
    - "587:587"  # SMTP Submission
    - "143:143"  # IMAP
    - "993:993"  # IMAPS
  environment:
    - ENABLE_SPAMASSASSIN=0
    - ENABLE_CLAMAV=0
```

**Option B : Serveur externe**
Modifiez les variables d'environnement Roundcube :
```yaml
ROUNDCUBE_DEFAULT_HOST: smtp.votre-serveur.com
ROUNDCUBE_SMTP_SERVER: smtp.votre-serveur.com
```

## Utilisation dans le frontend

### Intégrer le composant

Dans `MailView.vue` ou `HomeView.vue` :
```vue
<template>
  <div>
    <h1>Ma boîte mail</h1>
    <MailAccess />
  </div>
</template>

<script>
import MailAccess from '@/components/MailAccess.vue'

export default {
  components: {
    MailAccess
  }
}
</script>
```

## Sécurité

### Tokens SSO
- ✅ Validité limitée (5 minutes)
- ✅ Usage unique (supprimé après validation)
- ✅ Liés à un utilisateur spécifique
- ✅ Hash MD5 avec secret et timestamp

### Sessions
- ✅ Vérification régulière de la validité
- ✅ Expiration après 24h
- ✅ Token unique par session

### Recommandations
1. Utilisez HTTPS en production
2. Changez la clé secrète dans `mail.go`
3. Implémentez un rate limiting sur les endpoints SSO
4. Loggez les tentatives d'accès suspectes

## Dépannage

### L'utilisateur n'arrive pas à se connecter
1. Vérifiez que la session backend est active : `POST /api/session/check`
2. Vérifiez les logs du backend pour le token SSO
3. Vérifiez que Roundcube peut accéder à `host.docker.internal:8080`

### Le token SSO est invalide
- Les tokens expirent après 5 minutes
- Les tokens sont à usage unique
- Vérifiez l'horloge système (timestamps)

### Roundcube ne se connecte pas au mail
- Vérifiez la configuration SMTP/IMAP
- Testez manuellement la connexion au serveur mail
- Vérifiez les logs Roundcube : `/var/www/html/logs/`

## Prochaines étapes

1. **Configurer un serveur mail** (Postfix + Dovecot ou Docker Mailserver)
2. **Créer les comptes mail** pour chaque utilisateur Office1789
3. **Synchroniser les comptes** : créer automatiquement un compte mail à l'inscription
4. **Tester le flux complet** de connexion

## Support

Pour toute question, consultez :
- Documentation Roundcube : https://github.com/roundcube/roundcubemail
- Docker Mailserver : https://docker-mailserver.github.io/docker-mailserver/latest/
