package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/recipients"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestGetRecipient(t *testing.T){
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	recipient := models.Recipient{
		RecipientID: "abcd1",
		Email: "abcd@gmail.com",
		PushToken: "abcd",
	}

	err := SeedOneRecipient(&recipient)
	if err != nil {
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
		Value: fmt.Sprintf("%v", recipient.ID),
	})
	recipients.GetRecipient(c)
	assert.Equal(t, http.StatusOK, w.Code)

	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	got = got["recipient_details"].(map[string]interface{})
	assert.Equal(t,float64(recipient.ID),got["id"])
	assert.Equal(t,recipient.RecipientID,got["recipient_id"])
	assert.Equal(t,recipient.Email,got["email"])
	assert.Equal(t,recipient.PushToken,got["push_token"])

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
		Value: fmt.Sprintf("%v", recipient.ID+1),
	})
	recipients.GetRecipient(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)	
}