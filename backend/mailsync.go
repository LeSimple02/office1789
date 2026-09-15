package main

import (
	"fmt"
	"os"
	"os/exec"
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
	// Reload Postfix
	cmd := exec.Command("docker", "exec", "mailserver", "postfix", "reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reload Postfix: %v", err)
	}
	fmt.Println("✅ Postfix reloaded successfully")
	return nil
}

// syncMatrixDomains adds custom domains to Matrix homeserver configuration
func syncMatrixDomains() error {
	// Query all verified custom domains
	rows, err := db.Query(`
		SELECT DISTINCT custom_domain 
		FROM (
			SELECT custom_domain FROM users WHERE domain_verified = TRUE AND custom_domain IS NOT NULL
			UNION
			SELECT custom_domain FROM organizations WHERE domain_verified = TRUE AND custom_domain IS NOT NULL
		) AS domains
	`)
	if err != nil {
		return fmt.Errorf("failed to query custom domains: %v", err)
	}
	defer rows.Close()

	domains := []string{"office1789.com"} // Always include default domain
	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err == nil {
			domains = append(domains, domain)
		}
	}

	fmt.Printf("🔄 Matrix: Found %d domains to configure\n", len(domains))
	
	// Note: Matrix domains are typically configured at startup via homeserver.yaml
	// For dynamic updates, you would need to:
	// 1. Update /data/homeserver.yaml with new server_name or virtual_hosts
	// 2. Restart synapse container
	// For now, we'll just log them
	fmt.Println("📝 Matrix domains to configure:", strings.Join(domains, ", "))
	fmt.Println("⚠️  Note: Matrix requires homeserver.yaml update and container restart")
	
	return nil
}

// syncRoundcubeDomains updates Roundcube configuration for custom domains
func syncRoundcubeDomains() error {
	// Query all verified custom domains
	rows, err := db.Query(`
		SELECT DISTINCT custom_domain 
		FROM (
			SELECT custom_domain FROM users WHERE domain_verified = TRUE AND custom_domain IS NOT NULL
			UNION
			SELECT custom_domain FROM organizations WHERE domain_verified = TRUE AND custom_domain IS NOT NULL
		) AS domains
	`)
	if err != nil {
		return fmt.Errorf("failed to query custom domains: %v", err)
	}
	defer rows.Close()

	domains := []string{"office1789.com"}
	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err == nil {
			domains = append(domains, domain)
		}
	}

	fmt.Printf("🔄 Roundcube: Found %d domains to configure\n", len(domains))
	
	// Roundcube identities are user-specific, so we just need to ensure
	// the mail server accepts these domains (already done via Postfix)
	fmt.Println("✅ Roundcube will automatically use custom domains via IMAP")
	
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

	// Create Postfix hash database
	cmd := exec.Command("docker", "exec", "mailserver", "postmap", virtualMailboxesFile)
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Warning: Could not run postmap: %v\n", err)
	}

	// Reload Postfix
	if err := reloadPostfixConfig(); err != nil {
		fmt.Printf("⚠️  Warning: Could not reload Postfix: %v\n", err)
	}

	// Sync Matrix domains (informational for now)
	if err := syncMatrixDomains(); err != nil {
		fmt.Printf("⚠️  Warning: Matrix sync failed: %v\n", err)
	}

	// Sync Roundcube (informational)
	if err := syncRoundcubeDomains(); err != nil {
		fmt.Printf("⚠️  Warning: Roundcube sync failed: %v\n", err)
	}

	fmt.Println("✅ Mail server synchronization completed")
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
