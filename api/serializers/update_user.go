package serializers

import "code.jtg.tools/ayush.singhal/notifications-microservice/db/models"

// ChangeEmailInfo serializer to bind request data
type ChangeEmailInfo struct {
	OldEmail string `json:"old_email" binding:"required"`
	NewEmail string `json:"new_email" binding:"required"`
}

// ChangeRoleInfo serializer to bind request data
type ChangeRoleInfo struct {
	Email string `json:"email" binding:"required"`
	Role  int    `json:"role" binding:"required"`
}

// ChangePasswordInfo serializer to bind request data
type ChangePasswordInfo struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

// ChangeEmailInfoToUserModel converts ChangeEmailInfo serializer to User model
func ChangeEmailInfoToUserModel(info *ChangeEmailInfo, user *models.User) {
	user.Email = info.NewEmail
	user.Verified = false
}

// RevertChangesToUserEmailModel converts back the info serializer to User model
func RevertChangesToUserEmailModel(info *ChangeEmailInfo, user *models.User) {
	user.Email = info.OldEmail
	user.Verified = true
}

// ChangeRoleInfoToUserModel converts ChangeEmailInfo serializer to User model
func ChangeRoleInfoToUserModel(info *ChangeRoleInfo, user *models.User) {
	user.Role = info.Role
}

// RevertChangesToUserRoleModel converts back the info serializer to User model
func RevertChangesToUserRoleModel(info *ChangeRoleInfo, user *models.User) {
	user.Role = info.Role
}

//ChangePasswordInfoToUserModel converts ChangePasswordInfo serializer to User model
func ChangePasswordInfoToUserModel(info *ChangePasswordInfo, user *models.User) {
	user.Password = info.NewPassword
}

//RevertChangesToUserPasswordModel converts back the info serializer to User model
func RevertChangesToUserPasswordModel(info *ChangePasswordInfo, user *models.User) {
	user.Password = info.OldPassword
}
