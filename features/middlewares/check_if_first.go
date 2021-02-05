package middlewares

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//TODO: Handle case when first user fails to validate the emailID

//CheckIfFirst middleware checks the if another user exists to avoid creation of other user directly
func CheckIfFirst() gin.HandlerFunc {
	return func(c *gin.Context) {

		err := users.GetFirstUser(&models.User{})
		if err == gorm.ErrRecordNotFound {
			c.Next()
		}
		if err != nil {
			log.Println(err)
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
