# Configuration Backend pour Test Mobile/Desktop

## 🚀 Accès depuis votre portable

### Option 1: Utiliser votre IP locale (réseau Wi-Fi)

1. **Trouvez votre IP locale** (ordinateur avec backend):
```powershell
ipconfig
# Cherchez "Adresse IPv4" sous votre adaptateur Wi-Fi
# Ex: 192.168.1.100
```

2. **Configurez le backend Go** pour écouter sur toutes les interfaces:

Dans `backend/main.go`, remplacez:
```go
r.Run(":8080")  // Écoute seulement localhost
```
Par:
```go
r.Run("0.0.0.0:8080")  // Écoute sur toutes les interfaces
```

3. **Autorisez le port dans le pare-feu Windows**:
```powershell
New-NetFirewallRule -DisplayName "Office1789 Backend" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow
```

4. **Configurez l'API URL dans votre app** (webfront2):

Dans `webfront2/src/stores/global.js` ou configuration API:
```javascript
// En dev: utilisez votre IP locale
const API_URL = 'http://192.168.1.100:8080'  // Remplacez par VOTRE IP

// Détection automatique:
const API_URL = window.location.protocol === 'file:' 
  ? 'http://192.168.1.100:8080'  // Cordova/Electron
  : 'http://localhost:8080'       // Web dev
```

5. **Testez depuis votre portable**:
```bash
# Sur le portable, vérifiez la connexion:
curl http://192.168.1.100:8080/api/health
```

### Option 2: Tunnel avec ngrok (accès Internet)

Si vous n'êtes pas sur le même réseau:

1. **Installez ngrok**:
```powershell
winget install ngrok.ngrok
```

2. **Lancez le tunnel**:
```bash
ngrok http 8080
```

3. **Utilisez l'URL fournie** (ex: `https://abc123.ngrok.io`) dans votre app

### Option 3: Backend cloud temporaire

Déployez sur un service gratuit:
- **Render.com** (backend Go)
- **Railway.app** (Go + PostgreSQL)
- **Fly.io** (conteneurs)

## 🔧 Configuration API dans l'app

### Cordova (Android)
Créez `cordova-app/www/config.js`:
```javascript
window.API_CONFIG = {
  baseURL: 'http://192.168.1.100:8080',  // Votre IP
  timeout: 10000
}
```

### Electron (Desktop)
Dans `electron-app/preload.js`:
```javascript
contextBridge.exposeInMainWorld('config', {
  apiURL: process.env.API_URL || 'http://localhost:8080'
})
```

Lancez avec:
```bash
API_URL=http://192.168.1.100:8080 npm run electron:dev
```

## 🔐 CORS Backend

Configurez CORS dans `backend/main.go`:
```go
import "github.com/gin-contrib/cors"

func main() {
    r := gin.Default()
    
    // CORS pour permettre app mobile/desktop
    r.Use(cors.New(cors.Config{
        AllowOrigins: []string{
            "http://localhost:*",
            "http://192.168.*.*:*",  // Réseau local
            "file://*",              // Cordova
            "capacitor://*",
            "app://*",               // Electron
        },
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
    }))
    
    // Routes...
}
```

## 📱 Test Rapide

**Terminal 1** (Backend):
```bash
cd backend
go run main.go
# Devrait afficher: Listening on 0.0.0.0:8080
```

**Terminal 2** (Vérif réseau):
```bash
# Depuis votre portable (même Wi-Fi):
curl http://VOTRE_IP:8080/api/health
```

**Terminal 3** (Build app):
```bash
cd webfront2
npm run cordova:android
# Installer l'APK sur le portable
```

## ⚠️ Points importants

1. **Même réseau Wi-Fi** : Ordinateur et portable doivent être connectés au même réseau
2. **IP dynamique** : L'IP locale peut changer, utilisez une IP statique dans votre routeur
3. **HTTPS** : Pour features avancées (géolocalisation, notifications), utilisez HTTPS/ngrok
4. **Environnements** : Créez des configs différentes (dev/prod) dans votre app
