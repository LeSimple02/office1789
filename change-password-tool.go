package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	_ "github.com/lib/pq"
)

// HashPassword génère un hash bcrypt
func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// Changer le mot de passe Mail
func changeMailPassword(username, newPassword string) error {
	email := username + "@office1789.local"
	
	cmdDel := exec.Command("docker", "exec", "mailserver", "setup", "email", "del", email)
	_ = cmdDel.Run()
	
	cmdAdd := exec.Command("docker", "exec", "mailserver", "setup", "email", "add", email, newPassword)
	var stderr bytes.Buffer
	cmdAdd.Stderr = &stderr
	
	err := cmdAdd.Run()
	if err != nil {
		return fmt.Errorf("failed to change mail password: %v - %s", err, stderr.String())
	}
	
	return nil
}

// Changer le mot de passe Matrix
func changeMatrixPassword(username, newPassword string) error {
	cmd := exec.Command("docker", "exec", "synapse",
		"reset-password",
		"-c", "/data/homeserver.yaml",
		username,
		newPassword)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to change matrix password: %v - %s", err, stderr.String())
	}
	
	return nil
}

func main() {
	fmt.Println("====================================")
	fmt.Println("Office1789 - Password Change Tool")
	fmt.Println("====================================")
	fmt.Println()
	
	// Lire le username
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	
	if username == "" {
		fmt.Println("Error: Username cannot be empty")
		os.Exit(1)
	}
	
	// Lire le nouveau mot de passe (sans affichage)
	fmt.Print("New password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nError reading password:", err)
		os.Exit(1)
	}
	newPassword := string(passwordBytes)
	fmt.Println() // Nouvelle ligne après le mot de passe masqué
	
	if newPassword == "" {
		fmt.Println("Error: Password cannot be empty")
		os.Exit(1)
	}
	
	// Confirmer le mot de passe
	fmt.Print("Confirm password: ")
	confirmBytes, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Println("\nError reading password:", err)
		os.Exit(1)
	}
	confirmPassword := string(confirmBytes)
	fmt.Println()
	
	if newPassword != confirmPassword {
		fmt.Println("Error: Passwords do not match")
		os.Exit(1)
	}
	
	fmt.Println()
	fmt.Println("Connecting to database...")
	
	// Connexion PostgreSQL
	connStr := "host=localhost port=5432 user=admin password=adminpassword dbname=office1789db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Database connection error:", err)
		os.Exit(1)
	}
	defer db.Close()
	
	// Vérifier que l'utilisateur existe
	var userID int
	err = db.QueryRow("SELECT user_id FROM Users WHERE username=$1", username).Scan(&userID)
	if err != nil {
		fmt.Println("Error: User not found")
		os.Exit(1)
	}
	
	fmt.Println("Changing passwords...")
	
	// 1. Office1789 (PostgreSQL)
	newHash := HashPassword(newPassword)
	_, err = db.Exec("UPDATE Users SET password_hash=$1 WHERE user_id=$2", newHash, userID)
	if err != nil {
		fmt.Println("Error updating Office1789 password:", err)
		os.Exit(1)
	}
	fmt.Println("✓ Office1789 password changed")
	
	// 2. Mail
	err = changeMailPassword(username, newPassword)
	if err != nil {
		fmt.Println("⚠ Mail password change failed:", err)
	} else {
		fmt.Println("✓ Mail password changed")
	}
	
	// 3. Matrix
	err = changeMatrixPassword(username, newPassword)
	if err != nil {
		fmt.Println("⚠ Matrix password change failed:", err)
	} else {
		fmt.Println("✓ Matrix password changed")
	}
	
	fmt.Println()
	fmt.Println("Password change completed!")
}
