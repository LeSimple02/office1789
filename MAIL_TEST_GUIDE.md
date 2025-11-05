# Configuration Mail Office1789 - Guide de Test

## État du Système

✅ **Mailserver** : Opérationnel (sans SSL pour développement local)
✅ **Roundcube** : Opérationnel (connecté à PostgreSQL et Mailserver)
✅ **Base de données** : Tables Roundcube créées
✅ **Comptes mail** : jean@office1789.local et matthis@office1789.local

## Architecture

```
Vue.js Frontend (localhost:5173)
    └─> MailView avec iframe
         └─> Roundcube (localhost:8081)
              ├─> PostgreSQL (sessions/contacts)
              └─> Mailserver IMAP/SMTP
                   ├─> IMAP: port 143 (connexion)
                   └─> SMTP: port 587 (envoi)
```

## Comptes de Test

### Utilisateur 1 : Jean
- **Email** : jean@office1789.local
- **Mot de passe** : password123
- **Serveur IMAP** : mailserver:143
- **Serveur SMTP** : mailserver:587

### Utilisateur 2 : Matthis
- **Email** : matthis@office1789.local
- **Mot de passe** : password123
- **Serveur IMAP** : mailserver:143
- **Serveur SMTP** : mailserver:587

## Test de Connexion

### Option 1 : Via l'application Office1789
1. Démarrez le frontend Vue.js : `cd webfront2; npm run dev`
2. Connectez-vous à Office1789
3. Cliquez sur la section Mail
4. L'iframe Roundcube devrait apparaître
5. Connectez-vous avec jean@office1789.local / password123

### Option 2 : Direct sur Roundcube
1. Ouvrez http://localhost:8081
2. Connectez-vous avec :
   - **Utilisateur** : jean@office1789.local
   - **Mot de passe** : password123
   - (Le serveur mailserver est automatiquement configuré)

## Vérifications

### Vérifier les services Docker
```powershell
docker ps
```

Vous devriez voir :
- `postgres_db` (port 5432)
- `mailserver` (ports 25, 143, 587, 993)
- `roundcube` (port 8081)

### Vérifier les logs du mailserver
```powershell
docker logs mailserver --tail 50
```

Recherchez : `mail.office1789.local is up and running`

### Vérifier les logs de Roundcube
```powershell
docker logs roundcube --tail 50
```

Recherchez des erreurs IMAP (il ne devrait plus y en avoir)

### Vérifier les tables Roundcube
```powershell
docker exec postgres_db psql -U robespierre -d office1789 -c "\dt roundcube*"
```

Devrait afficher 8 tables : roundcube_users, roundcube_session, etc.

## Résolution des Problèmes Précédents

### ❌ Problème : "Erreur de connexion au stockage"
**Cause** : Mailserver en redémarrage constant à cause des certificats SSL manquants
**Solution** : Désactivé SSL avec `SSL_TYPE=` et `DOVECOT_TLS=no` dans docker-compose.yml

### ❌ Problème : Tables Roundcube manquantes
**Cause** : Schema non initialisé
**Solution** : Créé init-roundcube.sql avec toutes les tables nécessaires

### ❌ Problème : Roundcube ne peut pas se connecter à IMAP
**Cause** : Configuration pointant vers localhost au lieu de mailserver
**Solution** : Mis à jour ROUNDCUBE_DEFAULT_HOST=mailserver dans Dockerfile

## Scripts Disponibles

### init-roundcube-db.ps1
Initialise les tables Roundcube dans PostgreSQL

### create-mail-accounts.ps1
Crée les comptes mail jean et matthis (déjà fait)

### setup-mail.ps1
Gestion complète des comptes mail (création, suppression, liste)

## Configuration Technique

### docker-compose.yml
- Mailserver sans SSL (développement local)
- Roundcube connecté à mailserver via network Docker
- PostgreSQL partagé entre Roundcube et Office1789

### backend/mailmanager.go
- Fonctions pour créer/supprimer des comptes mail via Docker exec
- Intégration avec subscribe.go pour création automatique

### webfront2/src/views/MailView.vue
- Iframe simple pointant vers localhost:8081
- Pas de SSO complexe nécessaire (connexion directe Roundcube)

## Test Complet

1. **Démarrer tous les services** :
   ```powershell
   cd docker
   docker-compose up -d
   ```

2. **Vérifier que tout fonctionne** :
   ```powershell
   docker ps
   # Tous les containers doivent être "Up"
   ```

3. **Tester Roundcube directement** :
   - Ouvrir http://localhost:8081
   - Se connecter avec jean@office1789.local / password123
   - Vous devriez voir la boîte de réception (vide au début)

4. **Tester via Office1789** :
   ```powershell
   cd webfront2
   npm run dev
   ```
   - Se connecter à Office1789
   - Aller dans la section Mail
   - L'iframe Roundcube devrait s'afficher

5. **Envoyer un mail de test** :
   - Avec jean@office1789.local, envoyer un mail à matthis@office1789.local
   - Se connecter avec matthis pour vérifier la réception

## Notes Importantes

- ⚠️ **SSL désactivé** : Configuration pour développement uniquement
- ⚠️ **Mot de passe simple** : password123 pour tous les comptes de test
- ⚠️ **Domaine local** : office1789.local ne fonctionne que dans Docker
- ✅ **SSO non nécessaire** : Connexion directe Roundcube pour simplicité

## Prochaines Étapes (Optionnel)

1. Implémenter le vrai SSO si besoin (backend/mail.go existe déjà)
2. Activer SSL pour production
3. Synchroniser les mots de passe Office1789 avec les mots de passe mail
4. Configurer des domaines réels pour envoi externe

## Support

Si vous rencontrez toujours des problèmes :

1. Vérifiez les logs : `docker logs mailserver` et `docker logs roundcube`
2. Redémarrez tous les services : `docker-compose restart`
3. Vérifiez la connectivité réseau Docker : `docker network ls`
