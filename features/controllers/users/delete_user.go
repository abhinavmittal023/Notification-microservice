package users

import (
	"log"
	"net/http"
	"strconv"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//DeleteUser Controller for /users/delete/:id route
func DeleteUser(c *gin.Context){
	val,_ := c.Get("role")
	if val != 2{
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userID,err := strconv.Atoi(c.Param("id"))
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("String Conversion Error")
		return
	}
	var user models.User
	err = users.GetUserWithID(&user,uint64(userID))
	if err == gorm.ErrRecordNotFound{
		c.JSON(http.StatusBadRequest, gin.H{"id_not_present":"Id not in database"})
		return
	}

	err = users.SoftDeleteUser(&user)
	if err!= nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_server_error":"Internal Server Error"})
		log.Println("Delete User Service Error")
		return
	}
	c.JSON(http.StatusOK,gin.H{"status":"ok"})
}