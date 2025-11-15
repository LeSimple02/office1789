#!/bin/bash
# Patch pour bypass la vérification CSRF pour les logins SSO

# Fichier à patcher
LOGIN_FILE="/var/www/html/program/actions/login.php"

if [ -f "$LOGIN_FILE" ]; then
    # Vérifier si le patch n'est pas déjà appliqué
    if ! grep -q "HTTP_X_SSO_LOGIN" "$LOGIN_FILE"; then
        # Chercher la ligne qui contient "check_request" et ajouter le bypass avant
        sed -i "/rcmail::check_request/i\\
        // Office1789 SSO: Bypass CSRF check for SSO logins\\
        if (isset(\$_SERVER['HTTP_X_SSO_LOGIN']) && \$_SERVER['HTTP_X_SSO_LOGIN'] === 'true') {\\
            error_log('[SSO-Patch] Bypass CSRF check for SSO login');\\
        } else {" "$LOGIN_FILE"
        
        # Ajouter la fermeture du else après le check_request
        sed -i "/rcmail::check_request.*(/a\\
        }" "$LOGIN_FILE"
        
        echo "Patch CSRF SSO appliqué à login.php"
    else
        echo "Patch CSRF SSO déjà appliqué"
    fi
fi
