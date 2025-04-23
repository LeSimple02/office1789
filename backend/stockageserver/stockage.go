import (
	
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"	
	"github.com/gin-contrib/cors"
		
)

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	r.Use(cors.New(config))
	r.POST("/api/drive", drivei)
	r.POST("/api/driveaction", drivea)
	r.Run() // listen and serve on 0.0.0.0:8080
}
