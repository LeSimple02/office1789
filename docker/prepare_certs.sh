#!/bin/bash
set -e

BASE_DIR="/home/debian/office1789/docker"
# Chemin par défaut Let's Encrypt (certbot standalone avec DOMAIN_BASE=office1789.com)
LE_DIR="/etc/letsencrypt/live/office1789.com"
SSL_DIR="$BASE_DIR/config/ssl"

echo "🔐 Préparation des certificats TLS..."

sudo mkdir -p "$SSL_DIR"

if [ -f "$LE_DIR/fullchain.pem" ] && [ -f "$LE_DIR/privkey.pem" ]; then
  echo "✅ Certificats Let's Encrypt trouvés dans $LE_DIR, copie vers $SSL_DIR..."

  echo "📋 Nettoyage ancien SSL..."
  sudo rm -rf "$SSL_DIR"/*
  sudo cp "$LE_DIR/fullchain.pem" "$SSL_DIR/cert.pem"
  sudo cp "$LE_DIR/privkey.pem"  "$SSL_DIR/key.pem"

  sudo chown -R debian:debian "$SSL_DIR"
  sudo chmod 644 "$SSL_DIR/cert.pem"
  sudo chmod 600 "$SSL_DIR/key.pem"

  echo "✅ Certificats Let's Encrypt copiés dans $SSL_DIR"
else
  echo "⚠️  Aucun certificat Let's Encrypt trouvé dans $LE_DIR"
  echo "   -> Génération d'un certificat auto-signé TEMPORAIRE pour le mailserver"

  # Nettoyage et génération auto-signée
  sudo rm -rf "$SSL_DIR"/*
  sudo openssl req -x509 -nodes -newkey rsa:2048 \
    -subj "/CN=office1789.com" \
    -keyout "$SSL_DIR/key.pem" \
    -out "$SSL_DIR/cert.pem" \
    -days 1

  sudo chown -R debian:debian "$SSL_DIR"
  sudo chmod 644 "$SSL_DIR/cert.pem"
  sudo chmod 600 "$SSL_DIR/key.pem"

  echo "✅ Certificat auto-signé généré dans $SSL_DIR"
fi

ls -lh "$SSL_DIR"