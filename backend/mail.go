package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Structure pour la requête SSO mail
type MailSSORequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

// Générer un token SSO pour Roundcube
func generateRoundcubeToken(username string) string {
	// Créer un hash unique basé sur username + timestamp + secret
	secret := "office1789_secret_key" // À changer par une vraie clé secrète
	data := fmt.Sprintf("%s:%d:%s", username, time.Now().Unix(), secret)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}

// Endpoint pour obtenir l'URL de connexion automatique à Roundcube
func getRoundcubeAuth(c *gin.Context) {
	var req MailSSORequest
	c.BindJSON(&req)

	// Vérifier que la session est valide
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid session"})
		return
	}

	// Générer un token SSO temporaire
	ssoToken := generateRoundcubeToken(req.Username)
	
	// Stocker le token SSO avec expiration courte (5 minutes)
	ssoTokens[ssoToken] = ssoTokenData{
		Username: req.Username,
		Expiry:   time.Now().Add(5 * time.Minute),
	}

	// Construire l'URL de connexion automatique
	// Format: http://localhost:8081/?_task=mail&_token=TOKEN&_user=USERNAME
	roundcubeURL := fmt.Sprintf("http://localhost:8081/?_task=mail&_autologin=1&_user=%s&_token=%s", 
		req.Username, ssoToken)

	c.JSON(http.StatusOK, gin.H{
		"url": roundcubeURL,
		"username": req.Username,
	})
}

// Endpoint de validation pour Roundcube (webhook)
func validateRoundcubeSSOToken(c *gin.Context) {
	token := c.Query("token")
	username := c.Query("user")

	if token == "" || username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
		return
	}

	// Vérifier le token SSO
	if ssoData, ok := ssoTokens[token]; ok {
		if ssoData.Username == username && time.Now().Before(ssoData.Expiry) {
			// Token valide
			delete(ssoTokens, token) // Supprimer après utilisation
			c.JSON(http.StatusOK, gin.H{
				"valid": true,
				"username": username,
				"email": username + "@office1789.local", // Adapter selon votre domaine
			})
			return
		}
	}

	c.JSON(http.StatusUnauthorized, gin.H{"valid": false})
}

// Vérifier si l'utilisateur est connecté (endpoint de santé de session)
func checkSessionStatus(c *gin.Context) {
	var req MailSSORequest
	c.BindJSON(&req)

	// Vérifier la session
	if session, ok := sessions[req.Token]; ok && session.Username == req.Username && time.Now().Before(session.expiry) {
		c.JSON(http.StatusOK, gin.H{
			"connected": true,
			"username": req.Username,
			"expiry": session.expiry,
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{
		"connected": false,
	})
}
