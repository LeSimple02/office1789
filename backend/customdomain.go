package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// CustomDomainRequest structure
type CustomDomainRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	Domain   string `json:"domain"`
}

// CustomDomainResponse structure
type CustomDomainResponse struct {
	Success            bool   `json:"success"`
	Message            string `json:"message"`
	Domain             string `json:"domain,omitempty"`
	Verified           bool   `json:"verified"`
	VerificationToken  string `json:"verification_token,omitempty"`
	VerificationRecord string `json:"verification_record,omitempty"` // DNS TXT record to create
}

// generateVerificationToken generates a random token for domain verification
func generateVerificationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// isValidDomain checks if domain format is valid
func isValidDomain(domain string) bool {
	// Basic domain validation: letters, numbers, hyphens, dots
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	return domainRegex.MatchString(domain) && len(domain) <= 253
}

// AddCustomDomain adds a custom domain to user/organization (Professional/Enterprise only)
func AddCustomDomain(c *gin.Context) {
	var req CustomDomainRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, CustomDomainResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate session
	sess, valid := validateSession(req.Token, req.Username)
	if !valid {
		c.JSON(http.StatusUnauthorized, CustomDomainResponse{
			Success: false,
			Message: "Invalid session",
		})
		return
	}

	// Check user's offer (only Professional and Enterprise can use custom domains)
	var nboffer int
	var accountType string
	var organizationID sql.NullInt64
	
	err := db.QueryRow(`
		SELECT nboffer, COALESCE(account_type, 'personal'), organization_id 
		FROM Users WHERE user_id=$1
	`, sess.UserID).Scan(&nboffer, &accountType, &organizationID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, CustomDomainResponse{
			Success: false,
			Message: "Failed to get user information",
		})
		return
	}

	if nboffer < 2 {
		c.JSON(http.StatusForbidden, CustomDomainResponse{
			Success:  false,
			Message:  "Custom domains are only available for Professional and Enterprise plans. Please upgrade your account.",
			Verified: false,
		})
		return
	}

	// Validate domain format
	domain := strings.ToLower(strings.TrimSpace(req.Domain))
	if !isValidDomain(domain) {
		c.JSON(http.StatusBadRequest, CustomDomainResponse{
			Success: false,
			Message: "Invalid domain format. Please enter a valid domain (e.g., company.com)",
		})
		return
	}

	// Check if domain is already in use
	var existingCount int
	err = db.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT custom_domain FROM Users WHERE custom_domain=$1 AND domain_verified=TRUE
			UNION
			SELECT custom_domain FROM Organizations WHERE custom_domain=$1 AND domain_verified=TRUE
		) AS domains
	`, domain).Scan(&existingCount)
	
	if err == nil && existingCount > 0 {
		c.JSON(http.StatusConflict, CustomDomainResponse{
			Success: false,
			Message: "This domain is already in use by another account",
		})
		return
	}

	// Generate verification token
	token, err := generateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, CustomDomainResponse{
			Success: false,
			Message: "Failed to generate verification token",
		})
		return
	}

	// Store domain based on account type
	if accountType == "organization_owner" && organizationID.Valid {
		// Store in Organizations table
		_, err = db.Exec(`
			UPDATE Organizations 
			SET custom_domain=$1, domain_verified=FALSE, domain_verification_token=$2
			WHERE organization_id=$3
		`, domain, token, organizationID.Int64)
	} else {
		// Store in Users table (individual Pro/Enterprise account)
		_, err = db.Exec(`
			UPDATE Users 
			SET custom_domain=$1, domain_verified=FALSE
			WHERE user_id=$2
		`, domain, sess.UserID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, CustomDomainResponse{
			Success: false,
			Message: "Failed to save custom domain: " + err.Error(),
		})
		return
	}

	verificationRecord := fmt.Sprintf("office1789-verification=%s", token)
	
	fmt.Printf("✅ Custom domain added: %s for user %s (ID: %d), verification token: %s\n",
		domain, req.Username, sess.UserID, token)

	c.JSON(http.StatusOK, CustomDomainResponse{
		Success:            true,
		Message:            "Custom domain added successfully. Please verify domain ownership by adding a TXT record to your DNS.",
		Domain:             domain,
		Verified:           false,
		VerificationToken:  token,
		VerificationRecord: verificationRecord,
	})
}

// VerifyCustomDomain verifies domain ownership via DNS TXT record
func VerifyCustomDomain(c *gin.Context) {
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

	// Get domain and token
	var domain string
	var verificationToken string
	var accountType string
	var organizationID sql.NullInt64
	var alreadyVerified bool
	
	err := db.QueryRow(`
		SELECT COALESCE(account_type, 'personal'), organization_id 
		FROM Users WHERE user_id=$1
	`, sess.UserID).Scan(&accountType, &organizationID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get user info"})
		return
	}

	// Get domain and verification token based on account type
	if accountType == "organization_owner" && organizationID.Valid {
		err = db.QueryRow(`
			SELECT COALESCE(custom_domain, ''), COALESCE(domain_verification_token, ''), COALESCE(domain_verified, FALSE)
			FROM Organizations WHERE organization_id=$1
		`, organizationID.Int64).Scan(&domain, &verificationToken, &alreadyVerified)
	} else {
		err = db.QueryRow(`
			SELECT COALESCE(custom_domain, ''), COALESCE(domain_verified, FALSE)
			FROM Users WHERE user_id=$1
		`, sess.UserID).Scan(&domain, &alreadyVerified)
		
		// For individual accounts, we need to get the token from organization or generate it
		if domain != "" && !alreadyVerified {
			// Use a simple verification method for individual accounts
			verificationToken = fmt.Sprintf("office1789-user-%d", sess.UserID)
		}
	}
	
	if err != nil || domain == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "No custom domain configured for this account",
		})
		return
	}

	if alreadyVerified {
		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "Domain already verified",
			"domain":   domain,
			"verified": true,
		})
		return
	}

	// Lookup DNS TXT records
	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Failed to lookup DNS records for %s. Please ensure the TXT record is properly configured.", domain),
			"error":   err.Error(),
		})
		return
	}

	// Check if verification token is present in TXT records
	expectedRecord := fmt.Sprintf("office1789-verification=%s", verificationToken)
	verified := false
	
	for _, record := range txtRecords {
		if strings.Contains(record, expectedRecord) || strings.Contains(record, verificationToken) {
			verified = true
			break
		}
	}

	if !verified {
		c.JSON(http.StatusBadRequest, gin.H{
			"success":               false,
			"message":               "Domain verification failed. TXT record not found or incorrect.",
			"expected_txt_record":   expectedRecord,
			"found_txt_records":     txtRecords,
			"verification_token":    verificationToken,
		})
		return
	}

	// Update domain as verified
	if accountType == "organization_owner" && organizationID.Valid {
		_, err = db.Exec(`
			UPDATE Organizations 
			SET domain_verified=TRUE
			WHERE organization_id=$1
		`, organizationID.Int64)
	} else {
		_, err = db.Exec(`
			UPDATE Users 
			SET domain_verified=TRUE
			WHERE user_id=$1
		`, sess.UserID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update verification status",
		})
		return
	}

	fmt.Printf("✅ Domain verified: %s for user %s (ID: %d)\n", domain, req.Username, sess.UserID)

	// Trigger mail server sync
	go func() {
		if err := SyncMailServerConfig(); err != nil {
			fmt.Printf("⚠️ Failed to sync mail server after domain verification: %v\n", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"message":  fmt.Sprintf("Domain %s successfully verified!", domain),
		"domain":   domain,
		"verified": true,
	})
}

// GetCustomDomainInfo returns custom domain information
func GetCustomDomainInfo(c *gin.Context) {
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

	var domain sql.NullString
	var verified bool
	var accountType string
	var organizationID sql.NullInt64
	var verificationToken sql.NullString
	
	err := db.QueryRow(`
		SELECT COALESCE(account_type, 'personal'), organization_id 
		FROM Users WHERE user_id=$1
	`, sess.UserID).Scan(&accountType, &organizationID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get user info"})
		return
	}

	// Get domain info based on account type
	if accountType == "organization_owner" && organizationID.Valid {
		err = db.QueryRow(`
			SELECT custom_domain, COALESCE(domain_verified, FALSE), domain_verification_token
			FROM Organizations WHERE organization_id=$1
		`, organizationID.Int64).Scan(&domain, &verified, &verificationToken)
	} else {
		err = db.QueryRow(`
			SELECT custom_domain, COALESCE(domain_verified, FALSE)
			FROM Users WHERE user_id=$1
		`, sess.UserID).Scan(&domain, &verified)
	}
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get domain info"})
		return
	}

	response := gin.H{
		"success":  true,
		"domain":   nil,
		"verified": false,
	}

	if domain.Valid && domain.String != "" {
		response["domain"] = domain.String
		response["verified"] = verified
		
		if !verified && verificationToken.Valid {
			response["verification_token"] = verificationToken.String
			response["verification_record"] = fmt.Sprintf("office1789-verification=%s", verificationToken.String)
		}
	}

	c.JSON(http.StatusOK, response)
}

// RemoveCustomDomain removes the custom domain
func RemoveCustomDomain(c *gin.Context) {
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

	var accountType string
	var organizationID sql.NullInt64
	
	err := db.QueryRow(`
		SELECT COALESCE(account_type, 'personal'), organization_id 
		FROM Users WHERE user_id=$1
	`, sess.UserID).Scan(&accountType, &organizationID)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get user info"})
		return
	}

	// Remove domain based on account type
	if accountType == "organization_owner" && organizationID.Valid {
		_, err = db.Exec(`
			UPDATE Organizations 
			SET custom_domain=NULL, domain_verified=FALSE, domain_verification_token=NULL
			WHERE organization_id=$1
		`, organizationID.Int64)
	} else {
		_, err = db.Exec(`
			UPDATE Users 
			SET custom_domain=NULL, domain_verified=FALSE
			WHERE user_id=$1
		`, sess.UserID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to remove domain"})
		return
	}

	fmt.Printf("✅ Custom domain removed for user %s (ID: %d)\n", req.Username, sess.UserID)

	// Trigger mail server sync
	go func() {
		if err := SyncMailServerConfig(); err != nil {
			fmt.Printf("⚠️ Failed to sync mail server after domain removal: %v\n", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Custom domain removed successfully",
	})
}
