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
    // Configuration Docker (copié depuis config.docker.inc.php pour éviter les problèmes de scope)
    $config['db_dsnw'] = 'sqlite:////var/roundcube/db/sqlite.db?mode=0646';
    $config['db_dsnr'] = '';
    $config['imap_host'] = 'localhost:143';
    $config['smtp_host'] = 'localhost:587';
    $config['temp_dir'] = '/tmp/roundcube-temp';
    $config['skin'] = 'elastic';
    
    // BLOQUER le changement de mot de passe dans Roundcube
    // Les utilisateurs doivent utiliser Office1789 pour synchroniser tous les services
    $config['dont_override'] = ['password'];
    
    // Plugins Office1789 - DÉFINIS EN DERNIER APRÈS TOUTES LES AUTRES CONFIGS
    $config['plugins'] = ['office1789_sso', 'office1789_darkmode', 'archive', 'zipdownload'];
