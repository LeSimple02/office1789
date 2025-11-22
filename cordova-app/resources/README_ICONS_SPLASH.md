# Icons & Splash Setup

Place your images here before running `cordova-res`.

## Recommended Sizes (Android)
Icons:
- drawable-ldpi: 36x36
- drawable-mdpi: 48x48
- drawable-hdpi: 72x72
- drawable-xhdpi: 96x96
- drawable-xxhdpi: 144x144
- drawable-xxxhdpi: 192x192

Splash (portrait):
- port-ldpi: 200x320
- port-mdpi: 320x480
- port-hdpi: 480x800
- port-xhdpi: 720x1280
- port-xxhdpi: 960x1600
- port-xxxhdpi: 1280x1920

Use a flat background color and centered logo for splash.

## Using cordova-res
Install dev dependency (already can be added):

```bash
npm i -D cordova-res
```

Generate assets (from `webfront2/cordova`):
```bash
npx cordova-res android --skip-config --copy
```

Then rebuild:
```bash
cd ../cordova
npx cordova build android
```

## Notes
- config.xml already contains placeholder <icon> and <splash> entries; ensure filenames match.
- Remove unused densities if not provided to speed up build.
- For adaptive icons (Android 8+), consider adding foreground/background layers.
