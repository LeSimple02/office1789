package main

import (
	
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"	
	"github.com/gin-contrib/cors"
	

	
)

func welcome(c *gin.Context){
	c.Writer.Write([]byte("Hello welcome to the backend of Office1789\n"))
}


type createGroup struct {
	Token string	`json:token`
	Username string `json:username`
	Participant []string `json:recepter`
	
	
}



func main() {
	Connectdb()
	r := gin.Default()
	
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	
	r.GET("/api/welcome", welcome)
	r.POST("/api/subscribe", Sub)
	r.POST("/api/connect", Connect)
	r.POST("/api/getinfop", Getinfop)
	r.POST("/api/changeinfo", ChangeI)
	r.POST("/api/chat/createconv", createConv)
	
	r.Run() // listen and serve on 0.0.0.0:8080
}
