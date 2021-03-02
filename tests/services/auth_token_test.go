package services

import (
	"log"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/authservice"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/go-playground/assert/v2"
)

func TestValidateToken(t *testing.T) {
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
		Verified:  false,
		Role:      2,
	}

	err = SeedOneUser(&user)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	tokenString, err := auth.GenerateValidationToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.ValidationToken)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	user1, err := authservice.ValidateToken(tokenString, "validation")
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	assert.Equal(t, user.ID, user1.ID)

	tokenString, err = auth.GenerateRefreshToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.RefreshToken)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	user1, err = authservice.ValidateToken(tokenString, "refresh")
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	assert.Equal(t, user.ID, user1.ID)

	tokenString, err = auth.GenerateAccessToken(uint64(user.ID), user.Role, configuration.GetResp().Token.ExpiryTime.AccessToken)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	user1, err = authservice.ValidateToken(tokenString, "access")
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	assert.Equal(t, user.ID, user1.ID)
	assert.Equal(t, user.Role, user1.Role)
}

func TestValidateTokenInvalid(t *testing.T) {
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

	tokenString, err := auth.GenerateValidationToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.ValidationToken)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	user1, err := authservice.ValidateToken(tokenString, "refresh")
	assert.Equal(t, err.Error(), "Invalid Token")
	assert.Equal(t, user1, nil)

	tokenString, err = auth.GenerateRefreshToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.RefreshToken)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	user1, err = authservice.ValidateToken(tokenString, "access")
	assert.Equal(t, err.Error(), "Invalid Token")
	assert.Equal(t, user1, nil)

	tokenString, err = auth.GenerateAccessToken(uint64(user.ID), user.Role, configuration.GetResp().Token.ExpiryTime.AccessToken)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	user1, err = authservice.ValidateToken(tokenString, "refresh")
	assert.Equal(t, err.Error(), "Invalid Token")
	assert.Equal(t, user1, nil)
}
