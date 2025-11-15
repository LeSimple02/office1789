#!/bin/bash
# Script de test SSO Matrix/Element pour Linux
# Office1789 - Test end-to-end

echo "============================================================"
echo "OFFICE1789 - Test SSO Matrix/Element (Linux)"
echo "============================================================"
echo ""

# Configuration
BACKEND_URL="http://localhost:8080"
MATRIX_URL="http://localhost:8008"
ELEMENT_URL="http://localhost:8083"

# Couleurs
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
GRAY='\033[0;37m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 1. Test backend API
echo -e "${YELLOW}1. Test Backend API...${NC}"
if curl -s -X POST "$BACKEND_URL/api/session/check" \
   -H "Content-Type: application/json" \
   -d '{"token":"test"}' > /dev/null 2>&1; then
    echo -e "   ${GREEN}OK Backend accessible${NC}"
else
    echo -e "   ${RED}ERROR Backend inaccessible - Demarrez le backend${NC}"
    exit 1
fi

# 2. Test Matrix API
echo -e "\n${YELLOW}2. Test Matrix API...${NC}"
if RESPONSE=$(curl -s "$MATRIX_URL/_matrix/client/versions"); then
    echo -e "   ${GREEN}OK Matrix API accessible${NC}"
    VERSIONS=$(echo "$RESPONSE" | grep -o '"v1\.[0-9]*"' | tail -3 | tr '\n' ',' | sed 's/,$//')
    echo -e "   ${GRAY}INFO Versions supportees: $VERSIONS${NC}"
else
    echo -e "   ${RED}ERROR Matrix API inaccessible - Demarrez Synapse${NC}"
    exit 1
fi

# 3. Test Element
echo -e "\n${YELLOW}3. Test Element...${NC}"
if RESPONSE=$(curl -s "$ELEMENT_URL"); then
    if echo "$RESPONSE" | grep -q "office1789-sso.js"; then
        echo -e "   ${GREEN}OK Element accessible avec plugin SSO${NC}"
    else
        echo -e "   ${YELLOW}WARN Element accessible SANS plugin SSO${NC}"
        echo -e "   ${GRAY}TIP Rebuild Element: docker-compose build element${NC}"
    fi
else
    echo -e "   ${RED}ERROR Element inaccessible - Demarrez le conteneur${NC}"
    exit 1
fi

# 4. Test script SSO présent
echo -e "\n${YELLOW}4. Test script SSO dans Element...${NC}"
if RESULT=$(docker exec element ls -la /app/office1789-sso.js 2>&1); then
    echo -e "   ${GREEN}OK Script SSO present${NC}"
    echo -e "   ${GRAY}INFO $RESULT${NC}"
else
    echo -e "   ${RED}ERROR Script SSO manquant${NC}"
    exit 1
fi

# 5. Test endpoint Matrix SSO
echo -e "\n${YELLOW}5. Test endpoint /api/matrix/sso...${NC}"
echo -e "   ${YELLOW}Necessaire: token de session valide${NC}"
echo -e "   ${GRAY}Test manuel:${NC}"
echo -e "      ${GRAY}1. Login sur Office1789${NC}"
echo -e "      ${GRAY}2. Aller dans Chat${NC}"
echo -e "      ${GRAY}3. Cliquer sur 'Ouvrir Element (Chat)'${NC}"

# Résumé
echo -e "\n============================================================"
echo -e "${GREEN}OK Tests preliminaires OK${NC}"
echo -e "============================================================"

echo -e "\n${YELLOW}Prochaines etapes:${NC}"
echo -e "   ${NC}1. Demarrer Office1789 frontend (npm run dev)${NC}"
echo -e "   ${NC}2. Login avec un compte Office1789${NC}"
echo -e "   ${NC}3. Aller dans 'Chat' (ChatView)${NC}"
echo -e "   ${NC}4. Cliquer sur 'Ouvrir Element (Chat)'${NC}"
echo -e "   ${NC}5. Verifier l'authentification automatique${NC}"
echo ""

echo -e "${CYAN}URLs:${NC}"
echo -e "   ${GRAY}Backend: $BACKEND_URL${NC}"
echo -e "   ${GRAY}Matrix:  $MATRIX_URL${NC}"
echo -e "   ${GRAY}Element: $ELEMENT_URL${NC}"
echo -e "\n${CYAN}Office1789 - SSO Matrix pret !${NC}\n"
