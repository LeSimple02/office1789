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

# Injecter le CSS personnalisé dans le template HTML de Roundcube
SKIN_META="/var/www/html/skins/elastic/templates/minimal.html"
if [ -f "$SKIN_META" ]; then
    if ! grep -q "office1789-custom.css" "$SKIN_META"; then
        sed -i 's|</head>|<link rel="stylesheet" type="text/css" href="/skins/elastic/styles/office1789-custom.css">\n</head>|' "$SKIN_META"
        echo "CSS Office1789 injecté dans minimal.html"
    fi
fi

SKIN_FULL="/var/www/html/skins/elastic/templates/mail.html"
if [ -f "$SKIN_FULL" ]; then
    if ! grep -q "office1789-custom.css" "$SKIN_FULL"; then
        sed -i 's|</head>|<link rel="stylesheet" type="text/css" href="/skins/elastic/styles/office1789-custom.css">\n</head>|' "$SKIN_FULL"
        echo "CSS Office1789 injecté dans mail.html"
    fi
fi

echo "Configuration locale activée avec priorité maximale !"
