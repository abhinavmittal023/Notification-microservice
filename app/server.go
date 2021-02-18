package app

import (
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/authorization"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/recipients"
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
	users.ChangeOwnPasswordRoute(ownInfoGroup)
	users.GetUserProfileRoute(ownInfoGroup)

	userGroup := v1.Group("/users", middlewares.AuthorizeJWT(), middlewares.CheckIfSystemAdmin())
	users.GetAllUsersRoute(userGroup)
	users.AddUserRoute(userGroup)
	users.DeleteUserRoute(userGroup)
	users.GetUserRoute(userGroup)
	users.UpdateUserRoute(userGroup)
	users.ChangeDifferentUserPasswordRoute(userGroup)

	recipientGroup := v1.Group("/recipients", middlewares.AuthorizeJWT(), middlewares.CheckIfSystemAdmin())
	recipients.AddUpdateRecipientRoute(recipientGroup)
	recipients.GetRecipientRoute(recipientGroup)
	recipients.GetAllRecipientRoute(recipientGroup)

	channelGroup := v1.Group("/channels", middlewares.AuthorizeJWT(), middlewares.CheckIfSystemAdmin())
	channels.AddChannelRoute(channelGroup)
	channels.GetAllChannelsRoute(channelGroup)
	channels.UpdateChannelRoute(channelGroup)
	channels.DeleteChannelRoute(channelGroup)
	channels.GetChannelRoute(channelGroup)

	err := router.Run(":" + configuration.GetResp().Server.Port)
	return errors.Wrap(err, "Unable to run server")
}
