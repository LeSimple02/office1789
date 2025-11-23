#!/bin/bash

# Script pour vérifier les tables dans les 3 bases PostgreSQL

echo "======================================================"
echo "🔍 VÉRIFICATION DES BASES DE DONNÉES PostgreSQL"
echo "======================================================"
echo ""

# Couleurs
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Fonction pour compter les tables
check_database() {
    local container=$1
    local db_name=$2
    local db_user=$3
    local title=$4
    
    echo -e "${YELLOW}📦 $title${NC}"
    echo "   Container: $container"
    echo "   Database: $db_name"
    echo ""
    
    # Vérifier si le conteneur existe et est en cours d'exécution
    if ! docker ps --format '{{.Names}}' | grep -q "^${container}$"; then
        echo -e "   ${RED}❌ Conteneur non trouvé ou arrêté${NC}"
        echo ""
        return
    fi
    
    # Lister les tables
    echo "   📋 Tables trouvées:"
    docker exec $container psql -U $db_user -d $db_name -c "\dt" 2>/dev/null | grep -E "^ " || echo "   (aucune table ou erreur de connexion)"
    
    # Compter les tables
    table_count=$(docker exec $container psql -U $db_user -d $db_name -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';" 2>/dev/null | tr -d ' ')
    
    if [ ! -z "$table_count" ] && [ "$table_count" -gt 0 ]; then
        echo -e "   ${GREEN}✅ $table_count table(s) trouvée(s)${NC}"
    else
        echo -e "   ${RED}❌ Aucune table ou erreur de connexion${NC}"
    fi
    
    echo ""
    echo "------------------------------------------------------"
    echo ""
}

# Vérifier postgres_db (base principale Office1789)
check_database "postgres_db" "office1789" "office1789user" "BASE PRINCIPALE (office1789)"

# Vérifier postgres_roundcube
check_database "postgres_roundcube" "roundcube" "roundcube" "BASE ROUNDCUBE (Webmail)"

# Vérifier postgres_synapse
check_database "postgres_synapse" "synapse" "synapse_user" "BASE SYNAPSE (Matrix)"

echo "======================================================"
echo "🏁 VÉRIFICATION TERMINÉE"
echo "======================================================"
