package users

import (
	"time"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/jinzhu/gorm"
)

// GetUserWithID gets the user with specified ID from the database, and returns error/nil
func GetUserWithID(userID uint64) (*models.User, error) {
	var user models.User
	res := db.Get().First(&user, userID)
	return &user, res.Error
}

// GetFirstUser gets the details of the first user in the database
// Use with signup guards only
func GetFirstUser() (*models.User, error) {
	var user models.User
	res := db.Get().First(&user)
	if res.Error == nil && !user.Verified && time.Duration((time.Now().Unix()-user.CreatedAt.Unix()))-3600*configuration.GetResp().Token.ExpiryTime.ValidationToken > 0 {
		if err := DeleteUser(&user); err != nil {
			return nil, err
		}
		return nil, gorm.ErrRecordNotFound
	}
	return &user, res.Error
}

// GetUserWithEmail gets the user with specified Email from the database
func GetUserWithEmail(email string) (*models.User, error) {
	var user models.User
	res := db.Get().Where("email = ?", email).First(&user)
	return &user, res.Error
}

// GetAllUsers gets all users from the database and returns []models.User,err
func GetAllUsers(pagination *serializers.Pagination, userFilter *filter.User) ([]models.User, error) {

	var users []models.User
	dbg := db.Get()
	tx := dbg.Model(&models.User{})

	if userFilter.ID != 0 {
		tx = tx.Where("id = ?", userFilter.ID)
	}
	if userFilter.FirstName != "" {
		tx = tx.Where("first_name = ?", userFilter.FirstName)
	}
	if userFilter.LastName != "" {
		tx = tx.Where("last_name = ?", userFilter.LastName)
	}
	if userFilter.Email != "" {
		tx = tx.Where("email = ?", userFilter.Email)
	}
	if userFilter.Verified > 0 {
		tx = tx.Where("verified = ?", true)
	} else if userFilter.Verified < 0 {
		tx = tx.Where("verified = ?", false)
	}
	if userFilter.Role != 0 {
		tx = tx.Where("role = ?", userFilter.Role)
	}

	res := tx.Offset(pagination.Offset).Limit(pagination.Limit).Find(&users)
	return users, res.Error
}
