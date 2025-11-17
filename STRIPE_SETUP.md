# Guide d'intégration Stripe pour Office1789

## 📋 Prérequis

1. Compte Stripe : [https://dashboard.stripe.com/register](https://dashboard.stripe.com/register)
2. Go SDK installé (déjà fait via `go get`)

## 🔧 Étape 1 : Installation du SDK Stripe

```bash
cd backend
go get github.com/stripe/stripe-go/v81
go mod tidy
```

## 🔑 Étape 2 : Configuration des clés API

### 2.1 Récupérer vos clés Stripe

1. Connectez-vous au [Dashboard Stripe](https://dashboard.stripe.com/)
2. Allez dans **Developers** → **API keys**
3. Copiez votre **Secret key** (commence par `sk_test_` en mode test)

### 2.2 Créer le fichier .env

```bash
cp .env.example .env
```

Modifiez `.env` avec vos vraies valeurs :

```env
STRIPE_SECRET_KEY=sk_test_VOTRE_CLE_SECRETE
STRIPE_WEBHOOK_SECRET=whsec_VOTRE_WEBHOOK_SECRET
FRONTEND_URL=http://localhost:5173
```

## 💰 Étape 3 : Créer les produits et prix dans Stripe

### 3.1 Via le Dashboard Stripe

1. Allez dans **Products** → **Add product**
2. Créez 3 produits :

#### Produit 1 : Standard
- **Name**: Office1789 Standard
- **Description**: 20GB storage, professional email, priority support
- **Pricing**: 
  - Type: Recurring
  - Amount: €5.00
  - Billing period: Monthly
- Copiez le **Price ID** généré (ex: `price_1ABC...`)

#### Produit 2 : Professional
- **Name**: Office1789 Professional
- **Description**: 100GB storage, custom domain, team features
- **Pricing**: €12.00/month
- Copiez le **Price ID**

#### Produit 3 : Enterprise
- **Name**: Office1789 Enterprise
- **Description**: 500GB storage, unlimited domains, 24/7 support
- **Pricing**: €49.00/month
- Copiez le **Price ID**

### 3.2 Ajoutez les Price IDs dans .env

```env
STRIPE_PRICE_STANDARD=price_1ABC...
STRIPE_PRICE_PROFESSIONAL=price_1DEF...
STRIPE_PRICE_ENTERPRISE=price_1GHI...
```

## 🎣 Étape 4 : Configurer les Webhooks

### 4.1 En développement (local)

Installez Stripe CLI :

**Windows (PowerShell en admin):**
```powershell
scoop install stripe
```

**Alternative - télécharger manuellement:**
[https://github.com/stripe/stripe-cli/releases](https://github.com/stripe/stripe-cli/releases)

**Démarrer le tunnel webhook:**
```bash
stripe listen --forward-to localhost:8080/api/stripe/webhook
```

Stripe CLI vous donnera un **webhook secret** (commence par `whsec_`). Ajoutez-le dans `.env` :

```env
STRIPE_WEBHOOK_SECRET=whsec_VOTRE_SECRET_LOCAL
```

### 4.2 En production

1. Dans le Dashboard Stripe : **Developers** → **Webhooks**
2. Cliquez **Add endpoint**
3. URL : `https://votre-domaine.com/api/stripe/webhook`
4. Sélectionnez les événements :
   - `checkout.session.completed`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
   - `invoice.payment_failed`
5. Copiez le **Signing secret** dans `.env`

## 🚀 Étape 5 : Démarrer le backend

```bash
cd backend
go run .
```

Vous devriez voir :
```
✅ Stripe initialized successfully
[GIN-debug] POST   /api/stripe/checkout       --> main.CreateCheckoutSession
[GIN-debug] POST   /api/stripe/webhook        --> main.StripeWebhook
```

## 🧪 Étape 6 : Tester l'intégration

### Cartes de test Stripe

- **Succès** : `4242 4242 4242 4242`
- **Échec** : `4000 0000 0000 0002`
- **3D Secure** : `4000 0027 6000 3184`

Date : n'importe quelle date future  
CVC : n'importe quel 3 chiffres

### Test du flux complet

1. Frontend : Cliquez sur "Passer à Standard"
2. Le modal s'ouvre et appelle `/api/stripe/checkout`
3. Vous êtes redirigé vers Stripe Checkout
4. Utilisez la carte test `4242 4242 4242 4242`
5. Après paiement, retour vers `/account/subscription-success`
6. Le webhook met à jour automatiquement `nboffer` dans la DB

### Vérifier les logs

**Backend :**
```bash
Stripe checkout session created for user matthis (ID: 1), Plan: 1, Session: cs_test_...
✅ Subscription activated: User matthis (ID: 1) upgraded to plan 1 via Stripe session cs_test_...
```

**Stripe CLI (webhook listener):**
```bash
2024-01-15 14:23:45   --> checkout.session.completed [evt_1ABC...]
2024-01-15 14:23:46  <--  [200] POST http://localhost:8080/api/stripe/webhook [evt_1ABC...]
```

## 📊 Vérification dans la DB

```sql
SELECT username, nboffer FROM Users WHERE username = 'matthis';
-- Résultat attendu : nboffer = 1 (Standard)
```

## 🔒 Sécurité

### Mode Production

1. **Utilisez les clés LIVE** (commence par `sk_live_`)
2. **HTTPS obligatoire** pour les webhooks
3. **Variables d'environnement** : Ne jamais commit `.env`
4. **Vérification des signatures** : Toujours valider `Stripe-Signature`

### .gitignore

Assurez-vous que `.env` est dans `.gitignore` :

```gitignore
# Environment variables
.env
.env.local
.env.production
```

## 🐛 Résolution de problèmes

### Erreur : "Stripe price ID not configured"

→ Vérifiez que `STRIPE_PRICE_STANDARD`, `STRIPE_PRICE_PROFESSIONAL`, `STRIPE_PRICE_ENTERPRISE` sont dans `.env`

### Erreur : "Invalid signature"

→ Vérifiez `STRIPE_WEBHOOK_SECRET` dans `.env`  
→ En local, utilisez `stripe listen --forward-to localhost:8080/api/stripe/webhook`

### Le webhook ne se déclenche pas

→ Vérifiez que Stripe CLI tourne : `stripe listen`  
→ En production, vérifiez l'URL du webhook dans le Dashboard

### L'abonnement ne se met pas à jour

→ Vérifiez les logs backend  
→ Vérifiez les événements dans Dashboard Stripe → Developers → Events  
→ Testez manuellement : `stripe trigger checkout.session.completed`

## 📚 Documentation Stripe

- [Documentation Go](https://stripe.com/docs/api?lang=go)
- [Checkout Sessions](https://stripe.com/docs/payments/checkout)
- [Webhooks](https://stripe.com/docs/webhooks)
- [Testing](https://stripe.com/docs/testing)

## 💡 Fonctionnalités avancées (optionnel)

### Gestion des abonnements

- **Annulation** : Permettre à l'utilisateur d'annuler son abonnement
- **Mise à niveau/Rétrogradation** : Changer de plan en cours d'abonnement
- **Historique des paiements** : Afficher les factures

### Notifications

- **Email de confirmation** : Envoyer un email après souscription
- **Alerte échec de paiement** : Notifier l'utilisateur si le paiement échoue
- **Rappel de renouvellement** : Email avant le renouvellement

### Métriques

- **Dashboard analytics** : Suivre les conversions
- **MRR (Monthly Recurring Revenue)** : Calculer le revenu mensuel
- **Taux de churn** : Analyser les annulations

## ✅ Checklist avant production

- [ ] Clés API LIVE configurées
- [ ] Webhooks en HTTPS configurés
- [ ] Tests de paiement réussis
- [ ] Tests de webhook réussis
- [ ] Gestion des erreurs implémentée
- [ ] Logs de sécurité activés
- [ ] Plan Free (gratuit) toujours disponible
- [ ] Politique de remboursement définie
- [ ] CGV mises à jour avec conditions d'abonnement
