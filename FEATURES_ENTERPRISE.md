# 🚀 Office1789 - Nouvelles Fonctionnalités Enterprise

## 📅 Date de mise à jour : Novembre 2025

---

## ✨ Fonctionnalités implémentées

### 1. 🏢 Système d'organisation multi-comptes

Permet aux utilisateurs Professional et Enterprise de créer et gérer plusieurs sous-comptes au sein d'une organisation.

#### Caractéristiques
- **Professional** : jusqu'à 3 sous-comptes
- **Enterprise** : jusqu'à 20 sous-comptes
- Gestion centralisée par le compte principal
- Création/suppression de sous-comptes
- Chaque sous-compte a son propre login
- Partage des ressources de l'organisation

#### Architecture
- Table `Organizations` : stockage des organisations
- Liens hiérarchiques via `parent_account_id` et `organization_id`
- Types de compte : `personal`, `organization_owner`, `organization_member`

#### APIs
- `POST /api/organization/create-subaccount` - Créer un sous-compte
- `POST /api/organization/members` - Lister les membres
- `POST /api/organization/delete-member` - Supprimer un membre

#### Interface utilisateur
- `/account/organization` - Panel de gestion
- Formulaire de création de sous-comptes
- Liste des membres avec actions
- Indicateurs de limites (X/Y comptes utilisés)

---

### 2. 🌐 Domaines personnalisés

Permet aux utilisateurs Pro/Enterprise d'utiliser leur propre nom de domaine pour les emails.

#### Caractéristiques
- Domaine personnalisé pour emails (ex: `contact@votreentreprise.com`)
- Vérification de propriété via DNS TXT
- Configuration DNS automatisée
- Support SPF, DMARC, MX
- Synchronisation automatique avec Postfix

#### Workflow
1. Utilisateur ajoute un domaine
2. Token de vérification généré
3. Utilisateur configure DNS (TXT, MX, SPF)
4. Vérification via lookup DNS
5. Activation automatique du domaine
6. Sync avec serveur mail

#### APIs
- `POST /api/domain/add` - Ajouter un domaine
- `POST /api/domain/verify` - Vérifier via DNS
- `POST /api/domain/info` - Info domaine actuel
- `POST /api/domain/remove` - Retirer le domaine

#### Synchronisation mail
- `mailsync.go` : Scheduler automatique (5 min)
- Génération de `postfix-virtual-domains.cf`
- Génération de `postfix-virtual-mailboxes.cf`
- Reload Postfix automatique

#### Interface utilisateur
- `/account/custom-domain` - Panel de configuration
- Instructions DNS détaillées
- Copier-coller du token
- Vérification en un clic
- Indicateurs de statut

---

## 🗄️ Modifications de la base de données

### Table `Organizations` (nouvelle)
```sql
CREATE TABLE Organizations (
    organization_id SERIAL PRIMARY KEY,
    organization_name VARCHAR(255) NOT NULL,
    owner_user_id INT,
    max_members INT DEFAULT 1,
    custom_domain VARCHAR(255),
    domain_verified BOOLEAN DEFAULT FALSE,
    domain_verification_token VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Table `Users` (modifications)
```sql
ALTER TABLE Users ADD COLUMN organization_id INT;
ALTER TABLE Users ADD COLUMN parent_account_id INT;
ALTER TABLE Users ADD COLUMN account_type VARCHAR(20) DEFAULT 'personal';
ALTER TABLE Users ADD COLUMN custom_domain VARCHAR(255);
ALTER TABLE Users ADD COLUMN domain_verified BOOLEAN DEFAULT FALSE;

-- Foreign keys
ALTER TABLE Users ADD FOREIGN KEY (organization_id) 
    REFERENCES Organizations(organization_id) ON DELETE SET NULL;
ALTER TABLE Users ADD FOREIGN KEY (parent_account_id) 
    REFERENCES Users(user_id) ON DELETE CASCADE;
```

---

## 📁 Nouveaux fichiers

### Backend (Go)
- `backend/organization.go` - Gestion des organisations et sous-comptes
- `backend/customdomain.go` - Gestion des domaines personnalisés
- `backend/mailsync.go` - Synchronisation avec serveur mail

### Frontend (Vue.js)
- `webfront2/src/components/OrganizationPanel.vue` - Interface organisations
- `webfront2/src/components/CustomDomainPanel.vue` - Interface domaines

### Configuration mail
- `docker/config/postfix-virtual-domains.cf` - Liste des domaines acceptés
- `docker/config/postfix-virtual-mailboxes.cf` - Mappings email

### Scripts
- `docker/configure-custom-domains.ps1` - Configuration Windows
- `docker/configure-custom-domains.sh` - Configuration Linux

### Documentation
- `CUSTOM_DOMAINS.md` - Guide complet des domaines personnalisés
- `CUSTOM_DOMAINS_TESTING.md` - Tests et validation

---

## 🎯 Limites par plan

| Fonctionnalité | Free | Standard | Professional | Enterprise |
|----------------|------|----------|--------------|------------|
| **Stockage** | 1GB | 50GB | 200GB | Illimité |
| **Partage fichiers** | ❌ | ❌ | ✅ (3 membres) | ✅ (20 membres) |
| **Sous-comptes** | ❌ | ❌ | 3 | 20 |
| **Domaine personnalisé** | ❌ | ❌ | ✅ | ✅ |
| **Support prioritaire** | ❌ | ⭐ | ⭐ | ⭐⭐ (24/7) |
| **API Access** | ❌ | ❌ | ✅ | ✅ |
| **SLA** | - | - | - | 99.9% |

---

## 🔧 Configuration requise

### Serveur mail
1. Exécuter le script de configuration :
   ```powershell
   cd docker
   .\configure-custom-domains.ps1
   ```

2. Vérifier la configuration :
   ```bash
   docker exec mailserver cat /tmp/docker-mailserver/postfix-virtual-domains.cf
   docker exec mailserver postfix reload
   ```

### Backend
1. Compiler avec les nouvelles dépendances :
   ```bash
   cd backend
   go build
   ```

2. Le scheduler mail sync démarre automatiquement au lancement

### Base de données
Les migrations sont déjà appliquées. Vérifier avec :
```sql
SELECT column_name FROM information_schema.columns 
WHERE table_name = 'organizations';
```

---

## 🚀 Déploiement

### Checklist pré-déploiement
- [x] Base de données migrée (colonnes ajoutées)
- [x] Backend compilé avec nouvelles routes
- [x] docker-mailserver configuré pour domaines virtuels
- [x] Interface utilisateur mise à jour
- [x] Tests effectués (cf. CUSTOM_DOMAINS_TESTING.md)

### Commandes de déploiement

```bash
# 1. Backend
cd backend
go build
./main.exe  # ou ./main sur Linux

# 2. Mail server
cd docker
.\configure-custom-domains.ps1  # Windows
# ou
./configure-custom-domains.sh   # Linux

# 3. Frontend (si nécessaire)
cd webfront2
npm run build
```

---

## 📊 Monitoring

### Métriques à surveiller

```sql
-- Nombre d'organisations actives
SELECT COUNT(*) FROM Organizations;

-- Nombre de sous-comptes par organisation
SELECT o.organization_name, COUNT(u.user_id) as members
FROM Organizations o
LEFT JOIN Users u ON u.organization_id = o.organization_id
GROUP BY o.organization_id, o.organization_name;

-- Domaines personnalisés vérifiés
SELECT COUNT(*) FROM Organizations WHERE domain_verified = TRUE
UNION ALL
SELECT COUNT(*) FROM Users WHERE domain_verified = TRUE;
```

### Logs importants

```bash
# Backend logs
# Rechercher : "✅ Sub-account created", "✅ Domain verified"

# Mail server logs
docker logs -f mailserver | grep postfix

# Sync scheduler
# Rechercher : "🔄 Synchronizing custom domains"
```

---

## 🔒 Sécurité

### Mesures implémentées
- ✅ Validation stricte des formats de domaine
- ✅ Vérification DNS obligatoire (TXT record)
- ✅ Un domaine = un seul compte (unicité en BDD)
- ✅ SPF configuré : `v=spf1 include:office1789.com ~all`
- ✅ DMARC configuré : `v=DMARC1; p=quarantine`
- ✅ Isolation des organisations (CASCADE delete)
- ✅ Tokens de vérification aléatoires (32 bytes hex)

### Recommandations
- Surveiller les tentatives de vérification DNS échouées
- Limiter le nombre de vérifications par domaine (rate limiting)
- Logger tous les ajouts/suppressions de domaines
- Backup régulier des fichiers Postfix .cf

---

## 📞 Support utilisateur

### Questions fréquentes

**Q: Comment ajouter un sous-compte ?**  
R: Compte → Gestion de l'organisation → Créer un nouveau sous-compte

**Q: Pourquoi mon domaine ne se vérifie pas ?**  
R: Attendez jusqu'à 48h pour la propagation DNS. Vérifiez avec `nslookup -type=TXT votredomaine.com`

**Q: Puis-je utiliser un sous-domaine ?**  
R: Oui, `mail.votreentreprise.com` est valide

**Q: Comment passer de Pro à Enterprise ?**  
R: Compte → Changer d'offre → Enterprise (20 comptes au lieu de 3)

---

## 🎉 Prochaines étapes

### Améliorations possibles
- [ ] Gestion des quotas de stockage par sous-compte
- [ ] Rapport d'utilisation pour organisations
- [ ] API pour gestion programmatique
- [ ] Webhooks pour événements (nouveau membre, domaine vérifié)
- [ ] Interface d'administration pour super-admin
- [ ] Statistiques d'utilisation des domaines personnalisés
- [ ] Support DKIM (signatures email)
- [ ] Auto-configuration email clients (Autodiscover)

### Optimisations
- [ ] Cache des domaines vérifiés (Redis)
- [ ] Rate limiting sur vérifications DNS
- [ ] Batch processing pour sync mail (si > 100 domaines)
- [ ] Monitoring Prometheus/Grafana

---

## 📚 Documentation complète

- **`CUSTOM_DOMAINS.md`** - Guide utilisateur détaillé
- **`CUSTOM_DOMAINS_TESTING.md`** - Tests et scénarios de validation
- **`backend/organization.go`** - Code source avec commentaires
- **`backend/customdomain.go`** - Code source avec commentaires
- **`backend/mailsync.go`** - Code source avec commentaires

---

**Version :** 1.0  
**Statut :** ✅ Production Ready  
**Testé :** ✅ Compilé sans erreurs  
**Documenté :** ✅ Complet  

🎊 **Félicitations ! Office1789 est maintenant prêt pour les entreprises !** 🎊
