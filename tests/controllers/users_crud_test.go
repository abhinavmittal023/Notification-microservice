package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/users"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestGetAllUsers(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password, err := hash.Message("test12--", configuration.GetResp().PasswordHash)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
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

	err = SeedUsers(&usersList)
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
	users.GetAllUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	for i := range usersList {
		assert.Equal(t, float64(usersList[len(usersList)-1-i].ID), got["users"].([]interface{})[i].(map[string]interface{})["user_id"])
		assert.Equal(t, (usersList[len(usersList)-1-i].Email), got["users"].([]interface{})[i].(map[string]interface{})["email"])
		assert.Equal(t, (usersList[len(usersList)-1-i].FirstName), got["users"].([]interface{})[i].(map[string]interface{})["first_name"])
	}
}

func TestGetAllUsersPagination(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	password, err := hash.Message("test12--", configuration.GetResp().PasswordHash)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
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

	err = SeedUsers(&usersList)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest("GET", "?offset=0&limit=1", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	users.GetAllUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["users"].([]interface{})), 1)
	assert.Equal(t, float64(usersList[len(usersList)-1].ID), got["users"].([]interface{})[0].(map[string]interface{})["user_id"])
	assert.Equal(t, (usersList[len(usersList)-1].Email), got["users"].([]interface{})[0].(map[string]interface{})["email"])
	assert.Equal(t, (usersList[len(usersList)-1].FirstName), got["users"].([]interface{})[0].(map[string]interface{})["first_name"])

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	req, err = http.NewRequest("GET", "?offset=1&limit=1", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	users.GetAllUsers(c)

	assert.Equal(t, http.StatusOK, w.Code)

	got = gin.H{}
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["users"].([]interface{})), 1)
	assert.Equal(t, float64(usersList[len(usersList)-2].ID), got["users"].([]interface{})[0].(map[string]interface{})["user_id"])
	assert.Equal(t, (usersList[len(usersList)-2].Email), got["users"].([]interface{})[0].(map[string]interface{})["email"])
	assert.Equal(t, (usersList[len(usersList)-2].FirstName), got["users"].([]interface{})[0].(map[string]interface{})["first_name"])
}

func TestAddUser(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []byte(`{"first_name": "test","last_name": "","email": "test@test.com","role": 1}`)

	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	users.AddUser(c)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	users.AddUser(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddUserInvalid(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(`{"first_name": "test","last_name": "","email": "test@test","password": "test12--","role": 1}`)

	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	users.AddUser(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	data = []byte(`{"first_name": "test","last_name": "","email": "test@test.com","password": "test12--","role": -1}`)

	req, err = http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	users.AddUser(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestDeleteUser(t *testing.T) {

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
		log.Println(err)
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest("DELETE", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", user.ID),
	})
	users.DeleteUser(c)

	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	req, err = http.NewRequest("DELETE", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", user.ID+1),
	})
	users.DeleteUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestGetUser(t *testing.T) {

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
		Key:   "id",
		Value: fmt.Sprintf("%v", user.ID),
	})
	users.GetUser(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, float64(user.ID), got["user_details"].(map[string]interface{})["user_id"])
	assert.Equal(t, (user.Email), got["user_details"].(map[string]interface{})["email"])
	assert.Equal(t, (user.FirstName), got["user_details"].(map[string]interface{})["first_name"])

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
		Key:   "id",
		Value: fmt.Sprintf("%v", user.ID+1),
	})
	users.GetUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestUpdateUser(t *testing.T) {

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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []byte(`{"email": "test1@test.com","role": 1}`)

	req, err := http.NewRequest("PUT", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", user.ID),
	})
	users.UpdateUser(c)

	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", user.ID+1),
	})
	users.UpdateUser(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateUserInvlaid(t *testing.T) {

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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data := []byte(`{"email": "test1@test","role": 1}`)

	req, err := http.NewRequest("PUT", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", user.ID),
	})
	users.UpdateUser(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	data = []byte(`{"email": "test1@test.com","role": -1}`)

	req, err = http.NewRequest("PUT", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", user.ID),
	})
	users.UpdateUser(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
