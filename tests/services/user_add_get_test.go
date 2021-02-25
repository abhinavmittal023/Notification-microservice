package services

import (
	"log"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/gorm"
)

func TestGetUserWithID(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password := hash.Message("test12--", configuration.GetResp().PasswordHash)
	user := models.User{
		FirstName: "test",
		Email:     "test@test.com",
		Password:  password,
		Verified:  true,
		Role:      2,
	}

	err := SeedOneUser(&user)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	gotUser, gotError := users.GetUserWithID(uint64(user.ID))

	assert.Equal(t, gotError, nil)
	assert.Equal(t, user.ID, gotUser.ID)
	assert.Equal(t, (user.Email), gotUser.Email)
	assert.Equal(t, (user.FirstName), gotUser.FirstName)

	gotUser, gotError = users.GetUserWithID(uint64(user.ID) + 1)
	assert.Equal(t, gotError, gorm.ErrRecordNotFound)

}

func TestGetFirstUser(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password := hash.Message("test12--", configuration.GetResp().PasswordHash)
	usersList := []models.User{
		{
			FirstName: "test1",
			Email:     "test1@test.com",
			Password:  password,
			Verified:  true,
			Role:      2,
		},
		{
			FirstName: "test2",
			Email:     "test2@test.com",
			Password:  password,
			Verified:  true,
			Role:      1,
		},
		{
			FirstName: "test3",
			Email:     "test3@test.com",
			Password:  password,
			Verified:  false,
			Role:      1,
		},
	}

	err := SeedUsers(&usersList)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	gotUser, gotError := users.GetFirstUser()

	assert.Equal(t, gotError, nil)
	assert.Equal(t, usersList[0].ID, gotUser.ID)
	assert.Equal(t, (usersList[0].Email), gotUser.Email)
	assert.Equal(t, (usersList[0].FirstName), gotUser.FirstName)

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	gotUser, gotError = users.GetFirstUser()

	assert.Equal(t, gotError, gorm.ErrRecordNotFound)

	user := models.User{
		FirstName: "test",
		Email:     "test@test.com",
		Password:  password,
		Verified:  false,
		Role:      2,
	}

	err = SeedOneUser(&user)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	gotUser, gotError = users.GetFirstUser()

	assert.Equal(t, gotError, nil)
	assert.Equal(t, user.ID, gotUser.ID)
	assert.Equal(t, (user.Email), gotUser.Email)
	assert.Equal(t, (user.FirstName), gotUser.FirstName)
}

func TestGetUserWithEmail(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password := hash.Message("test12--", configuration.GetResp().PasswordHash)
	usersList := []models.User{
		{
			FirstName: "test1",
			Email:     "test1@test.com",
			Password:  password,
			Verified:  true,
			Role:      2,
		},
		{
			FirstName: "test2",
			Email:     "test2@test.com",
			Password:  password,
			Verified:  true,
			Role:      1,
		},
		{
			FirstName: "test3",
			Email:     "test3@test.com",
			Password:  password,
			Verified:  false,
			Role:      1,
		},
	}

	err := SeedUsers(&usersList)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	for i := range usersList {
		gotUser, gotError := users.GetUserWithEmail(usersList[i].Email)
		assert.Equal(t, gotError, nil)
		assert.Equal(t, gotUser.ID, usersList[i].ID)
		assert.Equal(t, gotUser.FirstName, usersList[i].FirstName)
	}

	_, gotError := users.GetUserWithEmail("wrongemail@test.com")
	assert.Equal(t, gotError, gorm.ErrRecordNotFound)
}

func TestGetAllUsers(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password := hash.Message("test12--", configuration.GetResp().PasswordHash)
	usersList := []models.User{
		{
			FirstName: "test1",
			Email:     "test1@test.com",
			Password:  password,
			Verified:  true,
			Role:      2,
		},
		{
			FirstName: "test2",
			Email:     "test2@test.com",
			Password:  password,
			Verified:  true,
			Role:      1,
		},
		{
			FirstName: "test3",
			Email:     "test3@test.com",
			Password:  password,
			Verified:  false,
			Role:      1,
		},
	}

	err := SeedUsers(&usersList)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	pagination := serializers.Pagination{
		Limit:  3,
		Offset: 0,
	}

	filter := filter.User{
		ID:        0,
		FirstName: "",
		LastName:  "",
		Email:     "",
		Verified:  0,
		Role:      0,
	}

	gotUser, gotError := users.GetAllUsers(&pagination, &filter)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, len(gotUser), 3)

	for i := range usersList {
		assert.Equal(t, usersList[len(usersList)-1-i].ID, gotUser[i].ID)
		assert.Equal(t, (usersList[len(usersList)-1-i].Email), gotUser[i].Email)
		assert.Equal(t, (usersList[len(usersList)-1-i].FirstName), gotUser[i].FirstName)
	}
}

func TestGetAllUsersFilters(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password := hash.Message("test12--", configuration.GetResp().PasswordHash)
	usersList := []models.User{
		{
			FirstName: "test1",
			Email:     "test1@test.com",
			Password:  password,
			Verified:  true,
			Role:      2,
		},
		{
			FirstName: "test2",
			Email:     "test2@test.com",
			Password:  password,
			Verified:  true,
			Role:      1,
		},
		{
			FirstName: "test3",
			Email:     "test3@test.com",
			Password:  password,
			Verified:  false,
			Role:      1,
		},
	}

	err := SeedUsers(&usersList)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	pagination := serializers.Pagination{
		Limit:  1,
		Offset: 0,
	}

	filter := filter.User{
		ID:        0,
		FirstName: "",
		LastName:  "",
		Email:     "",
		Verified:  0,
		Role:      1,
	}

	gotUser, gotError := users.GetAllUsers(&pagination, &filter)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, len(gotUser), 1)

	assert.Equal(t, usersList[2].ID, gotUser[0].ID)
	assert.Equal(t, (usersList[2].Email), gotUser[0].Email)
	assert.Equal(t, (usersList[2].FirstName), gotUser[0].FirstName)
	assert.Equal(t, (usersList[2].Role), gotUser[0].Role)
}

func TestCreateUser(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password := hash.Message("test12--", configuration.GetResp().PasswordHash)
	user := models.User{
		FirstName: "test1",
		Email:     "test1@test.com",
		Password:  password,
		Verified:  true,
		Role:      2,
	}

	gotError := users.CreateUser(&user)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, user.ID, uint(1))

}
