package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// MailAuthRequest représente une demande d'authentification mail
type MailAuthRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	Password string `json:"password"`
}

// MailAuthResponse représente la réponse d'authentification mail
type MailAuthResponse struct {
	URL   string `json:"url"`
	Error string `json:"error,omitempty"`
}

// GenerateMailSSO génère un token SSO pour Roundcube AVEC vérification du mot de passe
func GenerateMailSSO(c *gin.Context) {
	var req MailAuthRequest
	
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("[SSO] Erreur bind JSON:", err)
		c.JSON(http.StatusBadRequest, MailAuthResponse{
			Error: "Requête invalide",
		})
		return
	}

	fmt.Printf("[SSO] Requête reçue - Username: %s, Token présent: %v\n", req.Username, req.Token != "")

	// Valider la session
	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		fmt.Println("[SSO] Session invalide")
		c.JSON(http.StatusUnauthorized, MailAuthResponse{
			Error: "Session invalide",
		})
		return
	}

	fmt.Printf("[SSO] Session valide - UserID: %d\n", session.UserID)

	// Vérifier le mot de passe
	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM Users WHERE user_id=$1", session.UserID).Scan(&storedHash)
	if err != nil {
		fmt.Println("[SSO] Erreur récupération hash:", err)
		c.JSON(http.StatusInternalServerError, MailAuthResponse{
			Error: "Erreur lors de la vérification du mot de passe",
		})
		return
	}

	if !CheckPasswordHash(req.Password, storedHash) {
		fmt.Println("[SSO] Mot de passe incorrect")
		c.JSON(http.StatusUnauthorized, MailAuthResponse{
			Error: "Mot de passe incorrect",
		})
		return
	}

	fmt.Println("[SSO] Mot de passe validé")

	// Récupérer l'email de l'utilisateur
	var email string
	err = db.QueryRow("SELECT Email FROM Users WHERE user_id=$1", session.UserID).Scan(&email)
	if err != nil {
		fmt.Println("[SSO] Erreur récupération email:", err)
		c.JSON(http.StatusInternalServerError, MailAuthResponse{
			Error: "Erreur lors de la récupération de l'email",
		})
		return
	}

	fmt.Printf("[SSO] Email récupéré: %s\n", email)

	// Générer le token SSO avec le mot de passe
	ssoToken := generateSSOToken(req.Username, email, req.Password)
	
	fmt.Printf("[SSO] Token généré (50 premiers char): %s...\n", ssoToken[:min(50, len(ssoToken))])

	// Construire l'URL de Roundcube avec le token SSO
	roundcubeURL := "http://localhost:8081/?sso_token=" + ssoToken

	fmt.Printf("[SSO] URL générée: %s\n", roundcubeURL[:min(100, len(roundcubeURL))])

	c.JSON(http.StatusOK, MailAuthResponse{
		URL: roundcubeURL,
	})
}

// GenerateMailSSOAuto génère un token SSO SANS vérification du mot de passe (vrai SSO)
func GenerateMailSSOAuto(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}
	
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("[SSO-Auto] Erreur bind JSON:", err)
		c.JSON(http.StatusBadRequest, MailAuthResponse{
			Error: "Requête invalide",
		})
		return
	}

	fmt.Printf("[SSO-Auto] Requête reçue - Username: %s\n", req.Username)

	// Valider la session Office1789
	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		fmt.Println("[SSO-Auto] Session invalide")
		c.JSON(http.StatusUnauthorized, MailAuthResponse{
			Error: "Session invalide",
		})
		return
	}

	fmt.Printf("[SSO-Auto] Session valide - UserID: %d\n", session.UserID)

	// Récupérer l'email de l'utilisateur
	var email string
	err := db.QueryRow("SELECT Email FROM Users WHERE user_id=$1", session.UserID).Scan(&email)
	if err != nil {
		fmt.Println("[SSO-Auto] Erreur récupération email:", err)
		c.JSON(http.StatusInternalServerError, MailAuthResponse{
			Error: "Email non trouvé",
		})
		return
	}

	fmt.Printf("[SSO-Auto] Email récupéré: %s\n", email)

	// Utiliser le mot de passe de la session (capturé au login)
	password := session.Password
	if password == "" {
		fmt.Println("[SSO-Auto] ERREUR: Session sans mot de passe - reconnectez-vous")
		c.JSON(http.StatusUnauthorized, MailAuthResponse{
			Error: "Veuillez vous reconnecter à Office1789",
		})
		return
	}

	// Générer le token SSO avec le mot de passe Office1789 (synchronisé avec mail)
	ssoToken := generateSSOToken(req.Username, email, password)
	
	fmt.Printf("[SSO-Auto] Token SSO généré avec authentification\n")

	// Construire l'URL de Roundcube avec le token SSO
	roundcubeURL := "http://localhost:8081/?sso_token=" + ssoToken

	c.JSON(http.StatusOK, MailAuthResponse{
		URL: roundcubeURL,
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// generateSSOToken génère un token SSO signé pour Roundcube
func generateSSOToken(username, email, password string) string {
	// Secret partagé avec le plugin Roundcube
	secret := "Office1789-SecretKey-ChangeInProduction"

	// Créer les claims AVEC le mot de passe pour authentifier au mail
	claims := map[string]interface{}{
		"username": username,
		"email":    email,
		"password": password, // Mot de passe inclus pour auth mail
		"exp":      time.Now().Add(5 * time.Minute).Unix(), // Expire dans 5 minutes
		"iat":      time.Now().Unix(),
	}

	// Encoder les claims en JSON puis en base64
	claimsJSON, _ := json.Marshal(claims)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Signer avec HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(claimsB64))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	// Retourner le token: claims.signature
	return claimsB64 + "." + signature
}

// CheckSession vérifie si une session est valide
func CheckSession(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"connected": false})
		return
	}

	_, valid := validateSession(req.Token, req.Username)
	c.JSON(http.StatusOK, gin.H{"connected": valid})
}
