package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/notifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestNotificationAPI(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	channelsList := []models.Channel{
		{
			Name:     "email",
			Type:     1,
			Priority: 1,
		},
		{
			Name:     "web",
			Type:     2,
			Priority: 2,
		},
		{
			Name:     "push",
			Type:     3,
			Priority: 3,
		},
	}
	err := SeedChannels(&channelsList)
	if err != nil {
		t.Fail()
	}

	recipientsList := []models.Recipient{
		{
			RecipientID: "abcd1",
			Email:       "abcd@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID: "abcd2",
			Email:       "abcd1@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID: "abcd3",
			Email:       "abcd2@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
	}

	err = SeedRecipients(&recipientsList)
	if err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(`{
		"notifications": {
		   "recipients" : ["abcd1","abcd2","abcd3"],
		   "priority": "high",
		   "title": "hello",
		   "body": "world 123456"
	   }
   }`)
	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	notifications.PostSendNotifications(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["email"].(map[string]interface{})["success"].([]interface{})), 3)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["web"].(map[string]interface{})["failure"].([]interface{})), 3)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["push"].(map[string]interface{})["failure"].([]interface{})), 3)

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	channelsList = []models.Channel{
		{
			Name:     "email",
			Type:     1,
			Priority: 1,
		},
		{
			Name:     "web",
			Type:     2,
			Priority: 2,
		},
		{
			Name:     "push",
			Type:     3,
			Priority: 3,
		},
	}
	err = SeedChannels(&channelsList)
	if err != nil {
		t.Fail()
	}

	recipientsList = []models.Recipient{
		{
			RecipientID: "abcd1",
			Email:       "abcd@gmail.com",
			PushToken:   "cocfhBKISHaVzgP1wJP4G5:APA91bHftEr1al_cM0BYTb3V4hWEhX1l3d1hwkilIbAQrD7mVCukXMWORwKvDhqerozC24qkYQF4gTQ4-sP4T4W0CA5zgQbmjdpZ6shVLkw3F9FgYFq1DtUtycjbt5-Dn0D03p3xdfp9",
			WebToken:    "dSgDQp8SpCoNWcqMR0yph4:APA91bEbTpWZgLLTzhL-Bk4jw-emjyGAg_HHkwuNhuqXYniqPmD5vBuk69lk-WCtkeYqm0c0dF3J07XZU2ZEaSbdXyfEadt2afEU9KMBl1Bv0v7dZvueDagUl__UsCK4rocevlCU_A1w",
		},
		{
			RecipientID: "abcd2",
			Email:       "abcd1@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID: "abcd3",
			Email:       "abcd2@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
	}

	err = SeedRecipients(&recipientsList)
	if err != nil {
		t.Fail()
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	data = []byte(`{
		"notifications": {
		   "recipients" : ["abcd1","abcd2","abcd3"],
		   "priority": "high",
		   "title": "hello",
		   "body": "world 123456"
	   }
   }`)
	req, err = http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	notifications.PostSendNotifications(c)

	assert.Equal(t, http.StatusOK, w.Code)
	got = gin.H{}
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["email"].(map[string]interface{})["success"].([]interface{})), 3)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["web"].(map[string]interface{})["success"].([]interface{})), 1)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["web"].(map[string]interface{})["failure"].([]interface{})), 2)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["push"].(map[string]interface{})["success"].([]interface{})), 1)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["push"].(map[string]interface{})["failure"].([]interface{})), 2)
}

func TestNotificationAPIWithIncorrectRecipientID(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	channelsList := []models.Channel{
		{
			Name:     "email",
			Type:     1,
			Priority: 1,
		},
		{
			Name:     "web",
			Type:     2,
			Priority: 2,
		},
		{
			Name:     "push",
			Type:     3,
			Priority: 3,
		},
	}
	err := SeedChannels(&channelsList)
	if err != nil {
		t.Fail()
	}

	recipientsList := []models.Recipient{
		{
			RecipientID: "abcd1",
			Email:       "abcd@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID: "abcd2",
			Email:       "abcd1@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID: "abcd3",
			Email:       "abcd2@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
	}

	err = SeedRecipients(&recipientsList)
	if err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(`{
		"notifications": {
		   "recipients" : ["abcd1","abcd2","abcd3", "abcd4"],
		   "priority": "high",
		   "title": "hello",
		   "body": "world 123456"
	   }
   }`)
	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	notifications.PostSendNotifications(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["email"].(map[string]interface{})["success"].([]interface{})), 3)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["web"].(map[string]interface{})["failure"].([]interface{})), 3)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["push"].(map[string]interface{})["failure"].([]interface{})), 3)
	assert.Equal(t, len(got["recipient_id_incorrect"].([]interface{})), 1)
}

func TestNotificationAPIWithPreferredChannelType(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	channelsList := []models.Channel{
		{
			Name:     "email",
			Type:     1,
			Priority: 1,
		},
		{
			Name:     "web",
			Type:     2,
			Priority: 2,
		},
		{
			Name:     "push",
			Type:     3,
			Priority: 3,
		},
	}
	err := SeedChannels(&channelsList)
	if err != nil {
		t.Fail()
	}

	recipientsList := []models.Recipient{
		{
			RecipientID: "abcd1",
			Email:       "abcd@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID: "abcd2",
			Email:       "abcd1@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID:          "abcd3",
			Email:                "abcd2@gmail.com",
			PushToken:            "abcd",
			WebToken:             "abcd",
			PreferredChannelType: 1,
		},
	}

	err = SeedRecipients(&recipientsList)
	if err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(`{
		"notifications": {
		   "recipients" : ["abcd1","abcd2","abcd3", "abcd4"],
		   "priority": "low",
		   "title": "hello",
		   "body": "world 123456"
	   }
   }`)
	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	notifications.PostSendNotifications(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["email"].(map[string]interface{})["success"].([]interface{})), 1)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["push"].(map[string]interface{})["failure"].([]interface{})), 3)
	assert.Equal(t, len(got["recipient_id_incorrect"].([]interface{})), 1)
}

func TestNotificationAPIWithPreferredChannelDeleted(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	channelsList := []models.Channel{
		{
			Name:     "web",
			Type:     2,
			Priority: 2,
		},
		{
			Name:     "push",
			Type:     3,
			Priority: 3,
		},
	}
	err := SeedChannels(&channelsList)
	if err != nil {
		t.Fail()
	}

	recipientsList := []models.Recipient{
		{
			RecipientID: "abcd1",
			Email:       "abcd@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID: "abcd2",
			Email:       "abcd1@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID:          "abcd3",
			Email:                "abcd2@gmail.com",
			PushToken:            "abcd",
			WebToken:             "abcd",
			PreferredChannelType: 1,
		},
	}

	err = SeedRecipients(&recipientsList)
	if err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(`{
		"notifications": {
		   "recipients" : ["abcd1","abcd2","abcd3", "abcd4"],
		   "priority": "low",
		   "title": "hello",
		   "body": "world 123456"
	   }
   }`)
	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	notifications.PostSendNotifications(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["Preferred_channel_deleted"].([]interface{})), 1)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["push"].(map[string]interface{})["failure"].([]interface{})), 3)
	assert.Equal(t, len(got["recipient_id_incorrect"].([]interface{})), 1)
}

func TestNotificationAPIWithChannelNotExist(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	channelsList := []models.Channel{
		{
			Name:     "email",
			Type:     1,
			Priority: 1,
		},
		{
			Name:     "web",
			Type:     2,
			Priority: 2,
		},
	}
	err := SeedChannels(&channelsList)
	if err != nil {
		t.Fail()
	}

	recipientsList := []models.Recipient{
		{
			RecipientID: "abcd1",
			Email:       "abcd@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID: "abcd2",
			Email:       "abcd1@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID:          "abcd3",
			Email:                "abcd2@gmail.com",
			PushToken:            "abcd",
			WebToken:             "abcd",
			PreferredChannelType: 1,
		},
	}

	err = SeedRecipients(&recipientsList)
	if err != nil {
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(`{
		"notifications": {
		   "recipients" : ["abcd1","abcd2","abcd3", "abcd4"],
		   "priority": "low",
		   "title": "hello",
		   "body": "world 123456"
	   }
   }`)
	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	notifications.PostSendNotifications(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["email"].(map[string]interface{})["success"].([]interface{})), 1)
	assert.Equal(t, len(got["recipient_id_incorrect"].([]interface{})), 1)

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	channelsList = []models.Channel{
		{
			Name:     "email",
			Type:     1,
			Priority: 1,
		},
		{
			Name:     "push",
			Type:     3,
			Priority: 3,
		},
	}
	err = SeedChannels(&channelsList)
	if err != nil {
		t.Fail()
	}

	recipientsList = []models.Recipient{
		{
			RecipientID: "abcd1",
			Email:       "abcd@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID: "abcd2",
			Email:       "abcd1@gmail.com",
			PushToken:   "abcd",
			WebToken:    "abcd",
		},
		{
			RecipientID:          "abcd3",
			Email:                "abcd2@gmail.com",
			PushToken:            "abcd",
			WebToken:             "abcd",
			PreferredChannelType: 1,
		},
	}

	err = SeedRecipients(&recipientsList)
	if err != nil {
		t.Fail()
	}

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	data = []byte(`{
		"notifications": {
		   "recipients" : ["abcd1","abcd2","abcd3", "abcd4"],
		   "priority": "medium",
		   "title": "hello",
		   "body": "world 123456"
	   }
   }`)
	req, err = http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	notifications.PostSendNotifications(c)

	assert.Equal(t, http.StatusOK, w.Code)
	got = gin.H{}
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["email"].(map[string]interface{})["success"].([]interface{})), 1)
	assert.Equal(t, len(got["notification_status"].(map[string]interface{})["push"].(map[string]interface{})["failure"].([]interface{})), 3)
	assert.Equal(t, len(got["recipient_id_incorrect"].([]interface{})), 1)
}
