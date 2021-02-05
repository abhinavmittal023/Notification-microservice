package features

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/middlewares"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/signup"
	"github.com/gin-gonic/gin"
)

//InitServer is used to initialize server routes
func InitServer() error {

	router := gin.Default()
	//setting the cors headers
	router.Use(middlewares.CorsHeaders())

	api := router.Group("/api")
	v1 := api.Group("/v1")
	healthCheck := v1.Group("/ping")

	// healthCheck contains the /ping Health Check Endpoint
	healthCheck.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	signup.CreateUser(v1)
	err := router.Run(":" + configuration.GetResp().Server.Port)
	return err
}
