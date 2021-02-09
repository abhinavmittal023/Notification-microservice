package serializers

import (
	"time"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

//UserInfo serializer to bind request data
type UserInfo struct {
	FirstName     string    `json:"first_name"`
	LastName      string	`json:"last_name,omitempty"`
	Email         string    `json:"email"`
	Verified		bool	`json:"verified"`
	Role			int		`json:"role"`
	CreatedAt		time.Time	`json:"created_at"`
	UpdatedAt		time.Time	`json:"updated_at"`
}

//UserModelToUserInfo converts the user model to UserInfo struct
func UserModelToUserInfo(info *UserInfo, user *models.User){
	info.FirstName = user.FirstName
	info.LastName = user.LastName
	info.Email = user.Email
	info.Verified = user.Verified
	info.Role = user.Role
	info.CreatedAt = user.CreatedAt
	info.UpdatedAt = user.UpdatedAt
}