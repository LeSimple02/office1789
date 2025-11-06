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
        
        // Si token présent ET pas encore authentifié
        if (!empty($sso_token) && !isset($_SESSION['user_id'])) {
            error_log('[SSO] Token détecté: ' . substr($sso_token, 0, 20) . '...');
            
            // Valider et décoder le token
            $claims = $this->validate_token($sso_token);
            
            if ($claims && isset($claims['username']) && isset($claims['email'])) {
                error_log('[SSO] Claims décodés pour ' . $claims['email']);
                
                // Vérifier l'expiration
                if (isset($claims['exp']) && $claims['exp'] > time()) {
                    error_log('[SSO] Token valide, expiration OK');
                    
                    // Stocker dans session TEMPORAIRE pour authenticate
                    $_SESSION['temp_sso_login'] = array(
                        'user' => $claims['email'],
                        'pass' => 'password123',
                        'host' => 'mailserver'
                    );
                    
                    // Préparer les variables POST pour Roundcube
                    $_POST['_task'] = 'login';
                    $_POST['_action'] = 'login';
                    $_POST['_user'] = $claims['email'];
                    $_POST['_pass'] = 'password123';
                    $_POST['_host'] = 'mailserver';
                    
                    // Forcer task=login et action=login
                    $args['task'] = 'login';
                    $args['action'] = 'login';
                    
                    error_log('[SSO] Login préparé pour ' . $claims['email']);
                } else {
                    error_log('[SSO] Token expiré (exp: ' . $claims['exp'] . ', now: ' . time() . ')');
                }
            } else {
                error_log('[SSO] Token invalide - validation échouée');
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
