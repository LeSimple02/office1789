# 🚀 Installation rapide de Stripe

## Installation du package Go

```powershell
cd backend
go get github.com/stripe/stripe-go/v81
go mod tidy
```

## Configuration

1. **Créer le fichier .env** (si pas déjà fait)
```powershell
cp .env.example .env
```

2. **Ajouter vos clés Stripe dans .env**
```env
STRIPE_SECRET_KEY=sk_test_votre_cle
STRIPE_WEBHOOK_SECRET=whsec_votre_secret
STRIPE_PRICE_STANDARD=price_xxx
STRIPE_PRICE_PROFESSIONAL=price_xxx
STRIPE_PRICE_ENTERPRISE=price_xxx
FRONTEND_URL=http://localhost:5173
```

## Test en local

1. **Démarrer le backend**
```powershell
cd backend
go run .
```

2. **Installer Stripe CLI** (pour les webhooks en local)
```powershell
# Via Scoop (recommandé)
scoop install stripe

# Ou télécharger depuis https://github.com/stripe/stripe-cli/releases
```

3. **Démarrer le tunnel webhook**
```powershell
stripe login
stripe listen --forward-to localhost:8080/api/stripe/webhook
```

4. **Dans un autre terminal, démarrer le frontend**
```powershell
cd webfront2
npm run dev
```

## Test du paiement

1. Aller sur http://localhost:5173/account/change
2. Cliquer sur "Passer à Standard" (ou autre plan payant)
3. Utiliser la carte de test : **4242 4242 4242 4242**
4. Date : n'importe quelle date future
5. CVC : n'importe quels 3 chiffres

Après le paiement, vous serez redirigé vers `/account/subscription-success` puis automatiquement vers `/account`.

## Vérification

Check les logs backend pour voir :
```
Stripe checkout session created for user xxx (ID: x), Plan: x, Session: cs_test_...
✅ Subscription activated: User xxx (ID: x) upgraded to plan x via Stripe session cs_test_...
```

Check la DB :
```sql
SELECT username, nboffer FROM Users WHERE username = 'votre_username';
```

## 📚 Documentation complète

Voir `STRIPE_SETUP.md` pour le guide complet incluant :
- Configuration des produits dans le Dashboard Stripe
- Webhooks en production
- Cartes de test
- Résolution de problèmes
- Checklist avant production
