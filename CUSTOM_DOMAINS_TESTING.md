# Test de la fonctionnalité des domaines personnalisés

## 🧪 Tests automatisés

### 1. Test de la base de données

```sql
-- Vérifier que les colonnes existent
SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'organizations' 
  AND column_name IN ('custom_domain', 'domain_verified', 'domain_verification_token');

SELECT column_name, data_type 
FROM information_schema.columns 
WHERE table_name = 'users' 
  AND column_name IN ('custom_domain', 'domain_verified');

-- Résultat attendu : 5 lignes au total
```

### 2. Test de l'API - Ajouter un domaine

```bash
# Test avec curl
curl -X POST http://localhost:8080/api/domain/add \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "token": "YOUR_SESSION_TOKEN",
    "domain": "example-test.com"
  }'

# Réponse attendue (nboffer < 2) :
# {
#   "success": false,
#   "message": "Custom domains are only available for Professional and Enterprise plans..."
# }

# Réponse attendue (nboffer >= 2) :
# {
#   "success": true,
#   "message": "Custom domain added successfully...",
#   "domain": "example-test.com",
#   "verified": false,
#   "verification_token": "abc123...",
#   "verification_record": "office1789-verification=abc123..."
# }
```

### 3. Test de vérification DNS

```powershell
# Windows PowerShell
Resolve-DnsName -Type TXT example-test.com

# Résultat attendu :
# Name                                     Type   TTL   Section    Strings
# ----                                     ----   ---   -------    -------
# example-test.com                         TXT    3600  Answer     {office1789-verification=abc123...}
```

### 4. Test de l'API - Vérifier le domaine

```bash
curl -X POST http://localhost:8080/api/domain/verify \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "token": "YOUR_SESSION_TOKEN"
  }'

# Réponse attendue (DNS non configuré) :
# {
#   "success": false,
#   "message": "Failed to lookup DNS records..."
# }

# Réponse attendue (DNS OK) :
# {
#   "success": true,
#   "message": "Domain example-test.com successfully verified!",
#   "domain": "example-test.com",
#   "verified": true
# }
```

### 5. Test de synchronisation mail

```bash
# Vérifier que le domaine apparaît dans la config Postfix
docker exec mailserver cat /tmp/docker-mailserver/postfix-virtual-domains.cf

# Résultat attendu :
# office1789.com
# example-test.com

# Vérifier les mappings de mailboxes
docker exec mailserver cat /tmp/docker-mailserver/postfix-virtual-mailboxes.cf

# Résultat attendu :
# testuser@example-test.com testuser@office1789.com
```

### 6. Test d'envoi/réception email

```bash
# Tester la résolution d'adresse
docker exec mailserver postmap -q "testuser@example-test.com" /tmp/docker-mailserver/postfix-virtual-mailboxes.cf

# Résultat attendu : testuser@office1789.com

# Envoyer un email test (avec un client mail)
# De : external@gmail.com
# À : testuser@example-test.com
# Sujet : Test domaine personnalisé

# Vérifier réception :
# 1. Se connecter à Roundcube avec testuser@office1789.com
# 2. Email doit apparaître dans la boîte de réception
```

## 📋 Checklist de déploiement

### Avant le déploiement

- [x] Base de données : colonnes `custom_domain`, `domain_verified` ajoutées
- [x] API Backend : routes `/api/domain/*` créées
- [x] Synchronisation mail : `mailsync.go` implémenté
- [x] Scripts de configuration : `configure-custom-domains.ps1` et `.sh`
- [x] Interface utilisateur : `CustomDomainPanel.vue`
- [x] Documentation : `CUSTOM_DOMAINS.md`

### Déploiement

- [ ] Exécuter migrations BDD si nécessaire
- [ ] Compiler le backend : `go build`
- [ ] Configurer docker-mailserver : `.\configure-custom-domains.ps1`
- [ ] Redémarrer le backend
- [ ] Tester avec un compte Professional/Enterprise
- [ ] Vérifier les logs : `docker logs mailserver`

### Tests post-déploiement

- [ ] Ajouter un domaine test via UI
- [ ] Vérifier génération du token
- [ ] Configurer DNS (enregistrement TXT)
- [ ] Vérifier le domaine via UI
- [ ] Vérifier sync Postfix (fichiers .cf)
- [ ] Tester envoi email vers le domaine personnalisé
- [ ] Vérifier réception dans Roundcube

## 🐛 Scénarios de test

### Scénario 1 : Utilisateur Free essaie d'ajouter un domaine

**Étapes :**
1. Se connecter avec un compte Free (nboffer = 0)
2. Aller dans Compte → Domaine personnalisé
3. Essayer d'ajouter `test.com`

**Résultat attendu :**
- Message d'erreur : "Custom domains are only available for Professional and Enterprise plans"
- Pas de domaine ajouté en BDD

### Scénario 2 : Utilisateur Pro ajoute et vérifie un domaine

**Étapes :**
1. Se connecter avec un compte Pro (nboffer = 2)
2. Aller dans Compte → Domaine personnalisé
3. Ajouter `mycompany.com`
4. Copier le token de vérification
5. Ajouter l'enregistrement TXT dans le DNS
6. Attendre propagation DNS (15 min)
7. Cliquer sur "Vérifier le domaine"

**Résultat attendu :**
- Domaine ajouté avec `domain_verified = FALSE`
- Token généré et affiché
- Après vérification : `domain_verified = TRUE`
- Badge vert ✅ dans l'UI
- Fichiers Postfix mis à jour automatiquement
- Email `username@mycompany.com` fonctionne

### Scénario 3 : Organisation avec 5 membres

**Étapes :**
1. Compte Enterprise crée une organisation
2. Ajoute 5 sous-comptes : alice, bob, charlie, david, eve
3. Ajoute domaine personnalisé `bigcorp.com`
4. Vérifie le domaine

**Résultat attendu :**
- Tous les membres peuvent utiliser `@bigcorp.com`
- Mappings créés :
  - alice@bigcorp.com → alice@office1789.com
  - bob@bigcorp.com → bob@office1789.com
  - charlie@bigcorp.com → charlie@office1789.com
  - david@bigcorp.com → david@office1789.com
  - eve@bigcorp.com → eve@office1789.com

### Scénario 4 : Retirer un domaine personnalisé

**Étapes :**
1. Utilisateur avec domaine vérifié `test.com`
2. Clique sur "Retirer le domaine"
3. Confirme la suppression

**Résultat attendu :**
- Domaine retiré de la BDD (`custom_domain = NULL`)
- Fichiers Postfix mis à jour (domaine supprimé)
- Email `user@test.com` ne fonctionne plus
- Email `user@office1789.com` continue de fonctionner

### Scénario 5 : Domaine déjà utilisé

**Étapes :**
1. Utilisateur A ajoute et vérifie `shared.com`
2. Utilisateur B essaie d'ajouter `shared.com`

**Résultat attendu :**
- Erreur : "This domain is already in use by another account"
- Domaine non ajouté pour utilisateur B

## 📊 Métriques à surveiller

### Backend
- Nombre de domaines personnalisés actifs
- Taux de vérification réussi (verified / total)
- Temps moyen de vérification DNS
- Erreurs de sync mail

```sql
-- Statistiques des domaines
SELECT 
  COUNT(*) as total_domains,
  SUM(CASE WHEN domain_verified THEN 1 ELSE 0 END) as verified_domains,
  SUM(CASE WHEN NOT domain_verified THEN 1 ELSE 0 END) as pending_domains
FROM (
  SELECT custom_domain, domain_verified FROM Users WHERE custom_domain IS NOT NULL
  UNION ALL
  SELECT custom_domain, domain_verified FROM Organizations WHERE custom_domain IS NOT NULL
) AS all_domains;
```

### Serveur mail
- Nombre de domaines dans postfix-virtual-domains.cf
- Nombre de mappings dans postfix-virtual-mailboxes.cf
- Logs Postfix : erreurs de routage
- Emails rejetés pour domaines non vérifiés

```bash
# Compter les domaines
docker exec mailserver wc -l /tmp/docker-mailserver/postfix-virtual-domains.cf

# Compter les mappings
docker exec mailserver wc -l /tmp/docker-mailserver/postfix-virtual-mailboxes.cf

# Logs Postfix (dernières 100 lignes)
docker logs mailserver | grep postfix | tail -n 100
```

## 🔒 Sécurité

### Vérifications de sécurité

- [x] Validation du format de domaine (regex)
- [x] Vérification DNS obligatoire avant activation
- [x] Un domaine = un seul compte (unicité)
- [x] Isolation des domaines par utilisateur/organisation
- [x] SPF/DMARC configurés pour anti-spoofing

### Tests de sécurité

```bash
# Test 1 : Format de domaine invalide
# Input : "invalid..domain"
# Attendu : Erreur "Invalid domain format"

# Test 2 : Domaine avec protocole
# Input : "http://example.com"
# Attendu : Erreur "Invalid domain format"

# Test 3 : Sous-domaine
# Input : "mail.example.com"
# Attendu : Accepté (c'est valide)

# Test 4 : Token brute force
# Input : Essayer 1000 tokens aléatoires
# Attendu : Tous refusés, domaine reste non vérifié
```

## ✅ Résultat final

Si tous les tests passent :
- ✅ Base de données configurée
- ✅ API fonctionnelle
- ✅ Vérification DNS opérationnelle
- ✅ Synchronisation Postfix automatique
- ✅ Interface utilisateur intuitive
- ✅ Documentation complète
- ✅ Sécurité assurée

**Prêt pour la production ! 🚀**
