#!/bin/bash
# Injecter le CSS Office1789 dans les templates Roundcube

CSS_LINK='<link rel="stylesheet" type="text/css" href="/skins/elastic/styles/office1789-custom.css">'

# Injecter dans mail.html
if [ -f "/var/www/html/skins/elastic/templates/mail.html" ]; then
    if ! grep -q "office1789-custom.css" "/var/www/html/skins/elastic/templates/mail.html"; then
        sed -i "s|</head>|${CSS_LINK}\n</head>|" /var/www/html/skins/elastic/templates/mail.html
        echo "CSS injecté dans mail.html"
    fi
fi

# Injecter dans login.html
if [ -f "/var/www/html/skins/elastic/templates/login.html" ]; then
    if ! grep -q "office1789-custom.css" "/var/www/html/skins/elastic/templates/login.html"; then
        sed -i "s|</head>|${CSS_LINK}\n</head>|" /var/www/html/skins/elastic/templates/login.html
        echo "CSS injecté dans login.html"
    fi
fi

echo "Injection CSS terminée"
