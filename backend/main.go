package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func welcome(c *gin.Context) {
	c.Writer.Write([]byte("Hello welcome to the backend of Office1789\n"))
}

type createGroup struct {
	Token       string   `json:token`
	Username    string   `json:username`
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
	r.POST("/api/getinfop", getinfop)
	r.POST("/api/changeinfo", ChangeI)
	r.POST("/api/chat/createconv", createConv)
	r.POST("/api/drive/getfiles", getfiles)
	r.POST("/api/drive/upload", uploadFile)
	r.GET("/api/drive/download", downloadFile)
	r.POST("/api/drive/rename", renameFile)
	r.POST("/api/drive/gettrash", getTrashFiles)
	r.POST("/api/drive/trash", moveToTrash)
	r.POST("/api/drive/delete", deletePermanent)
	r.POST("/api/drive/deletePermanent", deletePermanent)
	r.POST("/api/drive/restore", restoreFile)
	r.POST("/api/drive/createFolder", createFolder)
	r.POST("/api/drive/moveFile", moveFile)
	r.POST("/api/drive/moveFolder", moveFolder)

	r.Run() // listen and serve on 0.0.0.0:8080
}
