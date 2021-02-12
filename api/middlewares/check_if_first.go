package middlewares

import (
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// TODO: Handle case when first user fails to validate the emailID

// CheckIfFirst middleware checks the if another user exists to avoid creation of other user directly
func CheckIfFirst() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Next() //Preflight Request
			return
		}
		_, err := users.GetFirstUser()
		if err == gorm.ErrRecordNotFound {
			c.Next()
			return
		} else if err != nil {
			log.Println(err.Error())
		}
		c.AbortWithStatus(http.StatusUnauthorized)
	}
}
