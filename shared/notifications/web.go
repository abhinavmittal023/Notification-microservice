package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"code.jtg.tools/ayush.singhal/notifications-microservice/constants"
	"github.com/pkg/errors"
)

// Web struct implements Notifications interface
type Web struct {
	To    string
	Title string
	Body  string
}

// NewNotification creates fills the values in the struct with the provided ones
func (web *Web) NewNotification(to string, title string, body string) {
	web.Body = body
	web.Title = title
	web.To = to
}

// SendNotification function sends a web notification to the specified deviceToken given the server key and title, body of the notification
func (web *Web) SendNotification() error {
	url := configuration.GetResp().WebNotification.URL

	var jsonStr = []byte(fmt.Sprintf(`{"notification": {
		"title": "%s", 
		"body": "%s"
		},
		"to" : "%s"}`, web.Title, web.Body, web.To))

	channel, err := channels.GetChannelWithType(uint(constants.ChannelIntType()[2]))
	if err != nil {
		log.Println(err.Error())
		return errors.Wrap(err, constants.Errors().InternalError)
	}

	var config serializers.WebConfig
	var serverKey string
	err = json.Unmarshal([]byte(channel.Configuration), &config)
	if err != nil || config.ServerKey == ""{
		serverKey = configuration.GetResp().WebNotification.ServerKey
	} else {
		serverKey = config.ServerKey
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", serverKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Declared an empty map interface
	var result map[string]interface{}

	if resp.Status != "200 OK" {
		return errors.New("Non 200 status received")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Error Reading Response")
	}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &result)

	if result["success"] != 1.0 {
		return errors.New("Notification Sending failed")
	}
	return nil
}
