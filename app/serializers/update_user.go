package serializers

import "code.jtg.tools/ayush.singhal/notifications-microservice/db/models"

// ChangeCredentialsInfo serializer to bind request data
type ChangeCredentialsInfo struct {
	Email string `json:"email" binding:"required"`
	Role  int    `json:"role" binding:"required"`
}

// ChangePasswordInfo serializer to bind request data
type ChangePasswordInfo struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangeCredentialsInfoToUserModel converts ChangeEmailInfo serializer to User model
func ChangeCredentialsInfoToUserModel(info *ChangeCredentialsInfo, user *models.User) {
	user.Email = info.Email
	user.Role = info.Role
}

// ChangePasswordInfoToUserModel converts ChangePasswordInfo serializer to User model
func ChangePasswordInfoToUserModel(info *ChangePasswordInfo, user *models.User) {
	user.Password = info.NewPassword
}
