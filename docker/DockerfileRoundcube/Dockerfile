FROM roundcube/roundcubemail:1.6.0-apache

# Copier le script d'entrypoint personnalisé
COPY entrypoint-wrapper.sh /entrypoint-wrapper.sh
RUN chmod +x /entrypoint-wrapper.sh

# Utiliser notre entrypoint qui wrap l'original
ENTRYPOINT ["/entrypoint-wrapper.sh"]
CMD ["apache2-foreground"]
