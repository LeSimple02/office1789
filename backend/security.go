package main

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Map en mémoire pour compatibilité (sera progressivement remplacé par DB)
var sessions = map[string]session{}

// Session structure basée sur user_id
type session struct {
	UserID   int
	Username string
	expiry   time.Time
}

type sessionSend struct {
	UserID   int       `json:"user_id"`
	Username string    `json:"Username"`
	Token    string    `json:"Token"`
	Expiry   time.Time `json:"Expiry"`
}

// Structure pour recevoir les requêtes de validation du frontend
type sessionVerify struct {
	Username string `json:"Username"`
	Token    string `json:"Token"`
}

// Créer une session en base de données ET en mémoire
func createSessionInDB(userID int, username string, token string, expiresAt time.Time) error {
	_, err := db.Exec(`
		INSERT INTO sessions (user_id, username, session_token, expiry)
		VALUES ($1, $2, $3, $4)
	`, userID, username, token, expiresAt)
	return err
}

// Récupérer une session depuis la base de données
func getSessionFromDB(token string) (*session, error) {
	var sess session
	var expiresAt time.Time
	
	err := db.QueryRow(`
		SELECT s.user_id, u.username, s.expiry
		FROM sessions s
		JOIN users u ON s.user_id = u.user_id
		WHERE s.session_token = $1 AND s.expiry > NOW()
	`, token).Scan(&sess.UserID, &sess.Username, &expiresAt)
	
	if err != nil {
		return nil, err
	}
	
	sess.expiry = expiresAt
	return &sess, nil
}

// Supprimer une session de la base de données
func deleteSessionFromDB(token string) error {
	_, err := db.Exec(`DELETE FROM sessions WHERE session_token = $1`, token)
	return err
}

// Nettoyer les sessions expirées
func cleanExpiredSessions() error {
	_, err := db.Exec(`DELETE FROM sessions WHERE expiry < NOW()`)
	return err
}

// Valider une session et retourner les informations utilisateur
func validateSession(token string, username string) (*session, bool) {
	if token == "" || username == "" {
		return nil, false
	}
	
	// Vérifier d'abord en mémoire (ancien système)
	if sess, ok := sessions[token]; ok && sess.Username == username {
		return &sess, true
	}
	
	// Sinon vérifier en DB (nouveau système)
	sess, err := getSessionFromDB(token)
	if err != nil || sess == nil {
		return nil, false
	}
	
	// Vérifier que le username correspond
	if sess.Username != username {
		return nil, false
	}
	
	// Mettre en cache en mémoire pour performance
	sessions[token] = *sess
	
	return sess, true
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
	}
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

