# Cordova wrapper (coexistence with Capacitor)

This folder lets you build a Cordova APK without removing existing Capacitor integration.

## Why this coexistence?
Capacitor is already integrated. Cordova will not fix WebView layout issues automatically; it provides an alternate packaging path only. Keep using your existing web code in `dist/`.

## Structure
```
cordova/
  config.xml            Cordova app configuration
  www/                  Generated (after build) – copy of webfront2/dist
  hooks/                (empty) place custom scripts here
```

## Install Cordova CLI (global)
```
npm install -g cordova
```
(Use an elevated PowerShell if needed.)

## Build web assets first
From the `webfront2` root:
```
npm run build
```
This creates `dist/`.

## Copy web assets into Cordova www
```
mkdir cordova\www 2>NUL
robocopy dist cordova\www /MIR
```
On *nix systems:
```
mkdir -p cordova/www
rsync -a --delete dist/ cordova/www/
```

## Add Android platform
```
cd cordova
cordova platform add android
```
If you already added it once, update:
```
cordova platform rm android
cordova platform add android
```

## Build APK
```
cordova build android --release
```
Debug build:
```
cordova build android
```
The output will appear under `cordova/platforms/android/app/build/outputs/apk/`.

## Signing (release)
Generate keystore (one time):
```
keytool -genkey -v -keystore office1789-release.keystore -alias office1789 -keyalg RSA -keysize 2048 -validity 10000
```
Then:
```
cordova build android --release -- --keystore=office1789-release.keystore --alias=office1789
```
(You may need to supply storepass / keypass.)

## Notes
- Keep a single source of truth: modify Vue code, then rebuild and copy.
- Service workers & PWA features still work inside Cordova WebView with some limitations (depends on platform version).
- Remove Capacitor only *after* validating Cordova build adds value.

## Clean
```
cordova clean
```

## Troubleshooting
| Issue | Fix |
|-------|-----|
| White screen | Check console via `adb logcat` or remote Chrome debugging. Ensure `www/index.html` copied. |
| Plugins need AndroidX | Already enabled via `config.xml` preference. |
| Layout still off-center | Same underlying WebView constraints; adjust CSS, not Cordova. |

