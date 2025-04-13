package main

import (
	
	"github.com/gin-gonic/gin"
	"strings"
	"net/http"
	"github.com/google/uuid"
	"time"
)


func Connect(c *gin.Context){
	var conn Connecti
	var hash string
	
	c.BindJSON(&conn)
	
	if strings.TrimSpace(conn.Username) != ""{
		rows := db.QueryRow("SELECT password_hash FROM Users WHERE username=$1", conn.Username)
		rows.Scan(&hash)
		passc := CheckPasswordHash(conn.Password, hash)
		
		
		if passc == true{
			sessionToken := uuid.NewString()
			expiresAtTime := time.Now().Add(120 * time.Second)
			sessions[sessionToken] = session {
				Username : conn.Username,
				expiry : expiresAtTime,
			}
			
			c.JSON(http.StatusOK,  sessionSend{sessionToken, conn.Username, expiresAtTime})
		} else{
			
			conn.Username = "no"
			c.JSON(http.StatusOK, conn)
		}
	}
	
	
	

}
