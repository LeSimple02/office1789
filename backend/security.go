package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/bcrypt"
	"sync"
	"time"
)

// Map en mémoire pour compatibilité (sera progressivement remplacé par DB)
var sessions = map[string]session{}
var sessionsMutex sync.RWMutex // Protection contre les accès concurrents

// Session structure basée sur user_id
type session struct {
	UserID   int
	Username string
	Password string // Mot de passe en clair (UNIQUEMENT en RAM, pour SSO)
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
	fmt.Println("[SESSION][DEBUG] createSessionInDB called")

	if db == nil {
		err := fmt.Errorf("db is not initialized")
		fmt.Printf("[SESSION][ERROR] %v\n", err)
		return err
	}

	// Récupérer le password depuis la RAM de façon sûre
	sessionsMutex.RLock()
	s, ok := sessions[token]
	sessionsMutex.RUnlock()

	password := ""
	if ok {
		password = s.Password
	}

	// Exécuter l'INSERT avec contexte et retry pour erreurs transitoires
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var lastErr error
	for attempt := 1; attempt <= 2; attempt++ {
		res, err := db.ExecContext(ctx, `
			INSERT INTO sessions (user_id, username, session_token, expiry, password_plain)
			VALUES ($1, $2, $3, $4, $5)
		`, userID, username, token, expiresAt, password)

		if err != nil {
			lastErr = err
			fmt.Printf("[SESSION][ERROR] Failed to insert session in DB (attempt %d): %v\n", attempt, err)
			fmt.Printf("[SESSION][DEBUG] userID=%v username=%v token=%v expiry=%v password_present=%v\n", userID, username, token, expiresAt, password != "")
			// small backoff before retry
			time.Sleep(150 * time.Millisecond)
			continue
		}

		// Log rows affected when possible
		if res != nil {
			if ra, errRows := res.RowsAffected(); errRows == nil {
				fmt.Printf("[SESSION][INFO] Session inserted in DB (attempt %d), rows affected=%d: userID=%v username=%v token=%v expiry=%v\n", attempt, ra, userID, username, token, expiresAt)
			} else {
				fmt.Printf("[SESSION][INFO] Session inserted in DB (attempt %d): userID=%v username=%v token=%v expiry=%v\n", attempt, userID, username, token, expiresAt)
			}
		}
		return nil
	}

	// Retourner la dernière erreur rencontrée
	return lastErr
}

// Récupérer une session depuis la base de données
func getSessionFromDB(token string) (*session, error) {
	var sess session
	var expiresAt time.Time
	var passwordPlain string
	
	err := db.QueryRow(`
		SELECT s.user_id, u.username, s.expiry, COALESCE(s.password_plain, '') as password
		FROM sessions s
		JOIN users u ON s.user_id = u.user_id
		WHERE s.session_token = $1 AND s.expiry > NOW()
	`, token).Scan(&sess.UserID, &sess.Username, &expiresAt, &passwordPlain)
	
	if err != nil {
		return nil, err
	}
	
	sess.expiry = expiresAt
	sess.Password = passwordPlain // Récupérer le password pour SSO
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
	
	// Vérifier d'abord en mémoire (ancien système) avec mutex
	sessionsMutex.RLock()
	sess, ok := sessions[token]
	sessionsMutex.RUnlock()
	
	if ok && sess.Username == username {
		return &sess, true
	}
	
	// Sinon vérifier en DB (nouveau système)
	sessFromDB, err := getSessionFromDB(token)
	if err != nil || sessFromDB == nil {
		return nil, false
	}
	
	// Vérifier que le username correspond
	if sessFromDB.Username != username {
		return nil, false
	}
	
	// Mettre en cache en mémoire pour performance avec mutex
	sessionsMutex.Lock()
	sessions[token] = *sessFromDB
	sessionsMutex.Unlock()
	
	return sessFromDB, true
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

// ============ CHIFFREMENT AES-256 POUR MAIL PASSWORD ============

// Clé secrète AES (32 bytes = 256 bits) - À CHANGER EN PRODUCTION
const AES_SECRET_KEY = "Office1789-AES-32BytesSecretKey!"

// EncryptPassword chiffre un mot de passe avec AES-256-GCM
func EncryptPassword(plaintext string) (string, error) {
	block, err := aes.NewCipher([]byte(AES_SECRET_KEY))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword déchiffre un mot de passe chiffré avec AES-256-GCM
func DecryptPassword(encrypted string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(AES_SECRET_KEY))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
