package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func welcome(c *gin.Context) {
	c.Writer.Write([]byte("Hello welcome to the backend of Office1789\n"))
}

func main() {
	Connectdb()
	r := gin.Default()

	// CORS configuration with credentials support
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8082", // OnlyOffice DocumentServer
		"http://localhost:5173", // Vite dev server
		"http://localhost:8081", // or your frontend production URL
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
	}
	config.AllowCredentials = true // Important: enable credentials
	r.Use(cors.New(config))

	r.GET("/api/welcome", welcome)
	r.POST("/api/subscribe", Sub)
	r.POST("/api/connect", Connect)
	r.POST("/api/getinfop", getinfop)
	r.POST("/api/changeinfo", ChangeI)
	r.POST("/api/chat/createconv", createConv)
	r.GET("/api/drive/download", downloadFile)
	r.POST("/api/drive/upload", uploadFile)
	r.POST("/api/drive/getfiles", getfiles)
	r.POST("/api/drive/gettrash", getTrashFiles)
	r.POST("/api/drive/createFolder", createFolder)
	r.POST("/api/drive/rename", renameFile)
	r.POST("/api/drive/delete", deletePermanent)
	r.POST("/api/drive/trash", moveToTrash)
	r.POST("/api/drive/restore", restoreFile)
	r.POST("/api/drive/moveFile", moveFile)
	r.POST("/api/drive/moveFolder", moveFolder)
	r.GET("/api/onlyoffice/config", onlyofficeConfig)
	r.POST("/api/onlyoffice/callback", onlyofficeCallback)
	r.POST("/api/drive/shareFile", createShareFile)
	r.POST("/api/drive/deactivateShareFile", deactivateShareFile)

	r.Run(":8080")
}
