<?php
define('INSTALL_PATH', realpath(__DIR__ . '/') . '/');
require_once INSTALL_PATH . 'program/include/iniset.php';

$rcmail = rcmail::get_instance();

header('Content-Type: text/plain; charset=utf-8');

echo "=== TEST TASK ET FILTRAGE ===\n\n";

echo "1. Task actuelle: " . ($rcmail->task ?? 'NON DÉFINI') . "\n";
echo "2. Action actuelle: " . ($rcmail->action ?? 'NON DÉFINI') . "\n";
echo "3. Paramètres GET: ";
print_r($_GET);
echo "\n4. Paramètres POST: ";
print_r($_POST);

echo "\n5. Plugin office1789_sso:\n";
if (file_exists(INSTALL_PATH . 'plugins/office1789_sso/office1789_sso.php')) {
    include_once(INSTALL_PATH . 'plugins/office1789_sso/office1789_sso.php');
    $test_plugin = new office1789_sso(null);
    echo "   task défini dans le plugin: " . $test_plugin->task . "\n";
    
    // Tester le filtre
    $current_task = $rcmail->task;
    $plugin_tasks = $test_plugin->task;
    
    if ($plugin_tasks && !preg_match('/^('.$plugin_tasks.')$/i', $current_task)) {
        echo "   ❌ PLUGIN FILTRÉ ! Task '$current_task' ne correspond pas au pattern '$plugin_tasks'\n";
    } else {
        echo "   ✓ Plugin correspondrait à la task actuelle\n";
    }
}

echo "\n6. Config plugins:\n";
print_r($rcmail->config->get('plugins'));
