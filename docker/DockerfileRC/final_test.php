<?php
define('INSTALL_PATH', '/var/www/html/');
require_once INSTALL_PATH . 'program/include/iniset.php';

$rcmail = rcmail::get_instance();

echo "=== DIAGNOSTIC FINAL ===\n\n";
echo "Plugins configurés: ";
$plugins = $rcmail->config->get('plugins');
if (empty($plugins)) {
    echo "VIDE OU NULL\n";
} else {
    print_r($plugins);
}

echo "\nPlugins chargés: ";
print_r($rcmail->plugins->loaded_plugins());

echo "\n\n=== CONTENU config.inc.php ===\n";
system("grep -n 'plugins' /var/www/html/config/config.inc.php | head -n 5");
