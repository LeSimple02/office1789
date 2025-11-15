#!/bin/bash
# Script pour synchroniser les mots de passe Office1789 avec le serveur mail

# Connexion à la base Office1789
DB_USERS=$(docker exec postgres_db psql -U robespierre -d office1789 -t -c "SELECT username, Email FROM Users WHERE Email LIKE '%@office1789.com';")

echo "=== Synchronisation des comptes mail avec Office1789 ==="
echo "$DB_USERS" | while IFS='|' read -r username email; do
    username=$(echo "$username" | xargs)
    email=$(echo "$email" | xargs)
    
    if [ ! -z "$username" ]; then
        echo "Traitement: $username -> $email"
        
        # Demander le mot de passe pour cet utilisateur
        read -sp "Mot de passe pour $username: " password
        echo
        
        # Générer le hash SHA512-CRYPT
        hash=$(docker exec mailserver doveadm pw -s SHA512-CRYPT -p "$password")
        
        # Mettre à jour postfix-accounts.cf
        docker exec mailserver bash -c "sed -i '/^${email}|/d' /tmp/docker-mailserver/postfix-accounts.cf"
        docker exec mailserver bash -c "echo '${email}|${hash}' >> /tmp/docker-mailserver/postfix-accounts.cf"
        
        echo "✓ Compte $email synchronisé"
    fi
done

echo "=== Redémarrage du serveur mail ==="
docker restart mailserver

echo "✓ Synchronisation terminée !"
