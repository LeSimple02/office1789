package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os/exec"
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
	TOTPCode  string `json:"totp_code"` // For 2FA authentication
}

// Créer un compte mail dans le mailserver
func createMailAccount(username, password string) error {
	email := username + "@office1789.local"
	
	// Commande Docker pour créer le compte mail
	cmd := exec.Command("docker", "exec", "mailserver", "setup", "email", "add", email, password)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create mail account: %v - %s", err, stderr.String())
	}
	
	return nil
}

// Créer un compte Matrix (Synapse)
func createMatrixAccount(username, password string) error {
	// Utiliser register_new_matrix_user CLI tool de Synapse
	// Cette commande crée un utilisateur via l'interface en ligne de commande
	cmd := exec.Command("docker", "exec", "synapse", 
		"register_new_matrix_user", 
		"http://localhost:8008",
		"-u", username,
		"-p", password,
		"-c", "/data/homeserver.yaml",
		"--no-admin")
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to create matrix account: %v - %s", err, stderr.String())
	}
	
	return nil
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

	// Garder le mot de passe en clair pour créer les comptes mail/Matrix
	plainPassword := subi.Password
	subi.Password = HashPassword(subi.Password)

	var userID int
	err := db.QueryRow("INSERT INTO Users (username, password_hash, email, phonenumber, nboffer, date_joined, last_login, domain) VALUES ($1, $2, $3, $4, $5, NOW(), NOW(), $6) RETURNING user_id", 
		subi.Username, subi.Password, subi.Email, subi.PhoneNumber, subi.Nboffer, "@office1789").Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Créer le compte mail automatiquement
	go func() {
		err := createMailAccount(subi.Username, plainPassword)
		if err != nil {
			fmt.Printf("Warning: Failed to create mail account for %s: %v\n", subi.Username, err)
		} else {
			fmt.Printf("Mail account created for %s@office1789.local\n", subi.Username)
		}
	}()

	// Créer le compte Matrix automatiquement
	go func() {
		err := createMatrixAccount(subi.Username, plainPassword)
		if err != nil {
			fmt.Printf("Warning: Failed to create Matrix account for %s: %v\n", subi.Username, err)
		} else {
			fmt.Printf("Matrix account created for @%s:office1789.com\n", subi.Username)
		}
	}()

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
