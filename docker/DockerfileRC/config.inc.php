<?php

/* Local configuration for Roundcube Webmail */

// ----------------------------------
// IMAP
// ----------------------------------
// The IMAP host (and optionally port number) chosen to perform the log-in.
// Leave blank to show a textbox at login, give a list of hosts
// to display a pulldown menu or set one host as string.
// Enter hostname with prefix ssl:// to use Implicit TLS, or use
// prefix tls:// to use STARTTLS.
// If port number is omitted it will be set to 993 (for ssl://) or 143 otherwise.
// Supported replacement variables:
// %n - hostname ($_SERVER['SERVER_NAME'])
// %t - hostname without the first part
// %d - domain (http hostname $_SERVER['HTTP_HOST'] without the first part)
// %s - domain name after the '@' from e-mail address provided at login screen
// For example %n = mail.domain.tld, %t = domain.tld
// WARNING: After hostname change update of mail_host column in users table is
//          required to match old user data records with the new host.
$config['imap_host'] = 'mailserver:143';

// ----------------------------------
// SQL DATABASE
// ----------------------------------
// Database connection string (DSN) for read+write operations
// Format (compatible with PEAR MDB2): db_provider://user:password@host/database
// Currently supported db_providers: mysql, pgsql, sqlite, mssql, sqlsrv, oracle
// For examples see http://pear.php.net/manual/en/package.database.mdb2.intro-dsn.php
// Note: for SQLite use absolute path (Linux): 'sqlite:////full/path/to/sqlite.db?mode=0646'
//       or (Windows): 'sqlite:///C:/full/path/to/sqlite.db'
// Note: Various drivers support various additional arguments for connection,
//       for Mysql: key, cipher, cert, capath, ca, verify_server_cert,
//       for Postgres: application_name, sslmode, sslcert, sslkey, sslrootcert, sslcrl, sslcompression, service.
//       e.g. 'mysql://roundcube:@localhost/roundcubemail?verify_server_cert=false'
$config['db_dsnw'] = 'pgsql://robespierre:guillotine@postgres_db:5432/office1789';

// you can define specific table (and sequence) names prefix
$config['db_prefix'] = 'roundcube_';

// ----------------------------------
// LOGGING/DEBUGGING
// ----------------------------------
// log driver:  'syslog', 'stdout' or 'file'.
$config['log_driver'] = 'stdout';

// Log sent messages to <log_dir>/sendmail.log or to syslog
$config['smtp_log'] = false;

// Log SQL queries to <log_dir>/sql.log or to syslog
$config['sql_debug'] = true;

// Log IMAP conversation to <log_dir>/imap.log or to syslog
$config['imap_debug'] = true;

// Log SMTP conversation to <log_dir>/smtp.log or to syslog
$config['smtp_debug'] = true;

// IMAP socket context options
// See http://php.net/manual/en/context.ssl.php
// The example below enables server certificate validation
//$config['imap_conn_options'] = [
//  'ssl'         => [
//     'verify_peer'  => true,
//     'verify_depth' => 3,
//     'cafile'       => '/etc/openssl/certs/ca.crt',
//   ],
// ];
// Note: These can be also specified as an array of options indexed by hostname
$config['imap_conn_options'] = array (
  'ssl' => 
  array (
    'verify_peer' => false,
    'verify_peer_name' => false,
    'allow_self_signed' => true,
  ),
);

// SMTP socket context options
// See http://php.net/manual/en/context.ssl.php
// The example below enables server certificate validation, and
// requires 'smtp_timeout' to be non zero.
// $config['smtp_conn_options'] = [
//   'ssl'         => [
//     'verify_peer'  => true,
//     'verify_depth' => 3,
//     'cafile'       => '/etc/openssl/certs/ca.crt',
//   ],
// ];
// Note: These can be also specified as an array of options indexed by hostname
$config['smtp_conn_options'] = array (
  'ssl' => 
  array (
    'verify_peer' => false,
    'verify_peer_name' => false,
    'allow_self_signed' => true,
  ),
);

// provide an URL where a user can get support for this Roundcube installation
// PLEASE DO NOT LINK TO THE ROUNDCUBE.NET WEBSITE HERE!
$config['support_url'] = '';

// Automatically register user in Roundcube database on successful (IMAP) logon.
// Set to false if only registered users should be allowed to the webmail.
// Note: If disabled you have to create records in Roundcube users table by yourself.
// Note: Roundcube does not manage/create users on a mail server.
$config['auto_create_user'] = false;

// This key is used for encrypting purposes, like storing of imap password
// in the session. For historical reasons it's called DES_key, but it's used
// with any configured cipher_method (see below).
// For the default cipher_method a required key length is 24 characters.
$config['des_key'] = '9wS3sh89WayiMORAxvYj5v0e';

// Automatically add this domain to user names for login
// Only for IMAP servers that require full e-mail addresses for login
// Specify an array with 'host' => 'domain' values to support multiple hosts
// Supported replacement variables:
// %h - user's IMAP hostname
// %n - hostname ($_SERVER['SERVER_NAME'])
// %t - hostname without the first part
// %d - domain (http hostname $_SERVER['HTTP_HOST'] without the first part)
// %z - IMAP domain (IMAP hostname without the first part)
// For example %n = mail.domain.tld, %t = domain.tld
$config['username_domain'] = 'office1789.local';

// Name your service. This is displayed on the login screen and in the window title
$config['product_name'] = 'Office1789 Mail';

// ----------------------------------
// PLUGINS
// ----------------------------------
// List of active plugins (in plugins/ directory)
$config['plugins'] = ['office1789_sso', 'archive', 'zipdownload'];

// the default locale setting (leave empty for auto-detection)
// RFC1766 formatted language name like en_US, de_DE, de_CH, fr_FR, pt_BR
$config['language'] = 'fr_FR';

// Set the spell checking engine. Possible values:
// - 'googie'  - the default (also used for connecting to Nox Spell Server, see 'spellcheck_uri' setting)
// - 'pspell'  - requires the PHP Pspell module and aspell installed
// - 'enchant' - requires the PHP Enchant module
// - 'atd'     - install your own After the Deadline server or check with the people at http://www.afterthedeadline.com before using their API
// Since Google shut down their public spell checking service, the default settings
// connect to http://spell.roundcube.net which is a hosted service provided by Roundcube.
// You can connect to any other googie-compliant service by setting 'spellcheck_uri' accordingly.
$config['spellcheck_engine'] = 'pspell';

// prefer displaying HTML messages
$config['prefer_html'] = false;

$config['zipdownload_selection'] = true;

$config['default_port'] = 143;

$config['smtp_port'] = 587;

$config['additional_header'] = '<link rel="stylesheet" href="/skins/elastic/styles/office1789-custom.css"><script src="/office1789-dark-sync.js"></script>';

