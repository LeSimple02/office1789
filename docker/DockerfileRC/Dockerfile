FROM roundcube/roundcubemail:1.6.0-apache

ENV ROUNDCUBE_DB_TYPE=pgsql \
    ROUNDCUBE_DB_HOST=postgres_db \
    ROUNDCUBE_DB_PORT=5432 \
    ROUNDCUBE_DB_NAME=office1789 \
    ROUNDCUBE_DB_USER=robespierre \
    ROUNDCUBE_DB_PASSWORD=guillotine \
    ROUNDCUBE_DEFAULT_HOST=mailserver \
    ROUNDCUBE_DEFAULT_PORT=143 \
    ROUNDCUBE_SMTP_SERVER=mailserver \
    ROUNDCUBE_SMTP_PORT=587 \
    ROUNDCUBE_PLUGINS=archive,zipdownload

# Modifier les en-têtes Apache pour autoriser l'iframe depuis Office1789
RUN echo 'Header always unset X-Frame-Options' >> /etc/apache2/conf-enabled/security.conf && \
    echo "Header always set Content-Security-Policy \"frame-ancestors 'self' http://localhost:5173 http://localhost:8080;\"" >> /etc/apache2/conf-enabled/security.conf && \
    a2enmod headers

# Copier le script d'initialisation
COPY init-config.sh /usr/local/bin/init-config.sh
RUN chmod +x /usr/local/bin/init-config.sh

# Copier le plugin SSO Office1789
COPY plugins/office1789_sso /var/www/html/plugins/office1789_sso

# Copier le handler SSO et le hook login qui seront inclus dans index.php par init-config.sh
COPY sso-handler.php /var/www/html/sso-handler.php
COPY sso-login-hook.php /var/www/html/sso-login-hook.php

# Copier le CSS personnalisé Office1789
COPY office1789-custom.css /var/www/html/skins/elastic/styles/office1789-custom.css

# Créer un wrapper pour l'entrypoint qui exécute notre script APRÈS l'entrypoint officiel
RUN echo '#!/bin/bash' > /usr/local/bin/custom-entrypoint.sh && \
    echo 'exec /docker-entrypoint.sh "$@" &' >> /usr/local/bin/custom-entrypoint.sh && \
    echo 'ENTRYPOINT_PID=$!' >> /usr/local/bin/custom-entrypoint.sh && \
    echo 'sleep 15' >> /usr/local/bin/custom-entrypoint.sh && \
    echo '/usr/local/bin/init-config.sh' >> /usr/local/bin/custom-entrypoint.sh && \
    echo 'wait $ENTRYPOINT_PID' >> /usr/local/bin/custom-entrypoint.sh && \
    chmod +x /usr/local/bin/custom-entrypoint.sh

EXPOSE 80
ENTRYPOINT ["/usr/local/bin/custom-entrypoint.sh"]
CMD ["apache2-foreground"]

