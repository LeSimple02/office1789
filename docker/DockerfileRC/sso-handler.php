<?php
/**
 * Office1789 SSO Handler
 * Ce fichier est inclus au début de index.php pour gérer l'authentification SSO
 */

$sso_secret = 'Office1789-SecretKey-ChangeInProduction';

// Fonction pour valider le token SSO
function validate_sso_token($token, $secret) {
    $parts = explode('.', $token);
    if (count($parts) !== 2) {
        return false;
    }
    
    list($claimsB64, $signature) = $parts;
    
    // Vérifier la signature
    $expected_signature = base64_encode(hash_hmac('sha256', $claimsB64, $secret, true));
    $expected_signature = str_replace(['+', '/'], ['-', '_'], rtrim($expected_signature, '='));
    
    if ($signature !== $expected_signature) {
        error_log('[SSO] Invalid signature');
        return false;
    }
    
    // Décoder les claims
    $claimsJSON = base64_decode(str_replace(['-', '_'], ['+', '/'], $claimsB64));
    $claims = json_decode($claimsJSON, true);
    
    if (!$claims) {
        error_log('[SSO] Invalid claims JSON');
        return false;
    }
    
    return $claims;
}

// Récupérer le token SSO depuis l'URL
$sso_token = isset($_GET['sso_token']) ? $_GET['sso_token'] : null;

// Si pas de token SSO = BLOQUER L'ACCÈS (sauf ressources statiques, POST SSO, et actions AJAX)
if (empty($sso_token)) {
    $request_method = isset($_SERVER['REQUEST_METHOD']) ? $_SERVER['REQUEST_METHOD'] : 'GET';
    $request_uri = isset($_SERVER['REQUEST_URI']) ? $_SERVER['REQUEST_URI'] : '';
    
    // Autoriser les ressources statiques (CSS, JS, images, fonts)
    $is_static_resource = (
        strpos($request_uri, '/skins/') !== false ||
        strpos($request_uri, '/plugins/') !== false ||
        strpos($request_uri, '/program/') !== false ||
        preg_match('/\.(css|js|png|jpg|jpeg|gif|svg|woff|woff2|ttf|ico)/', $request_uri)
    );
    
    // Autoriser les requêtes POST (pour le formulaire de login auto-submit SSO uniquement)
    $is_post_sso = ($request_method === 'POST' && isset($_POST['_action']) && $_POST['_action'] === 'login');
    
    // Autoriser les requêtes AJAX une fois connecté (task=mail, task=settings, etc.)
    // MAIS BLOQUER task=login (page de login) et task=logout
    $is_authenticated_task = (
        isset($_GET['_task']) && 
        $_GET['_task'] !== 'login' && 
        $_GET['_task'] !== 'logout' &&
        isset($_COOKIE['roundcube_sessid'])
    );
    
    // BLOQUER tout le reste (y compris GET / direct, ?_task=login, ?_task=logout)
    if (!$is_static_resource && !$is_post_sso && !$is_authenticated_task) {
        error_log('[SSO] Accès bloqué - pas de token SSO (URI: ' . $request_uri . ', Method: ' . $request_method . ')');
        header('HTTP/1.1 403 Forbidden');
        echo '<!DOCTYPE html><html><head><title>Accès refusé</title><style>
            body { font-family: Arial; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
            .box { background: white; padding: 40px; border-radius: 10px; box-shadow: 0 10px 40px rgba(0,0,0,0.3); text-align: center; max-width: 500px; }
            h1 { color: #e74c3c; margin: 0 0 20px 0; }
            p { color: #555; line-height: 1.6; }
            a { display: inline-block; margin-top: 20px; padding: 12px 30px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; text-decoration: none; border-radius: 5px; }
            </style></head><body>
            <div class="box">
                <h1>🔒 Accès Refusé</h1>
                <p>L\'accès direct à Roundcube est désactivé.</p>
                <p>Veuillez vous connecter via <strong>Office1789</strong>.</p>
                <a href="http://localhost:5173">Retour à Office1789</a>
            </div>
            </body></html>';
        exit;
    }
}

// Si token présent = AUTO-LOGIN
if (!empty($sso_token)) {
    error_log('[SSO] Token détecté: ' . substr($sso_token, 0, 20) . '... (length: ' . strlen($sso_token) . ')');
    
    $claims = validate_sso_token($sso_token, $sso_secret);
    
    if ($claims) {
        error_log('[SSO] Token valide - Claims: ' . json_encode(array_keys($claims)));
        
        if (isset($claims['email']) && isset($claims['password'])) {
            // Vérifier expiration
            if (isset($claims['exp']) && $claims['exp'] > time()) {
                error_log('[SSO] Token non expiré pour ' . $claims['email'] . ' (nonce: ' . ($claims['nonce'] ?? 'N/A') . ')');
                
                // IMPORTANT: Forcer un logout propre en supprimant tous les cookies de session Roundcube
                // Cela évite les conflits avec d'anciennes sessions existantes
                if (isset($_COOKIE['roundcube_sessid'])) {
                    setcookie('roundcube_sessid', '', time() - 3600, '/');
                    unset($_COOKIE['roundcube_sessid']);
                    error_log('[SSO] Cookie session Roundcube supprimé pour nouveau login propre');
                }
                
                // À ce stade, $RCMAIL (instance rcmail) est disponible car ce fichier
                // est inclus APRÈS iniset.php
                // Générer un token CSRF valide via l'API Roundcube
                
                $csrf_token = '';
                if (class_exists('rcmail') && method_exists('rcmail', 'get_instance')) {
                    $rcmail = rcmail::get_instance();
                    if ($rcmail && method_exists($rcmail, 'get_request_token')) {
                        $csrf_token = $rcmail->get_request_token();
                        error_log('[SSO] Token CSRF généré: ' . substr($csrf_token, 0, 10) . '...');
                    }
                }
                
                // Si pas de token CSRF, utiliser rcube_utils
                if (empty($csrf_token) && class_exists('rcube_utils')) {
                    $csrf_token = rcube_utils::request_token();
                    error_log('[SSO] Token CSRF via rcube_utils: ' . substr($csrf_token, 0, 10) . '...');
                }
                
                error_log('[SSO] Affichage formulaire auto-submit pour ' . $claims['email']);
            
            echo '<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Connexion SSO Office1789...</title>
    <style>
        body { 
            font-family: Arial, sans-serif; 
            display: flex; 
            justify-content: center; 
            align-items: center; 
            height: 100vh; 
            margin: 0; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); 
        }
        .loading {
            background: white;
            padding: 40px;
            border-radius: 10px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.3);
            text-align: center;
        }
        .spinner {
            border: 4px solid #f3f3f3;
            border-top: 4px solid #667eea;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 0 auto 20px;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="loading">
        <div class="spinner"></div>
        <h2>🔐 Connexion automatique en cours...</h2>
        <p>Authentification de ' . htmlspecialchars($claims['email'], ENT_QUOTES, 'UTF-8') . '</p>
    </div>
    <form id="sso_form" method="POST" action="/" style="display:none;">
        <input type="hidden" name="_task" value="login" />
        <input type="hidden" name="_action" value="login" />
        <input type="hidden" name="_token" value="' . htmlspecialchars($csrf_token, ENT_QUOTES, 'UTF-8') . '" />
        <input type="hidden" name="_user" value="' . htmlspecialchars($claims['email'], ENT_QUOTES, 'UTF-8') . '" />
        <input type="hidden" name="_pass" value="' . htmlspecialchars($claims['password'], ENT_QUOTES, 'UTF-8') . '" />
        <input type="hidden" name="_timezone" value="Europe/Paris" />
    </form>
    <script>
        // Soumettre le formulaire après un court délai
        setTimeout(function() {
            document.getElementById("sso_form").submit();
        }, 100);
    </script>
</body>
</html>';
            exit;
        } else {
            error_log('[SSO] Token expiré');
            header('HTTP/1.1 403 Forbidden');
            echo 'Token SSO expiré. Reconnectez-vous via Office1789.';
            exit;
        }
        } else {
            error_log('[SSO] Token invalide');
            header('HTTP/1.1 403 Forbidden');
            echo 'Token SSO invalide. Reconnectez-vous via Office1789.';
            exit;
        }
    }
}
