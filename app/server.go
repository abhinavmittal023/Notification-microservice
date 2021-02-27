package app

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/authorization"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/middlewares"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// InitServer is used to initialize server routes
func InitServer() error {
	router := gin.Default()
	// setting the cors headers
	router.Use(middlewares.CorsHeaders())

	v1 := router.Group("/api/v1")

	healthCheck := v1.Group("/health-check")

	// healthCheck contains the /health-check Health Check Endpoint
	healthCheck.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	firstSignUp := v1.Group("/signup", middlewares.CheckIfFirst())
	authorization.SignUpRoute(firstSignUp)

	authToken := v1.Group("/auth")
	auth.RefreshAccessTokenRoute(authToken)
	auth.ValidateEmailRoute(authToken)
	auth.CheckIfFirstRoute(authToken)

	loginGroup := v1.Group("/login", middlewares.CheckIfLogged())
	authorization.SignInRoute(loginGroup)

	ownInfoGroup := v1.Group("/profile", middlewares.AuthorizeJWT())
	users.GetUserProfileRoute(ownInfoGroup)
	users.ChangeOwnPasswordRoute(ownInfoGroup)

	UserGroup := v1.Group("/users", middlewares.AuthorizeJWT(), middlewares.CheckIfSystemAdmin())
	users.AddUserRoute(UserGroup)
	users.ChangeUserCredentialsRoute(UserGroup)
	users.DeleteUserRoute(UserGroup)
	users.GetUserRoute(UserGroup)
	users.GetAllUsersRoute(UserGroup)

	err := router.Run(":" + configuration.GetResp().Server.Port)
	return errors.Wrap(err, "Unable to run server")
}
