package controllers

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/authorization"
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

	c.Header("Content-Type", "application/json; charset=utf-8")

	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	authorization.SignUp(c)

	log.Printf("%s", w.Body)
	assert.Equal(t, 200, w.Code)

}
