# Utilisation de l'image officielle PostgreSQL
FROM postgres:16

ENV POSTGRES_USER=robespierre
ENV POSTGRES_PASSWORD=guillotine
ENV POSTGRES_DB=office1789

# Copie du script d'initialisation dans le répertoire d'initialisation de PostgreSQL
COPY backend/office1789-db/init.sql /docker-entrypoint-initdb.d/init.sql

# Définir les permissions (optionnel mais recommandé)
RUN chmod 755 /docker-entrypoint-initdb.d/init.sql
