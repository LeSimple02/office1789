package main


import (
	"database/sql"
	"fmt"
)




const (
	host = "localhost"
	port = 5432
	user = "robespierre"
	password = "guillotine"
	dbname = "office1789"
)


var db *sql.DB

func Connectdb() {
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
