package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/auth"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/authorization"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestSignup(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []byte(`{"first_name": "test","last_name": "","email": "test@test.com","password": "test12--"}`)

	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	authorization.SignUp(c)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestSignupInvalid(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []byte(`{"first_name": "test","last_name": "","email": "test@test","password": "test12--"}`)

	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	authorization.SignUp(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestCheckIfFirst(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
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
	auth.CheckIfFirst(c)
	assert.Equal(t, http.StatusOK, w.Code)

	password := hash.Message("test12--", configuration.GetResp().PasswordHash)
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

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Request = req
	auth.CheckIfFirst(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

}

func TestSignIn(t *testing.T) {
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

	if err := SeedOneUser(&user); err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []byte(`{"email": "test@test.com","password": "test12--"}`)

	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req
	authorization.SignIn(c)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	data = []byte(`{"email": "test@test.com","password": ""}`)

	req, err = http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req
	authorization.SignIn(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	data = []byte(`{"email": "test@test.com","password": "test12"}`)

	req, err = http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req
	authorization.SignIn(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	data = []byte(`{"email": "","password": "test12--"}`)

	req, err = http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req
	authorization.SignIn(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestChangePassword(t *testing.T) {

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

	if err := SeedOneUser(&user); err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user_id", user.ID)

	data := []byte(`{"old_password": "test12","new_password": "test12.."}`)

	req, err := http.NewRequest("PUT", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req

	users.ChangePassword(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Set("user_id", user.ID)

	data = []byte(`{"old_password": "test12--","new_password": "test12.."}`)

	req, err = http.NewRequest("PUT", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req

	users.ChangePassword(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserProfile(t *testing.T) {

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

	if err := SeedOneUser(&user); err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	c.Set("user_id", user.ID+1)

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req

	users.GetUserProfile(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	c.Set("user_id", user.ID)

	req, err = http.NewRequest("GET", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req

	users.GetUserProfile(c)
	assert.Equal(t, http.StatusOK, w.Code)

	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, float64(user.ID), got["user_id"])
	assert.Equal(t, user.FirstName, got["first_name"])
	assert.Equal(t, user.Email, got["email"])
}
