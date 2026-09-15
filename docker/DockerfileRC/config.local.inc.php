// Configuration locale Office1789
// Ce fichier écrase les configurations par défaut

// Ajouter le plugin SSO aux plugins existants (ne pas écraser)
if (!isset($config['plugins'])) {
    $config['plugins'] = [];
}
if (!in_array('office1789_sso', $config['plugins'])) {
    array_unshift($config['plugins'], 'office1789_sso'); // Ajouter en premier
}

// Configuration du domaine email
$config['username_domain'] = 'office1789.com';

// DÉSACTIVER la vérification CSRF pour permettre le SSO auto-login
$config['request_token_time'] = 0;

// Configuration IMAP/SMTP avec STARTTLS
$config['imap_conn_options'] = array(
    'ssl' => array(
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ),
);

// Activer STARTTLS pour IMAP
$config['default_host'] = 'tls://mailserver';
$config['default_port'] = 143;

$config['smtp_conn_options'] = array(
    'ssl' => array(
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ),
);

// Debug (désactiver en production)
$config['log_logins'] = true;
$config['log_session'] = true;
$config['sql_debug'] = true;
$config['imap_debug'] = true;
$config['smtp_debug'] = true;
$config['session_debug'] = true;
// Désactiver la vérification SSL IMAP/SMTP (utile pour dev/local)
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

<?php
// Configuration locale Office1789
// Ce fichier écrase les configurations par défaut

// Ajouter le plugin SSO aux plugins existants (ne pas écraser)
if (!isset($config['plugins'])) {
    $config['plugins'] = [];
}
if (!in_array('office1789_sso', $config['plugins'])) {
    array_unshift($config['plugins'], 'office1789_sso'); // Ajouter en premier
}

// Configuration du domaine email
$config['username_domain'] = 'office1789.com';

// DÉSACTIVER la vérification CSRF pour permettre le SSO auto-login
$config['request_token_time'] = 0;

// Configuration IMAP/SMTP avec STARTTLS
$config['imap_conn_options'] = array(
    'ssl' => array(
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ),
);

// Activer STARTTLS pour IMAP
$config['default_host'] = 'tls://mailserver';
$config['default_port'] = 143;

$config['smtp_conn_options'] = array(
    'ssl' => array(
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ),
);

// Debug (désactiver en production)
$config['log_logins'] = true;
$config['log_session'] = true;
$config['sql_debug'] = true;
$config['imap_debug'] = true;
$config['smtp_debug'] = true;
$config['session_debug'] = true;
