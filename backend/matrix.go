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

// GenerateMatrixSSOAuto génère un token SSO pour Element/Matrix sans redemander le mot de passe
func GenerateMatrixSSOAuto(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}

	if err := c.BindJSON(&req); err != nil {
		fmt.Println("[Matrix-SSO] Erreur bind JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide"})
		return
	}

	fmt.Printf("[Matrix-SSO] Requête reçue - Username: %s\n", req.Username)

	// Valider la session Office1789
	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		fmt.Println("[Matrix-SSO] Session invalide")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalide"})
		return
	}

	fmt.Printf("[Matrix-SSO] Session valide - UserID: %d\n", session.UserID)

	// Récupérer le password depuis la table Users (source unique de vérité)
	var encryptedPassword string
	err := db.QueryRow("SELECT COALESCE(mail_password, '') FROM users WHERE user_id=$1", session.UserID).Scan(&encryptedPassword)
	if err != nil {
		fmt.Printf("[Matrix-SSO] ❌ Erreur récupération password: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
		return
	}

	fmt.Printf("[Matrix-SSO] 🔑 Password chiffré récupéré depuis DB, présent: %v\n", encryptedPassword != "")

	// Vérifier que le password existe
	// Gestion migration : mail_password NULL ou ancien hash bcrypt
	if encryptedPassword == "" {
		fmt.Printf("[Matrix-SSO] ⚠️ mail_password NULL pour %s - Migration nécessaire\n", req.Username)
		c.JSON(http.StatusPreconditionFailed, gin.H{
			"error": "Votre compte nécessite une migration. Veuillez changer votre mot de passe dans 'Mon compte' pour activer l'accès Mail/Matrix.",
		})
		return
	}

	// Déchiffrer le mot de passe (AES-256-GCM)
	userPassword, err := DecryptPassword(encryptedPassword)
	if err != nil {
		// Si déchiffrement échoue = ancien hash bcrypt non migré
		fmt.Printf("[Matrix-SSO] ⚠️ Déchiffrement impossible pour %s - Ancien hash détecté\n", req.Username)
		c.JSON(http.StatusPreconditionFailed, gin.H{
			"error": "Votre compte utilise l'ancien système. Changez votre mot de passe dans 'Mon compte' pour réactiver Mail/Matrix.",
		})
		return
	}

	// Générer le token SSO Matrix
	matrixUserID := fmt.Sprintf("@%s:office1789.com", req.Username)
	ssoToken := generateMatrixSSOToken(req.Username, matrixUserID, userPassword)

	fmt.Printf("[Matrix-SSO] Token SSO généré pour %s\n", matrixUserID)

	// Construire l'URL d'Element avec le token SSO
	elementURL := fmt.Sprintf("http://localhost:8083/?sso_token=%s", ssoToken)

	c.JSON(http.StatusOK, gin.H{
		"url": elementURL,
	})
}

// generateMatrixSSOToken génère un token SSO signé pour Element/Matrix
func generateMatrixSSOToken(username, matrixUserID, password string) string {
	// Secret partagé avec le plugin Element
	secret := "Office1789-Matrix-SecretKey-ChangeInProduction"

	// Créer les claims AVEC le mot de passe pour authentifier à Matrix
	claims := map[string]interface{}{
		"username":     username,               // Username Office1789
		"matrixUserId": matrixUserID,           // @username:office1789.com
		"password":     password,               // Mot de passe inclus pour auth Matrix
		"exp":          time.Now().Add(5 * time.Minute).Unix(), // Expire dans 5 minutes
		"iat":          time.Now().Unix(),
		"nonce":        time.Now().UnixNano(), // Nonce unique pour éviter replay attack
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
