#!/bin/bash
set -e

# Attendre que Roundcube soit initialisé
if [ -d "/var/www/html/config" ]; then
    echo "Copie de la configuration personnalisée..."
    cp /tmp/custom-config.inc.php /var/www/html/config/config.inc.php
    chmod 644 /var/www/html/config/config.inc.php
    echo "Configuration personnalisée appliquée !"
fi

# Lancer la commande par défaut
exec "$@"
