package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"code.jtg.tools/ayush.singhal/notifications-microservice/configuration"
	"github.com/pkg/errors"
)

// PushNotification function sends a push notification to the specified deviceToken given the server key and title, body of the notification
func PushNotification(deviceToken string, serverKey string, title string, notificationBody string) (bool, error) {
	url := configuration.GetResp().PushNotification.URL

	var jsonStr = []byte(fmt.Sprintf(`{"notification": {
		"title": "%s", 
		"body": "%s"
		},
		"to" : "%s"}`, title, notificationBody, deviceToken))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", fmt.Sprintf("key=%s", serverKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Declared an empty map interface
	var result map[string]interface{}

	if resp.Status != "200 OK" {
		return false, fmt.Errorf("Non 200 status received")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "Error Reading Response")
	}

	// Unmarshal or Decode the JSON to the interface.
	json.Unmarshal(body, &result)

	if result["success"] != 1.0 {
		return false, nil
	}
	return true, nil
}
