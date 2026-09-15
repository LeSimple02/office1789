package main

import(
	"github.com/gin-gonic/gin"
)


func createConv(c *gin.Context){

	/*
	var cre createGroup
	var groupid int;
	
	c.BindJSON(&cre)
	
	session, valid := validateSession(cre.Token, cre.Username)
	if valid {
	
		rows := db.QueryRow("INSERT INTO groups(subject) DEFAULT VALUES RETURNING group_id;")
		rows.Scan(&groupid)
		
		db.QueryRow("INSERT INTO participants(group_id, user_id) VALUES($1, $2);", groupid, session.UserID)
		for _, par := range cre.Participant { 
			db.QueryRow("INSERT INTO participants(group_id, user_id) VALUES($1, $2);", groupid, par)
		}
		
		//c.Writer.Write([]byte(strings.Join(recepUL, ",")))
	} 
	*/
	
}



