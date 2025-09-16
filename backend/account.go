package main

import (
	"net/http"
	"strings"
	"time"
	//"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChangeIn struct {
	Username     string `json:username`
	LastUsername string `json:lastusername`
	Password     string `json:password`
	Email        string `json:email`
	Domain       string `json:email`
	PhoneNumber  string `json:phonenumber`
	Nboffer      int    `json:nboffer`
	Token        string `json:token`
}

func drivei(c *gin.Context) {

	var verif sessionSend
	var infop Subscribe

	c.BindJSON(&verif)

	if sessions[verif.Token].Username == verif.Username {
		rows := db.QueryRow("SELECT domain, nboffer, date_joined, last_login, phonenumber, email FROM Users WHERE username=$1", verif.Username)
		rows.Scan(&infop.Domain, &infop.Nboffer, &infop.DateJoined, &infop.LastLogin, &infop.PhoneNumber, &infop.Email)

		c.JSON(http.StatusOK, infop)

	} else {
		infop.Username = "no"
		c.JSON(http.StatusOK, infop)
	}
	

}

func getinfop(c *gin.Context) {

	var verif sessionSend
	var infop Subscribe

	c.ShouldBindJSON(&verif)
	

	if sessions[verif.Token].Username == verif.Username && verif.Username != ""{
		rows := db.QueryRow("SELECT domain, nboffer, date_joined, last_login, phonenumber, email FROM Users WHERE username=$1", verif.Username)
		rows.Scan(&infop.Domain, &infop.Nboffer, &infop.DateJoined, &infop.LastLogin, &infop.PhoneNumber, &infop.Email)
		
		c.JSON(http.StatusOK, infop)
		return

	} else {
		infop.Username = "no"
		
		c.JSON(http.StatusOK, infop)
		
	}

}

func ChangeI(c *gin.Context) {
	var count int
	var cha ChangeIn
	var infova vInfo

	c.BindJSON(&cha)

	if sessions[cha.Token].Username == cha.LastUsername {

		if strings.TrimSpace(cha.Username) != "" {
			rows := db.QueryRow("SELECT count(*) FROM Users WHERE username=$1", cha.Username)
			rows.Scan(&count)

			if count > 0 {
				infova.Username = "no"
			}

		}
		if strings.TrimSpace(cha.Email) != "" {
			rows := db.QueryRow("SELECT count(*) FROM Users WHERE email=$1", cha.Email)
			rows.Scan(&count)
			if count > 0 {
				infova.Email = "no"
			}

		}
		if strings.TrimSpace(cha.PhoneNumber) != "" {
			rows := db.QueryRow("SELECT count(*) FROM Users WHERE phonenumber=$1", cha.PhoneNumber)
			rows.Scan(&count)
			if count > 0 {
				infova.Phone = "no"
			}
		}

		if infova.Phone == "no" || infova.Username == "no" || infova.Email == "no" {
			c.JSON(http.StatusOK, infova)
		} else {

			if strings.TrimSpace(cha.LastUsername) != "" {

				if cha.Password != "" {
					cha.Password = HashPassword(cha.Password)
					db.Exec("UPDATE Users SET password_hash=$1 WHERE username=$2", cha.Password, cha.LastUsername)
				}

				if cha.Email != "" {
					db.Exec("UPDATE Users SET email=$1 WHERE username=$2", cha.Email, cha.LastUsername)
				}
				if cha.PhoneNumber != "" {
					db.Exec("UPDATE Users SET phonenumber=$1 WHERE username=$2", cha.PhoneNumber, cha.LastUsername)
				}
				if cha.Username != "" {
					db.Exec("UPDATE Users SET username=$1 WHERE username=$2", cha.Username, cha.LastUsername)
				}

				sessionToken := uuid.NewString()
				expiresAtTime := time.Now().Add(120 * time.Second)
				if cha.Username != "" {
					sessions[sessionToken] = session{
						Username: cha.Username,
						expiry:   expiresAtTime,
					}
					c.JSON(http.StatusOK, sessionSend{sessionToken, cha.Username, expiresAtTime})
				} else {
					c.JSON(http.StatusOK, infova)
				}

			}

		}
	}
}
