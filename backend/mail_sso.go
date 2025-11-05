package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Secret partagé entre Go et Roundcube pour signer les tokens
const SSO_SECRET = "Office1789-SecretKey-ChangeInProduction"

type SSOTokenClaims struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	ExpiresAt int64  `json:"exp"`
}

// Générer un token SSO pour Roundcube
func GenerateMailSSOToken(c *gin.Context) {
	sessionToken := c.GetHeader("Authorization")
	if sessionToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing session token"})
		return
	}

	// Vérifier la session
	sess, exists := sessions[sessionToken]
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
		return
	}

	if time.Now().After(sess.expiry) {
		delete(sessions, sessionToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session expired"})
		return
	}

	// Créer le token SSO
	claims := SSOTokenClaims{
		Username:  sess.Username,
		Email:     sess.Username + "@office1789.local",
		ExpiresAt: time.Now().Add(5 * time.Minute).Unix(), // Token valide 5 minutes
	}

	// Encoder en JSON
	claimsJSON, _ := json.Marshal(claims)
	claimsB64 := base64.URLEncoding.EncodeToString(claimsJSON)

	// Signer avec HMAC-SHA256
	h := hmac.New(sha256.New, []byte(SSO_SECRET))
	h.Write([]byte(claimsB64))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Token final = claims.signature
	token := claimsB64 + "." + signature

	c.JSON(http.StatusOK, gin.H{
		"sso_token": token,
		"username":  sess.Username,
		"email":     sess.Username + "@office1789.local",
	})
}
