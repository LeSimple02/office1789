import fs from 'fs';
import path from 'path';

const distDir = path.resolve('dist');
const targetDir = path.resolve('..', 'cordova-app', 'www');

function copyRecursive(src, dest) {
  const stat = fs.statSync(src);
  if (stat.isDirectory()) {
    if (!fs.existsSync(dest)) fs.mkdirSync(dest, { recursive: true });
    for (const entry of fs.readdirSync(src)) {
      copyRecursive(path.join(src, entry), path.join(dest, entry));
    }
  } else {
    fs.copyFileSync(src, dest);
  }
}

function removeIfExists(p) {
  if (fs.existsSync(p)) {
    for (const entry of fs.readdirSync(p)) {
      const full = path.join(p, entry);
      const stat = fs.lstatSync(full);
      if (stat.isDirectory()) {
        removeIfExists(full);
      } else {
        fs.unlinkSync(full);
      }
    }
    fs.rmdirSync(p);
  }
}

if (!fs.existsSync(distDir)) {
  console.error('dist directory not found. Run build first.');
  process.exit(1);
}

removeIfExists(targetDir);
fs.mkdirSync(targetDir, { recursive: true });
copyRecursive(distDir, targetDir);

console.log('Copied dist -> ../cordova-app/www');
