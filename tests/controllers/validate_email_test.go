package controllers

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	sharedAuth "code.jtg.tools/ayush.singhal/notifications-microservice/shared/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestValidateEmail(t *testing.T) {
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

	if err := SeedOneUser(&user); err != nil {
		t.Fail()
	}

	tokenString, err := sharedAuth.GenerateValidationToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.ValidationToken)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "token",
		Value: fmt.Sprintf("%s", tokenString),
	})
	auth.ValidateEmail(c)
	assert.Equal(t, http.StatusFound, w.Code)

	tokenString = tokenString[:len(tokenString)-1]

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	req, err = http.NewRequest("GET", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "token",
		Value: fmt.Sprintf("%s", tokenString),
	})
	auth.ValidateEmail(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	tokenString, err = sharedAuth.GenerateValidationToken(uint64(user.ID)+1, configuration.GetResp().Token.ExpiryTime.ValidationToken)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	req, err = http.NewRequest("GET", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "token",
		Value: fmt.Sprintf("%s", tokenString),
	})
	auth.ValidateEmail(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
