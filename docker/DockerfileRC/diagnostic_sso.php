<?php
/**
 * Diagnostic SSO Office1789
 * Accessible via http://localhost:8081/diagnostic_sso.php
 */

// Charger Roundcube
define('INSTALL_PATH', realpath(__DIR__ . '/') . '/');
require_once INSTALL_PATH . 'program/include/iniset.php';

// Obtenir l'instance Roundcube
$rcmail = rcmail::get_instance();

header('Content-Type: text/html; charset=utf-8');
?>
<!DOCTYPE html>
<html>
<head>
    <title>Diagnostic SSO Office1789</title>
    <style>
        body { font-family: monospace; background: #1e1e1e; color: #d4d4d4; padding: 20px; }
        .ok { color: #4ec9b0; }
        .error { color: #f48771; }
        .warn { color: #dcdcaa; }
        pre { background: #2d2d2d; padding: 10px; border-radius: 5px; }
        h2 { color: #569cd6; }
    </style>
</head>
<body>
    <h1>🔍 Diagnostic SSO Office1789</h1>
    
    <h2>1. Configuration Plugins</h2>
    <pre><?php 
        echo "Variable d'environnement ROUNDCUBE_PLUGINS:\n";
        echo getenv('ROUNDCUBE_PLUGINS') ?: 'Non définie';
        echo "\n\n\$config['plugins'] depuis config.inc.php:\n";
        print_r($rcmail->config->get('plugins'));
    ?></pre>
    
    <h2>2. Plugins Chargés par Roundcube</h2>
    <pre><?php 
        $loaded = $rcmail->plugins->loaded_plugins();
        if (empty($loaded)) {
            echo '<span class="error">✗ AUCUN plugin chargé !</span>';
        } else {
            echo '<span class="ok">✓ ' . count($loaded) . ' plugins chargés:</span>' . "\n";
            print_r($loaded);
        }
    ?></pre>
    
    <h2>3. Plugin office1789_sso</h2>
    <pre><?php 
        if (in_array('office1789_sso', $loaded)) {
            echo '<span class="ok">✓ office1789_sso EST CHARGÉ</span>';
        } else {
            echo '<span class="error">✗ office1789_sso N\'EST PAS CHARGÉ</span>';
            echo "\n\nVérifications:\n";
            $plugin_file = INSTALL_PATH . 'plugins/office1789_sso/office1789_sso.php';
            if (file_exists($plugin_file)) {
                echo "✓ Fichier existe: $plugin_file\n";
                echo "  Taille: " . filesize($plugin_file) . " bytes\n";
                echo "  Modifié: " . date('Y-m-d H:i:s', filemtime($plugin_file)) . "\n";
            } else {
                echo "✗ Fichier INTROUVABLE: $plugin_file\n";
            }
        }
    ?></pre>
    
    <h2>4. Test Token SSO</h2>
    <pre><?php 
        $sso_token = $_GET['sso_token'] ?? 'Aucun token dans l\'URL';
        echo "Token reçu: " . htmlspecialchars($sso_token);
        
        if ($sso_token && $sso_token !== 'Aucun token dans l\'URL') {
            echo "\n\nDécodage du token:\n";
            $parts = explode('.', $sso_token);
            if (count($parts) === 2) {
                $claims_b64 = $parts[0];
                $signature_b64 = $parts[1];
                $claims_json = base64_decode($claims_b64);
                echo "Claims: $claims_json\n";
                $claims = json_decode($claims_json, true);
                print_r($claims);
            } else {
                echo "Format invalide (attendu: claims.signature)";
            }
        }
    ?></pre>
    
    <h2>5. État Session</h2>
    <pre><?php 
        if ($rcmail->user->ID) {
            echo '<span class="ok">✓ Utilisateur connecté: ' . $rcmail->user->get_username() . '</span>';
        } else {
            echo '<span class="warn">⚠ Aucun utilisateur connecté</span>';
        }
        
        echo "\n\nSession:\n";
        print_r($_SESSION);
    ?></pre>
    
</body>
</html>
