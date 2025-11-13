package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Connect(c *gin.Context) {
	var conn Connecti
	var hash string
	var userID int

	c.BindJSON(&conn)

	if strings.TrimSpace(conn.Username) != "" {
		// Récupérer l'ID utilisateur et le hash du mot de passe
		rows := db.QueryRow("SELECT user_id, password_hash FROM Users WHERE username=$1", conn.Username)
		err := rows.Scan(&userID, &hash)
		
		if err != nil {
			conn.Username = "no"
			c.JSON(http.StatusOK, conn)
			return
		}
		
		passc := CheckPasswordHash(conn.Password, hash)

		if passc == true {
			// Nettoyer les anciennes sessions expirées
			cleanExpiredSessions()
			
			// Créer une nouvelle session
			sessionToken := uuid.NewString()
			expiresAtTime := time.Now().Add(24 * time.Hour)
			
			// Sauvegarder en mémoire (ancien système)
			sessions[sessionToken] = session{
				UserID:   userID,
				Username: conn.Username,
				expiry:   expiresAtTime,
			}
			
			// Sauvegarder en base de données (nouveau système)
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
		} else {
			conn.Username = "no"
			c.JSON(http.StatusOK, conn)
		}
	}
}
