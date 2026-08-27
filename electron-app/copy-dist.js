#!/usr/bin/env node
/**
 * Copy Vite build output (dist) to Electron www directory
 */
const fs = require('fs')
const path = require('path')

const sourceDir = path.join(__dirname, '..', 'webfront2', 'dist')
const targetDir = path.join(__dirname, 'www')

function copyRecursive(src, dest) {
  if (fs.existsSync(dest)) {
    fs.rmSync(dest, { recursive: true, force: true })
  }
  
  fs.mkdirSync(dest, { recursive: true })
  
  const entries = fs.readdirSync(src, { withFileTypes: true })
  
  for (const entry of entries) {
    const srcPath = path.join(src, entry.name)
    const destPath = path.join(dest, entry.name)
    
    if (entry.isDirectory()) {
      copyRecursive(srcPath, destPath)
    } else {
      fs.copyFileSync(srcPath, destPath)
    }
  }
}

console.log('Copying dist -> electron-app/www')
copyRecursive(sourceDir, targetDir)
console.log('✓ Copy complete')
