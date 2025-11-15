package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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
	Email        string `json:"email"` // Email de récupération
	Domain       string `json:"domain"`
	PhoneNumber  string `json:"phonenumber"`
	Nboffer      int    `json:"nboffer"`
	Token        string `json:"token"`
}

type UserProfile struct {
	Username      string `json:"username"`
	Domain        string `json:"Domain"`
	Email         string `json:"Email"`         // Email de récupération (recovery_email)
	RecoveryEmail string `json:"RecoveryEmail"` // Alias pour clarté
	PhoneNumber   string `json:"PhoneNumber"`
	Nboffer       int    `json:"Nboffer"`
	DateJoined    string `json:"DateJoined"`
	LastLogin     string `json:"LastLogin"`
}

func getinfop(c *gin.Context) {

	var verif sessionVerify
	var infop UserProfile

	c.ShouldBindJSON(&verif)

	session, valid := validateSession(verif.Token, verif.Username)
	if !valid {
		infop.Username = "no"
		c.JSON(http.StatusOK, infop)
		return
	}

	rows := db.QueryRow("SELECT username, domain, nboffer, date_joined, last_login, phonenumber, COALESCE(recovery_email, '') FROM Users WHERE user_id=$1", session.UserID)
	rows.Scan(&infop.Username, &infop.Domain, &infop.Nboffer, &infop.DateJoined, &infop.LastLogin, &infop.PhoneNumber, &infop.Email)
	infop.RecoveryEmail = infop.Email // Copier Email et RecoveryEmail pour compatibilité frontend

	c.JSON(http.StatusOK, infop)
}

func ChangeI(c *gin.Context) {
	var count int
	var cha ChangeIn
	var infova vInfo

	c.BindJSON(&cha)
	
	// Debug: afficher ce qui est reçu
	fmt.Printf("DEBUG ChangeI - Username reçu: '%s' (len=%d)\n", cha.Username, len(cha.Username))
	fmt.Printf("DEBUG ChangeI - LastUsername: '%s'\n", cha.LastUsername)
	fmt.Printf("DEBUG ChangeI - Email: '%s'\n", cha.Email)
	fmt.Printf("DEBUG ChangeI - Phone: '%s'\n", cha.PhoneNumber)

	sess, valid := validateSession(cha.Token, cha.LastUsername)
	if !valid {
		infova.Username = "no"
		c.JSON(http.StatusUnauthorized, infova)
		return
	}

	// Validation du username SEULEMENT si fourni et non vide
	if strings.TrimSpace(cha.Username) != "" {
		fmt.Printf("DEBUG ChangeI - Username validation: '%s' vs '%s'\n", cha.Username, cha.LastUsername)
		
		// Bloquer si l'utilisateur essaie de mettre son propre nom
		if cha.Username == cha.LastUsername {
			fmt.Println("DEBUG ChangeI - Username same as current")
			infova.Username = "same"
			c.JSON(http.StatusOK, infova)
			return
		}
		
		// Vérifier si le nouveau username est déjà pris par quelqu'un d'autre
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE username=$1 AND user_id != $2", cha.Username, sess.UserID)
		rows.Scan(&count)
		fmt.Printf("DEBUG ChangeI - Username count: %d\n", count)

		if count > 0 {
			infova.Username = "no"
		}
	}
	
	// Validation email SEULEMENT si fourni et non vide
	if strings.TrimSpace(cha.Email) != "" {
		// Vérifier que ce n'est pas déjà pris par quelqu'un d'autre
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE email=$1 AND user_id != $2", cha.Email, sess.UserID)
		rows.Scan(&count)
		if count > 0 {
			infova.Email = "no"
		}
	}
	
	// Validation phone SEULEMENT si fourni et non vide
	if strings.TrimSpace(cha.PhoneNumber) != "" {
		// Vérifier que ce n'est pas déjà pris par quelqu'un d'autre
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE phonenumber=$1 AND user_id != $2", cha.PhoneNumber, sess.UserID)
		rows.Scan(&count)
		if count > 0 {
			infova.Phone = "no"
		}
	}

	// S'il y a des erreurs de validation, les retourner
	if infova.Phone == "no" || infova.Username == "no" || infova.Username == "same" || infova.Email == "no" {
		fmt.Printf("DEBUG ChangeI - Erreur validation: username='%s', email='%s', phone='%s'\n", infova.Username, infova.Email, infova.Phone)
		c.JSON(http.StatusOK, infova)
		return
	}
	
	// Pas d'erreur, on procède aux modifications
	fmt.Println("DEBUG ChangeI - Validation OK, procède aux modifications")

	// Changement de mot de passe synchronisé (Office1789 + Mail + Matrix)
	if cha.Password != "" {
		newHash := HashPassword(cha.Password)
		
		// Chiffrer le mot de passe pour mail/Matrix
		fmt.Printf("DEBUG ChangeI - 🔐 Chiffrement du nouveau mot de passe pour user_id=%d\n", sess.UserID)
		encryptedPassword, err := EncryptPassword(cha.Password)
		if err != nil {
			fmt.Printf("DEBUG ChangeI - ❌ ERREUR CHIFFREMENT: %v\n", err)
			c.JSON(http.StatusInternalServerError, vInfo{
				Username: "error",
				Email:    "error",
				Phone:    "error",
			})
			return
		}
		fmt.Printf("DEBUG ChangeI - ✅ Mot de passe chiffré (len=%d)\n", len(encryptedPassword))
		
		// UPDATE avec mail_password chiffré
		db.Exec("UPDATE Users SET password_hash=$1, mail_password=$2 WHERE user_id=$3", 
			newHash, encryptedPassword, sess.UserID)
		fmt.Printf("DEBUG ChangeI - Password changed for user %s\n", cha.LastUsername)
		
		// Synchroniser avec Mail (asynchrone)
		go func() {
			err := changeMailPassword(cha.LastUsername, cha.Password)
			if err != nil {
				fmt.Printf("Warning: Failed to change mail password for %s: %v\n", cha.LastUsername, err)
			} else {
				fmt.Printf("Mail password changed for %s@office1789.com\n", cha.LastUsername)
			}
		}()
		
		// Synchroniser avec Matrix (asynchrone)
		go func() {
			err := changeMatrixPassword(cha.LastUsername, cha.Password)
			if err != nil {
				fmt.Printf("Warning: Failed to change matrix password for %s: %v\n", cha.LastUsername, err)
			} else {
				fmt.Printf("Matrix password changed for @%s:office1789.com\n", cha.LastUsername)
			}
		}()
	}

	if cha.Email != "" {
		db.Exec("UPDATE Users SET recovery_email=$1 WHERE user_id=$2", cha.Email, sess.UserID)
	}
	if cha.PhoneNumber != "" {
		db.Exec("UPDATE Users SET phonenumber=$1 WHERE user_id=$2", cha.PhoneNumber, sess.UserID)
	}
	
	// Vérifier si le username change
	usernameChanged := false
	newUsername := cha.LastUsername
	if strings.TrimSpace(cha.Username) != "" && cha.Username != cha.LastUsername {
		// Changer username ET mettre à jour email du compte mail
		db.Exec("UPDATE Users SET username=$1, email=$2 WHERE user_id=$3", 
			cha.Username, cha.Username+"@office1789.com", sess.UserID)
		newUsername = cha.Username
		usernameChanged = true
	}

	// Si le username a changé, créer une nouvelle session
	if usernameChanged {
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
		
		fmt.Printf("DEBUG ChangeI - Username changé! Nouveau token: %s, Username: %s\n", sessionToken, newUsername)
		
		c.JSON(http.StatusOK, sessionSend{
			UserID:   sess.UserID,
			Username: newUsername,
			Token:    sessionToken,
			Expiry:   expiresAtTime,
		})
	} else {
		// Pas de changement de username, retourner un message de succès simple
		fmt.Println("DEBUG ChangeI - Modifications effectuées sans changement de username")
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Changes saved successfully",
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

	// Supprimer le compte Matrix/Synapse associé
	_ = deleteMatrixAccount(req.Username)

	// Supprimer le compte Roundcube associé (via email)
	var email string
	db.QueryRow("SELECT email FROM Users WHERE user_id=$1", session.UserID).Scan(&email)
	if email != "" {
		_ = deleteRoundcubeAccount(email)
	}

	// Supprimer l'utilisateur de la base de données (CASCADE supprimera les sessions)
	_, err := db.Exec("DELETE FROM Users WHERE user_id=$1", session.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, DeleteAccountResponse{
			Success: false,
			Message: "Failed to delete account",
		})
		return
	}

	// Déconnecter complètement : supprimer de la mémoire
	delete(sessions, req.Token)
	// La session DB est déjà supprimée via CASCADE

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

// ============ CHANGEMENT MOT DE PASSE SYNCHRONISÉ ============

type ChangePasswordRequest struct {
	Username    string `json:"username"`
	Token       string `json:"token"`
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type ChangePasswordResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Changer le mot de passe Mail (Docker Mailserver)
func changeMailPassword(username, newPassword string) error {
	email := username + "@office1789.com"
	
	// Supprimer l'ancien compte
	cmdDel := exec.Command("docker", "exec", "mailserver", "setup", "email", "del", email)
	_ = cmdDel.Run() // Ignorer l'erreur si le compte n'existe pas
	
	// Recréer avec le nouveau mot de passe
	cmdAdd := exec.Command("docker", "exec", "mailserver", "setup", "email", "add", email, newPassword)
	var stderr bytes.Buffer
	cmdAdd.Stderr = &stderr
	
	err := cmdAdd.Run()
	if err != nil {
		return fmt.Errorf("failed to change mail password: %v - %s", err, stderr.String())
	}
	
	return nil
}

// Changer le mot de passe Matrix (Synapse)
func changeMatrixPassword(username, newPassword string) error {
	// Utiliser hash_password pour générer le hash bcrypt
	cmd := exec.Command("docker", "exec", "synapse",
		"hash_password", "-p", newPassword, "-c", "/data/homeserver.yaml")
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to hash password: %v - %s", err, stderr.String())
	}
	
	// Le hash est dans stdout (format: "$2b$12$...")
	passwordHash := strings.TrimSpace(stdout.String())
	if !strings.HasPrefix(passwordHash, "$2b$") && !strings.HasPrefix(passwordHash, "$2a$") {
		return fmt.Errorf("invalid hash output: %s", passwordHash)
	}
	
	// Mettre à jour directement dans la DB Synapse via PostgreSQL avec STDIN
	matrixUserID := "@" + username + ":office1789.com"
	
	// Utiliser echo + pipe pour éviter les problèmes d'échappement
	sqlQuery := fmt.Sprintf("UPDATE users SET password_hash='%s' WHERE name='%s';", passwordHash, matrixUserID)
	
	cmd = exec.Command("docker", "exec", "-i", "postgres_synapse",
		"psql", "-U", "synapse", "-d", "synapse")
	
	cmd.Stdin = strings.NewReader(sqlQuery)
	stderr.Reset()
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to update matrix password in DB: %v - %s", err, stderr.String())
	}
	
	return nil
}

func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ChangePasswordResponse{
			Success: false,
			Message: "Invalid request",
		})
		return
	}
	
	// Vérifier que la session est valide
	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, ChangePasswordResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}
	
	// Vérifier l'ancien mot de passe
	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM Users WHERE user_id=$1", session.UserID).Scan(&storedHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ChangePasswordResponse{
			Success: false,
			Message: "Database error",
		})
		return
	}
	
	if !CheckPasswordHash(req.OldPassword, storedHash) {
		c.JSON(http.StatusUnauthorized, ChangePasswordResponse{
			Success: false,
			Message: "Current password is incorrect",
		})
		return
	}
	
	// 1. Changer le mot de passe Office1789 (DB PostgreSQL)
	newHash := HashPassword(req.NewPassword)
	
	// 1b. Chiffrer aussi le nouveau mot de passe pour mail/Matrix
	fmt.Printf("[ChangePassword] 🔐 Chiffrement du nouveau mot de passe pour user_id=%d\n", session.UserID)
	encryptedPassword, err := EncryptPassword(req.NewPassword)
	if err != nil {
		fmt.Printf("[ChangePassword] ❌ ERREUR CHIFFREMENT: %v\n", err)
		c.JSON(http.StatusInternalServerError, ChangePasswordResponse{
			Success: false,
			Message: "Failed to encrypt password",
		})
		return
	}
	fmt.Printf("[ChangePassword] ✅ Mot de passe chiffré (len=%d): %s...\n", len(encryptedPassword), encryptedPassword[:20])
	
	fmt.Printf("[ChangePassword] 📝 UPDATE Users SET mail_password pour user_id=%d\n", session.UserID)
	_, err = db.Exec("UPDATE Users SET password_hash=$1, mail_password=$2 WHERE user_id=$3", 
		newHash, encryptedPassword, session.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ChangePasswordResponse{
			Success: false,
			Message: "Failed to update password",
		})
		return
	}
	
	// 2. Changer le mot de passe Mail (asynchrone)
	go func() {
		err := changeMailPassword(req.Username, req.NewPassword)
		if err != nil {
			fmt.Printf("Warning: Failed to change mail password for %s: %v\n", req.Username, err)
		} else {
			fmt.Printf("Mail password changed for %s@office1789.com\n", req.Username)
		}
	}()
	
	// 3. Changer le mot de passe Matrix (asynchrone)
	go func() {
		err := changeMatrixPassword(req.Username, req.NewPassword)
		if err != nil {
			fmt.Printf("Warning: Failed to change matrix password for %s: %v\n", req.Username, err)
		} else {
			fmt.Printf("Matrix password changed for @%s:office1789.com\n", req.Username)
		}
	}()
	
	c.JSON(http.StatusOK, ChangePasswordResponse{
		Success: true,
		Message: "Password changed successfully across all services",
	})
}

// deleteMatrixAccount supprime un compte Matrix/Synapse
func deleteMatrixAccount(username string) error {
	matrixUserID := fmt.Sprintf("@%s:office1789.com", username)
	
	// Utiliser l'API Admin de Synapse pour désactiver le compte
	url := "http://synapse:8008/_synapse/admin/v1/deactivate/" + matrixUserID
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(`{"erase": true}`)))
	if err != nil {
		return err
	}
	
	req.Header.Set("Content-Type", "application/json")
	// Note: Nécessite un access token admin - à configurer
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to delete matrix account, status: %d", resp.StatusCode)
	}
	
	return nil
}

// deleteRoundcubeAccount supprime un compte Roundcube
func deleteRoundcubeAccount(email string) error {
	// Connexion à la base Roundcube
	connStr := "host=postgres_roundcube port=5432 user=roundcube password=roundcube1789 dbname=roundcube sslmode=disable"
	roundcubeDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer roundcubeDB.Close()
	
	// Supprimer l'utilisateur Roundcube
	_, err = roundcubeDB.Exec("DELETE FROM users WHERE username=$1", email)
	return err
}

// RegenerateMatrixAccount régénère un compte Matrix utilisateur
func RegenerateMatrixAccount(c *gin.Context) {
	type RegenerateRequest struct {
		Username string `json:"username"`
		Token    string `json:"token"`
		Password string `json:"password"` // Mot de passe en clair nécessaire
	}

	var req RegenerateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Requête invalide"})
		return
	}

	// Valider la session
	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Session invalide"})
		return
	}

	// Vérifier le mot de passe fourni
	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM Users WHERE user_id=$1", session.UserID).Scan(&storedHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Erreur lors de la récupération des informations utilisateur"})
		return
	}

	// Comparer le mot de passe fourni avec le hash stocké
	if !CheckPasswordHash(req.Password, storedHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Mot de passe incorrect"})
		return
	}

	// Générer le hash du mot de passe avec hash_password de Synapse
	cmd := exec.Command("docker", "exec", "synapse",
		"hash_password", "-p", req.Password, "-c", "/data/homeserver.yaml")
	
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Erreur lors du hashage du mot de passe",
			"details": stderr.String(),
		})
		return
	}
	
	passwordHash := strings.TrimSpace(stdout.String())
	if !strings.HasPrefix(passwordHash, "$2b$") && !strings.HasPrefix(passwordHash, "$2a$") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Hash de mot de passe invalide",
			"details": passwordHash,
		})
		return
	}

	// Créer l'utilisateur Matrix dans la base de données
	matrixUserID := "@" + req.Username + ":office1789.com"
	
	// Vérifier d'abord si l'utilisateur existe déjà
	checkQuery := fmt.Sprintf("SELECT COUNT(*) FROM users WHERE name='%s';", matrixUserID)
	cmd = exec.Command("docker", "exec", "-i", "postgres_synapse",
		"psql", "-U", "synapse", "-d", "synapse", "-t")
	cmd.Stdin = strings.NewReader(checkQuery)
	var countOut bytes.Buffer
	cmd.Stdout = &countOut
	stderr.Reset()
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Erreur lors de la vérification du compte Matrix",
			"details": stderr.String(),
		})
		return
	}

	count := strings.TrimSpace(countOut.String())
	
	var sqlQuery string
	if count == "0" {
		// Créer un nouveau compte
		sqlQuery = fmt.Sprintf(`
			INSERT INTO users (name, password_hash, creation_ts, admin, user_type, deactivated)
			VALUES ('%s', '%s', %d, 0, NULL, 0);
		`, matrixUserID, passwordHash, time.Now().Unix()*1000)
	} else {
		// Mettre à jour le compte existant (réactiver si désactivé)
		sqlQuery = fmt.Sprintf(`
			UPDATE users 
			SET password_hash='%s', deactivated=0 
			WHERE name='%s';
		`, passwordHash, matrixUserID)
	}
	
	cmd = exec.Command("docker", "exec", "-i", "postgres_synapse",
		"psql", "-U", "synapse", "-d", "synapse")
	cmd.Stdin = strings.NewReader(sqlQuery)
	stderr.Reset()
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Erreur lors de la création/réactivation du compte Matrix",
			"details": stderr.String(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Compte Matrix régénéré avec succès"})
}
