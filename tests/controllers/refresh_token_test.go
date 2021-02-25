package controllers

import (
	"bytes"
	"encoding/json"
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

func TestRefeshTokenInvalid(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(`{"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsIlJvbGUiOjAsIlRva2VuVHlwZSI6InJlZnJlc2giLCJleHAiOjE2MTQyMjk4NzYsImlhdCI6MTYxNDIxOTA3NiwiaXNzIjoibm90aWZpY2F0aW9uLW1pY3Jvc2VydmljZSJ9.tqv61EE5rwSi0gA2pMQe-TMhqYk3GUPB76vmaXaGHdo"}`)
	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req

	auth.RefreshAccessToken(c)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefreshToken(t *testing.T) {
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

	refreshToken, err := sharedAuth.GenerateRefreshToken(uint64(user.ID), configuration.GetResp().Token.ExpiryTime.RefreshToken)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(fmt.Sprintf(`{"refresh_token":"%s"}`, refreshToken))
	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	c.Request = req

	auth.RefreshAccessToken(c)
	assert.Equal(t, http.StatusOK, w.Code)

	got := gin.H{}

	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	_, err = sharedAuth.ValidateToken(got["access_token"].(string))
	if err != nil {
		t.Fatal(err)
	}
}
