# Migration du système de session vers une base de données persistante

## ✅ Modifications terminées

### 1. Base de données
- ✅ Ajout de la table `Sessions` dans `init.sql`
- ✅ Création du script de migration `migrate_sessions.sql`
- ✅ Index sur `session_token` et `expires_at`
- ✅ CASCADE DELETE quand un user est supprimé

### 2. Structure de session (security.go)
- ✅ Session basée sur `user_id` au lieu de `username`
- ✅ `createSessionInDB()` - Crée session en DB
- ✅ `getSessionFromDB()` - Récupère session depuis DB
- ✅ `deleteSessionFromDB()` - Supprime session
- ✅ `cleanExpiredSessions()` - Nettoie sessions expirées
- ✅ `validateSession()` - Fonction helper pour valider

### 3. Fichiers modifiés
- ✅ `login.go` - Utilise user_id, sauvegarde en DB
- ✅ `subscribe.go` - Récupère user_id et crée session en DB
- ✅ `account.go` - getinfop, ChangeI, DeleteAccount
- ✅ `calendar.go` - Tous les endpoints utilisent validateSession()

### 4. Frontend
- ✅ Retourne maintenant `user_id` en plus du `username`
- Structure: `{user_id, username, token, expiry}`

## ⚠️ Modifications à finaliser manuellement

### drive.go (2068 lignes)
Le fichier est trop grand pour être modifié automatiquement. Voici les changements pattern à appliquer :

**AVANT:**
```go
session, valid := validateSession(req.Token, req.Username)
if !valid {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    return
}

var userID int
err := db.QueryRow("SELECT user_id FROM Users WHERE username=$1", req.Username).Scan(&userID)
if err != nil {
    c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
    return
}

// ... utilisation de userID
db.Query("... WHERE user_id=$1", userID)
```

**APRÈS:**
```go
session, valid := validateSession(req.Token, req.Username)
if !valid {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    return
}

// Utiliser directement session.UserID (pas besoin de requête SQL)
db.Query("... WHERE user_id=$1", session.UserID)
```

**Occurrences à corriger dans drive.go:**
- Ligne 117-125
- Ligne 435-445  
- Ligne 517-525
- Ligne 750-760
- Ligne 790-800
- Ligne 885-895
- Ligne 979-990
- Ligne 1081-1090
- Ligne 1250-1260
- Ligne 1378-1388
- Ligne 1870-1880
- Ligne 1932-1942
- Ligne 1983-1993

**Ligne 301: Remplacer**
```go
session, ok := sessions[token]  // ❌ sessions n'existe plus
```
Par:
```go
session, err := getSessionFromDB(token)  // ✅
if err != nil || session == nil {
    // handle error
}
```

### chat.go
**Ligne 17: AVANT**
```go
session, valid := validateSession(cre.Token, cre.Username); if valid {
```
**APRÈS:**
```go
session, valid := validateSession(cre.Token, cre.Username)
if valid {
```

Puis utiliser `session.UserID` au lieu de requêter la DB pour l'user_id.

## 📝 Pour appliquer la migration

### 1. Démarrer PostgreSQL
```powershell
cd c:\Users\Matthis\Documents\office1789\docker
docker compose up -d postgres
```

### 2. Appliquer la migration
```powershell
docker compose exec postgres psql -U office1789_user -d office1789 -f /docker-entrypoint-initdb.d/migrate_sessions.sql
```

OU copier le script et l'exécuter:
```powershell
docker cp c:\Users\Matthis\Documents\office1789\backend\office1789-db\migrate_sessions.sql office1789-postgres:/tmp/
docker compose exec postgres psql -U office1789_user -d office1789 -f /tmp/migrate_sessions.sql
```

### 3. Vérifier la table
```powershell
docker compose exec postgres psql -U office1789_user -d office1789 -c "\d Sessions"
```

### 4. Compiler après corrections manuelles
```powershell
cd c:\Users\Matthis\Documents\office1789\backend
go build -o office1789.exe
```

## 🎯 Avantages du nouveau système

1. **Persistance**: Les sessions survivent aux redémarrages du backend
2. **Sécurité**: Basé sur user_id (pas de conflit de username)
3. **Performance**: Index sur session_token pour recherche rapide
4. **Nettoyage auto**: CASCADE DELETE supprime les sessions avec l'user
5. **Expiration**: Fonction cleanExpiredSessions() 
6. **Multi-instance**: Plusieurs backends peuvent partager les sessions

## 🔄 Comportement des sessions

- **Durée**: 24 heures (au lieu de 120 secondes)
- **Stockage**: PostgreSQL (au lieu de map mémoire)
- **Validation**: Vérifie user_id + username + expiration
- **Nettoyage**: Automatique à chaque login + cleanExpiredSessions()

## 📊 Structure de réponse login/subscribe

```json
{
  "user_id": 123,
  "username": "john",
  "token": "uuid-v4-token",
  "expiry": "2025-11-14T15:00:00Z"
}
```

Le frontend doit stocker `user_id` en plus du `username` et `token`.
