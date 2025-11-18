package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// PasswordResetRequest structure pour demande de réinitialisation
type PasswordResetRequest struct {
	Identifier string `json:"identifier"` // username ou email
}

// PasswordResetConfirm structure pour confirmer la réinitialisation
type PasswordResetConfirm struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// RequestPasswordReset génère un token et envoie l'email de réinitialisation
func RequestPasswordReset(c *gin.Context) {
	var req PasswordResetRequest

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Requête invalide",
		})
		return
	}

	if req.Identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Veuillez entrer votre identifiant ou email",
		})
		return
	}

	// Rechercher l'utilisateur par username ou recovery_email
	var userID int
	var username string
	var recoveryEmail string

	// Essayer d'abord par username
	err := db.QueryRow(`
		SELECT user_id, username, COALESCE(recovery_email, '') 
		FROM Users 
		WHERE username = $1
	`, req.Identifier).Scan(&userID, &username, &recoveryEmail)

	// Si pas trouvé par username, essayer par recovery_email
	if err != nil {
		err = db.QueryRow(`
			SELECT user_id, username, recovery_email 
			FROM Users 
			WHERE recovery_email = $1 AND recovery_email IS NOT NULL AND recovery_email != ''
		`, req.Identifier).Scan(&userID, &username, &recoveryEmail)

		if err != nil {
			// Pour des raisons de sécurité, on ne dit pas si l'utilisateur existe ou non
			fmt.Printf("User not found for identifier: %s\n", req.Identifier)
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Si cet utilisateur existe, un email a été envoyé",
			})
			return
		}
	}

	// Vérifier que l'utilisateur a un recovery_email
	if recoveryEmail == "" {
		fmt.Printf("User %s has no recovery email configured\n", username)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Si cet utilisateur existe, un email a été envoyé",
		})
		return
	}

	// Générer un token unique
	token := generateResetToken()

	// Sauvegarder le token dans la DB avec expiration (1 heure)
	expiresAt := time.Now().Add(1 * time.Hour)

	_, err = db.Exec(`
		INSERT INTO password_reset_tokens (user_id, token, expires_at, used)
		VALUES ($1, $2, $3, false)
		ON CONFLICT (user_id) 
		DO UPDATE SET token = $2, expires_at = $3, used = false, created_at = NOW()
	`, userID, token, expiresAt)

	if err != nil {
		fmt.Printf("Error saving reset token: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Erreur lors de la génération du lien",
		})
		return
	}

	// Envoyer l'email
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	resetLink := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	err = sendResetEmail(recoveryEmail, username, resetLink)
	if err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Erreur lors de l'envoi de l'email",
		})
		return
	}

	fmt.Printf("Password reset email sent to %s for user %s\n", email, username)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Un email de réinitialisation a été envoyé",
	})
}

// ResetPassword réinitialise le mot de passe avec le token
func ResetPassword(c *gin.Context) {
	var req PasswordResetConfirm

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Requête invalide",
		})
		return
	}

	if req.Token == "" || req.NewPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Token et nouveau mot de passe requis",
		})
		return
	}

	// Valider le token
	var userID int
	var expiresAt time.Time
	var used bool

	err := db.QueryRow(`
		SELECT user_id, expires_at, used 
		FROM password_reset_tokens 
		WHERE token = $1
	`, req.Token).Scan(&userID, &expiresAt, &used)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Token invalide ou expiré",
		})
		return
	}

	// Vérifier si le token a déjà été utilisé
	if used {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Ce lien a déjà été utilisé",
		})
		return
	}

	// Vérifier si le token a expiré
	if time.Now().After(expiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Ce lien a expiré. Veuillez faire une nouvelle demande",
		})
		return
	}

	// Récupérer le username pour la synchronisation Mail/Matrix
	var username string
	err = db.QueryRow("SELECT username FROM Users WHERE user_id=$1", userID).Scan(&username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Erreur lors de la mise à jour",
		})
		return
	}

	// Hasher le nouveau mot de passe
	hashedPassword, err := HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Erreur lors du hashage",
		})
		return
	}

	// Chiffrer le mot de passe pour Mail/Matrix
	encryptedPassword, err := EncryptPassword(req.NewPassword)
	if err != nil {
		fmt.Printf("Error encrypting password for Mail/Matrix: %v\n", err)
		encryptedPassword = "" // Continue sans Mail/Matrix en cas d'erreur
	}

	// Mettre à jour le mot de passe dans la DB
	_, err = db.Exec(`
		UPDATE Users 
		SET password_hash = $1, mail_password = $2 
		WHERE user_id = $3
	`, hashedPassword, encryptedPassword, userID)

	if err != nil {
		fmt.Printf("Error updating password: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Erreur lors de la mise à jour du mot de passe",
		})
		return
	}

	// Synchroniser avec le serveur Mail (postfix-accounts.cf)
	err = UpdateMailPassword(username, req.NewPassword)
	if err != nil {
		fmt.Printf("Warning: Failed to sync mail password for %s: %v\n", username, err)
		// Continue quand même, le mot de passe Office1789 est mis à jour
	}

	// Marquer le token comme utilisé
	_, err = db.Exec(`
		UPDATE password_reset_tokens 
		SET used = true 
		WHERE token = $1
	`, req.Token)

	if err != nil {
		fmt.Printf("Error marking token as used: %v\n", err)
		// Non bloquant
	}

	fmt.Printf("Password reset successful for user %s (ID: %d)\n", username, userID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Mot de passe réinitialisé avec succès",
	})
}

// generateResetToken génère un token aléatoire sécurisé
func generateResetToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// sendResetEmail envoie l'email de réinitialisation
func sendResetEmail(to, username, resetLink string) error {
	// Configuration SMTP (à adapter selon votre serveur)
	smtpHost := "localhost"
	smtpPort := "25"
	from := "noreply@office1789.com"

	// Message HTML
	subject := "Réinitialisation de votre mot de passe Office1789"
	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: 'Segoe UI', Arial, sans-serif; background-color: #f5f5f7; margin: 0; padding: 20px; }
        .container { max-width: 600px; margin: 0 auto; background: white; border-radius: 16px; box-shadow: 0 4px 12px rgba(0,0,0,0.1); overflow: hidden; }
        .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 40px 20px; text-align: center; }
        .header h1 { margin: 0; font-size: 28px; font-weight: 700; }
        .content { padding: 40px 30px; }
        .content p { color: #333; line-height: 1.6; margin: 16px 0; }
        .button { display: inline-block; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 16px 40px; border-radius: 12px; text-decoration: none; font-weight: 600; margin: 24px 0; }
        .link { color: #667eea; word-break: break-all; font-size: 14px; }
        .footer { background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }
        .warning { background: #fff3cd; border-left: 4px solid #ffc107; padding: 12px; margin: 20px 0; border-radius: 4px; color: #856404; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Office1789</h1>
            <p>Réinitialisation de mot de passe</p>
        </div>
        <div class="content">
            <p>Bonjour <strong>%s</strong>,</p>
            <p>Vous avez demandé la réinitialisation de votre mot de passe Office1789.</p>
            <p>Cliquez sur le bouton ci-dessous pour créer un nouveau mot de passe :</p>
            <p style="text-align: center;">
                <a href="%s" class="button">Réinitialiser mon mot de passe</a>
            </p>
            <p>Ou copiez ce lien dans votre navigateur :</p>
            <p class="link">%s</p>
            <div class="warning">
                <strong>⚠️ Important :</strong> Ce lien est valable pendant 1 heure et ne peut être utilisé qu'une seule fois.
            </div>
            <p>Si vous n'avez pas demandé cette réinitialisation, ignorez cet email. Votre mot de passe restera inchangé.</p>
        </div>
        <div class="footer">
            <p>© 2025 Office1789 - Solution collaborative française</p>
            <p>Cet email a été envoyé automatiquement, merci de ne pas y répondre.</p>
        </div>
    </div>
</body>
</html>
`, username, resetLink, resetLink)

	// Construire le message SMTP
	message := fmt.Sprintf("From: Office1789 <%s>\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// Envoyer l'email (sans authentification pour localhost)
	addr := smtpHost + ":" + smtpPort
	err := smtp.SendMail(addr, nil, from, []string{to}, []byte(message))

	return err
}
