package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	virtualDomainsFile   = "./docker/config/postfix-virtual-domains.cf"
	virtualMailboxesFile = "./docker/config/postfix-virtual-mailboxes.cf"
	configDir            = "./docker/config"
)

// ensureConfigDir creates the config directory if it doesn't exist
func ensureConfigDir() error {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}
	}
	return nil
}

// syncCustomDomains synchronizes verified custom domains with Postfix configuration
func syncCustomDomains() error {
	// Ensure config directory exists
	if err := ensureConfigDir(); err != nil {
		return err
	}

	// Get all verified custom domains from database
	rows, err := db.Query(`
		SELECT DISTINCT custom_domain 
		FROM (
			SELECT custom_domain FROM Users 
			WHERE custom_domain IS NOT NULL AND domain_verified = TRUE
			UNION
			SELECT custom_domain FROM Organizations 
			WHERE custom_domain IS NOT NULL AND domain_verified = TRUE
		) AS all_domains
		WHERE custom_domain != ''
		ORDER BY custom_domain
	`)
	
	if err != nil {
		return fmt.Errorf("failed to query custom domains: %v", err)
	}
	defer rows.Close()

	// Build domains list
	domains := []string{"office1789.com"} // Default domain always present
	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err == nil && domain != "" {
			domains = append(domains, domain)
		}
	}

	// Write to virtual domains file
	domainsContent := "# Virtual domains configuration for custom domains\n"
	domainsContent += "# Format: domain (one per line)\n"
	domainsContent += "# This file is automatically updated by Office1789 when users add custom domains\n\n"
	domainsContent += strings.Join(domains, "\n") + "\n"

	if err := os.WriteFile(virtualDomainsFile, []byte(domainsContent), 0644); err != nil {
		return fmt.Errorf("failed to write virtual domains file: %v", err)
	}

	fmt.Printf("✅ Synced %d custom domains to Postfix\n", len(domains)-1)
	return nil
}

// syncVirtualMailboxes creates email aliases for custom domains
func syncVirtualMailboxes() error {
	// Get all users with verified custom domains (individual accounts)
	userRows, err := db.Query(`
		SELECT u.username, u.custom_domain, u.email
		FROM Users u
		WHERE u.custom_domain IS NOT NULL 
		  AND u.custom_domain != '' 
		  AND u.domain_verified = TRUE
		  AND u.account_type IN ('personal', 'organization_owner')
	`)
	
	if err != nil {
		return fmt.Errorf("failed to query users with custom domains: %v", err)
	}
	defer userRows.Close()

	mailboxMappings := []string{}
	
	// Individual accounts
	for userRows.Next() {
		var username, customDomain, officeEmail string
		if err := userRows.Scan(&username, &customDomain, &officeEmail); err == nil {
			// Map username@customdomain.com -> username@office1789.com
			customEmail := fmt.Sprintf("%s@%s", username, customDomain)
			mapping := fmt.Sprintf("%s %s", customEmail, officeEmail)
			mailboxMappings = append(mailboxMappings, mapping)
		}
	}

	// Get organization members with parent's custom domain
	orgRows, err := db.Query(`
		SELECT u.username, u.email, o.custom_domain
		FROM Users u
		JOIN Organizations o ON u.organization_id = o.organization_id
		WHERE o.custom_domain IS NOT NULL 
		  AND o.custom_domain != ''
		  AND o.domain_verified = TRUE
		  AND u.account_type = 'organization_member'
	`)
	
	if err == nil {
		defer orgRows.Close()
		for orgRows.Next() {
			var username, officeEmail, customDomain string
			if err := orgRows.Scan(&username, &officeEmail, &customDomain); err == nil {
				// Map username@org-customdomain.com -> username@office1789.com
				customEmail := fmt.Sprintf("%s@%s", username, customDomain)
				mapping := fmt.Sprintf("%s %s", customEmail, officeEmail)
				mailboxMappings = append(mailboxMappings, mapping)
			}
		}
	}

	// Write to virtual mailboxes file
	mailboxesContent := "# Virtual mailbox mappings for custom domains\n"
	mailboxesContent += "# Format: email@custom-domain.com username@office1789.com\n"
	mailboxesContent += "# This file maps custom domain emails to actual Office1789 accounts\n"
	mailboxesContent += "# Automatically updated by Office1789 API\n\n"
	
	if len(mailboxMappings) > 0 {
		mailboxesContent += strings.Join(mailboxMappings, "\n") + "\n"
	}

	if err := os.WriteFile(virtualMailboxesFile, []byte(mailboxesContent), 0644); err != nil {
		return fmt.Errorf("failed to write virtual mailboxes file: %v", err)
	}

	fmt.Printf("✅ Synced %d virtual mailbox mappings\n", len(mailboxMappings))
	return nil
}

// reloadPostfixConfig reloads Postfix configuration in mailserver container
func reloadPostfixConfig() error {
	// Note: This would execute: docker exec mailserver postfix reload
	// For now, we'll just log it - implement actual docker exec if needed
	fmt.Println("📮 Postfix configuration files updated. Run: docker exec mailserver postfix reload")
	return nil
}

// SyncMailServerConfig synchronizes all custom domain configurations with mail server
func SyncMailServerConfig() error {
	fmt.Println("🔄 Synchronizing custom domains with mail server...")
	
	if err := syncCustomDomains(); err != nil {
		return fmt.Errorf("failed to sync domains: %v", err)
	}
	
	if err := syncVirtualMailboxes(); err != nil {
		return fmt.Errorf("failed to sync mailboxes: %v", err)
	}
	
	if err := reloadPostfixConfig(); err != nil {
		return fmt.Errorf("failed to reload Postfix: %v", err)
	}
	
	return nil
}

// StartMailSyncScheduler starts a background goroutine that periodically syncs mail config
func StartMailSyncScheduler() {
	go func() {
		// Initial sync
		time.Sleep(5 * time.Second)
		if err := SyncMailServerConfig(); err != nil {
			fmt.Printf("❌ Initial mail sync failed: %v\n", err)
		}
		
		// Periodic sync every 5 minutes
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			if err := SyncMailServerConfig(); err != nil {
				fmt.Printf("❌ Periodic mail sync failed: %v\n", err)
			}
		}
	}()
	
	fmt.Println("✅ Mail sync scheduler started (every 5 minutes)")
}
