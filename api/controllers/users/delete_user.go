package users

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/api/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DeleteUserRoute is used to delete users from database
func DeleteUserRoute(router *gin.RouterGroup) {
	router.DELETE("/delete/:id", DeleteUser)
}

// DeleteUser Controller for /users/delete/:id route
func DeleteUser(c *gin.Context) {
	val, _ := c.Get(constants.Role)
	if val != 2 {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("String Conversion Error")
		return
	}
	user, err := users.GetUserWithID(uint64(userID))
	if err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not in database"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("GetUserWithID service error")
		return
	}

	err = users.SoftDeleteUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		log.Println("Delete User Service Error")
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
