package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/recipients"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestGetRecipient(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	recipient := models.Recipient{
		RecipientID: "abcd1",
		Email:       "abcd@gmail.com",
		PushToken:   "abcd",
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
	assert.Equal(t, float64(recipient.ID), got["id"])
	assert.Equal(t, recipient.RecipientID, got["recipient_id"])
	assert.Equal(t, recipient.Email, got["email"])
	assert.Equal(t, recipient.PushToken, got["push_token"])

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

func TestGetRecipientWithChannel(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	channel := models.Channel{
		Name:     "email",
		Type:     1,
		Priority: 1,
	}
	if err := SeedOneChannel(&channel); err != nil {
		log.Println(err)
		t.Fail()
	}

	recipient := models.Recipient{
		RecipientID:          "abcd1",
		Email:                "abcd@gmail.com",
		PushToken:            "abcd",
		PreferredChannelType: 1,
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

	recipientDetails := got["recipient_details"].(map[string]interface{})
	preferredChannel := got["preferred_channel"].(map[string]interface{})
	assert.Equal(t, float64(recipient.ID), recipientDetails["id"])
	assert.Equal(t, recipient.RecipientID, recipientDetails["recipient_id"])
	assert.Equal(t, recipient.Email, recipientDetails["email"])
	assert.Equal(t, recipient.PushToken, recipientDetails["push_token"])

	assert.Equal(t, float64(channel.ID), preferredChannel["id"])
	assert.Equal(t, channel.Name, preferredChannel["name"])
	assert.Equal(t, float64(channel.Priority), preferredChannel["priority"])
	assert.Equal(t, float64(channel.Type), preferredChannel["type"])
}

func TestGetRecipientWithDeletedChannel(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	recipient := models.Recipient{
		RecipientID:          "abcd1",
		Email:                "abcd@gmail.com",
		PushToken:            "abcd",
		PreferredChannelType: 1,
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

	recipientDetails := got["recipient_details"].(map[string]interface{})
	assert.Equal(t, float64(recipient.ID), recipientDetails["id"])
	assert.Equal(t, recipient.RecipientID, recipientDetails["recipient_id"])
	assert.Equal(t, recipient.Email, recipientDetails["email"])
	assert.Equal(t, recipient.PushToken, recipientDetails["push_token"])

	assert.Equal(t, "Preferred Channel Email was Deleted", got["warning"])
}

func TestGetAllRecipients(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	recipientsList := []models.Recipient{
		{
			RecipientID: "abcd1",
			Email:       "abcd@gmail.com",
			PushToken:   "abcd",
		},
		{
			RecipientID: "abcd2",
			Email:       "abcd1@gmail.com",
			PushToken:   "abcd",
		},
		{
			RecipientID: "abcd3",
			Email:       "abcd2@gmail.com",
			PushToken:   "abcd",
		},
	}

	err := SeedRecipients(&recipientsList)
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
	recipients.GetAllRecipient(c)
	assert.Equal(t, http.StatusOK, w.Code)

	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	recipientRecords := got["recipient_records"].([]interface{})

	for i := range recipientRecords {
		assert.Equal(t, float64(recipientsList[len(recipientsList)-1-i].ID), recipientRecords[i].(map[string]interface{})["id"])
		assert.Equal(t, (recipientsList[len(recipientsList)-1-i].Email), recipientRecords[i].(map[string]interface{})["email"])
		assert.Equal(t, recipientsList[len(recipientsList)-1-i].RecipientID, recipientRecords[i].(map[string]interface{})["recipient_id"])
		assert.Equal(t, recipientsList[len(recipientsList)-1-i].PushToken, recipientRecords[i].(map[string]interface{})["push_token"])
	}
}

func TestAddUpdateRecipient(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	file, err := os.Open("recipients.csv")
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	fi, err := file.Stat()
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	boundary := writer.Boundary()
	part, err := writer.CreateFormFile("recipients", fi.Name())
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	part.Write(fileContents)
	err = writer.Close()
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data;boundary=%s", boundary))
	req.Header.Set("Content-Length", "1204")
	req.Header.Set("Content-Disposition", `form-dataname="recipients"; filename="recipients.csv"`)
	c.Request = req
	recipients.AddUpdateRecipient(c)
	assert.Equal(t, http.StatusOK, w.Code)
}
