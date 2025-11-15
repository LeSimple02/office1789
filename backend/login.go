package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type loginError struct {
	Username string `json:"Username"`
}

func Connect(c *gin.Context) {
	var conn Connecti
	var hash string
	var userID int

	c.BindJSON(&conn)

	if strings.TrimSpace(conn.Username) == "" {
		c.JSON(http.StatusOK, loginError{Username: "no"})
		return
	}

	// Récupérer l'ID utilisateur et le hash du mot de passe
	rows := db.QueryRow("SELECT user_id, password_hash FROM Users WHERE username=$1", conn.Username)
	err := rows.Scan(&userID, &hash)
	
	if err != nil {
		c.JSON(http.StatusOK, loginError{Username: "no"})
		return
	}
	
	passc := CheckPasswordHash(conn.Password, hash)

	if passc != true {
		c.JSON(http.StatusOK, loginError{Username: "no"})
		return
	}

	// Check if 2FA is enabled for this user
	if Check2FARequired(userID) {
		// 2FA is required - check if TOTP code provided
		if conn.TOTPCode == "" {
			// Return special response indicating 2FA is required
			c.JSON(http.StatusOK, gin.H{
				"require_2fa": true,
				"user_id":     userID,
				"username":    conn.Username,
			})
			return
		}

		// Validate TOTP code
		if !ValidateTOTPLogin(userID, conn.TOTPCode) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Code 2FA invalide"})
			return
		}
	}

	// Nettoyer les anciennes sessions expirées
	cleanExpiredSessions()
	
	// Créer une nouvelle session
	sessionToken := uuid.NewString()
	expiresAtTime := time.Now().Add(24 * time.Hour)
	
	// Sauvegarder en mémoire
	sessions[sessionToken] = session{
		UserID:   userID,
		Username: conn.Username,
		expiry:   expiresAtTime,
	}
	
	// Sauvegarder en base de données
	err = createSessionInDB(userID, conn.Username, sessionToken, expiresAtTime)
	if err != nil {
		// Log error but continue (fallback to memory-only session)
	}
	
	// Mettre à jour last_login
	db.Exec("UPDATE Users SET last_login = NOW() WHERE user_id = $1", userID)

	c.JSON(http.StatusOK, sessionSend{
		UserID:   userID,
		Username: conn.Username,
		Token:    sessionToken,
		Expiry:   expiresAtTime,
	})
}
