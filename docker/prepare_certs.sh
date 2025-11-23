#!/bin/bash
set -e

BASE_DIR="/home/debian/office1789/docker"
LE_DIR="/etc/letsencrypt/live/office1789.com-0001"
SSL_DIR="$BASE_DIR/config/ssl"

echo "🔐 Préparation des certificats TLS..."

if [ ! -f "$LE_DIR/fullchain.pem" ] || [ ! -f "$LE_DIR/privkey.pem" ]; then
  echo "❌ Certificats Let's Encrypt introuvables dans $LE_DIR"
  exit 1
fi

sudo mkdir -p "$SSL_DIR"

echo "📋 Nettoyage ancien SSL..."
sudo rm -rf "$SSL_DIR"/*
sudo cp "$LE_DIR/fullchain.pem" "$SSL_DIR/cert.pem"
sudo cp "$LE_DIR/privkey.pem"  "$SSL_DIR/key.pem"

sudo chown -R debian:debian "$SSL_DIR"
sudo chmod 644 "$SSL_DIR/cert.pem"
sudo chmod 600 "$SSL_DIR/key.pem"

echo "✅ Certificats copiés dans $SSL_DIR"
ls -lh "$SSL_DIR"