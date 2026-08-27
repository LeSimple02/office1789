<?php
    $config['plugins'] = [];
    $config['log_driver'] = 'stdout';
    $config['zipdownload_selection'] = true;
    $config['des_key'] = 'JIWIG94wMqHgAmFGK1wwoUzh';
    $config['enable_spellcheck'] = true;
    $config['spellcheck_engine'] = 'pspell';
    include(__DIR__ . '/config.docker.inc.php');
    $config['imap_conn_options'] = [
        'ssl' => [
            'verify_peer' => false,
            'verify_peer_name' => false,
            'allow_self_signed' => true,
        ],
    ];
    $config['smtp_conn_options'] = [
        'ssl' => [
            'verify_peer' => false,
            'verify_peer_name' => false,
            'allow_self_signed' => true,
        ],
    ];
    
