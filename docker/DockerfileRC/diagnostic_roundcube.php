<?php
echo "=== DIAGNOSTIC ROUNDCUBE ===\n\n";

// 1. Vérifier la config
echo "1. CONFIGURATION\n";
include("/var/www/html/config/config.inc.php");
echo "Plugins configurés dans config.inc.php: ";
print_r($config['plugins']);
echo "\n";

// 2. Vérifier que Roundcube peut se charger
echo "2. CHARGEMENT ROUNDCUBE\n";
define('INSTALL_PATH', '/var/www/html/');
if (file_exists(INSTALL_PATH . 'program/include/iniset.php')) {
    echo "✓ iniset.php existe\n";
    include_once(INSTALL_PATH . 'program/include/iniset.php');
    echo "✓ iniset.php chargé\n";
} else {
    echo "✗ iniset.php introuvable\n";
    exit(1);
}

// 3. Vérifier rcmail
echo "\n3. INSTANCE RCMAIL\n";
try {
    $rcmail = rcmail::get_instance();
    echo "✓ rcmail instance créée\n";
    
    // 4. Vérifier les plugins chargés
    echo "\n4. PLUGINS CHARGÉS PAR ROUNDCUBE\n";
    $loaded_plugins = $rcmail->plugins->loaded_plugins();
    if (empty($loaded_plugins)) {
        echo "✗ AUCUN plugin chargé !\n";
    } else {
        echo "Plugins chargés: ";
        print_r($loaded_plugins);
    }
    
    // 5. Vérifier si office1789_sso est actif
    echo "\n5. PLUGIN office1789_sso\n";
    if (in_array('office1789_sso', $loaded_plugins)) {
        echo "✓ office1789_sso EST chargé\n";
    } else {
        echo "✗ office1789_sso N'EST PAS chargé\n";
        echo "   → Vérifions pourquoi...\n";
        if (file_exists('/var/www/html/plugins/office1789_sso/office1789_sso.php')) {
            echo "   ✓ Fichier plugin existe\n";
            $syntax = shell_exec('php -l /var/www/html/plugins/office1789_sso/office1789_sso.php 2>&1');
            echo "   " . trim($syntax) . "\n";
        } else {
            echo "   ✗ Fichier plugin INTROUVABLE\n";
        }
    }
    
} catch (Exception $e) {
    echo "✗ Erreur: " . $e->getMessage() . "\n";
}
