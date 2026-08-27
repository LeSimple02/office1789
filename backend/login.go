package main

import (
	"fmt"
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
	// Debug logs: début
	fmt.Println("[LOGIN][DEBUG] Début de la fonction Connect")

	var conn Connecti
	var hash string
	var userID int

	if err := c.BindJSON(&conn); err != nil {
		fmt.Printf("[LOGIN][ERROR] Failed to parse JSON payload: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_payload"})
		return
	}
	fmt.Printf("[LOGIN][DEBUG] Payload reçu: username=%v, password=(redacted), totp_code=%v\n", conn.Username, conn.TOTPCode)

	if strings.TrimSpace(conn.Username) == "" {
		fmt.Println("[LOGIN][DEBUG] Username vide, retour erreur")
		c.JSON(http.StatusOK, loginError{Username: "no"})
		return
	}

	// Récupérer l'ID utilisateur et le hash du mot de passe
	rows := db.QueryRow("SELECT user_id, password_hash FROM Users WHERE username=$1", conn.Username)
	err := rows.Scan(&userID, &hash)

	if err != nil {
		fmt.Printf("[LOGIN][DEBUG] Utilisateur non trouvé ou erreur SQL: %v\n", err)
		c.JSON(http.StatusOK, loginError{Username: "no"})
		return
	}

	passc := CheckPasswordHash(conn.Password, hash)

	if passc != true {
		fmt.Println("[LOGIN][DEBUG] Mot de passe incorrect")
		c.JSON(http.StatusOK, loginError{Username: "no"})
		return
	}

	// Check if 2FA is enabled for this user
	if Check2FARequired(userID) {
		// 2FA is required - check if TOTP code provided
		if conn.TOTPCode == "" {
			fmt.Println("[LOGIN][DEBUG] 2FA requis mais code non fourni")
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
			fmt.Println("[LOGIN][DEBUG] Code 2FA invalide")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Code 2FA invalide"})
			return
		}
	}

	// Nettoyer les anciennes sessions expirées
	fmt.Println("[LOGIN][DEBUG] Nettoyage des sessions expirées")
	cleanExpiredSessions()

	// Créer une nouvelle session
	fmt.Println("[LOGIN][DEBUG] Création d'une nouvelle session (token, RAM, DB)")
	sessionToken := uuid.NewString()
	expiresAtTime := time.Now().Add(24 * time.Hour)

	// Sauvegarder en mémoire avec mutex (mot de passe inclus pour SSO)
	sessionsMutex.Lock()
	sessions[sessionToken] = session{
		UserID:   userID,
		Username: conn.Username,
		Password: conn.Password, // Stocké en RAM uniquement pour SSO
		expiry:   expiresAtTime,
	}
	sessionsMutex.Unlock()

	// Sauvegarder en base de données
	fmt.Println("[LOGIN][DEBUG] Appel createSessionInDB depuis Connect (login)")
	err = createSessionInDB(userID, conn.Username, sessionToken, expiresAtTime)
	if err != nil {
		fmt.Printf("[LOGIN][ERROR] createSessionInDB depuis Connect (first attempt): %v\n", err)
		// retry once
		err2 := createSessionInDB(userID, conn.Username, sessionToken, expiresAtTime)
		if err2 != nil {
			fmt.Printf("[LOGIN][ERROR] createSessionInDB depuis Connect (second attempt): %v\n", err2)
			// Cleanup RAM session to avoid inconsistent state
			sessionsMutex.Lock()
			delete(sessions, sessionToken)
			sessionsMutex.Unlock()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server_error"})
			return
		}
	}

	// Mettre à jour last_login
	fmt.Println("[LOGIN][DEBUG] Mise à jour du champ last_login")
	db.Exec("UPDATE Users SET last_login = NOW() WHERE user_id = $1", userID)

	fmt.Printf("[LOGIN][DEBUG] Fin login: userID=%v, username=%v, token=%v, expiry=%v\n", userID, conn.Username, sessionToken, expiresAtTime)

	c.JSON(http.StatusOK, sessionSend{
		UserID:   userID,
		Username: conn.Username,
		Token:    sessionToken,
		Expiry:   expiresAtTime,
	})
}
