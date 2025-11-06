<?php
require_once '/var/www/html/program/lib/Roundcube/bootstrap.php';

error_log('[TEST] Début du test de chargement du plugin SSO');

$rcmail = rcmail::get_instance();
error_log('[TEST] Instance Roundcube créée');

// Charger le plugin
$rcmail->plugins->load_plugin('office1789_sso');
error_log('[TEST] load_plugin appelé');

// Vérifier si la classe existe
if (class_exists('office1789_sso')) {
    error_log('[TEST] SUCCESS: Classe office1789_sso existe');
} else {
    error_log('[TEST] ERROR: Classe office1789_sso N\'EXISTE PAS');
}

// Vérifier les plugins chargés
$loaded = $rcmail->plugins->loaded_plugins();
error_log('[TEST] Plugins chargés: ' . implode(', ', $loaded));

echo "Test terminé\n";
