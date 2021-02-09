package features

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/middlewares"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/signin"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/signup"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/routes/users"
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

	loginGroup := v1.Group("/login",middlewares.CheckIfLogged())
	signin.SignIn(loginGroup)

	userGroup := v1.Group("/users",middlewares.AuthorizeJWT())
	users.AddUser(userGroup)
	users.ChangeUserEmail(userGroup)
	users.ChangeUserRole(userGroup)
	users.ChangeUserPassword(userGroup)
	users.DeleteUser(userGroup)
	users.GetUser(userGroup)
	users.GetAllUsers(userGroup)

	err := router.Run(":" + configuration.GetResp().Server.Port)
	return err
}
