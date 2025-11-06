<?php
define('INSTALL_PATH', realpath(__DIR__ . '/') . '/');
require_once INSTALL_PATH . 'program/include/iniset.php';

$rcmail = rcmail::get_instance();

header('Content-Type: text/plain; charset=utf-8');

echo "=== TEST CHARGEMENT PLUGINS WEB ===\n\n";

echo "1. \$config['plugins'] depuis rcmail->config:\n";
$plugins_config = $rcmail->config->get('plugins');
print_r($plugins_config);

echo "\n2. Plugins chargés par rcmail->plugins:\n";
$loaded = $rcmail->plugins->loaded_plugins();
print_r($loaded);

echo "\n3. Contenu de config.inc.php:\n";
$config_content = file_get_contents(INSTALL_PATH . 'config/config.inc.php');
preg_match('/\$config\[\'plugins\'\]\s*=\s*\[(.*?)\];/s', $config_content, $matches);
echo "Ligne trouvée: " . ($matches[0] ?? 'NON TROUVÉ') . "\n";

echo "\n4. Test manuel d'inclusion:\n";
$config = array(); // Utiliser $config, pas $test_config !
include(INSTALL_PATH . 'config/config.inc.php');
echo "Plugins après include manuel: ";
print_r($config['plugins'] ?? 'NON DÉFINI');
