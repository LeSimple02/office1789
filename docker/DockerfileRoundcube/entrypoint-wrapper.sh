#!/bin/sh
set -e

# Run the original entrypoint
/docker-entrypoint.sh "$@" &
ENTRYPOINT_PID=$!

# Wait for config files to be generated
sleep 3

# Inject IMAP and SMTP connection options to disable certificate verification
if [ -f /var/www/html/config/config.inc.php ]; then
    cat >> /var/www/html/config/config.inc.php <<'PHPEOF'

// Custom: Disable SSL verification for self-signed certs
$config['imap_conn_options'] = array(
    'ssl' => array(
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ),
);
$config['smtp_conn_options'] = array(
    'ssl' => array(
        'verify_peer' => false,
        'verify_peer_name' => false,
        'allow_self_signed' => true,
    ),
);
PHPEOF
    echo "✓ Injected SSL conn_options into Roundcube config"
fi

# Wait for the original entrypoint to finish
wait $ENTRYPOINT_PID
