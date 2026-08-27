#!/bin/bash
# Entrypoint personnalisé pour Roundcube avec SSO et injection de config
set -e

# Copier la config custom si présente
if [ -f /var/roundcube/config-wrapper.inc.php ]; then
    echo "[Entrypoint] Injection de la config custom Roundcube..."
    cp /var/roundcube/config-wrapper.inc.php /var/www/html/config/config.inc.php
fi

# Copier le plugin SSO si présent
if [ -d /var/roundcube/plugins/office1789_sso ]; then
    echo "[Entrypoint] Injection du plugin SSO..."
    cp -r /var/roundcube/plugins/office1789_sso /var/www/html/plugins/
fi

# Copier le script auto-dark-mode si présent
if [ -f /var/roundcube/auto-dark-mode.js ]; then
    echo "[Entrypoint] Injection du script auto-dark-mode..."
    cp /var/roundcube/auto-dark-mode.js /var/www/html/skins/elastic/js/
fi

# Copier le logo custom si présent
if [ -f /var/roundcube/office1789_logo.png ]; then
    echo "[Entrypoint] Injection du logo custom..."
    cp /var/roundcube/office1789_logo.png /var/www/html/skins/elastic/images/
fi

# Patch DocumentRoot si besoin (pour images officielles)
if [ -f /etc/apache2/sites-available/000-default.conf ]; then
    sed -i 's|/var/www/html/public_html|/var/www/html|g' /etc/apache2/sites-available/000-default.conf
fi

# Lancer Apache en foreground
echo "[Entrypoint] Démarrage d'Apache..."
exec apache2-foreground#!/bin/sh
set -e

# Run the original entrypoint
/docker-entrypoint.sh "$@" &
ENTRYPOINT_PID=$!

# Wait for config files to be generated
sleep 3

# Inject IMAP and SMTP connection options to disable certificate verification
if [ -f /var/www/html/config/config.inc.php ]; then
    cat >> /var/www/html/config/config.inc.php <<'PHPEOF'

// Custom: Disable SSL verification for self-signed certs
$config['imap_conn_options'] = array(
    'ssl' => array(
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ),
);
$config['smtp_conn_options'] = array(
    'ssl' => array(
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ),
);
PHPEOF
    echo "✓ Injected SSL conn_options into Roundcube config"
fi

# Wait for the original entrypoint to finish
wait $ENTRYPOINT_PID

# Start Apache in foreground so the container stays up

# Fix DocumentRoot if needed (replace /var/www/html/public_html by /var/www/html)
sed -i 's|/var/www/html/public_html|/var/www/html|g' /etc/apache2/sites-enabled/000-default.conf

exec /usr/sbin/apache2ctl -D FOREGROUND
