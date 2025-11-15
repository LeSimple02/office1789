<?php
/**
 * Office1789 SSO Plugin
 * 
 * Permet l'authentification automatique via token SSO généré par le backend Go
 */

class office1789_sso extends rcube_plugin
{
    public $task = '.*';
    private $secret = 'Office1789-SecretKey-ChangeInProduction'; // Même secret que backend Go

    function init()
    {
        error_log('[SSO] Plugin: Initialisation');
        $this->add_hook('startup', array($this, 'startup'));
        $this->add_hook('authenticate', array($this, 'authenticate'));
        $this->add_hook('login_after', array($this, 'login_after'));
    }

    function startup($args)
    {
        $rcmail = rcmail::get_instance();
        
        // Récupérer le token SSO depuis l'URL
        $sso_token = rcube_utils::get_input_value('sso_token', rcube_utils::INPUT_GPC);
        
        // SÉCURITÉ : Bloquer TOUT accès sans token SSO valide (sauf si déjà authentifié)
        if (empty($sso_token) && !isset($_SESSION['user_id'])) {
            error_log('[SSO] SÉCURITÉ - Accès refusé sans token SSO valide');
            
            // Afficher une page d'erreur propre
            header('HTTP/1.1 403 Forbidden');
            echo '<!DOCTYPE html><html><head><title>Accès refusé</title><style>
                body { font-family: Arial, sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
                .error-box { background: white; padding: 40px; border-radius: 10px; box-shadow: 0 10px 40px rgba(0,0,0,0.3); text-align: center; max-width: 500px; }
                h1 { color: #e74c3c; margin: 0 0 20px 0; }
                p { color: #555; line-height: 1.6; }
                a { display: inline-block; margin-top: 20px; padding: 12px 30px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; text-decoration: none; border-radius: 5px; }
                </style></head><body>
                <div class="error-box">
                    <h1>🔒 Accès Refusé</h1>
                    <p>L\'accès direct à Roundcube est désactivé pour des raisons de sécurité.</p>
                    <p>Veuillez vous connecter via <strong>Office1789</strong> pour accéder à votre messagerie.</p>
                    <a href="http://localhost:5173">Retour à Office1789</a>
                </div>
                </body></html>';
            exit;
        }
        
        // Si token présent ET pas encore authentifié
        if (!empty($sso_token) && !isset($_SESSION['user_id'])) {
            error_log('[SSO] Token détecté: ' . substr($sso_token, 0, 20) . '...');
            
            // Valider et décoder le token
            $claims = $this->validate_token($sso_token);
            
            if ($claims && isset($claims['username']) && isset($claims['email'])) {
                error_log('[SSO] Claims décodés pour ' . $claims['email']);
                error_log('[SSO] Claims complets: ' . json_encode($claims));
                
                // Vérifier l'expiration
                if (isset($claims['exp']) && $claims['exp'] > time()) {
                    error_log('[SSO] Token valide, expiration OK (exp: ' . $claims['exp'] . ', now: ' . time() . ')');
                    
                    // Stocker dans session TEMPORAIRE pour authenticate
                    // Utiliser l'email complet et le mot de passe fourni dans le token
                    $password = isset($claims['password']) ? $claims['password'] : '';
                    
                    $_SESSION['temp_sso_login'] = array(
                        'user' => $claims['email'],
                        'pass' => $password,
                        'host' => 'ssl://mailserver:993'
                    );
                    
                    error_log('[SSO] Session SSO créée, affichage formulaire auto-submit');
                    
                    // Afficher un formulaire qui se soumet automatiquement
                    echo '<!DOCTYPE html>
<html>
<head>
    <title>Connexion SSO Office1789...</title>
</head>
<body>
    <form id="sso_form" method="POST" action="/">
        <input type="hidden" name="_task" value="login" />
        <input type="hidden" name="_action" value="login" />
        <input type="hidden" name="_user" value="' . htmlspecialchars($claims['email']) . '" />
        <input type="hidden" name="_pass" value="' . htmlspecialchars($password) . '" />
        <input type="hidden" name="_host" value="ssl://mailserver:993" />
        <p style="text-align:center; font-family:Arial; margin-top:100px;">
            🔐 Connexion automatique en cours...<br><br>
            <span style="color:#666;">Si la page ne se charge pas, <a href="#" onclick="document.getElementById(\'sso_form\').submit(); return false;">cliquez ici</a></span>
        </p>
    </form>
    <script>
        document.getElementById("sso_form").submit();
    </script>
</body>
</html>';
                    exit;
                } else {
                    error_log('[SSO] Token expiré (exp: ' . ($claims['exp'] ?? 'N/A') . ', now: ' . time() . ')');
                    header('HTTP/1.1 403 Forbidden');
                    echo 'Token SSO expiré. Veuillez vous reconnecter via Office1789.';
                    exit;
                }
            } else {
                error_log('[SSO] Token invalide - validation échouée');
                header('HTTP/1.1 403 Forbidden');
                echo 'Token SSO invalide. Veuillez vous reconnecter via Office1789.';
                exit;
            }
        } else if (!empty($sso_token) && isset($_SESSION['user_id'])) {
            error_log('[SSO] Token présent mais utilisateur déjà authentifié (user_id: ' . $_SESSION['user_id'] . ')');
        }
        
        return $args;
    }

    function authenticate($args)
    {
        error_log('[SSO] Hook authenticate appelé');
        
        // Si c'est une connexion SSO (données stockées temporairement)
        if (isset($_SESSION['temp_sso_login'])) {
            $sso_data = $_SESSION['temp_sso_login'];
            
            error_log('[SSO] Authentification SSO pour ' . $sso_data['user']);
            
            $args['user'] = $sso_data['user'];
            $args['pass'] = $sso_data['pass'];
            $args['host'] = $sso_data['host'];
            $args['valid'] = true;
            $args['cookiecheck'] = false;
            
            // Nettoyer immédiatement
            unset($_SESSION['temp_sso_login']);
            
            error_log('[SSO] Auto-login configuré et session nettoyée');
        } else {
            error_log('[SSO] Authentification normale (pas de SSO)');
        }
        
        return $args;
    }
    
    function login_after($args)
    {
        error_log('[SSO] Hook login_after - Login réussi pour user_id: ' . (isset($_SESSION['user_id']) ? $_SESSION['user_id'] : 'N/A'));
        return $args;
    }

    private function validate_token($token)
    {
        // Format: claims.signature
        $parts = explode('.', $token);
        if (count($parts) !== 2) {
            return false;
        }
        
        list($claimsB64, $signature) = $parts;
        
        // Vérifier la signature
        $expected_signature = base64_encode(hash_hmac('sha256', $claimsB64, $this->secret, true));
        $expected_signature = str_replace(['+', '/'], ['-', '_'], rtrim($expected_signature, '='));
        
        if ($signature !== $expected_signature) {
            error_log('[SSO] Invalid signature - Expected: ' . $expected_signature . ', Got: ' . $signature);
            return false;
        }
        
        // Décoder les claims
        $claimsJSON = base64_decode(str_replace(['-', '_'], ['+', '/'], $claimsB64));
        $claims = json_decode($claimsJSON, true);
        
        if (!$claims) {
            error_log('[SSO] Invalid claims JSON: ' . $claimsJSON);
            return false;
        }
        
        return $claims;
    }
}
