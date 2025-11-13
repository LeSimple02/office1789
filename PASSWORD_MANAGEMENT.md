# Gestion des Mots de Passe Office1789

## 🔐 Système Unifié

Office1789 utilise un **système de mot de passe unifié** pour tous les services :
- **Office1789** (Authentification principale)
- **Mail** (Roundcube + Docker Mailserver)
- **Matrix** (Synapse / Element)

## 🚫 Restrictions

### Changement de mot de passe **BLOQUÉ** dans :
- ❌ **Roundcube** : Option désactivée via `$config['dont_override'] = ['password']`
- ❌ **Clients Matrix** (Element, etc.) : Désactivé via `password_config: enabled: false`

### Changement de mot de passe **AUTORISÉ** uniquement via :
- ✅ **Interface Office1789** : Via l'endpoint `/api/account/change-password`
- ✅ **Script CLI** : `change-password.ps1` (pour admins)

## 📝 Utilisation

### Pour les utilisateurs

1. Se connecter à Office1789
2. Aller dans **Paramètres du compte**
3. Utiliser le formulaire de changement de mot de passe
4. Le nouveau mot de passe sera automatiquement synchronisé sur :
   - Office1789 (PostgreSQL)
   - Mail (Docker Mailserver)
   - Matrix (Synapse)

### Pour les administrateurs (CLI)

```powershell
.\change-password.ps1
```

Le script vous demandera :
- Username
- Nouveau mot de passe
- Confirmation

## 🔧 Détails Techniques

### Endpoint API

**POST** `/api/account/change-password`

**Body :**
```json
{
  "username": "jean",
  "token": "session-token",
  "oldPassword": "ancien-mot-de-passe",
  "newPassword": "nouveau-mot-de-passe"
}
```

**Response :**
```json
{
  "success": true,
  "message": "Password changed successfully across all services"
}
```

### Processus de changement

1. **Vérification** de la session et de l'ancien mot de passe
2. **Mise à jour Office1789** : Hash bcrypt dans PostgreSQL
3. **Mise à jour Mail** (asynchrone) : 
   - `docker exec mailserver setup email del user@office1789.local`
   - `docker exec mailserver setup email add user@office1789.local newpass`
4. **Mise à jour Matrix** (asynchrone) :
   - `docker exec synapse reset-password -c /data/homeserver.yaml username newpass`

### Sécurité

- ✅ Mot de passe **toujours hashé** (bcrypt) dans PostgreSQL
- ✅ Mot de passe en clair **uniquement** pendant la requête HTTP
- ✅ Pas de stockage du mot de passe en clair
- ✅ Validation de l'ancien mot de passe avant changement
- ✅ Session token requis pour authentification

## 📁 Fichiers concernés

### Backend
- `backend/account.go` : Fonction `ChangePassword()`
- `backend/main.go` : Route `/api/account/change-password`
- `backend/subscribe.go` : Fonctions helper (`changeMailPassword`, `changeMatrixPassword`)

### Configuration
- `docker/DockerfileRC/config.inc.php` : Blocage Roundcube
- `docker/synapse/conf/homeserver.yaml` : Blocage Matrix

### Scripts
- `change-password-tool.go` : CLI standalone
- `change-password.ps1` : Wrapper PowerShell

## 🎯 Avantages

1. **Simplicité** : Un seul mot de passe pour tous les services
2. **Sécurité** : Point de changement centralisé et auditable
3. **Synchronisation automatique** : Pas de désynchronisation possible
4. **Administration facilitée** : Script CLI pour support technique

## ⚠️ Notes importantes

- Le changement de mot de passe dans `account.go` via `ChangeI()` est **désactivé**
- Seul l'endpoint dédié `/api/account/change-password` fonctionne
- Les erreurs Mail/Matrix sont **loggées** mais ne bloquent pas la mise à jour Office1789
- Les comptes sont créés automatiquement lors de l'inscription (voir `subscribe.go`)
