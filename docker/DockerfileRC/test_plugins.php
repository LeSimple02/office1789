<?php
include("/var/www/html/config/config.inc.php");
echo "=== PLUGINS CONFIGURÉS ===\n";
print_r($config['plugins']);
echo "\n=== FICHIERS PLUGINS ===\n";
$plugin_dir = '/var/www/html/plugins/office1789_sso';
if (is_dir($plugin_dir)) {
    echo "✓ Plugin office1789_sso existe\n";
    if (file_exists($plugin_dir . '/office1789_sso.php')) {
        echo "✓ Fichier office1789_sso.php existe (" . filesize($plugin_dir . '/office1789_sso.php') . " bytes)\n";
    }
}
$plugin_dir = '/var/www/html/plugins/office1789_darkmode';
if (is_dir($plugin_dir)) {
    echo "✓ Plugin office1789_darkmode existe\n";
    if (file_exists($plugin_dir . '/office1789_darkmode.php')) {
        echo "✓ Fichier office1789_darkmode.php existe (" . filesize($plugin_dir . '/office1789_darkmode.php') . " bytes)\n";
    }
}
