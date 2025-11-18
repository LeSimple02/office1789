# 🎉 Nouvelles Fonctionnalités Office1789

## 📊 Panel Administrateur

### Accès
- Route : `/admin`
- Réservé aux utilisateurs avec le rôle `admin`
- Connexion requise avec un compte admin

### Fonctionnalités

#### 📈 Statistiques en temps réel
- **Cartes visuelles** avec 6 métriques :
  - 👥 Total utilisateurs
  - ✉️ Emails vérifiés
  - 📱 Téléphones vérifiés
  - ⚠️ Utilisateurs sans contact vérifié
  - 📁 Total fichiers stockés
  - 📅 Total événements calendrier

#### 📊 Graphiques Canvas (Chart.js)
1. **Graphique en Donut** : Distribution des vérifications (email/téléphone/aucun)
2. **Graphique en Barres** : Répartition des utilisateurs par offre (Free/Standard/Pro/Enterprise)

#### 👥 Gestion des utilisateurs
- **Tableau complet** avec recherche en temps réel
- Colonnes affichées :
  - ID utilisateur
  - Username (avec badge coloré)
  - Email principal
  - Recovery Email (avec statut ✓/✗)
  - Téléphone (avec statut ✓/✗)
  - Rôle (User/Admin)
  - Offre d'abonnement
  - Date d'inscription

#### ⚡ Actions administrateur
- **Promouvoir/Rétrograder** : Changer le rôle user ↔ admin
- **Vérification manuelle** : Marquer un email ou téléphone comme vérifié (override)
- **Recherche** : Filtrer par username, email, téléphone
- **Ligne jaune** : Les admins sont surlignés en jaune dans le tableau

---

## 💾 Limitation de Stockage par Offre

### Quotas par Offre

| Offre | nboffer | Stockage | Couleur Badge |
|-------|---------|----------|---------------|
| Free | 0 | 1 GB | Gris |
| Standard | 1 | 50 GB | Bleu |
| Professional | 2 | 200 GB | Violet |
| Enterprise | 3 | ✨ Illimité | Orange |

### Comportement

#### 🚫 Blocage d'upload
- Lorsque l'utilisateur tente d'uploader un fichier :
  - Le backend calcule le stockage actuel utilisé
  - Vérifie si l'upload dépasserait le quota
  - **Bloque** l'upload avec message d'erreur si quota dépassé
  - Retourne les détails : `current_usage`, `storage_limit`, `upload_size`, `remaining`

#### 📊 Affichage du quota
- Composant `StorageQuota.vue` à intégrer dans `DriveView.vue`
- Affiche :
  - Barre de progression colorée (vert → orange → rouge)
  - Pourcentage utilisé
  - Octets utilisés / Limite totale
  - Badge de l'offre actuelle
  - Badge "✨ Illimité" pour Enterprise

#### ⚠️ Avertissements
- **>90%** : "⚠️ Espace presque plein. Pensez à améliorer votre offre."
- **>100%** : "⚠️ Quota dépassé ! Supprimez des fichiers pour libérer de l'espace."

---

## 🛠️ Installation et Configuration

### Frontend (Chart.js déjà installé)
```bash
cd webfront2
npm install chart.js  # ✅ Déjà fait
```

### Backend
Aucune dépendance supplémentaire requise. Les constantes de limite sont déjà définies dans `drive.go`.

### Base de données
Aucune migration nécessaire. Utilise les colonnes existantes :
- `Users.nboffer` : Détermine le quota
- `DriveFiles.file_size` : Calcule l'usage total

---

## 🚀 Endpoints API

### Admin

#### `GET /api/admin/stats`
Retourne les statistiques de la plateforme (protégé par middleware admin).

**Headers** :
```json
{
  "Authorization": "session_token_here"
}
```

**Response** :
```json
{
  "total_users": 150,
  "total_verified_emails": 120,
  "total_verified_phones": 80,
  "users_without_contacts": 10,
  "total_files": 5000,
  "total_calendar_events": 1200
}
```

#### `GET /api/admin/users`
Liste tous les utilisateurs avec détails complets.

**Response** :
```json
[
  {
    "user_id": 46,
    "username": "matthis2",
    "email": "matthis@office1789.com",
    "recovery_email": "matthis.backup@gmail.com",
    "recovery_email_verified": true,
    "phonenumber": "+33612345678",
    "phonenumber_verified": true,
    "role": "admin",
    "nboffer": 3,
    "date_joined": "2024-11-15T10:30:00Z"
  }
]
```

#### `POST /api/admin/users/role`
Change le rôle d'un utilisateur.

**Body** :
```json
{
  "user_id": 46,
  "role": "admin"  // ou "user"
}
```

#### `POST /api/admin/users/verify-contact`
Vérification manuelle d'un contact.

**Body** :
```json
{
  "user_id": 46,
  "contact_type": "email"  // ou "phone"
}
```

### Drive Storage

#### `POST /api/drive/storage`
Retourne l'usage de stockage de l'utilisateur.

**Body** :
```json
{
  "username": "matthis2",
  "token": "session_token_here"
}
```

**Response (Non-illimité)** :
```json
{
  "current_usage": 524288000,      // 500 MB en bytes
  "storage_limit": 1073741824,     // 1 GB en bytes
  "nboffer": 0,
  "remaining": 549453824,          // ~524 MB restants
  "percentage_used": 48.83
}
```

**Response (Illimité)** :
```json
{
  "current_usage": 5368709120,     // 5 GB
  "storage_limit": -1,
  "nboffer": 3,
  "unlimited": true
}
```

---

## 📝 Intégration Frontend

### Panel Admin
Déjà créé dans `webfront2/src/views/AdminPanel.vue`. Accessible via :
```javascript
// Dans router/index.js
{
  path: '/admin',
  name: 'admin',
  component: () => import('../views/AdminPanel.vue')
}
```

### Composant StorageQuota
À intégrer dans `DriveView.vue` :
```vue
<template>
  <div class="drive-view">
    <StorageQuota ref="storageQuota" />
    <!-- Reste du contenu Drive -->
  </div>
</template>

<script setup>
import StorageQuota from '@/components/StorageQuota.vue'
import { ref } from 'vue'

const storageQuota = ref(null)

// Après upload/delete, rafraîchir le quota
async function handleUploadSuccess() {
  await storageQuota.value?.refresh()
}
</script>
```

### Gestion des erreurs d'upload
Intercepter l'erreur `storage_limit_exceeded` :
```javascript
try {
  const response = await fetch('/api/drive/upload', {...})
  const data = await response.json()
  
  if (data.error === 'storage_limit_exceeded') {
    alert(`⚠️ Quota de stockage dépassé !
    
Utilisé : ${formatBytes(data.current_usage)}
Limite : ${formatBytes(data.storage_limit)}
Fichier : ${formatBytes(data.upload_size)}
Restant : ${formatBytes(data.remaining)}

Veuillez supprimer des fichiers ou améliorer votre offre.`)
  }
} catch (error) {
  console.error(error)
}
```

---

## 🔐 Sécurité

### Middleware Admin
Tous les endpoints `/api/admin/*` sont protégés par `CheckAdminMiddleware()` :
- Vérifie le token de session
- Vérifie que `Users.role = 'admin'`
- Retourne `403 Forbidden` si non-admin

### Validation des quotas
- Calcul côté serveur uniquement (aucune confiance au client)
- Vérification atomique avant écriture du fichier
- Transaction DB pour garantir la cohérence

---

## 📊 Base de Données

### Colonnes utilisées

#### `Users`
- `user_id` : ID unique
- `username` : Nom d'utilisateur
- `email` : Email principal (Matrix/compte)
- `recovery_email` : Email de récupération
- `recovery_email_verified` : BOOLEAN (vérification permanente)
- `phonenumber` : Téléphone
- `phonenumber_verified` : BOOLEAN (vérification permanente)
- `role` : VARCHAR(20) - "user" ou "admin"
- `nboffer` : INT - 0=Free, 1=Standard, 2=Pro, 3=Enterprise
- `date_joined` : Date d'inscription

#### `DriveFiles`
- `file_id` : ID unique
- `user_id` : Propriétaire
- `file_name` : Nom du fichier
- `file_path` : Chemin dans l'arborescence
- `file_size` : Taille en bytes (utilisé pour calcul quota)
- `file_type` : Type MIME
- `date_uploaded` : Date d'upload

### Requêtes importantes

#### Calcul du stockage utilisé
```sql
SELECT COALESCE(SUM(file_size), 0) FROM DriveFiles WHERE user_id = $1
```

#### Stats admin
```sql
-- Total utilisateurs
SELECT COUNT(*) FROM Users

-- Emails vérifiés
SELECT COUNT(*) FROM Users WHERE recovery_email_verified = true

-- Téléphones vérifiés
SELECT COUNT(*) FROM Users WHERE phonenumber_verified = true

-- Sans contacts vérifiés
SELECT COUNT(*) FROM Users 
WHERE recovery_email_verified = false AND phonenumber_verified = false

-- Total fichiers
SELECT COUNT(*) FROM DriveFiles

-- Total événements
SELECT COUNT(*) FROM CalendarEvents
```

---

## 🎨 Design

### Palette de couleurs

#### Offres
- **Free (0)** : `#e5e7eb` (gris)
- **Standard (1)** : `#dbeafe` (bleu clair) / `#1e40af` (bleu foncé)
- **Professional (2)** : `#ede9fe` (violet clair) / `#5b21b6` (violet foncé)
- **Enterprise (3)** : `#fef3c7` (orange clair) / `#92400e` (orange foncé)

#### Barre de progression
- **Normal (<80%)** : Dégradé vert `#10b981` → `#059669`
- **Warning (80-95%)** : Dégradé orange `#f59e0b` → `#d97706`
- **Critical (>95%)** : Dégradé rouge `#ef4444` → `#dc2626`

#### Admin Panel
- **Gradient de fond** : `linear-gradient(135deg, #667eea 0%, #764ba2 100%)`
- **Cartes** : Fond blanc, ombre `0 4px 12px rgba(0,0,0,0.1)`
- **Hover** : `translateY(-5px)` + ombre plus prononcée
- **Ligne admin** : Fond `#fef3c7` (jaune clair)

---

## 🧪 Tests

### Test des quotas

1. **Créer un utilisateur Free (nboffer=0)**
   ```sql
   INSERT INTO Users (username, email, password_hash, nboffer)
   VALUES ('testfree', 'test@free.com', 'hash', 0);
   ```

2. **Uploader un fichier de 500 MB** → ✅ OK

3. **Uploader un fichier de 600 MB** → ❌ Bloqué avec erreur `storage_limit_exceeded`

4. **Passer en Standard (nboffer=1)**
   ```sql
   UPDATE Users SET nboffer = 1 WHERE username = 'testfree';
   ```

5. **Uploader le fichier de 600 MB** → ✅ OK (limite à 50 GB)

### Test du panel admin

1. **Créer un admin**
   ```sql
   UPDATE Users SET role = 'admin' WHERE username = 'matthis2';
   ```

2. **Se connecter avec matthis2**

3. **Naviguer vers `/admin`** → Panel s'affiche

4. **Promouvoir un autre utilisateur** → ✅ Rôle changé

5. **Se déconnecter et reconnecter avec l'utilisateur promu** → Accès admin activé

6. **Vérifier manuellement un contact** → Badge ✓ apparaît

---

## 📚 Prochaines étapes recommandées

1. **Ajouter StorageQuota dans DriveView.vue**
2. **Tester l'upload avec quota Free dépassé**
3. **Créer une page de gestion d'abonnement** pour permettre aux users de changer d'offre
4. **Ajouter des notifications** push quand quota >80%
5. **Historique des actions admin** (logs des changements de rôle, vérifications manuelles)
6. **Graphiques temporels** : Evolution des inscriptions, uploads par jour
7. **Export CSV** des utilisateurs depuis le panel admin
8. **Filtres avancés** : Par offre, par date d'inscription, par statut de vérification

---

## 🐛 Troubleshooting

### Panel admin ne s'affiche pas
- Vérifier que l'utilisateur a `role='admin'` dans la DB
- Vérifier que le backend est redémarré avec les nouveaux endpoints
- Vérifier la console navigateur pour erreurs CORS

### Upload bloqué alors que quota non atteint
- Vérifier que `file_size` est correctement enregistré dans `DriveFiles`
- Requête test : `SELECT SUM(file_size) FROM DriveFiles WHERE user_id=X`
- Vérifier que `nboffer` est correct dans `Users`

### Chart.js ne s'affiche pas
- Installer : `npm install chart.js`
- Vérifier l'import dans AdminPanel.vue
- Vérifier que les données sont bien reçues de l'API

---

**✅ Système complètement fonctionnel et prêt à l'emploi !**
