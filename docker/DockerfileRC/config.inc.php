<?php
    // Configuration Office1789
    global $config; // FORCE GLOBAL SCOPE
    if (!isset($config) || !is_array($config)) {
        $config = array();
    }
    
    $config['log_driver'] = 'stdout';
    $config['zipdownload_selection'] = true;
    $config['des_key'] = '1Zpms5JFDU17jVORhV0xWkGB';
    $config['enable_spellcheck'] = true;
    $config['spellcheck_engine'] = 'pspell';
    
    // Plugins Office1789 - DÉFINIS DIRECTEMENT ICI
    // Note: Ne PAS utiliser include() car le scope ne fonctionne pas avec Roundcube
    // Configuration Docker - Base de données PostgreSQL
    $config['db_dsnw'] = 'pgsql://roundcube:roundcube1789@postgres_roundcube/roundcube';
    $config['db_dsnr'] = '';
    
    // Configuration IMAP - Port 993 avec SSL/TLS (IMAPS)
    $config['imap_host'] = 'ssl://mailserver:993';
    $config['imap_conn_options'] = array(
        'ssl' => array(
            'verify_peer'       => false,
            'verify_peer_name'  => false,
            'allow_self_signed' => true,
        ),
    );
    $config['imap_timeout'] = 30;
    
    // Configuration SMTP - Port 587 avec STARTTLS
    $config['smtp_host'] = 'tls://mailserver:587';
    $config['smtp_conn_options'] = array(
        'ssl' => array(
            'verify_peer'       => false,
            'verify_peer_name'  => false,
            'allow_self_signed' => true,
        ),
    );
    
    $config['temp_dir'] = '/tmp/roundcube-temp';
    $config['skin'] = 'elastic';
    
    // BLOQUER le changement de mot de passe dans Roundcube
    // Les utilisateurs doivent utiliser Office1789 pour synchroniser tous les services
    $config['dont_override'] = ['password'];
    
    // Plugins Office1789 - DÉFINIS EN DERNIER APRÈS TOUTES LES AUTRES CONFIGS
    $config['plugins'] = ['office1789_sso', 'office1789_darkmode', 'archive', 'zipdownload'];
