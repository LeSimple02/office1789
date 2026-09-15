#!/bin/bash
# Script pour configurer Roundcube avec notre config locale

CONFIG_FILE="/var/www/html/config/config.inc.php"
LOCAL_CONFIG="/var/www/html/config/config.local.inc.php"

# Attendre que Roundcube soit initialisé
while [ ! -f "$CONFIG_FILE" ]; do
    echo "Attente de l'initialisation de Roundcube..."
    sleep 2
done

echo "Réécriture de config.inc.php pour charger config.local.inc.php en dernier..."
# Réécrire complètement config.inc.php pour inclure config.docker D'ABORD, puis config.local qui écrase
cat > "$CONFIG_FILE" << 'EOFCONFIG'
<?php
    // Charger d'abord les paramètres par défaut de Roundcube
    $config['plugins'] = [];
    $config['log_driver'] = 'stdout';
    $config['zipdownload_selection'] = true;
    $config['des_key'] = '1Zpms5JFDU17jVORhV0xWkGB';
    $config['enable_spellcheck'] = true;
    $config['spellcheck_engine'] = 'pspell';

    // Charger config.docker.inc.php EN PREMIER (sera écrasé)
    if (file_exists(__DIR__ . '/config.docker.inc.php')) {
        include(__DIR__ . '/config.docker.inc.php');
    }

    // Puis charger notre config locale QUI ÉCRASE config.docker (priorité maximale)
    if (file_exists(__DIR__ . '/config.local.inc.php')) {
        include(__DIR__ . '/config.local.inc.php');
    }
EOFCONFIG

echo "Configuration locale activée avec priorité maximale !"

# Patcher index.php pour inclure le handler SSO juste après iniset.php
INDEX_FILE="/var/www/html/index.php"
if [ -f "$INDEX_FILE" ]; then
    if ! grep -q "sso-handler.php" "$INDEX_FILE"; then
        sed -i "/require_once 'program\/include\/iniset.php';/a\\\n// Office1789 SSO Handler\\nrequire_once __DIR__ . '\/sso-handler.php';" "$INDEX_FILE"
        echo "Handler SSO ajouté à index.php"
    else
        echo "Handler SSO déjà présent dans index.php"
    fi
fi
