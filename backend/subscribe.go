package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Subscribe struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"Email"`
	Domain      string `json:"Domain"`
	PhoneNumber string `json:"PhoneNumber"`
	Nboffer     int    `json:"Nboffer"`
	DateJoined  string `json:"DateJoined"`
	LastLogin   string `json:"LastLogin"`
}

type vInfo struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

type Connecti struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	LastLogin string `json:"lastlogin"`
}

func Sub(c *gin.Context) {

	var count int
	var subi Subscribe
	var infova vInfo

	c.BindJSON(&subi)

	// Initialiser à vide
	infova.Username = ""
	infova.Email = ""
	infova.Phone = ""
	
	// Vérifier username
	if strings.TrimSpace(subi.Username) != "" {
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE username=$1", subi.Username)
		rows.Scan(&count)

		if count > 0 {
			infova.Username = "no"
		}
	}
	
	// Vérifier email
	if strings.TrimSpace(subi.Email) != "" {
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE email=$1", subi.Email)
		rows.Scan(&count)
		if count > 0 {
			infova.Email = "no"
		}
	}
	
	// Vérifier phonenumber
	if strings.TrimSpace(subi.PhoneNumber) != "" {
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE phonenumber=$1", subi.PhoneNumber)
		rows.Scan(&count)
		if count > 0 {
			infova.Phone = "no"
		}
	}

	// Si un champ est déjà pris, retourner les erreurs
	if infova.Phone == "no" || infova.Username == "no" || infova.Email == "no" {
		c.JSON(http.StatusOK, infova)
		return
	}

	// Tous les champs sont valides, créer l'utilisateur
	if strings.TrimSpace(subi.Username) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}

	subi.Password = HashPassword(subi.Password)

	var userID int
	err := db.QueryRow("INSERT INTO Users (username, password_hash, email, phonenumber, nboffer, date_joined, last_login, domain) VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), $6) RETURNING user_id", 
		subi.Username, subi.Password, subi.Email, subi.PhoneNumber, subi.Nboffer, "@office1789").Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	sessionToken := uuid.NewString()
	expiresAtTime := time.Now().Add(24 * time.Hour)
	
	// Créer session en mémoire
	sessions[sessionToken] = session{
		UserID:   userID,
		Username: subi.Username,
		expiry:   expiresAtTime,
	}
	
	// Créer session en DB
	_ = createSessionInDB(userID, subi.Username, sessionToken, expiresAtTime)

	c.JSON(http.StatusOK, sessionSend{
		UserID:   userID,
		Username: subi.Username,
		Token:    sessionToken,
		Expiry:   expiresAtTime,
	})
}
