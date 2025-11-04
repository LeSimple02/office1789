<?php

/* Local configuration for Roundcube Webmail */

// ----------------------------------
// IMAP
// ----------------------------------
$config['imap_host'] = getenv('ROUNDCUBE_DEFAULT_HOST') ?: 'localhost:143';
$config['imap_auth_type'] = null;
$config['imap_delimiter'] = null;

// ----------------------------------
// SMTP
// ----------------------------------
$config['smtp_host'] = getenv('ROUNDCUBE_SMTP_SERVER') ?: 'localhost:587';
$config['smtp_user'] = '%u';
$config['smtp_pass'] = '%p';
$config['smtp_auth_type'] = null;

// ----------------------------------
// System
// ----------------------------------
$config['des_key'] = 'office1789_roundcube_des_key_change_me';
$config['product_name'] = 'Office1789 Mail';
$config['useragent'] = 'Office1789 Roundcube Webmail';
$config['support_url'] = '';

// ----------------------------------
// Plugins
// ----------------------------------
$config['plugins'] = array_filter(array_unique(array_merge(
    $config['plugins'],
    ['archive', 'zipdownload']
)));

// ----------------------------------
// Security - IMPORTANT for iframe
// ----------------------------------
// Autoriser l'embedding dans un iframe depuis votre application
$config['x_frame_options'] = 'ALLOW-FROM http://localhost:5173';

// Alternative plus permissive (à utiliser en développement uniquement)
// $config['x_frame_options'] = 'SAMEORIGIN';

// Désactiver la protection CSRF pour l'iframe (à utiliser avec précaution)
$config['session_samesite'] = 'None';
$config['session_secure'] = false; // Mettre à true en production avec HTTPS

// Autoriser les cookies tiers (nécessaire pour l'iframe)
$config['session_domain'] = '';

// ----------------------------------
// Database
// ----------------------------------
$config['db_dsnw'] = sprintf(
    '%s://%s:%s@%s:%s/%s',
    getenv('ROUNDCUBE_DB_TYPE') ?: 'pgsql',
    getenv('ROUNDCUBE_DB_USER') ?: 'robespierre',
    getenv('ROUNDCUBE_DB_PASSWORD') ?: 'guillotine',
    getenv('ROUNDCUBE_DB_HOST') ?: 'postgres_db',
    getenv('ROUNDCUBE_DB_PORT') ?: '5432',
    getenv('ROUNDCUBE_DB_NAME') ?: 'office1789'
);

// ----------------------------------
// Logging/Debugging
// ----------------------------------
$config['log_driver'] = 'file';
$config['log_dir'] = '/var/log/roundcube/';
$config['session_debug'] = false;
$config['sql_debug'] = false;
$config['imap_debug'] = false;
$config['smtp_debug'] = false;

// ----------------------------------
// Interface
// ----------------------------------
$config['language'] = 'fr_FR';
$config['enable_installer'] = false;
$config['auto_create_user'] = true;

