package signin

import (
	"encoding/json"
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/features/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

//SignIn Controller for /signin route
func SignIn(c *gin.Context){
	var info serializers.LoginInfo
	if c.BindJSON(&info) != nil {
		c.JSON(http.StatusBadRequest,gin.H{"required":"Email,Password are required"})
		return
	}
	
	var user models.User
	err := users.GetUserWithEmail(&user,info.Email)
	if err == gorm.ErrRecordNotFound{
		c.JSON(http.StatusUnauthorized, gin.H{"email_not_present":"EmailId not in database"})
		return
	}
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"internal_error":"Internal Server Error"})
		return
	}

	if !user.Verified{
		c.JSON(http.StatusUnauthorized, gin.H{"email_not_verified":"EmailId not verified"})
		return
	}

	if !hash.Validate(info.Password,user.Password,configuration.GetResp().PasswordHash){
		c.JSON(http.StatusUnauthorized, gin.H{"incorrect_password":"Passwords mismatch"})
		return
	}
	
	info.FirstName = user.FirstName
	info.LastName = user.LastName
	info.Password = ""
	var token serializers.RefreshToken

	token.AccessToken,err = auth.GenerateAccessToken(uint64(user.ID),user.Role,configuration.GetResp().Token.ExpiryTime.AccessToken)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"access_token_generation_error":"Access Token not generated"})
		return
	}
	token.RefreshToken,err = auth.GenerateRefreshToken(uint64(user.ID),configuration.GetResp().Token.ExpiryTime.RefreshToken)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"refresh_token_generation_error":"Refresh Token not generated"})
		return
	}

	js, err := json.Marshal(&serializers.LoginResponse{
		LoginInfo: info,
		RefreshToken: token,
	})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"signin_json_marshal_error": "JSON marshalling error"})
		return
	}
	c.Data(http.StatusOK, "application/json", js)
}
