# Office1789 - Project Structure

Office1789 est maintenant organisé en 3 sous-projets distincts :

## 📁 Structure

```
office1789/
├── webfront2/          # Application Vue 3 principale
├── cordova-app/        # Application mobile Android (Cordova)
└── electron-app/       # Application desktop (Electron)
```

## 🌐 webfront2 (Web App)

Application Vue 3 + Vite principale.

```bash
cd webfront2
npm install
npm run dev          # Dev server
npm run build        # Build production
```

## 📱 cordova-app (Mobile Android)

Application mobile basée sur Cordova.

### Build & Install
```bash
# Depuis webfront2, build et copie vers cordova-app
cd webfront2
npm run cordova:android

# Ou directement depuis cordova-app
cd cordova-app
npx cordova build android

# Install sur device
adb install -r platforms/android/app/build/outputs/apk/debug/app-debug.apk
```

### Splash Screen
- Logo Office1789 centré sur fond blanc
- Proportions respectées (pas d'étirement)
- Durée : 1.5 secondes avec fade
- Généré via Python : `python create-splash-minimal.py`

## 🖥️ electron-app (Desktop)

Application desktop multi-plateforme.

### Développement
```bash
cd electron-app
npm install
npm run electron:dev   # Build Vue + Launch Electron
```

### Build Production
```bash
npm run electron:build  # Crée les packages Windows/Mac/Linux
```

Les exécutables seront dans `electron-app/dist/`.

## 🔄 Workflow de développement

1. **Développer** dans `webfront2/`
2. **Build mobile** : `npm run cordova:android` depuis webfront2
3. **Build desktop** : `npm run electron:dev` depuis electron-app

## 📝 Notes

- **Base URL** : Vue app configuré avec `base: './'` pour compatibilité Cordova/Electron
- **Routing** : Hash mode (`/#/`) utilisé sous `file://` protocol
- **i18n** : Caractères spéciaux échappés (`{'@'}`, `{'['}`, `{']'}`)
