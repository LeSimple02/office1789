# 🚀 Guide de démarrage rapide - Fonctionnalités Enterprise

## Pour les administrateurs système

### ⚡ Installation rapide (5 minutes)

```powershell
# 1. Aller dans le dossier backend
cd C:\Users\Matthis\Documents\office1789\backend

# 2. Compiler le backend
go build

# 3. Configurer le serveur mail
cd ..\docker
.\configure-custom-domains.ps1

# 4. Démarrer le backend
cd ..\backend
.\main.exe

# ✅ C'est prêt !
```

### 🔍 Vérification rapide

```bash
# Vérifier que tout fonctionne
docker ps  # Tous les conteneurs doivent être "Up"
docker exec mailserver postfix status  # Doit dire "running"
docker exec mailserver cat /tmp/docker-mailserver/postfix-virtual-domains.cf  # Doit contenir office1789.com
```

---

## Pour les utilisateurs finaux

### 🏢 Créer une organisation (2 minutes)

1. **Se connecter** à Office1789
2. Aller dans **Compte** → **Gestion de l'organisation**
3. Remplir le formulaire :
   - Nom d'utilisateur du sous-compte
   - Mot de passe (min 8 caractères)
   - Email de récupération (optionnel)
   - Nom de l'organisation (première fois uniquement)
4. Cliquer sur **Créer le sous-compte**

**Résultat :** Un nouveau compte est créé avec l'email `nouveaucompte@office1789.com`

### 🌐 Ajouter un domaine personnalisé (10 minutes + délai DNS)

#### Étape 1 : Ajouter le domaine (1 minute)
1. Aller dans **Compte** → **Domaine personnalisé**
2. Entrer votre domaine : `votreentreprise.com`
3. Cliquer sur **Ajouter le domaine**
4. **Copier le token** affiché (ex: `office1789-verification=abc123...`)

#### Étape 2 : Configurer le DNS (5 minutes)
Aller chez votre registrar (OVH, Gandi, Cloudflare...) et ajouter ces enregistrements :

| Type | Nom | Valeur | Priorité |
|------|-----|--------|----------|
| TXT | @ | `office1789-verification=VOTRE_TOKEN` | - |
| MX | @ | `mail.office1789.com` | 10 |
| TXT | @ | `v=spf1 include:office1789.com ~all` | - |
| TXT | _dmarc | `v=DMARC1; p=quarantine` | - |

**Important :** Remplacez `VOTRE_TOKEN` par le token copié à l'étape 1.

#### Étape 3 : Vérifier (5 minutes à 48h)
1. Attendre la propagation DNS (15 min minimum)
2. Vérifier avec : `nslookup -type=TXT votreentreprise.com`
3. Dans Office1789, cliquer sur **Vérifier le domaine**
4. Si ✅ apparaît, c'est bon !

**Résultat :** Vous pouvez maintenant utiliser `contact@votreentreprise.com` au lieu de `contact@office1789.com`

---

## 💡 Cas d'usage

### Cas 1 : Petite entreprise (Plan Professional)
**Besoin :** 3 employés + domaine `petiteboutique.fr`

```
1. Patron crée 3 sous-comptes :
   - alice@office1789.com
   - bob@office1789.com  
   - charlie@office1789.com

2. Patron ajoute le domaine petiteboutique.fr

3. Tous peuvent maintenant utiliser :
   - alice@petiteboutique.fr
   - bob@petiteboutique.fr
   - charlie@petiteboutique.fr
```

**Coût :** 12€/mois (plan Professional)

### Cas 2 : Moyenne entreprise (Plan Enterprise)
**Besoin :** 15 employés + domaine `monentreprise.com`

```
1. Admin crée 15 sous-comptes
2. Admin configure monentreprise.com
3. Chaque employé a son email @monentreprise.com
4. Partage de fichiers entre tous les membres
5. Support 24/7 prioritaire
```

**Coût :** 50€/mois (plan Enterprise)

### Cas 3 : Freelance avec plusieurs projets
**Besoin :** Domaines multiples pour différents clients

```
Plan Professional avec 1 domaine principal
OU
Plan Enterprise si besoins de plusieurs sous-comptes
```

---

## 🆘 Dépannage express

### ❌ "Custom domains are only available for Professional..."
**Solution :** Passez au plan Professional ou Enterprise dans **Compte** → **Changer d'offre**

### ❌ "Domain verification failed"
**Solutions :**
1. Attendre plus longtemps (jusqu'à 48h)
2. Vérifier le DNS : `nslookup -type=TXT votredomaine.com`
3. Vérifier que le token est exact (copier-coller)
4. Essayer avec un autre serveur DNS : `nslookup -type=TXT votredomaine.com 8.8.8.8`

### ❌ "This domain is already in use"
**Solution :** Ce domaine est déjà vérifié par un autre compte. Utilisez un autre domaine ou contactez le support.

### ❌ "Sub-account limit reached"
**Solution :** 
- Plan Professional : limite de 3 → Passez à Enterprise (20)
- Plan Enterprise : limite de 20 → Contactez le support pour un plan custom

### ❌ Les emails n'arrivent pas
**Vérifications :**
1. Domaine vérifié ? (badge ✅ dans l'interface)
2. Enregistrement MX configuré ? `nslookup -type=MX votredomaine.com`
3. Postfix rechargé ? `docker exec mailserver postfix reload`
4. Logs : `docker logs mailserver | tail -n 50`

---

## 📊 Comprendre les limites

```
Free (0€)     : 1 compte, 1GB, pas de domaine, pas de partage
Standard (5€) : 1 compte, 50GB, pas de domaine, pas de partage
Professional (12€) : 1 compte + 3 sous-comptes, 200GB, 1 domaine, partage 3 membres
Enterprise (50€)   : 1 compte + 20 sous-comptes, illimité, 1 domaine, partage 20 membres
```

---

## 🎓 Formation rapide (5 minutes)

### Pour le compte principal (organization_owner)
**Vous pouvez :**
- ✅ Créer des sous-comptes (limité par votre plan)
- ✅ Supprimer des sous-comptes
- ✅ Ajouter un domaine personnalisé
- ✅ Gérer toute l'organisation
- ✅ Partager des fichiers avec tous les membres

**Vous ne pouvez pas :**
- ❌ Supprimer votre propre compte (owner) sans supprimer l'organisation
- ❌ Transférer la propriété à un sous-compte (pour l'instant)

### Pour les sous-comptes (organization_member)
**Vous pouvez :**
- ✅ Vous connecter avec votre propre email/password
- ✅ Utiliser toutes les fonctionnalités Office1789 (mail, drive, chat, calendrier)
- ✅ Utiliser le domaine personnalisé de l'organisation
- ✅ Recevoir des fichiers partagés

**Vous ne pouvez pas :**
- ❌ Créer d'autres sous-comptes
- ❌ Modifier le domaine personnalisé
- ❌ Gérer l'organisation

---

## 🔐 Sécurité - Bonnes pratiques

### ✅ À FAIRE
- Utiliser des mots de passe forts pour tous les sous-comptes
- Activer la 2FA pour le compte principal
- Vérifier régulièrement les membres actifs
- Configurer SPF et DMARC pour éviter le spam
- Backup régulier des données importantes

### ❌ À ÉVITER
- Partager les identifiants de sous-comptes
- Utiliser le même mot de passe pour tous
- Oublier de retirer les membres qui quittent l'entreprise
- Négliger la vérification du domaine
- Ignorer les emails de sécurité

---

## 📞 Obtenir de l'aide

### Support technique
- **Professional** : Support prioritaire ⭐ (réponse sous 24h)
- **Enterprise** : Support 24/7 ⭐⭐ (réponse sous 2h)
- **Email** : support@office1789.com
- **Documentation** : Voir `CUSTOM_DOMAINS.md`

### Communauté
- **Forum** : https://community.office1789.com
- **Discord** : https://discord.gg/office1789
- **GitHub Issues** : Pour reporter des bugs

---

## 🎉 Félicitations !

Vous êtes maintenant prêt à utiliser Office1789 avec :
- ✅ Votre propre domaine personnalisé
- ✅ Plusieurs sous-comptes pour votre équipe
- ✅ Tous les outils collaboratifs (mail, drive, chat, calendrier)
- ✅ Support prioritaire

**Besoin d'aide ?** Consultez `CUSTOM_DOMAINS.md` ou contactez le support ! 🚀
