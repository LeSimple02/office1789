package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CalendarEvent struct {
	EventID          int       `json:"id"`
	UserID           int       `json:"user_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Start            time.Time `json:"start"`
	End              time.Time `json:"end"`
	Location         string    `json:"location"`
	DateCreated      time.Time `json:"date_created"`
}

type EventRequest struct {
	Username    string `json:"username"`
	Token       string `json:"token"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Start       string `json:"start"`
	End         string `json:"end"`
	Location    string `json:"location"`
}

type EventDeleteRequest struct {
	Username string `json:"username"`
	Token    string `json:"token"`
	EventID  string `json:"event_id"`
}

// GetCalendarEvents récupère tous les événements d'un utilisateur
func GetCalendarEvents(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Token    string `json:"token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Vérifier la session
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Récupérer l'ID utilisateur
	var userID int
	err := db.QueryRow("SELECT user_id FROM Users WHERE username = $1", req.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	// Récupérer les événements
	rows, err := db.Query(`
		SELECT event_id, user_id, event_title, event_description, event_start, event_end, location, date_created
		FROM CalendarEvents
		WHERE user_id = $1
		ORDER BY event_start ASC
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}
	defer rows.Close()

	var events []CalendarEvent
	for rows.Next() {
		var event CalendarEvent
		var description, location sql.NullString
		err := rows.Scan(
			&event.EventID,
			&event.UserID,
			&event.Title,
			&description,
			&event.Start,
			&event.End,
			&location,
			&event.DateCreated,
		)
		if err != nil {
			continue
		}
		if description.Valid {
			event.Description = description.String
		}
		if location.Valid {
			event.Location = location.String
		}
		events = append(events, event)
	}

	if events == nil {
		events = []CalendarEvent{}
	}

	c.JSON(http.StatusOK, gin.H{"events": events})
}

// CreateCalendarEvent crée un nouvel événement
func CreateCalendarEvent(c *gin.Context) {
	var req EventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Vérifier la session
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Récupérer l'ID utilisateur
	var userID int
	err := db.QueryRow("SELECT user_id FROM Users WHERE username = $1", req.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	// Parser les dates
	startTime, err := time.Parse(time.RFC3339, req.Start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time"})
		return
	}

	endTime, err := time.Parse(time.RFC3339, req.End)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time"})
		return
	}

	// Insérer l'événement
	var eventID int
	err = db.QueryRow(`
		INSERT INTO CalendarEvents (user_id, event_title, event_description, event_start, event_end, location)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING event_id
	`, userID, req.Title, req.Description, startTime, endTime, req.Location).Scan(&eventID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Event created successfully",
		"event_id": eventID,
	})
}

// DeleteCalendarEvent supprime un événement
func DeleteCalendarEvent(c *gin.Context) {
	var req EventDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Vérifier la session
	if session, ok := sessions[req.Token]; !ok || session.Username != req.Username || req.Username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Récupérer l'ID utilisateur
	var userID int
	err := db.QueryRow("SELECT user_id FROM Users WHERE username = $1", req.Username).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user"})
		return
	}

	// Vérifier que l'événement appartient à l'utilisateur avant de le supprimer
	var ownerID int
	err = db.QueryRow("SELECT user_id FROM CalendarEvents WHERE event_id = $1", req.EventID).Scan(&ownerID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this event"})
		return
	}

	// Supprimer l'événement
	_, err = db.Exec("DELETE FROM CalendarEvents WHERE event_id = $1", req.EventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}
