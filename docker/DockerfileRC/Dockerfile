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

# Copier le fichier de config personnalisé si besoin
# COPY config.inc.php /var/www/html/config/config.inc.php

EXPOSE 80
CMD ["apache2-foreground"]

