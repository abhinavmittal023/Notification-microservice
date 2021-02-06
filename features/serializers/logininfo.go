package serializers

import "code.jtg.tools/ayush.singhal/notifications-microservice/db/models"

//LoginInfo serializer to bind request data
type LoginInfo struct {
	FirstName	string		`json:"first_name" binding:"required"`
	LastName	string		`json:"last_name"`
	Email        string      `json:"email" binding:"required"`
	Password     string      `json:"password,omitempty" binding:"required"`
	Role		int			`json:"role"`
}

//LoginInfoToUserModel converts LoginInfo serializer to User model
func LoginInfoToUserModel(info LoginInfo,user *models.User){
	user.FirstName = info.FirstName
	user.LastName = info.LastName
	user.Email = info.Email
	user.Password = info.Password
	user.Verified = false
	user.Role = info.Role
}