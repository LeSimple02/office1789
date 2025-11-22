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

// Usernames réservés pour le système
var reservedUsernames = []string{
	"noreply",
	"no-reply",
	"admin",
	"administrator",
	"postmaster",
	"hostmaster",
	"webmaster",
	"root",
	"system",
	"support",
	"info",
	"contact",
	"abuse",
	"security",
	"mailer-daemon",
	"daemon",
	"office1789",
	"office",
}

// Vérifier si un username est réservé
func isReservedUsername(username string) bool {
	usernameLower := strings.ToLower(strings.TrimSpace(username))
	for _, reserved := range reservedUsernames {
		if usernameLower == reserved {
			return true
		}
	}
	return false
}

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
	email := username + "@office1789.com"
	
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
		// Vérifier si le username est réservé
		if isReservedUsername(subi.Username) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "reserved_username",
				"message": "Ce nom d'utilisateur est réservé et ne peut pas être utilisé",
			})
			return
		}
		
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE username=$1", subi.Username)
		rows.Scan(&count)

		if count > 0 {
			infova.Username = "no"
		}
	}
	
	// Vérifier recovery_email (optionnel)
	if strings.TrimSpace(subi.Email) != "" {
		// Vérifier que l'email a été vérifié dans les 30 dernières minutes
		if !CheckVerificationStatus(subi.Email, "email") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "email_not_verified",
				"message": "Vous devez vérifier votre email avant de créer un compte",
			})
			return
		}
		
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE recovery_email=$1", subi.Email)
		rows.Scan(&count)
		if count > 0 {
			infova.Email = "no"
		}
	}
	
	// Vérifier phonenumber
	if strings.TrimSpace(subi.PhoneNumber) != "" {
		// Vérifier que le téléphone a été vérifié dans les 30 dernières minutes
		if !CheckVerificationStatus(subi.PhoneNumber, "phone") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "phone_not_verified",
				"message": "Vous devez vérifier votre numéro avant de créer un compte",
			})
			return
		}
		
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

	// NOTE: subi.Email = email de récupération (optionnel, peut être vide ou différent)
	// L'email du compte mail sera TOUJOURS username@office1789.com (construit côté SSO)

	// Garder le mot de passe en clair pour créer les comptes mail/Matrix
	plainPassword := subi.Password
	subi.Password = HashPassword(subi.Password)

	// Chiffrer le mot de passe pour mail/Matrix (stockage sécurisé mais réversible)
	encryptedMailPassword, err := EncryptPassword(plainPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}

	// Construire l'email du compte mail (toujours username@office1789.com)
	mailAddress := subi.Username + "@office1789.com"
	
	// Déterminer si les contacts sont vérifiés
	emailVerified := false
	phoneVerified := false
	if strings.TrimSpace(subi.Email) != "" {
		emailVerified = true // Déjà vérifié via CheckVerificationStatus plus haut
	}
	if strings.TrimSpace(subi.PhoneNumber) != "" {
		phoneVerified = true // Déjà vérifié via CheckVerificationStatus plus haut
	}
	
	var userID int
	err = db.QueryRow("INSERT INTO Users (username, password_hash, email, recovery_email, recovery_email_verified, phonenumber, phonenumber_verified, nboffer, date_joined, last_login, domain, mail_password, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW(), $9, $10, 'user') RETURNING user_id", 
		subi.Username, subi.Password, mailAddress, subi.Email, emailVerified, subi.PhoneNumber, phoneVerified, subi.Nboffer, "@office1789", encryptedMailPassword).Scan(&userID)

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
			fmt.Printf("Mail account created for %s@office1789.com\n", subi.Username)
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
	
	// Créer session en mémoire avec mutex
	sessionsMutex.Lock()
	sessions[sessionToken] = session{
		UserID:   userID,
		Username: subi.Username,
		expiry:   expiresAtTime,
	}
	sessionsMutex.Unlock()
	
	// Créer session en DB
	_ = createSessionInDB(userID, subi.Username, sessionToken, expiresAtTime)

	c.JSON(http.StatusOK, sessionSend{
		UserID:   userID,
		Username: subi.Username,
		Token:    sessionToken,
		Expiry:   expiresAtTime,
	})
}
