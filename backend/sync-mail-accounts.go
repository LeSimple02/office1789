package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/lib/pq"
)

// Script pour créer les comptes mail pour tous les utilisateurs existants
func main() {
	// Connexion à la base de données
	connStr := "host=localhost port=5432 user=robespierre password=guillotine dbname=office1789 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erreur de connexion à la base de données:", err)
	}
	defer db.Close()

	// Tester la connexion
	err = db.Ping()
	if err != nil {
		log.Fatal("Impossible de se connecter à la base de données:", err)
	}

	fmt.Println("✅ Connecté à la base de données PostgreSQL")
	fmt.Println("")

	// Récupérer tous les utilisateurs
	rows, err := db.Query("SELECT username FROM Users")
	if err != nil {
		log.Fatal("Erreur lors de la récupération des utilisateurs:", err)
	}
	defer rows.Close()

	var users []string
	for rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			log.Println("Erreur de scan:", err)
			continue
		}
		users = append(users, username)
	}

	if len(users) == 0 {
		fmt.Println("❌ Aucun utilisateur trouvé dans la base de données")
		os.Exit(0)
	}

	fmt.Printf("📋 %d utilisateur(s) trouvé(s) dans la base de données\n", len(users))
	fmt.Println("")

	// Demander le mot de passe par défaut
	fmt.Println("⚠️  Les comptes mail vont être créés pour tous les utilisateurs")
	fmt.Print("Entrez le mot de passe par défaut à utiliser (ou appuyez sur Entrée pour utiliser 'office1789'): ")
	
	var defaultPassword string
	fmt.Scanln(&defaultPassword)
	
	if defaultPassword == "" {
		defaultPassword = "office1789"
	}

	fmt.Println("")
	fmt.Println("🔧 Création des comptes mail...")
	fmt.Println("")

	successCount := 0
	errorCount := 0

	for _, username := range users {
		email := fmt.Sprintf("%s@office1789.local", username)
		fmt.Printf("Création de %s... ", email)

		err := createMailAccountScript(username, defaultPassword)
		if err != nil {
			fmt.Printf("❌ ERREUR: %v\n", err)
			errorCount++
		} else {
			fmt.Println("✅")
			successCount++
		}
	}

	fmt.Println("")
	fmt.Println("========================================")
	fmt.Printf("✅ Comptes créés avec succès: %d\n", successCount)
	if errorCount > 0 {
		fmt.Printf("❌ Erreurs: %d\n", errorCount)
	}
	fmt.Println("========================================")
	fmt.Println("")
	fmt.Println("💡 Les utilisateurs peuvent maintenant se connecter à Roundcube avec:")
	fmt.Println("   Email: username@office1789.local")
	fmt.Printf("   Mot de passe: %s\n", defaultPassword)
	fmt.Println("   URL: http://localhost:8081")
}

// Fonction helper pour créer un compte mail (copie de mailmanager.go)
func createMailAccountScript(username, password string) error {
	email := fmt.Sprintf("%s@office1789.local", username)
	
	// Import nécessaire
	cmd := exec.Command("docker", "exec", "mailserver", "setup", "email", "add", email, password)
	
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	
	err := cmd.Run()
	if err != nil {
		// Si l'erreur dit que le compte existe déjà, ce n'est pas grave
		if strings.Contains(stderr.String(), "already exists") || strings.Contains(stderr.String(), "existe déjà") {
			return nil
		}
		return fmt.Errorf("%v - %s", err, stderr.String())
	}
	
	return nil
}
