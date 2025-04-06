package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"database/sql"
	"time"
	"fmt"
	"strings"
	 "golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"	
	"github.com/gin-contrib/cors"
	"github.com/google/uuid"

	
)

func welcome(c *gin.Context){
	c.Writer.Write([]byte("Hello welcome to the backend of Office1789\n"))
}



const (
	host = "localhost"
	port = 5432
	user = "robespierre"
	password = "guillotine"
	dbname = "office1789"
)

type Subscribe struct{
	Username string `json:username`
	Password string `json:password`
	Email string `json:email`
	Domain string `json:email`
	PhoneNumber string `json:phonenumber`
	Nboffer int `json:nboffer`
	DateJoined string `json:datejoined`
	LastLogin string `json:lastlogin`
}

type ChangeIn struct{
	Username string `json:username`
	LastUsername string `json:lastusername`
	Password string `json:password`
	Email string `json:email`
	Domain string `json:email`
	PhoneNumber string `json:phonenumber`
	Nboffer int `json:nboffer`
	Token string	`json:token`
	
}

type Connecti struct{
	Username string `json:username`
	Password string `json:password`
	LastLogin string `json:lastlogin`
}


func HashPassword(password string) (string) {
    bytes, err:= bcrypt.GenerateFromPassword([]byte(password), 14)
    if(err != nil){
    }
    return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}


var sessions = map[string]session{}

type session struct {
	Username string
	expiry time .Time
}

type sessionSend struct {
	Token string	`json:token`
	Username string `json:username`
	Expiry time .Time
}


type createGroup struct {
	Token string	`json:token`
	Username string `json:username`
	Patcipant []string `json:recepter`
	
	
}

type vInfo struct {
	Email string	`json:email`
	Username string `json:username`
	Phone string `json:phone`
}

func changeI(c *gin.Context){
	var count int
	var cha ChangeIn
	var infova vInfo
	
	c.BindJSON(&cha)
	
	fmt.Println(cha)
	if(sessions[cha.Token].Username == cha.LastUsername){
	
		if strings.TrimSpace(cha.Username) != ""{
			rows := db.QueryRow("SELECT count(*) FROM Users WHERE username=$1", cha.Username)
			fmt.Println(cha.Username)
			rows.Scan(&count)
			fmt.Println(count)
			
			if count>0{
				infova.Username = "no"
			}
			
		}
		if strings.TrimSpace(cha.Email) != ""{
			rows := db.QueryRow("SELECT count(*) FROM Users WHERE email=$1", cha.Email)
			rows.Scan(&count)
			if count>0{
				infova.Email = "no"
			}
			
		}
		if strings.TrimSpace(cha.PhoneNumber) !="" {
			rows := db.QueryRow("SELECT count(*) FROM Users WHERE phonenumber=$1", cha.PhoneNumber)
			rows.Scan(&count)
			if count>0 {
				infova.Phone = "no"
			}
		}
			
		if (infova.Phone=="no" || infova.Username=="no" || infova.Email=="no"){
			c.JSON(http.StatusOK, infova)
		} else{
			
				
				
				if strings.TrimSpace(cha.LastUsername) != ""{
			
					if(cha.Password != ""){
						cha.Password = HashPassword(cha.Password)
						db.Exec("UPDATE Users SET password_hash=$1 WHERE username=$2", cha.Password, cha.LastUsername)
					}
			
					if(cha.Email != ""){
						db.Exec("UPDATE Users SET email=$1 WHERE username=$2", cha.Email, cha.LastUsername)
					}
					if(cha.PhoneNumber != ""){
						db.Exec("UPDATE Users SET phonenumber=$1 WHERE username=$2", cha.PhoneNumber, cha.LastUsername)
					}
					if(cha.Username != ""){
						db.Exec("UPDATE Users SET username=$1 WHERE username=$2", cha.Username, cha.LastUsername)
					}
					
					
					sessionToken := uuid.NewString()
					expiresAtTime := time.Now().Add(120 * time.Second)
					if(cha.Username != ""){
						sessions[sessionToken] = session {
						Username : cha.Username,
						expiry : expiresAtTime,
						}
						c.JSON(http.StatusOK,  sessionSend{sessionToken, cha.Username, expiresAtTime})
					}else{
						c.JSON(http.StatusOK, infova)
					}
					
				}
				
				
			}	
		}
}


func connect(c *gin.Context){
	var conn Connecti
	var hash string
	
	c.BindJSON(&conn)
	
	fmt.Println(conn.Password)
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




func createConv(c *gin.Context){

	var cre createGroup
	var groupid int;
	
	c.BindJSON(&verif)
	
	
	if(sessions[cre.Token].Username == cre.Username){
	
		rows := db.QueryRow("INSERT INTO groups(subject) DEFAULT VALUES RETURNING group_id;")
		rows.Scan(&groupid)
		
		rows := db.QueryRow("INSERT INTO participants(group_id, user_id) VALUES($1, 3);", groupid, cre.Username)
		for par in cre.Participant{ 
			rows := db.QueryRow("INSERT INTO participants(group_id, user_id) VALUES($1, 3);", groupid, par)
		}
		
		//c.Writer.Write([]byte(strings.Join(recepUL, ",")))
}


func sub(c *gin.Context){
	
	var count int;
	var subi Subscribe
	var infova vInfo
	
	c.BindJSON(&subi)
	
	
	
	if strings.TrimSpace(subi.Username) != ""{
		rows := db.QueryRow("SELECT count(*) FROM Users WHERE username=$1", subi.Username)
		fmt.Println(subi.Username)
		rows.Scan(&count)
		fmt.Println(count)
		
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
		
		
		
				result, err := db.Exec("INSERT INTO Users (username, password_hash, email, phonenumber, nboffer, date_joined, last_login, domain) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", subi.Username, subi.Password, subi.Email, subi.PhoneNumber, subi.Nboffer, subi.DateJoined, subi.LastLogin, "@office1789")
				if err != nil {
			  	panic(err)
				}
				fmt.Println(result)
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


func getinfop(c * gin.Context){

	var verif sessionSend;
	var infop Subscribe;
	
	c.BindJSON(&verif)

	if sessions[verif.Token].Username == verif.Username{
		rows := db.QueryRow("SELECT domain, nboffer, date_joined, last_login, phonenumber, email FROM Users WHERE username=$1", verif.Username)
		rows.Scan(&infop.Domain, &infop.Nboffer, &infop.DateJoined, &infop.LastLogin, &infop.PhoneNumber, &infop.Email)
		
		c.JSON(http.StatusOK, infop)
	
	}else{
		infop.Username = "no"
		c.JSON(http.StatusOK, infop)
	}

}

var db *sql.DB

func connectdb() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
    	var err error
    	db, err = sql.Open("postgres", psqlInfo)
    	
  	if err != nil {
    	panic(err)
  	}
    	
    
	err = db.Ping()
    	
    	if err != nil {
  	panic(err)
	}
	fmt.Println("Successfully connected!")
	
	

}

func main() {
	connectdb()
	r := gin.Default()
	
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	
	r.GET("/api/welcome", welcome)
	r.POST("/api/subscribe", sub)
	r.POST("/api/connect", connect)
	r.POST("/api/getinfop", getinfop)
	r.POST("/api/changeinfo", changeI)
	r.POST("/api/chat/createconv", createConv)
	
	r.Run() // listen and serve on 0.0.0.0:8080
}
