package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChangeIn struct {
	Username     string `json:"username"`
	LastUsername string `json:"lastusername"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	Domain       string `json:"domain"`
	PhoneNumber  string `json:"phonenumber"`
	Nboffer      int    `json:"nboffer"`
	Token        string `json:"token"`
}

func getinfop(c *gin.Context) {

	var verif sessionVerify
	var infop Subscribe

	c.ShouldBindJSON(&verif)

	session, valid := validateSession(verif.Token, verif.Username)
	if !valid {
		infop.Username = "no"
		c.JSON(http.StatusOK, infop)
		return
	}

	rows := db.QueryRow("SELECT domain, nboffer, date_joined, last_login, phonenumber, email FROM Users WHERE user_id=$1", session.UserID)
	rows.Scan(&infop.Domain, &infop.Nboffer, &infop.DateJoined, &infop.LastLogin, &infop.PhoneNumber, &infop.Email)

	c.JSON(http.StatusOK, infop)
}

func ChangeI(c *gin.Context) {
	var count int
	var cha ChangeIn
	var infova vInfo

	c.BindJSON(&cha)

	sess, valid := validateSession(cha.Token, cha.LastUsername)
	if !valid {
		infova.Username = "no"
		c.JSON(http.StatusUnauthorized, infova)
		return
	}

	if strings.TrimSpace(cha.Username) != "" {
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE username=$1", cha.Username)
		rows.Scan(&count)

		if count > 0 {
			infova.Username = "no"
		}

	}
	if strings.TrimSpace(cha.Email) != "" {
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE email=$1", cha.Email)
		rows.Scan(&count)
		if count > 0 {
			infova.Email = "no"
		}

	}
	if strings.TrimSpace(cha.PhoneNumber) != "" {
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE phonenumber=$1", cha.PhoneNumber)
		rows.Scan(&count)
		if count > 0 {
			infova.Phone = "no"
		}
	}

	if infova.Phone == "no" || infova.Username == "no" || infova.Email == "no" {
		c.JSON(http.StatusOK, infova)
	} else {

		if cha.Password != "" {
			cha.Password = HashPassword(cha.Password)
			db.Exec("UPDATE Users SET password_hash=$1 WHERE user_id=$2", cha.Password, sess.UserID)
		}

		if cha.Email != "" {
			db.Exec("UPDATE Users SET email=$1 WHERE user_id=$2", cha.Email, sess.UserID)
		}
		if cha.PhoneNumber != "" {
			db.Exec("UPDATE Users SET phonenumber=$1 WHERE user_id=$2", cha.PhoneNumber, sess.UserID)
		}
		
		newUsername := cha.LastUsername
		if cha.Username != "" {
			db.Exec("UPDATE Users SET username=$1 WHERE user_id=$2", cha.Username, sess.UserID)
			newUsername = cha.Username
		}

		// Créer nouvelle session
		sessionToken := uuid.NewString()
		expiresAtTime := time.Now().Add(24 * time.Hour)
		
		// Supprimer l'ancienne session
		delete(sessions, cha.Token)
		deleteSessionFromDB(cha.Token)
		
		// Créer nouvelle session en mémoire
		sessions[sessionToken] = session{
			UserID:   sess.UserID,
			Username: newUsername,
			expiry:   expiresAtTime,
		}
		
		// Créer nouvelle session en DB
		_ = createSessionInDB(sess.UserID, newUsername, sessionToken, expiresAtTime)
		
		c.JSON(http.StatusOK, sessionSend{
			UserID:   sess.UserID,
			Username: newUsername,
			Token:    sessionToken,
			Expiry:   expiresAtTime,
		})
	}
}

type DeleteAccountRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type DeleteAccountResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func DeleteAccount(c *gin.Context) {
	var req DeleteAccountRequest
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, DeleteAccountResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}

	// Vérifier que la session est valide
	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, DeleteAccountResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Supprimer tous les fichiers/uploads de l'utilisateur
	uploadsPath := "./uploads/" + req.Username
	// Ignorer l'erreur si le dossier n'existe pas
	_ = removeDirectory(uploadsPath)

	// Supprimer l'avatar de l'utilisateur
	avatarPath := "./avatars/" + req.Username + ".png"
	_ = removeFile(avatarPath)

	// Supprimer l'utilisateur de la base de données (CASCADE supprimera les sessions)
	_, err := db.Exec("DELETE FROM Users WHERE user_id=$1", session.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DeleteAccountResponse{
			Success: false,
			Message: "Failed to delete account",
		})
		return
	}

	c.JSON(http.StatusOK, DeleteAccountResponse{
		Success: true,
		Message: "Account deleted successfully",
	})
}

// Helper functions for file/directory removal
func removeFile(path string) error {
	return os.Remove(path)
}

func removeDirectory(path string) error {
	return os.RemoveAll(path)
}

