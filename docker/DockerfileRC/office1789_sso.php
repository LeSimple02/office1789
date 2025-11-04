<?php
/**
 * Office1789 SSO Plugin for Roundcube
 * 
 * This plugin allows automatic login via SSO token from Office1789 backend
 */

class office1789_sso extends rcube_plugin
{
    public $task = 'login|mail';
    private $backend_url = 'http://host.docker.internal:8080';

    function init()
    {
        $this->add_hook('startup', array($this, 'startup'));
        $this->add_hook('authenticate', array($this, 'authenticate'));
    }

    function startup($args)
    {
        $rcmail = rcmail::get_instance();

        // Check if auto-login is requested
        if ($args['task'] == 'mail' && 
            !empty($_GET['_autologin']) && 
            !empty($_GET['_token']) && 
            !empty($_GET['_user'])) {
            
            $token = rcube_utils::get_input_value('_token', rcube_utils::INPUT_GET);
            $user = rcube_utils::get_input_value('_user', rcube_utils::INPUT_GET);

            // Validate token with backend
            if ($this->validate_sso_token($token, $user)) {
                // Create session
                $_POST['_user'] = $user;
                $_POST['_pass'] = $this->get_user_password($user);
                $_POST['_task'] = 'mail';
                
                // Force authentication
                $args['task'] = 'login';
                $args['action'] = 'login';
            }
        }

        return $args;
    }

    function authenticate($args)
    {
        // Allow authentication via SSO
        if (!empty($_GET['_autologin'])) {
            $args['valid'] = true;
            $args['abort'] = false;
        }

        return $args;
    }

    private function validate_sso_token($token, $user)
    {
        // Call backend API to validate token
        $url = $this->backend_url . '/api/mail/validate-sso?token=' . urlencode($token) . '&user=' . urlencode($user);
        
        $ch = curl_init($url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch, CURLOPT_TIMEOUT, 5);
        
        $response = curl_exec($ch);
        $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        curl_close($ch);

        if ($http_code == 200) {
            $data = json_decode($response, true);
            return isset($data['valid']) && $data['valid'] === true;
        }

        return false;
    }

    private function get_user_password($user)
    {
        // For SSO, we can use a dummy password or integrate with your auth system
        // In production, you should integrate this with your actual mail server authentication
        return 'office1789_sso_' . md5($user);
    }
}
