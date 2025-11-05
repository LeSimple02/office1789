package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Créer un compte mail pour un utilisateur
func createMailAccount(username, password string) error {
	// Format de l'email : username@office1789.local
	email := fmt.Sprintf("%s@office1789.local", username)
	
	// Vérifier si le compte existe déjà
	checkCmd := exec.Command("docker", "exec", "mailserver", "setup", "email", "list")
	output, err := checkCmd.Output()
	if err != nil {
		return fmt.Errorf("erreur lors de la vérification des comptes mail: %v", err)
	}
	
	// Si le compte existe déjà, ne rien faire
	if strings.Contains(string(output), email) {
		fmt.Printf("⚠️  Compte mail %s existe déjà\n", email)
		return nil
	}
	
	// Créer le compte mail avec docker exec
	cmd := exec.Command("docker", "exec", "mailserver", "setup", "email", "add", email, password)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("erreur lors de la création du compte mail %s: %v - %s", email, err, stderr.String())
	}
	
	fmt.Printf("✅ Compte mail créé: %s\n", email)
	return nil
}

// Supprimer un compte mail
func deleteMailAccount(username string) error {
	email := fmt.Sprintf("%s@office1789.local", username)
	
	cmd := exec.Command("docker", "exec", "mailserver", "setup", "email", "del", email)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("erreur lors de la suppression du compte mail %s: %v - %s", email, err, stderr.String())
	}
	
	fmt.Printf("✅ Compte mail supprimé: %s\n", email)
	return nil
}

// Mettre à jour le mot de passe d'un compte mail
func updateMailPassword(username, newPassword string) error {
	email := fmt.Sprintf("%s@office1789.local", username)
	
	cmd := exec.Command("docker", "exec", "mailserver", "setup", "email", "update", email, newPassword)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("erreur lors de la mise à jour du mot de passe mail %s: %v - %s", email, err, stderr.String())
	}
	
	fmt.Printf("✅ Mot de passe mail mis à jour: %s\n", email)
	return nil
}

// Lister tous les comptes mail
func listMailAccounts() ([]string, error) {
	cmd := exec.Command("docker", "exec", "mailserver", "setup", "email", "list")
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la liste des comptes mail: %v", err)
	}
	
	// Parser la sortie pour extraire les emails
	lines := strings.Split(string(output), "\n")
	var emails []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "@office1789.local") {
			// Extraire l'email de la ligne
			parts := strings.Fields(line)
			if len(parts) > 0 {
				emails = append(emails, parts[0])
			}
		}
	}
	
	return emails, nil
}
