<?php
/**
 * Office1789 SSO Plugin
 * 
 * Permet l'authentification automatique via token SSO généré par le backend Go
 */

class office1789_sso extends rcube_plugin
{
    public $task = 'login|mail';
    private $secret = 'Office1789-SecretKey-ChangeInProduction'; // Même secret que backend Go

    function init()
    {
        rcube::write_log('console', 'SSO Plugin: Initialisation');
        $this->add_hook('startup', array($this, 'startup'));
        $this->add_hook('authenticate', array($this, 'authenticate'));
    }

    function startup($args)
    {
        $rcmail = rcmail::get_instance();
        
        // Récupérer le token SSO depuis l'URL
        $sso_token = rcube_utils::get_input_value('sso_token', rcube_utils::INPUT_GPC);
        
        // Si token présent ET pas encore authentifié
        if (!empty($sso_token) && !$_SESSION['user_id']) {
            rcube::write_log('console', 'SSO: Token détecté: ' . substr($sso_token, 0, 20) . '...');
            
            // Valider et décoder le token
            $claims = $this->validate_token($sso_token);
            
            if ($claims && isset($claims['username']) && isset($claims['email'])) {
                rcube::write_log('console', 'SSO: Claims décodés pour ' . $claims['email']);
                
                // Vérifier l'expiration
                if (isset($claims['exp']) && $claims['exp'] > time()) {
                    rcube::write_log('console', 'SSO: Token valide, expiration OK');
                    
                    // Stocker dans session
                    $_SESSION['sso_validated'] = true;
                    $_SESSION['sso_email'] = $claims['email'];
                    $_SESSION['sso_pass'] = 'password123';
                    $_SESSION['sso_host'] = 'mailserver';
                    
                    // Forcer le login
                    $_POST['_user'] = $claims['email'];
                    $_POST['_pass'] = 'password123';
                    $_POST['_host'] = 'mailserver';
                    $_POST['_task'] = 'login';
                    $_POST['_action'] = 'login';
                    
                    $args['task'] = 'login';
                    $args['action'] = 'login';
                    
                    rcube::write_log('console', 'SSO: Variables POST définies, déclenchement login');
                } else {
                    rcube::write_log('errors', 'SSO: Token expiré (exp: ' . $claims['exp'] . ', now: ' . time() . ')');
                }
            } else {
                rcube::write_log('errors', 'SSO: Token invalide - validation échouée');
            }
        } else if (!empty($sso_token)) {
            rcube::write_log('console', 'SSO: Token présent mais utilisateur déjà authentifié');
        }
        
        return $args;
    }

    function authenticate($args)
    {
        rcube::write_log('console', 'SSO: Hook authenticate appelé');
        rcube::write_log('console', 'SSO: Session sso_validated = ' . (isset($_SESSION['sso_validated']) ? 'true' : 'false'));
        
        // Si c'est une connexion SSO validée
        if (isset($_SESSION['sso_validated']) && $_SESSION['sso_validated'] === true) {
            rcube::write_log('console', 'SSO: Authentification SSO détectée');
            
            $args['user'] = $_SESSION['sso_email'];
            $args['pass'] = $_SESSION['sso_pass'];
            $args['host'] = $_SESSION['sso_host'];
            $args['valid'] = true;
            $args['cookiecheck'] = false;
            
            rcube::write_log('console', 'SSO: Auto-login configuré pour ' . $args['user'] . ' sur ' . $args['host']);
            
            // Nettoyer la session SSO
            unset($_SESSION['sso_validated']);
            unset($_SESSION['sso_email']);
            unset($_SESSION['sso_pass']);
            unset($_SESSION['sso_host']);
        } else {
            rcube::write_log('console', 'SSO: Pas d\'authentification SSO, passage normal');
        }
        
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
            rcube::write_log('errors', 'SSO: Invalid signature');
            return false;
        }
        
        // Décoder les claims
        $claimsJSON = base64_decode(str_replace(['-', '_'], ['+', '/'], $claimsB64));
        $claims = json_decode($claimsJSON, true);
        
        if (!$claims) {
            rcube::write_log('errors', 'SSO: Invalid claims JSON');
            return false;
        }
        
        return $claims;
    }
}
