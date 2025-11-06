<?php
error_reporting(E_ALL);
ini_set('display_errors', 1);

echo "=== TEST CHARGEMENT PLUGINS ROUNDCUBE ===\n\n";

// Charger la config
include('/var/www/html/config/config.inc.php');

echo "Plugins configurés:\n";
print_r($config['plugins']);

echo "\n\n=== TEST CHARGEMENT CLASSE PLUGIN ===\n";

// Charger les classes de base Roundcube nécessaires
require_once('/var/www/html/program/lib/Roundcube/rcube_plugin.php');

// Tenter de charger le plugin
$plugin_file = '/var/www/html/plugins/office1789_sso/office1789_sso.php';

if (file_exists($plugin_file)) {
    echo "✓ Fichier plugin existe: $plugin_file\n";
    require_once($plugin_file);
    
    if (class_exists('office1789_sso')) {
        echo "✓ Classe office1789_sso chargée avec succès\n";
        
        // Tenter d'instancier
        try {
            $plugin = new office1789_sso(null);
            echo "✓ Plugin instancié\n";
            
            // Appeler init() manuellement
            echo "\n=== APPEL DE init() ===\n";
            $plugin->init();
            echo "\n✓ init() appelé - vérifiez les logs [SSO] ci-dessus\n";
            
        } catch (Exception $e) {
            echo "✗ Erreur instanciation: " . $e->getMessage() . "\n";
        }
    } else {
        echo "✗ Classe office1789_sso NON trouvée après require\n";
    }
} else {
    echo "✗ Fichier plugin N'EXISTE PAS: $plugin_file\n";
}
