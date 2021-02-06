package features

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/middlewares"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/signup"
	"github.com/gin-gonic/gin"
)

//InitServer is used to initialize server routes
func InitServer() error {

	router := gin.Default()
	//setting the cors headers
	router.Use(middlewares.CorsHeaders())

	v1 := router.Group("/api/v1")

	healthCheck := v1.Group("/ping")

	// healthCheck contains the /ping Health Check Endpoint
	healthCheck.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	firstSignUp := v1.Group("/signup",middlewares.CheckIfFirst())
	signup.SignUp(firstSignUp)

	authorization := v1.Group("/auth")
	auth.Auth(authorization)

	err := router.Run(":" + configuration.GetResp().Server.Port)
	return err
}
