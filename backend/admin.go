package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminStatsResponse représente les statistiques pour l'admin
type AdminStatsResponse struct {
	TotalUsers           int `json:"total_users"`
	TotalVerifiedEmails  int `json:"total_verified_emails"`
	TotalVerifiedPhones  int `json:"total_verified_phones"`
	UsersWithoutContacts int `json:"users_without_contacts"`
	TotalFiles           int `json:"total_files"`
	TotalCalendarEvents  int `json:"total_calendar_events"`
}

// AdminUserListResponse représente la liste des utilisateurs
type AdminUserListResponse struct {
	UserID              int    `json:"user_id"`
	Username            string `json:"username"`
	Email               string `json:"email"`
	RecoveryEmail       string `json:"recovery_email"`
	RecoveryEmailVerified bool `json:"recovery_email_verified"`
	PhoneNumber         string `json:"phonenumber"`
	PhoneNumberVerified bool   `json:"phonenumber_verified"`
	Role                string `json:"role"`
	Nboffer             int    `json:"nboffer"`
	DateJoined          string `json:"date_joined"`
	LastLogin           string `json:"last_login"`
}

// CheckAdminMiddleware vérifie si l'utilisateur est admin
func CheckAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization token"})
			c.Abort()
			return
		}

		sess, exists := sessions[tokenString]
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session"})
			c.Abort()
			return
		}

		// Vérifier le rôle admin
		var role string
		err := db.QueryRow("SELECT role FROM Users WHERE user_id=$1", sess.UserID).Scan(&role)
		if err != nil || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Set("userID", sess.UserID)
		c.Set("username", sess.Username)
		c.Next()
	}
}

// GetAdminStats retourne les statistiques globales
func GetAdminStats(c *gin.Context) {
	var stats AdminStatsResponse

	// Total users
	db.QueryRow("SELECT COUNT(*) FROM Users").Scan(&stats.TotalUsers)

	// Users with verified emails
	db.QueryRow("SELECT COUNT(*) FROM Users WHERE recovery_email_verified = true").Scan(&stats.TotalVerifiedEmails)

	// Users with verified phones
	db.QueryRow("SELECT COUNT(*) FROM Users WHERE phonenumber_verified = true").Scan(&stats.TotalVerifiedPhones)

	// Users without verified contacts
	db.QueryRow("SELECT COUNT(*) FROM Users WHERE recovery_email_verified = false AND phonenumber_verified = false").Scan(&stats.UsersWithoutContacts)

	// Total files
	db.QueryRow("SELECT COUNT(*) FROM DriveFiles").Scan(&stats.TotalFiles)

	// Total calendar events
	db.QueryRow("SELECT COUNT(*) FROM CalendarEvents").Scan(&stats.TotalCalendarEvents)

	c.JSON(http.StatusOK, stats)
}

// GetAllUsers retourne la liste de tous les utilisateurs
func GetAllUsers(c *gin.Context) {
	rows, err := db.Query(`
		SELECT user_id, username, email, COALESCE(recovery_email, ''), 
		       COALESCE(recovery_email_verified, false),
		       COALESCE(phonenumber, ''), COALESCE(phonenumber_verified, false),
		       COALESCE(role, 'user'), COALESCE(nboffer, 0),
		       date_joined, COALESCE(last_login, date_joined)
		FROM Users
		ORDER BY date_joined DESC
	`)

	if err != nil {
		fmt.Printf("Error fetching users: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []AdminUserListResponse
	for rows.Next() {
		var user AdminUserListResponse
		err := rows.Scan(
			&user.UserID,
			&user.Username,
			&user.Email,
			&user.RecoveryEmail,
			&user.RecoveryEmailVerified,
			&user.PhoneNumber,
			&user.PhoneNumberVerified,
			&user.Role,
			&user.Nboffer,
			&user.DateJoined,
			&user.LastLogin,
		)
		if err != nil {
			fmt.Printf("Error scanning user: %v\n", err)
			continue
		}
		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUserRole met à jour le rôle d'un utilisateur
type UpdateUserRoleRequest struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"` // 'user' ou 'admin'
}

func UpdateUserRole(c *gin.Context) {
	var req UpdateUserRoleRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Valider le rôle
	if req.Role != "user" && req.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be 'user' or 'admin'"})
		return
	}

	// Mettre à jour
	_, err := db.Exec("UPDATE Users SET role=$1 WHERE user_id=$2", req.Role, req.UserID)
	if err != nil {
		fmt.Printf("Error updating role: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role updated successfully",
	})
}

// VerifyUserContact permet à l'admin de marquer manuellement un contact comme vérifié
type VerifyUserContactRequest struct {
	UserID      int    `json:"user_id"`
	ContactType string `json:"contact_type"` // 'email' ou 'phone'
}

func VerifyUserContact(c *gin.Context) {
	var req VerifyUserContactRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if req.ContactType == "email" {
		_, err := db.Exec("UPDATE Users SET recovery_email_verified=true WHERE user_id=$1", req.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
			return
		}
	} else if req.ContactType == "phone" {
		_, err := db.Exec("UPDATE Users SET phonenumber_verified=true WHERE user_id=$1", req.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify phone"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contact_type must be 'email' or 'phone'"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Contact verified successfully",
	})
}
