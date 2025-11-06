<?php
/**
 * Office1789 SSO Plugin
 * Authentification automatique via token SSO du backend Go
 */

class office1789_sso extends rcube_plugin
{
    public $task = '.*'; // S'exécuter sur toutes les tâches
    private $secret = 'Office1789-SecretKey-ChangeInProduction';

    function init()
    {
        error_log('[SSO] Plugin: Initialisation du plugin office1789_sso');
        $this->add_hook('startup', array($this, 'startup'));
        $this->add_hook('authenticate', array($this, 'authenticate'));
        $this->add_hook('login_after', array($this, 'login_after'));
    }

    function startup($args)
    {
        error_log('[SSO] Hook startup appelé - Task: ' . $args['task']);
        
        $rcmail = rcmail::get_instance();
        
        // Si déjà connecté, ne rien faire
        if ($rcmail->user->ID) {
            error_log('[SSO] Utilisateur déjà connecté (ID: ' . $rcmail->user->ID . ')');
            return $args;
        }

        // Récupérer le token SSO depuis l'URL
        $sso_token = rcube_utils::get_input_value('sso_token', rcube_utils::INPUT_GPC);
        
        if (!empty($sso_token)) {
            error_log('[SSO] Token détecté: ' . substr($sso_token, 0, 20) . '...');
            
            // Valider et décoder le token
            $claims = $this->validate_token($sso_token);
            
            if ($claims && isset($claims['username']) && isset($claims['email'])) {
                error_log('[SSO] Claims décodés - username: ' . $claims['username'] . ', email: ' . $claims['email']);
                
                // Vérifier l'expiration
                if (isset($claims['exp']) && $claims['exp'] > time()) {
                    error_log('[SSO] Token valide, expiration: ' . date('Y-m-d H:i:s', $claims['exp']));
                    
                    // Stocker dans session temporaire
                    $_SESSION['temp_sso_login'] = array(
                        'user' => $claims['email'],
                        'host' => 'mailserver:143'
                    );
                    
                    // Préparer la connexion
                    $_POST['_user'] = $claims['email'];
                    $_POST['_pass'] = 'password123'; // Mot de passe IMAP réel
                    $_POST['_task'] = 'mail';
                    $_POST['_action'] = 'login';
                    
                    $args['task'] = 'login';
                    $args['action'] = 'login';
                    
                    error_log('[SSO] Login préparé pour: ' . $claims['email']);
                } else {
                    error_log('[SSO] Token expiré');
                }
            } else {
                error_log('[SSO] Claims invalides ou incomplets');
            }
        }
        
        return $args;
    }

    function authenticate($args)
    {
        error_log('[SSO] Hook authenticate appelé - User: ' . ($args['user'] ?? 'N/A'));
        
        // Si connexion SSO en cours
        if (isset($_SESSION['temp_sso_login'])) {
            $sso_data = $_SESSION['temp_sso_login'];
            error_log('[SSO] Authentification SSO pour: ' . $sso_data['user']);
            
            $args['user'] = $sso_data['user'];
            $args['pass'] = 'password123'; // Mot de passe IMAP réel
            $args['host'] = $sso_data['host'];
            $args['valid'] = true;
            $args['cookiecheck'] = false;
            
            // Nettoyer la session temporaire
            unset($_SESSION['temp_sso_login']);
            
            error_log('[SSO] Authentification SSO validée');
        }
        
        return $args;
    }

    function login_after($args)
    {
        error_log('[SSO] Hook login_after - Login réussi pour user_id: ' . ($args['user_id'] ?? 'N/A'));
        return $args;
    }

    private function validate_token($token)
    {
        // Format: base64(claims).base64(signature)
        $parts = explode('.', $token);
        if (count($parts) !== 2) {
            error_log('[SSO] Format de token invalide (pas 2 parties)');
            return false;
        }

        list($claimsB64, $signatureB64) = $parts;

        // Vérifier la signature HMAC
        $expected_signature = base64_encode(hash_hmac('sha256', $claimsB64, $this->secret, true));
        $expected_signature = str_replace(['+', '/', '='], ['-', '_', ''], $expected_signature);

        if ($signatureB64 !== $expected_signature) {
            error_log('[SSO] Signature invalide');
            error_log('[SSO] Attendu: ' . $expected_signature);
            error_log('[SSO] Reçu: ' . $signatureB64);
            return false;
        }

        // Décoder les claims
        $claimsJSON = base64_decode($claimsB64);
        $claims = json_decode($claimsJSON, true);

        if (!$claims) {
            error_log('[SSO] Erreur décodage JSON claims');
            return false;
        }

        return $claims;
    }
}
