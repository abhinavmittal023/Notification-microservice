package serializers

import (
	"strings"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// ChangeCredentialsInfo serializer to bind request data
type ChangeCredentialsInfo struct {
	ID    uint64 `json:"-"`
	Email string `json:"email" binding:"required"`
	Role  int    `json:"role" binding:"required"`
}

// ChangePasswordInfo serializer to bind request data
type ChangePasswordInfo struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangeCredentialsInfoToUserModel converts ChangeEmailInfo serializer to User model
func ChangeCredentialsInfoToUserModel(info *ChangeCredentialsInfo, user *models.User) {
	user.Email = strings.ToLower(info.Email)
	user.Role = int(info.Role)
}

// ChangePasswordInfoToUserModel converts ChangePasswordInfo serializer to User model
func ChangePasswordInfoToUserModel(info *ChangePasswordInfo, user *models.User) {
	user.Password = info.NewPassword
}
