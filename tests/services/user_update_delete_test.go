package services

import (
	"log"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/gorm"
)

func TestPatchUser(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password, err := hash.Message("test12--", configuration.GetResp().PasswordHash)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	user := models.User{
		FirstName: "test",
		Email:     "test@test.com",
		Password:  password,
		Verified:  true,
		Role:      2,
	}

	err = SeedOneUser(&user)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	user.FirstName = "changed"
	user.LastName = "name"

	gotError := users.PatchUser(&user)
	assert.Equal(t, gotError, nil)

	newUser := models.User{}
	err = db.Get().First(&newUser).Error
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	assert.Equal(t, newUser.ID, uint(1))
	assert.Equal(t, newUser.Email, "test@test.com")
	assert.Equal(t, newUser.LastName, "name")
	assert.NotEqual(t, newUser.CreatedAt, newUser.UpdatedAt)
}

func TestSoftDeleteUser(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password, err := hash.Message("test12--", configuration.GetResp().PasswordHash)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	user := models.User{
		FirstName: "test",
		Email:     "test@test.com",
		Password:  password,
		Verified:  true,
		Role:      2,
	}

	err = SeedOneUser(&user)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	gotError := users.DeleteUser(&user)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, db.Get().First(models.User{}, user.ID).Error, gorm.ErrRecordNotFound)
}
