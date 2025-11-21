# Domaines Personnalisés - Guide de Configuration

## 🌐 Vue d'ensemble

Office1789 permet aux utilisateurs **Professional** et **Enterprise** d'utiliser leur propre nom de domaine pour leurs adresses email. Par exemple, au lieu de `contact@office1789.com`, vous pouvez utiliser `contact@votreentreprise.com`.

## 📋 Prérequis

- Plan **Professional** (3 comptes) ou **Enterprise** (20 comptes)
- Accès à la configuration DNS de votre domaine
- Domaine valide (ex: `votreentreprise.com`)

## 🔧 Configuration côté serveur

### 1. Installer le support des domaines personnalisés

**Sur Windows (PowerShell) :**
```powershell
cd docker
.\configure-custom-domains.ps1
```

**Sur Linux :**
```bash
cd docker
chmod +x configure-custom-domains.sh
./configure-custom-domains.sh
```

Ce script :
- Configure Postfix pour accepter les domaines virtuels
- Met en place les mappings d'adresses email
- Recharge automatiquement la configuration

### 2. Vérification

Vérifier que la configuration est appliquée :
```bash
docker exec mailserver cat /tmp/docker-mailserver/postfix-virtual-domains.cf
docker exec mailserver cat /tmp/docker-mailserver/postfix-virtual-mailboxes.cf
```

## 👤 Utilisation par l'utilisateur

### Étape 1 : Ajouter un domaine

1. Connectez-vous à Office1789
2. Allez dans **Compte** → **Domaine personnalisé**
3. Entrez votre nom de domaine (ex: `monentreprise.com`)
4. Cliquez sur **Ajouter le domaine**

Un token de vérification sera généré automatiquement.

### Étape 2 : Vérifier la propriété du domaine

#### Configuration DNS requise

Ajoutez les enregistrements DNS suivants chez votre registrar (OVH, Gandi, Cloudflare, etc.) :

| Type | Nom/Host | Valeur | Priorité |
|------|----------|--------|----------|
| **TXT** | `@` ou domaine principal | `office1789-verification=VOTRE_TOKEN` | - |
| **MX** | `@` | `mail.office1789.com` | 10 |
| **TXT** | `@` | `v=spf1 include:office1789.com ~all` | - |
| **TXT** | `_dmarc` | `v=DMARC1; p=quarantine; rua=mailto:postmaster@office1789.com` | - |

**Important :** Remplacez `VOTRE_TOKEN` par le token affiché dans l'interface Office1789.

#### Exemples de configuration

**OVH :**
- Type : `TXT`
- Sous-domaine : (vide)
- Valeur : `office1789-verification=abc123...`

**Cloudflare :**
- Type : `TXT`
- Name : `@`
- Content : `office1789-verification=abc123...`

**Gandi :**
- Type : `TXT`
- Nom : `@`
- Valeur : `office1789-verification=abc123...`

### Étape 3 : Vérifier le domaine

1. Attendez la propagation DNS (5 minutes à 48 heures)
2. Cliquez sur **Vérifier le domaine** dans Office1789
3. Si la vérification réussit, votre domaine est activé ! ✅

### Étape 4 : Utiliser vos nouvelles adresses

Une fois vérifié, vos emails fonctionnent automatiquement :

**Compte individuel :**
- `votreusername@votreentreprise.com` → redirige vers `votreusername@office1789.com`

**Organisation :**
- Tous les membres de l'organisation peuvent utiliser le domaine
- `membre1@votreentreprise.com`
- `membre2@votreentreprise.com`
- etc.

## 🔄 Synchronisation automatique

Le backend Office1789 synchronise automatiquement les domaines avec le serveur mail :

- **Au démarrage** : sync initial après 5 secondes
- **Périodique** : sync toutes les 5 minutes
- **Lors des changements** :
  - Ajout de domaine
  - Vérification de domaine
  - Suppression de domaine

### Synchronisation manuelle

Si besoin, déclencher manuellement :

```bash
# Recharger la configuration Postfix
docker exec mailserver postfix reload

# Vérifier les logs
docker logs mailserver | tail -n 50
```

## 🏢 Fonctionnalités par plan

| Fonctionnalité | Free | Standard | Professional | Enterprise |
|----------------|------|----------|--------------|------------|
| Domaine personnalisé | ❌ | ❌ | ✅ | ✅ |
| Sous-comptes | ❌ | ❌ | 3 | 20 |
| Support prioritaire | ❌ | ⭐ | ⭐ | ⭐⭐ |
| Stockage | 1GB | 50GB | 200GB | Illimité |

## 🔍 Dépannage

### Le domaine ne se vérifie pas

**Causes possibles :**
1. **Propagation DNS incomplète** : Attendez jusqu'à 48h
2. **Enregistrement TXT incorrect** : Vérifiez que le token est exact
3. **Problème de registrar** : Certains registrars ont des délais plus longs

**Vérifier manuellement :**
```bash
# Windows PowerShell
Resolve-DnsName -Type TXT votreentreprise.com

# Linux/Mac
dig TXT votreentreprise.com
nslookup -type=TXT votreentreprise.com
```

Vous devez voir : `office1789-verification=VOTRE_TOKEN`

### Les emails n'arrivent pas

**Vérifications :**
1. Domaine vérifié (badge vert ✅ dans l'interface)
2. Enregistrement MX configuré : `mail.office1789.com` priorité 10
3. Postfix rechargé : `docker exec mailserver postfix reload`
4. Logs du serveur : `docker logs mailserver`

**Tester la réception :**
```bash
# Vérifier les mappings
docker exec mailserver postmap -q "contact@votreentreprise.com" /tmp/docker-mailserver/postfix-virtual-mailboxes.cf
```

### Le domaine est déjà utilisé

**Erreur :** "This domain is already in use by another account"

**Solution :** Ce domaine est déjà vérifié par un autre utilisateur. Chaque domaine ne peut être utilisé que par un seul compte Office1789.

## 📊 Architecture technique

```
┌─────────────────┐
│  Utilisateur    │
│  ajoute domaine │
└────────┬────────┘
         │
         ▼
┌─────────────────────────┐
│  Office1789 Backend     │
│  - Génère token         │
│  - Stocke en BDD        │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  Utilisateur configure  │
│  DNS TXT, MX, SPF       │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  Backend vérifie DNS    │
│  - net.LookupTXT()      │
│  - Compare token        │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  SyncMailServerConfig() │
│  - Génère .cf files     │
│  - Postfix reload       │
└────────┬────────────────┘
         │
         ▼
┌─────────────────────────┐
│  docker-mailserver      │
│  - Virtual domains      │
│  - Virtual mailboxes    │
│  - Emails fonctionnent  │
└─────────────────────────┘
```

## 📁 Fichiers de configuration

### Backend (Go)

- **`backend/customdomain.go`** : API domaines personnalisés
  - `POST /api/domain/add` : Ajouter un domaine
  - `POST /api/domain/verify` : Vérifier via DNS
  - `POST /api/domain/info` : Info domaine actuel
  - `POST /api/domain/remove` : Retirer le domaine

- **`backend/mailsync.go`** : Synchronisation avec Postfix
  - `syncCustomDomains()` : Génère postfix-virtual-domains.cf
  - `syncVirtualMailboxes()` : Génère postfix-virtual-mailboxes.cf
  - `SyncMailServerConfig()` : Sync complet
  - `StartMailSyncScheduler()` : Scheduler automatique (5 min)

### Base de données

```sql
-- Table Organizations
custom_domain VARCHAR(255)
domain_verified BOOLEAN
domain_verification_token VARCHAR(255)

-- Table Users
custom_domain VARCHAR(255)
domain_verified BOOLEAN
```

### Mail Server

- **`docker/config/postfix-virtual-domains.cf`** : Liste des domaines acceptés
- **`docker/config/postfix-virtual-mailboxes.cf`** : Mappings email → compte

## 🚀 Déploiement en production

### Checklist

- [ ] Exécuter `configure-custom-domains.ps1` (Windows) ou `.sh` (Linux)
- [ ] Vérifier que Postfix accepte les domaines virtuels
- [ ] Tester avec un domaine test
- [ ] Configurer le DNS du domaine principal si nécessaire
- [ ] Surveiller les logs : `docker logs -f mailserver`
- [ ] Documenter pour les utilisateurs finaux

### Monitoring

```bash
# Vérifier les domaines actifs
docker exec mailserver cat /tmp/docker-mailserver/postfix-virtual-domains.cf

# Voir les mappings
docker exec mailserver cat /tmp/docker-mailserver/postfix-virtual-mailboxes.cf

# Logs Postfix en temps réel
docker logs -f mailserver | grep postfix

# Tester un mapping
docker exec mailserver postmap -q "test@customdomain.com" /tmp/docker-mailserver/postfix-virtual-mailboxes.cf
```

## 📞 Support

Pour toute question :
- **Professional** : Support prioritaire ⭐
- **Enterprise** : Support prioritaire 24/7 ⭐⭐
- **Email** : support@office1789.com
- **Documentation** : https://docs.office1789.com

---

**Dernière mise à jour :** Novembre 2025  
**Version :** 1.0  
**Statut :** Production-ready ✅
