package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// SendVerificationCodeRequest représente une demande d'envoi de code
type SendVerificationCodeRequest struct {
	Contact string `json:"contact"` // Email ou téléphone
	Type    string `json:"type"`    // "email" ou "phone"
}

// VerifyCodeRequest représente une demande de vérification de code
type VerifyCodeRequest struct {
	Contact string `json:"contact"` // Email ou téléphone
	Code    string `json:"code"`    // Code à 6 chiffres
	Type    string `json:"type"`    // "email" ou "phone"
}

// SendVerificationCode envoie un code de vérification par email ou SMS
func SendVerificationCode(c *gin.Context) {
	var req SendVerificationCodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Données invalides",
		})
		return
	}

	// Valider le type
	if req.Type != "email" && req.Type != "phone" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Type invalide (email ou phone)",
		})
		return
	}

	// Vérifier que le contact n'est pas déjà utilisé par un autre utilisateur
	var count int
	if req.Type == "email" {
		db.QueryRow("SELECT COUNT(*) FROM Users WHERE recovery_email=$1", req.Contact).Scan(&count)
	} else {
		db.QueryRow("SELECT COUNT(*) FROM Users WHERE phonenumber=$1", req.Contact).Scan(&count)
	}

	if count > 0 {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Ce " + map[string]string{"email": "email", "phone": "numéro"}[req.Type] + " est déjà utilisé",
		})
		return
	}

	// Générer un code à 6 chiffres
	code := generateVerificationCode()

	// Sauvegarder dans la DB avec expiration (10 minutes)
	expiresAt := time.Now().Add(10 * time.Minute)

	// Invalider les anciens codes pour ce contact
	_, err := db.Exec(`
		UPDATE verification_codes 
		SET verified = true 
		WHERE contact = $1 AND verified = false
	`, req.Contact)

	if err != nil {
		fmt.Printf("Error invalidating old codes: %v\n", err)
	}

	// Créer un nouveau code
	_, err = db.Exec(`
		INSERT INTO verification_codes (contact, code, type, expires_at)
		VALUES ($1, $2, $3, $4)
	`, req.Contact, code, req.Type, expiresAt)

	if err != nil {
		fmt.Printf("Error saving verification code: %v\n", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "Erreur lors de la génération du code",
		})
		return
	}

	// Envoyer le code
	if req.Type == "email" {
		err = sendVerificationEmail(req.Contact, code)
		if err != nil {
			fmt.Printf("Error sending verification email: %v\n", err)
			c.JSON(500, gin.H{
				"success": false,
				"message": "Erreur lors de l'envoi de l'email",
			})
			return
		}
	} else {
		// Pour le SMS, on pourrait intégrer un service comme Twilio
		// Pour l'instant, on log juste le code
		fmt.Printf("SMS Verification code for %s: %s\n", req.Contact, code)
		// TODO: Intégrer un service SMS
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Code de vérification envoyé",
		"expires_in": 600, // 10 minutes en secondes
	})
}

// VerifyCode vérifie un code de vérification
func VerifyCode(c *gin.Context) {
	var req VerifyCodeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Données invalides",
		})
		return
	}

	// Chercher le code dans la DB
	var id int
	var expiresAt time.Time
	var verified bool

	err := db.QueryRow(`
		SELECT id, expires_at, verified
		FROM verification_codes
		WHERE contact = $1 AND code = $2 AND type = $3
		ORDER BY created_at DESC
		LIMIT 1
	`, req.Contact, req.Code, req.Type).Scan(&id, &expiresAt, &verified)

	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Code invalide",
		})
		return
	}

	// Vérifier si le code a déjà été utilisé
	if verified {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Ce code a déjà été utilisé",
		})
		return
	}

	// Vérifier si le code a expiré
	if time.Now().After(expiresAt) {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Ce code a expiré",
		})
		return
	}

	// Marquer le code comme vérifié
	_, err = db.Exec(`
		UPDATE verification_codes 
		SET verified = true 
		WHERE id = $1
	`, id)

	if err != nil {
		fmt.Printf("Error marking code as verified: %v\n", err)
		c.JSON(500, gin.H{
			"success": false,
			"message": "Erreur lors de la vérification",
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Code vérifié avec succès",
		"verified": true,
	})
}

// CheckVerificationStatus vérifie si un contact a été vérifié récemment (dans les 30 minutes)
func CheckVerificationStatus(contact string, contactType string) bool {
	var count int
	thirtyMinutesAgo := time.Now().Add(-30 * time.Minute)

	err := db.QueryRow(`
		SELECT COUNT(*) 
		FROM verification_codes 
		WHERE contact = $1 
		AND type = $2 
		AND verified = true 
		AND created_at > $3
	`, contact, contactType, thirtyMinutesAgo).Scan(&count)

	if err != nil {
		fmt.Printf("Error checking verification status: %v\n", err)
		return false
	}

	return count > 0
}

// generateVerificationCode génère un code aléatoire à 6 chiffres
func generateVerificationCode() string {
	max := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		// Fallback sur time-based si erreur
		return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

// sendVerificationEmail envoie un email avec le code de vérification
func sendVerificationEmail(to, code string) error {
	smtpHost := "localhost"
	smtpPort := "25"
	from := "noreply@office1789.com"

	subject := "Code de vérification Office1789"
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
        .content { padding: 40px 30px; text-align: center; }
        .content p { color: #333; line-height: 1.6; margin: 16px 0; }
        .code-box { background: #f0f0f5; border: 3px dashed #667eea; border-radius: 12px; padding: 30px; margin: 30px 0; }
        .code { font-size: 48px; font-weight: bold; letter-spacing: 10px; color: #667eea; font-family: 'Courier New', monospace; }
        .footer { background: #f8f9fa; padding: 20px; text-align: center; color: #666; font-size: 12px; }
        .warning { background: #fff3cd; border-left: 4px solid #ffc107; padding: 12px; margin: 20px 0; border-radius: 4px; color: #856404; text-align: left; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Office1789</h1>
            <p>Code de Vérification</p>
        </div>
        <div class="content">
            <p>Voici votre code de vérification :</p>
            <div class="code-box">
                <div class="code">%s</div>
            </div>
            <p>Entrez ce code dans l'application pour confirmer votre adresse email.</p>
            <div class="warning">
                <strong>⚠️ Important :</strong> Ce code expire dans 10 minutes et ne peut être utilisé qu'une seule fois.
            </div>
            <p>Si vous n'avez pas demandé ce code, ignorez cet email.</p>
        </div>
        <div class="footer">
            <p>© 2025 Office1789 - Solution collaborative française</p>
            <p>Cet email a été envoyé automatiquement, merci de ne pas y répondre.</p>
        </div>
    </div>
</body>
</html>
`, code)

	message := fmt.Sprintf("From: Office1789 <%s>\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	addr := smtpHost + ":" + smtpPort
	err := smtp.SendMail(addr, nil, from, []string{to}, []byte(message))

	return err
}
