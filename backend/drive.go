package main

/*
import (
	
	"github.com/gin-gonic/gin"
	"strings"
	"net/http"
	"github.com/google/uuid"
	"time"
)




func Sub(c *gin.Context){
	
	var count int;
	var subi Subscribe
	var infova vInfo
	
	c.BindJSON(&subi)
	
	
	
	if strings.TrimSpace(subi.Username) != ""{
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE username=$1", subi.Username)
		rows.Scan(&count)
		
		if count>0{
			infova.Username = "no"
		}else{
			infova.Username = "yes"
		}
		
	}
	if strings.TrimSpace(subi.Email) != ""{
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE email=$1", subi.Email)
		rows.Scan(&count)
		if count>0{
			infova.Email = "no"
		}
		
	}
	if strings.TrimSpace(subi.PhoneNumber) !="" {
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE phonenumber=$1", subi.PhoneNumber)
		rows.Scan(&count)
		if count>0 {
			infova.Phone = "no"
		}
	}
		
	if (infova.Phone=="no" || infova.Username=="no" || infova.Email=="no"){
		c.JSON(http.StatusOK, infova)
	} else if(infova.Username=="yes"){
		rows := db.QueryRow("SELECT count(*) INTO Users WHERE username=$1 OR email=$2 OR phonenumber=3$)", subi.Username, subi.Email, subi.PhoneNumber)
		rows.Scan(&count)
		if count==0{
		
			if strings.TrimSpace(subi.Username) != ""{
		
				subi.Password = HashPassword(subi.Password)
		
				db.Exec("INSERT INTO Users (username, password_hash, email, phonenumber, nboffer, date_joined, last_login, domain) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", subi.Username, subi.Password, subi.Email, subi.PhoneNumber, subi.Nboffer, subi.DateJoined, subi.LastLogin, "@office1789")
				
				sessionToken := uuid.NewString()
				expiresAtTime := time.Now().Add(120 * time.Second)
				sessions[sessionToken] = session {
				Username : subi.Username,
				expiry : expiresAtTime,
			}
			
			c.JSON(http.StatusOK,  sessionSend{sessionToken, subi.Username, expiresAtTime})
			}
		
		}
		
	}
}
*/

