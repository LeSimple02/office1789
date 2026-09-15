package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Organization limits per offer
const (
	SubAccountsStandard     = 0  // Cannot create sub-accounts
	SubAccountsProfessional = 3  // Can have 3 sub-accounts
	SubAccountsEnterprise   = 20 // Can have 20 sub-accounts
)

// getSubAccountsLimit returns the number of sub-accounts allowed based on user's offer
func getSubAccountsLimit(nboffer int) int {
	switch nboffer {
	case 0, 1:
		return SubAccountsStandard // Free and Standard: no sub-accounts
	case 2:
		return SubAccountsProfessional
	case 3:
		return SubAccountsEnterprise
	default:
		return SubAccountsStandard
	}
}

// CreateSubAccountRequest structure
type CreateSubAccountRequest struct {
	Username         string `json:"username"`
	Token            string `json:"token"`
	SubUsername      string `json:"sub_username"`
	SubPassword      string `json:"sub_password"`
	SubEmail         string `json:"sub_email"`         // Optional recovery email
	SubPhoneNumber   string `json:"sub_phone_number"`  // Optional
	OrganizationName string `json:"organization_name"` // Optional, for first sub-account
}

// SubAccountResponse structure
type SubAccountResponse struct {
	Success          bool   `json:"success"`
	Message          string `json:"message"`
	SubAccountID     int    `json:"sub_account_id,omitempty"`
	SubAccountEmail  string `json:"sub_account_email,omitempty"`
	OrganizationID   int    `json:"organization_id,omitempty"`
	OrganizationName string `json:"organization_name,omitempty"`
	CurrentCount     int    `json:"current_count,omitempty"`
	MaxAllowed       int    `json:"max_allowed,omitempty"`
}

// CreateSubAccount creates a sub-account for Professional/Enterprise users
func CreateSubAccount(c *gin.Context) {
	var req CreateSubAccountRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, SubAccountResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate session
	sess, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, SubAccountResponse{
			Success: false,
			Message: "Invalid session",
		})
		return
	}

	// Check if sub-username is reserved
	if isReservedUsername(req.SubUsername) {
		c.JSON(http.StatusBadRequest, SubAccountResponse{
			Success: false,
			Message: "Ce nom d'utilisateur est réservé et ne peut pas être utilisé",
		})
		return
	}

	// Get user's offer and check sub-account limit
	var nboffer int
	var accountType string
	var organizationID sql.NullInt64
	err := db.QueryRow(`
		SELECT nboffer, COALESCE(account_type, 'personal'), organization_id 
		FROM Users WHERE user_id=$1
	`, sess.UserID).Scan(&nboffer, &accountType, &organizationID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, SubAccountResponse{
			Success: false,
			Message: "Failed to get user information",
		})
		return
	}

	maxSubAccounts := getSubAccountsLimit(nboffer)
	if maxSubAccounts == 0 {
		c.JSON(http.StatusForbidden, SubAccountResponse{
			Success:      false,
			Message:      "Sub-accounts feature is only available for Professional and Enterprise plans. Please upgrade your account.",
			MaxAllowed:   maxSubAccounts,
			CurrentCount: 0,
		})
		return
	}

	// Count current sub-accounts
	var currentSubAccounts int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM Users 
		WHERE parent_account_id=$1
	`, sess.UserID).Scan(&currentSubAccounts)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, SubAccountResponse{
			Success: false,
			Message: "Failed to count existing sub-accounts",
		})
		return
	}

	if currentSubAccounts >= maxSubAccounts {
		c.JSON(http.StatusForbidden, SubAccountResponse{
			Success:      false,
			Message:      fmt.Sprintf("Sub-account limit reached (%d/%d). Upgrade to create more sub-accounts.", currentSubAccounts, maxSubAccounts),
			CurrentCount: currentSubAccounts,
			MaxAllowed:   maxSubAccounts,
		})
		return
	}

	// Check if sub_username already exists
	var exists int
	err = db.QueryRow("SELECT COUNT(*) FROM Users WHERE username=$1", req.SubUsername).Scan(&exists)
	if err != nil || exists > 0 {
		c.JSON(http.StatusConflict, SubAccountResponse{
			Success: false,
			Message: "Username already taken",
		})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.SubPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SubAccountResponse{
			Success: false,
			Message: "Failed to hash password",
		})
		return
	}

	// Create or get organization
	var orgID int
	var orgName string
	
	if !organizationID.Valid {
		// First sub-account: create organization
		if req.OrganizationName == "" {
			req.OrganizationName = req.Username + " Organization"
		}
		
		err = db.QueryRow(`
			INSERT INTO Organizations (organization_name, owner_user_id, max_members, created_at)
			VALUES ($1, $2, $3, NOW())
			RETURNING organization_id
		`, req.OrganizationName, sess.UserID, maxSubAccounts).Scan(&orgID)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, SubAccountResponse{
				Success: false,
				Message: "Failed to create organization: " + err.Error(),
			})
			return
		}

		// Update parent account to link to organization
		_, err = db.Exec(`
			UPDATE Users 
			SET organization_id=$1, account_type='organization_owner' 
			WHERE user_id=$2
		`, orgID, sess.UserID)
		
		if err != nil {
			c.JSON(http.StatusInternalServerError, SubAccountResponse{
				Success: false,
				Message: "Failed to update parent account",
			})
			return
		}
		
		orgName = req.OrganizationName
	} else {
		// Use existing organization
		orgID = int(organizationID.Int64)
		err = db.QueryRow("SELECT organization_name FROM Organizations WHERE organization_id=$1", orgID).Scan(&orgName)
		if err != nil {
			orgName = "Organization"
		}
	}

	// Create sub-account
	subAccountEmail := req.SubUsername + "@office1789.com"
	var subAccountID int
	
	// Sub-accounts inherit the parent's plan to access organization features
	// They share the same storage and features as the organization
	var parentOffer int
	err = db.QueryRow(`SELECT nboffer FROM Users WHERE user_id=$1`, sess.UserID).Scan(&parentOffer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, SubAccountResponse{
			Success: false,
			Message: "Failed to get parent plan: " + err.Error(),
		})
		return
	}
	
	err = db.QueryRow(`
		INSERT INTO Users (
			username, password_hash, email, recovery_email, phonenumber,
			nboffer, domain, organization_id, parent_account_id, account_type,
			date_joined, last_login
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING user_id
	`, req.SubUsername, string(hashedPassword), subAccountEmail, req.SubEmail, 
	   req.SubPhoneNumber, parentOffer, "@office1789", orgID, sess.UserID, "organization_member").Scan(&subAccountID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, SubAccountResponse{
			Success: false,
			Message: "Failed to create sub-account: " + err.Error(),
		})
		return
	}

	fmt.Printf("✅ Sub-account created: %s (ID: %d) under parent %s (ID: %d), Organization: %s (ID: %d)\n",
		req.SubUsername, subAccountID, req.Username, sess.UserID, orgName, orgID)

	c.JSON(http.StatusOK, SubAccountResponse{
		Success:          true,
		Message:          "Sub-account created successfully",
		SubAccountID:     subAccountID,
		SubAccountEmail:  subAccountEmail,
		OrganizationID:   orgID,
		OrganizationName: orgName,
		CurrentCount:     currentSubAccounts + 1,
		MaxAllowed:       maxSubAccounts,
	})
}

// OrganizationMember structure
type OrganizationMember struct {
	UserID        int       `json:"user_id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	AccountType   string    `json:"account_type"`
	DateJoined    time.Time `json:"date_joined"`
	LastLogin     time.Time `json:"last_login"`
	RecoveryEmail string    `json:"recovery_email,omitempty"`
	PhoneNumber   string    `json:"phone_number,omitempty"`
}

// GetOrganizationMembers returns all members of the organization
func GetOrganizationMembers(c *gin.Context) {
	type reqT struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}
	
	var req reqT
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	sess, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid session"})
		return
	}

	// Get organization info
	var orgID sql.NullInt64
	var orgName string
	var accountType string
	var maxMembers int
	
	err := db.QueryRow(`
		SELECT u.organization_id, u.account_type, COALESCE(o.organization_name, ''), COALESCE(o.max_members, 0)
		FROM Users u
		LEFT JOIN Organizations o ON u.organization_id = o.organization_id
		WHERE u.user_id=$1
	`, sess.UserID).Scan(&orgID, &accountType, &orgName, &maxMembers)
	
	if err != nil || !orgID.Valid {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false, 
			"message": "No organization found for this account",
		})
		return
	}

	// Get all members of the organization
	rows, err := db.Query(`
		SELECT user_id, username, email, account_type, date_joined, last_login,
		       COALESCE(recovery_email, ''), COALESCE(phonenumber, '')
		FROM Users
		WHERE organization_id=$1
		ORDER BY account_type DESC, date_joined ASC
	`, orgID.Int64)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get members"})
		return
	}
	defer rows.Close()

	var members []OrganizationMember
	for rows.Next() {
		var m OrganizationMember
		err := rows.Scan(&m.UserID, &m.Username, &m.Email, &m.AccountType, 
			&m.DateJoined, &m.LastLogin, &m.RecoveryEmail, &m.PhoneNumber)
		if err == nil {
			members = append(members, m)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success":           true,
		"organization_id":   orgID.Int64,
		"organization_name": orgName,
		"max_members":       maxMembers,
		"current_count":     len(members),
		"members":           members,
	})
}

// DeleteSubAccount deletes a sub-account (only owner can do this)
func DeleteSubAccount(c *gin.Context) {
	type reqT struct {
		Username      string `json:"username"`
		Token         string `json:"token"`
		SubAccountID  int    `json:"sub_account_id"`
	}
	
	var req reqT
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request"})
		return
	}

	sess, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid session"})
		return
	}

	// Verify that the sub-account belongs to this parent
	var parentID int
	err := db.QueryRow("SELECT parent_account_id FROM Users WHERE user_id=$1", req.SubAccountID).Scan(&parentID)
	
	if err != nil || parentID != sess.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "You can only delete your own sub-accounts",
		})
		return
	}

	// Delete the sub-account (CASCADE will handle related data)
	_, err = db.Exec("DELETE FROM Users WHERE user_id=$1", req.SubAccountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete sub-account",
		})
		return
	}

	fmt.Printf("✅ Sub-account deleted: ID %d by parent %s (ID: %d)\n", req.SubAccountID, req.Username, sess.UserID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Sub-account deleted successfully",
	})
}
