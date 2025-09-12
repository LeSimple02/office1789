# Utilisation de l'image officielle PostgreSQL
FROM postgres:16

# Copie du script d'initialisation dans le répertoire d'initialisation de PostgreSQL
COPY init.sql /docker-entrypoint-initdb.d/init.sql

# Définir les permissions (optionnel mais recommandé)
RUN chmod 755 /docker-entrypoint-initdb.d/init.sql
