#!/bin/bash
# Script pour configurer Roundcube avec notre config locale

CONFIG_FILE="/var/www/html/config/config.inc.php"
LOCAL_CONFIG="/var/www/html/config/config.local.inc.php"

# Attendre que Roundcube soit initialisÃ©
while [ ! -f "$CONFIG_FILE" ]; do
    echo "Attente de l'initialisation de Roundcube..."
    sleep 2
done

echo "RÃ©Ã©criture de config.inc.php pour charger config.local.inc.php en dernier..."

# RÃ©Ã©crire complÃ¨tement config.inc.php pour inclure config.docker D'ABORD, puis config.local qui Ã©crase
cat > "$CONFIG_FILE" << 'EOFCONFIG'
<?php
    // Charger d'abord les paramÃ¨tres par dÃ©faut de Roundcube
    $config['plugins'] = [];
    $config['log_driver'] = 'stdout';
    $config['zipdownload_selection'] = true;
    $config['des_key'] = '1Zpms5JFDU17jVORhV0xWkGB';
    $config['enable_spellcheck'] = true;
    $config['spellcheck_engine'] = 'pspell';
    
    // Charger config.docker.inc.php EN PREMIER (sera Ã©crasÃ©)
    if (file_exists(__DIR__ . '/config.docker.inc.php')) {
        include(__DIR__ . '/config.docker.inc.php');
    }
    
    // Puis charger notre config locale QUI Ã‰CRASE config.docker (prioritÃ© maximale)
    if (file_exists(__DIR__ . '/config.local.inc.php')) {
        include(__DIR__ . '/config.local.inc.php');
    }
EOFCONFIG

# Modifier index.php pour inclure le handler SSO au début
INDEX_FILE="/var/www/html/index.php"
if [ -f "$INDEX_FILE" ]; then
    # Vérifier si SSO handler n'est pas déjà inclus
    if ! grep -q "sso-handler.php" "$INDEX_FILE"; then
        # Insérer après la ligne "require_once 'program/include/iniset.php';"
        sed -i "/require_once 'program\/include\/iniset.php';/a\\
\\
// Office1789 SSO Handler\\
require_once __DIR__ . '/sso-handler.php';" "$INDEX_FILE"
        echo "Handler SSO ajouté à index.php"
    fi
fi

# Injecter le CSS personnalisé dans le template HTML de Roundcube
SKIN_META="/var/www/html/skins/elastic/templates/minimal.html"
if [ -f "$SKIN_META" ]; then
    if ! grep -q "office1789-custom.css" "$SKIN_META"; then
        sed -i 's|</head>|<link rel="stylesheet" type="text/css" href="/skins/elastic/styles/office1789-custom.css">\n</head>|' "$SKIN_META"
        echo "CSS Office1789 injectÃ© dans minimal.html"
    fi
fi

SKIN_FULL="/var/www/html/skins/elastic/templates/mail.html"
if [ -f "$SKIN_FULL" ]; then
    if ! grep -q "office1789-custom.css" "$SKIN_FULL"; then
        sed -i 's|</head>|<link rel="stylesheet" type="text/css" href="/skins/elastic/styles/office1789-custom.css">\n</head>|' "$SKIN_FULL"
        echo "CSS Office1789 injectÃ© dans mail.html"
    fi
fi

echo "Configuration locale activÃ©e avec prioritÃ© maximale !"
