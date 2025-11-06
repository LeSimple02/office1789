#!/bin/bash
# Commenter la ligne plugins dans defaults.inc.php
sed -i "845s/.*/\/\/ \\\$config['plugins'] = []; \/\/ Commenté - défini dans config.inc.php/" /var/www/html/config/defaults.inc.php
echo "Ligne 845 modifiée:"
sed -n '845p' /var/www/html/config/defaults.inc.php
