package benchmark

import (
	"fmt"
	"log"
	"sync"
	"testing"

	apimessage "code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/api_message"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/notifications"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"code.jtg.tools/ayush.singhal/notifications-microservice/shared/hash"
	sendNotification "code.jtg.tools/ayush.singhal/notifications-microservice/shared/notifications"
)

func BenchmarkNotificationsRecipients10000(b *testing.B) {

	recipientList := []models.Recipient{}
	var maxRecipient = 10000
	for i := 1; i <= maxRecipient; i++ {
		recipient := models.Recipient{
			RecipientID: fmt.Sprintf("%v", i),
			Email:       fmt.Sprintf("test%v@test.com", i),
			PushToken:   hash.GenerateSecureToken(10),
			WebToken:    hash.GenerateSecureToken(10),
		}
		recipientList = append(recipientList, recipient)
	}

	notification := models.Notification{
		Title:    "test",
		Priority: 3,
		Body:     "test",
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

	for n := 0; n < b.N; n++ {
		openAPI := apimessage.OpenAPI{
			NotificationStatus:          make(map[string]apimessage.OpenAPIChannel),
			RecipientIDIncorrect:        []string{},
			PreferredChannelTypeDeleted: []string{},
			RepeatedID:                  []string{},
		}
		errorChan := make(chan error)
		stopChan := make(chan bool)
		var wg sync.WaitGroup

		wg.Add(1)
		go notifications.SendToRecipients(channelsList, recipientList, &openAPI, errorChan, stopChan, notification, sendNotification.MockNotification{}, &wg)
		err := <-errorChan
		if err != nil {
			close(stopChan)
			wg.Wait()
			log.Println(err.Error())
			b.Fail()
		}
		wg.Wait()
	}
}

func BenchmarkNotificationsRecipients1000(b *testing.B) {

	recipientList := []models.Recipient{}
	var maxRecipient = 1000
	for i := 1; i <= maxRecipient; i++ {
		recipient := models.Recipient{
			RecipientID: fmt.Sprintf("%v", i),
			Email:       fmt.Sprintf("test%v@test.com", i),
			PushToken:   hash.GenerateSecureToken(10),
			WebToken:    hash.GenerateSecureToken(10),
		}
		recipientList = append(recipientList, recipient)
	}

	notification := models.Notification{
		Title:    "test",
		Priority: 3,
		Body:     "test",
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

	for n := 0; n < b.N; n++ {
		openAPI := apimessage.OpenAPI{
			NotificationStatus:          make(map[string]apimessage.OpenAPIChannel),
			RecipientIDIncorrect:        []string{},
			PreferredChannelTypeDeleted: []string{},
			RepeatedID:                  []string{},
		}
		errorChan := make(chan error)
		stopChan := make(chan bool)
		var wg sync.WaitGroup

		wg.Add(1)
		go notifications.SendToRecipients(channelsList, recipientList, &openAPI, errorChan, stopChan, notification, sendNotification.MockNotification{}, &wg)
		err := <-errorChan
		if err != nil {
			close(stopChan)
			wg.Wait()
			log.Println(err.Error())
			b.Fail()
		}
		wg.Wait()
	}
}

func BenchmarkNotificationsRecipients100(b *testing.B) {

	recipientList := []models.Recipient{}
	var maxRecipient = 100
	for i := 1; i <= maxRecipient; i++ {
		recipient := models.Recipient{
			RecipientID: fmt.Sprintf("%v", i),
			Email:       fmt.Sprintf("test%v@test.com", i),
			PushToken:   hash.GenerateSecureToken(10),
			WebToken:    hash.GenerateSecureToken(10),
		}
		recipientList = append(recipientList, recipient)
	}

	notification := models.Notification{
		Title:    "test",
		Priority: 3,
		Body:     "test",
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

	for n := 0; n < b.N; n++ {
		openAPI := apimessage.OpenAPI{
			NotificationStatus:          make(map[string]apimessage.OpenAPIChannel),
			RecipientIDIncorrect:        []string{},
			PreferredChannelTypeDeleted: []string{},
			RepeatedID:                  []string{},
		}
		errorChan := make(chan error)
		stopChan := make(chan bool)
		var wg sync.WaitGroup

		wg.Add(1)
		go notifications.SendToRecipients(channelsList, recipientList, &openAPI, errorChan, stopChan, notification, sendNotification.MockNotification{}, &wg)
		err := <-errorChan
		if err != nil {
			close(stopChan)
			wg.Wait()
			log.Println(err.Error())
			b.Fail()
		}
		wg.Wait()
	}
}

func BenchmarkNotifications10000(b *testing.B) {

	recipientList := []models.Recipient{}
	var maxRecipient = 10
	var maxNotifications = 10000
	for i := 1; i <= maxRecipient; i++ {
		recipient := models.Recipient{
			RecipientID: fmt.Sprintf("%v", i),
			Email:       fmt.Sprintf("test%v@test.com", i),
			PushToken:   hash.GenerateSecureToken(10),
			WebToken:    hash.GenerateSecureToken(10),
		}
		recipientList = append(recipientList, recipient)
	}

	notification := models.Notification{
		Title:    "test",
		Priority: 3,
		Body:     "test",
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

	for n := 0; n < b.N; n++ {
		for not := 0; not < maxNotifications; not++ {
			openAPI := apimessage.OpenAPI{
				NotificationStatus:          make(map[string]apimessage.OpenAPIChannel),
				RecipientIDIncorrect:        []string{},
				PreferredChannelTypeDeleted: []string{},
				RepeatedID:                  []string{},
			}
			errorChan := make(chan error)
			stopChan := make(chan bool)
			var wg sync.WaitGroup

			wg.Add(1)
			go notifications.SendToRecipients(channelsList, recipientList, &openAPI, errorChan, stopChan, notification, sendNotification.MockNotification{}, &wg)
			err := <-errorChan
			if err != nil {
				close(stopChan)
				wg.Wait()
				log.Println(err.Error())
				b.Fail()
			}
			wg.Wait()
		}
	}
}
