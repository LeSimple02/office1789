<?php
echo "=== TEST CONFIGURATION PLUGINS ===\n\n";

// Simuler exactement ce que fait Roundcube
$config = array();

echo "1. Avant tout include:\n";
echo "   \$config['plugins'] = "; var_dump($config['plugins'] ?? 'non défini');

echo "\n2. Include config.inc.php:\n";
include("/var/www/html/config/config.inc.php");
echo "   \$config['plugins'] = "; print_r($config['plugins']);

echo "\n3. Vérifier si config.local.inc.php existe:\n";
if (file_exists('/var/www/html/config/config.local.inc.php')) {
    echo "   ✓ config.local.inc.php existe\n";
    echo "   Contenu:\n";
    echo file_get_contents('/var/www/html/config/config.local.inc.php');
} else {
    echo "   ✗ config.local.inc.php INTROUVABLE\n";
}

echo "\n4. Valeur finale de \$config['plugins']:\n";
print_r($config['plugins']);
