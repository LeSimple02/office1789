package main

import (
	"bytes"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"image/png"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type TOTPRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	Password string `json:"password,omitempty"`
	Code     string `json:"code,omitempty"`
}

// Enable2FA generates a TOTP secret and returns QR code
func Enable2FA(c *gin.Context) {
	var req TOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalide"})
		return
	}

	userID := session.UserID
	username := req.Username

	// Check if 2FA is already enabled
	var existing bool
	err := db.QueryRow("SELECT enabled FROM user_totp WHERE user_id = $1", userID).Scan(&existing)
	if err == nil && existing {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA déjà activé"})
		return
	}

	// Generate TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Office1789",
		AccountName: username,
		SecretSize:  32,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur génération secret"})
		return
	}

	// Generate QR code
	var buf bytes.Buffer
	img, err := key.Image(256, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur génération QR code"})
		return
	}
	png.Encode(&buf, img)

	// Generate backup codes (8 codes of 8 characters each)
	backupCodes := make([]string, 8)
	hashedBackupCodes := make([]string, 8)
	for i := 0; i < 8; i++ {
		code := generateBackupCode()
		backupCodes[i] = code
		hashed, _ := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		hashedBackupCodes[i] = string(hashed)
	}

	// Save to database (enabled = false until verified)
	_, err = db.Exec(`
		INSERT INTO user_totp (user_id, secret, enabled, backup_codes, created_at)
		VALUES ($1, $2, false, $3, NOW())
		ON CONFLICT (user_id) DO UPDATE 
		SET secret = $2, backup_codes = $3, enabled = false, created_at = NOW()
	`, userID, key.Secret(), pq.Array(hashedBackupCodes))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur sauvegarde"})
		return
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	c.JSON(http.StatusOK, gin.H{
		"qr_code":      qrCodeBase64,
		"secret":       key.Secret(),
		"backup_codes": backupCodes,
	})
}

// Verify2FA validates the TOTP code and enables 2FA
func Verify2FA(c *gin.Context) {
	var req TOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalide"})
		return
	}

	userID := session.UserID

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code requis"})
		return
	}

	// Get secret from database
	var secret string
	err := db.QueryRow("SELECT secret FROM user_totp WHERE user_id = $1", userID).Scan(&secret)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA non configuré"})
		return
	}

	// Validate TOTP code
	valid = totp.Validate(req.Code, secret)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code invalide"})
		return
	}

	// Enable 2FA
	_, err = db.Exec("UPDATE user_totp SET enabled = true, last_used = NOW() WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur activation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA activé avec succès"})
}

// Disable2FA disables 2FA after password verification
func Disable2FA(c *gin.Context) {
	var req TOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalide"})
		return
	}

	userID := session.UserID

	if req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mot de passe requis"})
		return
	}

	// Verify password
	var passwordHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE user_id = $1", userID).Scan(&passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur vérification"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
		return
	}

	// Delete 2FA configuration
	_, err = db.Exec("DELETE FROM user_totp WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur désactivation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA désactivé"})
}

// Get2FAStatus returns the current 2FA status
func Get2FAStatus(c *gin.Context) {
	var req TOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalide"})
		return
	}

	userID := session.UserID

	var enabled bool
	err := db.QueryRow("SELECT enabled FROM user_totp WHERE user_id = $1", userID).Scan(&enabled)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"enabled": false})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur récupération statut"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enabled": enabled})
}

// RegenerateBackupCodes generates new backup codes
func RegenerateBackupCodes(c *gin.Context) {
	var req TOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	session, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Session invalide"})
		return
	}

	userID := session.UserID

	// Check if 2FA is enabled
	var enabled bool
	err := db.QueryRow("SELECT enabled FROM user_totp WHERE user_id = $1", userID).Scan(&enabled)
	if err != nil || !enabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "2FA non activé"})
		return
	}

	// Generate new backup codes
	backupCodes := make([]string, 8)
	hashedBackupCodes := make([]string, 8)
	for i := 0; i < 8; i++ {
		code := generateBackupCode()
		backupCodes[i] = code
		hashed, _ := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
		hashedBackupCodes[i] = string(hashed)
	}

	// Update database
	_, err = db.Exec("UPDATE user_totp SET backup_codes = $1 WHERE user_id = $2", pq.Array(hashedBackupCodes), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur génération codes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"backup_codes": backupCodes,
	})
}

// generateBackupCode generates a random 8-character backup code
func generateBackupCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	rand.Read(b)
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

// ValidateTOTPLogin validates TOTP code during login
func ValidateTOTPLogin(userID int, code string) bool {
	var secret string
	var enabled bool
	var backupCodesArray []string

	err := db.QueryRow("SELECT secret, enabled, backup_codes FROM user_totp WHERE user_id = $1", userID).
		Scan(&secret, &enabled, pq.Array(&backupCodesArray))

	if err != nil || !enabled {
		return false
	}

	// First, try TOTP validation
	if totp.Validate(code, secret) {
		// Update last_used timestamp
		db.Exec("UPDATE user_totp SET last_used = NOW() WHERE user_id = $1", userID)
		return true
	}

	// If TOTP fails, check backup codes
	for i, hashedCode := range backupCodesArray {
		if bcrypt.CompareHashAndPassword([]byte(hashedCode), []byte(code)) == nil {
			// Backup code is valid, remove it from the array
			backupCodesArray = append(backupCodesArray[:i], backupCodesArray[i+1:]...)
			db.Exec("UPDATE user_totp SET backup_codes = $1, last_used = NOW() WHERE user_id = $2",
				pq.Array(backupCodesArray), userID)
			return true
		}
	}

	return false
}

// Check2FARequired checks if user has 2FA enabled (used during login)
func Check2FARequired(userID int) bool {
	var enabled bool
	err := db.QueryRow("SELECT enabled FROM user_totp WHERE user_id = $1", userID).Scan(&enabled)
	return err == nil && enabled
}
