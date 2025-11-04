FROM roundcube/roundcubemail:1.6.0-apache

ENV ROUNDCUBE_DB_TYPE=pgsql \
    ROUNDCUBE_DB_HOST=postgres_db \
    ROUNDCUBE_DB_PORT=5432 \
    ROUNDCUBE_DB_NAME=office1789 \
    ROUNDCUBE_DB_USER=robespierre \
    ROUNDCUBE_DB_PASSWORD=guillotine \
    ROUNDCUBE_DEFAULT_HOST=mail.localhost \
    ROUNDCUBE_SMTP_SERVER=mail.localhost \
    ROUNDCUBE_PLUGINS=archive,zipdownload

# Copier la configuration personnalisée pour autoriser l'iframe
COPY config.inc.php /var/www/html/config/config.inc.php

# Modifier les en-têtes Apache pour autoriser l'iframe
RUN echo 'Header always unset X-Frame-Options' >> /etc/apache2/conf-enabled/security.conf && \
    echo 'Header set Content-Security-Policy "frame-ancestors http://localhost:5173 http://localhost:8080;"' >> /etc/apache2/conf-enabled/security.conf && \
    a2enmod headers

EXPOSE 80
CMD ["apache2-foreground"]

