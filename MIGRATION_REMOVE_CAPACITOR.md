# Removal of Capacitor

This project has been migrated to use Cordova only.

## Removed Items
- Dependencies: @capacitor/android, @capacitor/core, @capacitor/cli, @ionic/cli
- Config file: webfront2/capacitor.config.ts
- (Pending) Native android/ folder from previous Capacitor setup

## Cordova Build
1. Run `npm install` in `webfront2/` to prune removed deps.
2. Build web assets: `npm run build`.
3. Copy into Cordova: `npm run cordova:copy`.
4. Build Android APK: `npm run cordova:android` (first run will add platform).

## Rollback (If Needed)
To restore Capacitor:
1. Re-add dependencies:
   ```bash
   npm install --save-dev @capacitor/core @capacitor/cli @capacitor/android @ionic/cli
   ```
2. Restore `capacitor.config.ts` (from version control) under `webfront2/`.
3. Recreate Android project:
   ```bash
   npx cap init
   npx cap add android
   npx cap sync
   ```

## Notes
- Source code had no direct Capacitor API usage, so removal is low risk.
- Verify environment variables and any plugin-specific functionality (none were configured).
- Delete the obsolete `android/` directory when you are certain rollback is not required.
