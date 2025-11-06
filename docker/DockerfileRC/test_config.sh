#!/bin/bash
echo "=== config.local.inc.php ==="
cat /var/www/html/config/config.local.inc.php

echo -e "\n\n=== Test inclusion ==="
php -r '$config = array(); include("/var/www/html/config/config.inc.php"); echo "Plugins après inclusion:\n"; print_r($config["plugins"]);'
