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
	Username     string `json:username`
	LastUsername string `json:lastusername`
	Password     string `json:password`
	Email        string `json:email`
	Domain       string `json:email`
	PhoneNumber  string `json:phonenumber`
	Nboffer      int    `json:nboffer`
	Token        string `json:token`
}

func getinfop(c *gin.Context) {

	var verif sessionSend
	var infop Subscribe

	c.ShouldBindJSON(&verif)

	if sessions[verif.Token].Username == verif.Username && verif.Username != "" {
		rows := db.QueryRow("SELECT domain, nboffer, date_joined, last_login, phonenumber, email FROM Users WHERE username=$1", verif.Username)
		rows.Scan(&infop.Domain, &infop.Nboffer, &infop.DateJoined, &infop.LastLogin, &infop.PhoneNumber, &infop.Email)

		c.JSON(http.StatusOK, infop)
		return

	} else {
		infop.Username = "no"

		c.JSON(http.StatusOK, infop)

	}

}

func ChangeI(c *gin.Context) {
	var count int
	var cha ChangeIn
	var infova vInfo

	c.BindJSON(&cha)

	if sessions[cha.Token].Username == cha.LastUsername {

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

			if strings.TrimSpace(cha.LastUsername) != "" {

				if cha.Password != "" {
					cha.Password = HashPassword(cha.Password)
					db.Exec("UPDATE Users SET password_hash=$1 WHERE username=$2", cha.Password, cha.LastUsername)
				}

				if cha.Email != "" {
					db.Exec("UPDATE Users SET email=$1 WHERE username=$2", cha.Email, cha.LastUsername)
				}
				if cha.PhoneNumber != "" {
					db.Exec("UPDATE Users SET phonenumber=$1 WHERE username=$2", cha.PhoneNumber, cha.LastUsername)
				}
				if cha.Username != "" {
					db.Exec("UPDATE Users SET username=$1 WHERE username=$2", cha.Username, cha.LastUsername)
				}

				sessionToken := uuid.NewString()
				expiresAtTime := time.Now().Add(120 * time.Second)
				if cha.Username != "" {
					sessions[sessionToken] = session{
						Username: cha.Username,
						expiry:   expiresAtTime,
					}
					c.JSON(http.StatusOK, sessionSend{cha.Username, sessionToken, expiresAtTime})
				} else {
					c.JSON(http.StatusOK, infova)
				}

			}

		}
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
	session, exists := sessions[req.Token]
	if !exists || session.Username != req.Username {
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

	// Supprimer l'utilisateur de la base de données
	_, err := db.Exec("DELETE FROM Users WHERE username=$1", req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DeleteAccountResponse{
			Success: false,
			Message: "Failed to delete account",
		})
		return
	}

	// Supprimer la session
	delete(sessions, req.Token)

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

