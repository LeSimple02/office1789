# 📧 Configuration du Serveur Mail - Office1789

## 🎯 Fonctionnement automatique

✨ **Bonne nouvelle !** Les comptes mail sont maintenant créés **automatiquement** lors de l'inscription d'un nouvel utilisateur !

- Quand un utilisateur s'inscrit avec le username `jean`, un compte mail `jean@office1789.local` est créé
- Le mot de passe du compte mail = mot de passe de connexion Office1789
- L'utilisateur peut immédiatement accéder à sa boîte mail via Roundcube

## 🚀 Installation rapide

### 1. Démarrer les services
```powershell
cd docker
docker-compose up -d
```

Attendez environ 30 secondes que le serveur mail démarre complètement.

### 2. Synchroniser les utilisateurs existants

Si vous avez déjà des utilisateurs dans votre base de données **AVANT** d'installer le serveur mail :

```powershell
.\sync-existing-users.ps1
```

Ce script va créer automatiquement les comptes mail pour tous les utilisateurs existants.

### 3. Créer des comptes mail manuellement (optionnel)

**Option A - Script interactif :**
```powershell
.\setup-mail.ps1
```
Choisissez l'option 4 pour créer les comptes par défaut.

**Option B - Commande manuelle:**
```powershell
# Créer un compte pour Jean
docker exec -it mailserver setup email add jean@office1789.local jean1789

# Créer un compte pour Matthis
docker exec -it mailserver setup email add matthis@office1789.local matthis1789
```

### 3. Se connecter à Roundcube

1. Ouvrez http://localhost:8081
2. Connectez-vous avec :
   - **Utilisateur:** `jean@office1789.local`
   - **Mot de passe:** `jean1789`

## 📋 Gestion des comptes

### Lister les comptes existants
```powershell
docker exec mailserver setup email list
```

### Ajouter un nouveau compte
```powershell
docker exec -it mailserver setup email add utilisateur@office1789.local motdepasse
```

### Supprimer un compte
```powershell
docker exec -it mailserver setup email del utilisateur@office1789.local
```

### Changer le mot de passe
```powershell
docker exec -it mailserver setup email update utilisateur@office1789.local nouveau_motdepasse
```

## 🔧 Configuration des services

### Serveur Mail (mailserver)
- **SMTP:** `mailserver:587` (dans Docker) ou `localhost:587` (depuis l'hôte)
- **IMAP:** `mailserver:143` (dans Docker) ou `localhost:143` (depuis l'hôte)
- **IMAPS:** `mailserver:993` (dans Docker) ou `localhost:993` (depuis l'hôte)
- **Domaine:** `office1789.local`

### Roundcube (webmail)
- **URL:** http://localhost:8081
- **Base de données:** PostgreSQL (office1789)
- **Configuration:** Automatique via variables d'environnement

## 🐛 Dépannage

### Le serveur mail ne démarre pas
```powershell
# Vérifier les logs
docker logs mailserver

# Redémarrer le container
docker restart mailserver
```

### Impossible de se connecter
1. Vérifiez que le compte existe :
   ```powershell
   docker exec mailserver setup email list
   ```

2. Vérifiez que le serveur mail est démarré :
   ```powershell
   docker ps | findstr mailserver
   ```

3. Testez la connexion IMAP :
   ```powershell
   telnet localhost 143
   ```

### Roundcube ne se connecte pas au serveur mail
1. Vérifiez que les containers sont sur le même réseau :
   ```powershell
   docker network inspect docker_default
   ```

2. Testez la connexion depuis Roundcube :
   ```powershell
   docker exec -it roundcube ping mailserver
   ```

### Réinitialiser complètement le serveur mail
```powershell
# Arrêter et supprimer les containers
docker-compose down

# Supprimer les volumes (⚠️ Cela supprimera tous les emails !)
docker volume rm docker_maildata docker_mailstate docker_maillogs

# Redémarrer
docker-compose up -d

# Recréer les comptes
.\setup-mail.ps1
```

## 📚 Ports utilisés

| Service | Port | Description |
|---------|------|-------------|
| PostgreSQL | 5432 | Base de données |
| Roundcube | 8081 | Interface webmail |
| OnlyOffice | 8082 | Éditeur de documents |
| SMTP | 25 | Envoi d'emails |
| SMTP Submission | 587 | Envoi sécurisé |
| IMAP | 143 | Réception d'emails |
| IMAPS | 993 | Réception sécurisée |
| Backend Go | 8080 | API Office1789 |

## 🔐 Sécurité

- Les mots de passe sont stockés chiffrés dans le serveur mail
- SSL auto-signé activé pour IMAPS/SMTPS
- En production, utilisez des certificats SSL valides (Let's Encrypt)

## 📖 Documentation

- [Docker Mailserver](https://docker-mailserver.github.io/docker-mailserver/latest/)
- [Roundcube](https://github.com/roundcube/roundcubemail/wiki)
- [Office1789](../README.md)
