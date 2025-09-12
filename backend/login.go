package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Connect(c *gin.Context) {
	var conn Connecti
	var hash string

	c.BindJSON(&conn)

	if strings.TrimSpace(conn.Username) != "" {
		rows := db.QueryRow("SELECT password_hash FROM Users WHERE username=$1", conn.Username)
		rows.Scan(&hash)
		passc := CheckPasswordHash(conn.Password, hash)

		if passc == true {
			sessionToken := uuid.NewString()
			expiresAtTime := time.Now().Add(120 * time.Second)
			sessions[sessionToken] = session{
				Username: conn.Username,
				expiry:   expiresAtTime,
			}

			c.JSON(http.StatusOK, sessionSend{conn.Username, sessionToken, expiresAtTime})
		} else {

			conn.Username = "no"
			c.JSON(http.StatusOK, conn)
		}
	}

}
