package app

import (
	"time"

	"github.com/MyriadFlow/gateway/api"
	"github.com/MyriadFlow/gateway/util/pkg/auth"
	"github.com/MyriadFlow/gateway/util/pkg/logwrapper"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"

	"github.com/MyriadFlow/gateway/config/dbconfig/dbinit"
	"github.com/MyriadFlow/gateway/config/envconfig"
	"github.com/gin-gonic/gin"
)

var GinApp *gin.Engine

func Init() {
	envconfig.InitEnvVars()
	gin.SetMode(envconfig.EnvVars.APP_MODE)
	auth.Init()
	logwrapper.Init()
	dbinit.Init()

	GinApp = gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	config.ExposeHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	// serve static files
	GinApp.Use(static.Serve("/", static.LocalFile("./web", false)))
	GinApp.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"status": 404, "message": "Invalid Endpoint Request"})
	})
	GinApp.Use(cors.New(config))
	api.ApplyRoutes(GinApp)
}
