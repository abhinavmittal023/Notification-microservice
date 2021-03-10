package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/controllers/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
)

func TestGetAllChannels(t *testing.T) {
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
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	channels.GetAllChannels(c)

	assert.Equal(t, http.StatusOK, w.Code)
	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	for i := range channelsList {
		assert.Equal(t, float64(channelsList[i].ID), got["channels"].([]interface{})[i].(map[string]interface{})["id"])
		assert.Equal(t, (channelsList[i].Name), got["channels"].([]interface{})[i].(map[string]interface{})["name"])
		assert.Equal(t, float64(channelsList[i].Priority), got["channels"].([]interface{})[i].(map[string]interface{})["priority"])
		assert.Equal(t, float64(channelsList[i].Type), got["channels"].([]interface{})[i].(map[string]interface{})["type"])
	}
}

func TestGetAllChannelsPagination(t *testing.T) {

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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req, err := http.NewRequest("GET", "?offset=0&limit=1", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	channels.GetAllChannels(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["channels"].([]interface{})), 1)
	assert.Equal(t, float64(channelsList[0].ID), got["channels"].([]interface{})[0].(map[string]interface{})["id"])
	assert.Equal(t, (channelsList[0].Name), got["channels"].([]interface{})[0].(map[string]interface{})["name"])
	assert.Equal(t, float64(channelsList[0].Type), got["channels"].([]interface{})[0].(map[string]interface{})["type"])
	assert.Equal(t, float64(channelsList[0].Priority), got["channels"].([]interface{})[0].(map[string]interface{})["priority"])

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	req, err = http.NewRequest("GET", "?offset=1&limit=1", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	channels.GetAllChannels(c)

	assert.Equal(t, http.StatusOK, w.Code)

	got = gin.H{}
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(got["channels"].([]interface{})), 1)
	assert.Equal(t, float64(channelsList[1].ID), got["channels"].([]interface{})[0].(map[string]interface{})["id"])
	assert.Equal(t, (channelsList[1].Name), got["channels"].([]interface{})[0].(map[string]interface{})["name"])
	assert.Equal(t, float64(channelsList[1].Type), got["channels"].([]interface{})[0].(map[string]interface{})["type"])
	assert.Equal(t, float64(channelsList[1].Priority), got["channels"].([]interface{})[0].(map[string]interface{})["priority"])
}

func TestAddChannel(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	data,err := json.Marshal(&serializers.ChannelInfo{
		Name: "email",
		Type: 1,
		Priority: 1,
		Configuration: `{ "email": "no-reply@notificationmicroservice.com" ,"password": "", "smtp_host": "127.0.0.1", "smtp_port": "1025" }`,
	})

	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	channels.AddChannel(c)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder() // For Channel with same provided type
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	channels.AddChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddChannelInvalid(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	// request format doesn't contain 'name' required field
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(`{"type": 1,"priority": 1}`)
	req, err := http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	channels.AddChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// type format incorrect
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	data = []byte(`{"name":"email","type": "1","priority": 1}`)
	req, err = http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	channels.AddChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// type greater than maximum type allowed
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	data = []byte(`{"name":"email","type": 4,"priority": 1}`)
	req, err = http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	channels.AddChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// type equal to 0 not
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	data = []byte(`{"name":"email","type": 0,"priority": 1}`)
	req, err = http.NewRequest("POST", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	channels.AddChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteChannel(t *testing.T) {
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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, err := http.NewRequest("Delete", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", channel.ID),
	})
	channels.DeleteChannel(c)
	assert.Equal(t, http.StatusOK, w.Code)

	// Try deleting already deleted channel

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", channel.ID),
	})
	channels.DeleteChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteChannelInvalid(t *testing.T) {
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

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req, err := http.NewRequest("Delete", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", channel.ID+1),
	})
	channels.DeleteChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Delete channel without id param

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	req, err = http.NewRequest("Delete", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	c.Request = req
	channels.DeleteChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetChannel(t *testing.T) {
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
		Value: fmt.Sprintf("%v", channelsList[0].ID),
	})
	channels.GetChannel(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var got gin.H
	err = json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, float64(channelsList[0].ID), got["channel_details"].(map[string]interface{})["id"])
	assert.Equal(t, channelsList[0].Name, got["channel_details"].(map[string]interface{})["name"])
	assert.Equal(t, float64(channelsList[0].Priority), got["channel_details"].(map[string]interface{})["priority"])
	assert.Equal(t, float64(channelsList[0].Type), got["channel_details"].(map[string]interface{})["type"])

	// Checking with invalid id
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
		Value: fmt.Sprintf("%v", channelsList[2].ID+1),
	})
	channels.GetChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Checking without id as param
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)

	req, err = http.NewRequest("GET", "", nil)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	channels.GetChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateChannel(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	channel := models.Channel{
		Name:     "email",
		Type:     1,
		Configuration: `{ "email": "no-reply@notificationmicroservice.com" ,"password": "", "smtp_host": "127.0.0.1", "smtp_port": "1025" }`,
		Priority: 1,
	}
	if err := SeedOneChannel(&channel); err != nil {
		log.Println(err)
		t.Fail()
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data,err := json.Marshal(&serializers.ChannelInfo{
		Name: "email",
		Type: 1,
		Configuration: `{ "email": "no-reply@notificationmicroservice.com" ,"password": "", "smtp_host": "127.0.0.1", "smtp_port": "1025" }`,
		Priority: 2,
	})
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	req, err := http.NewRequest("Put", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", channel.ID),
	})
	channels.UpdateChannel(c)
	assert.Equal(t, http.StatusOK, w.Code)

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", channel.ID+1),
	})
	channels.UpdateChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateChannelInvalid(t *testing.T) {
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

	// type greater than max value
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	data := []byte(`{"name": "email","type": 4, "priority": 1}`)
	req, err := http.NewRequest("Put", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", channel.ID),
	})
	channels.UpdateChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// priority greater than max value
	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	data = []byte(`{"name": "email","type": 2, "priority": 4}`)
	req, err = http.NewRequest("Put", "", bytes.NewReader(data))
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}
	c.Request = req
	c.Params = append(c.Params, gin.Param{
		Key:   "id",
		Value: fmt.Sprintf("%v", channel.ID),
	})
	channels.UpdateChannel(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
