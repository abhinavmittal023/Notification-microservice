package services

import (
	"log"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/recipients"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/gorm"
)

func TestGetRecipientWithID(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	recipient := models.Recipient{
		RecipientID: "1",
		Email:       "test@test.com",
		PushToken:   "test1234",
	}

	err := SeedOneRecipient(&recipient)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	gotRecipient, gotError := recipients.GetRecipientWithID(uint64(recipient.ID))

	assert.Equal(t, gotError, nil)
	assert.Equal(t, recipient.ID, gotRecipient.ID)
	assert.Equal(t, (recipient.Email), gotRecipient.Email)
	assert.Equal(t, recipient.RecipientID, gotRecipient.RecipientID)

	_, gotError = recipients.GetRecipientWithID(uint64(recipient.ID) + 1)
	assert.Equal(t, gotError, gorm.ErrRecordNotFound)

}

func TestGetRecipientWithRecipientID(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	recipient := models.Recipient{
		RecipientID: "1",
		Email:       "test@test.com",
		PushToken:   "test1234",
	}

	err := SeedOneRecipient(&recipient)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	gotRecipient, gotError := recipients.GetRecipientWithRecipientID(recipient.RecipientID)

	assert.Equal(t, gotError, nil)
	assert.Equal(t, recipient.ID, gotRecipient.ID)
	assert.Equal(t, (recipient.Email), gotRecipient.Email)
	assert.Equal(t, recipient.RecipientID, gotRecipient.RecipientID)

	_, gotError = recipients.GetRecipientWithRecipientID("2")
	assert.Equal(t, gotError, gorm.ErrRecordNotFound)

}

func TestGetAllRecipients(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	recipientList := []models.Recipient{
		{
			RecipientID: "1",
			Email:       "test1@test.com",
			PushToken:   "test1234",
		},
		{
			RecipientID: "2",
			Email:       "test2@test.com",
			WebToken:    "test1234",
		},
		{
			RecipientID: "3",
			WebToken:    "test5678",
			PushToken:   "test1234",
		},
	}

	err := SeedRecipients(&recipientList)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	pagination := serializers.Pagination{
		Limit:  3,
		Offset: 0,
	}

	filter := filter.Recipient{
		RecipientID:          "",
		Email:                0,
		PushToken:            0,
		WebToken:             0,
		PreferredChannelType: 0,
	}

	gotRecipient, gotError := recipients.GetAllRecipients(&pagination, &filter)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, len(gotRecipient), 3)

	for i := range recipientList {
		assert.Equal(t, recipientList[i].ID, gotRecipient[i].ID)
		assert.Equal(t, (recipientList[i].RecipientID), gotRecipient[i].RecipientID)
		assert.Equal(t, (recipientList[i].Email), gotRecipient[i].Email)
		assert.Equal(t, (recipientList[i].WebToken), gotRecipient[i].WebToken)
		assert.Equal(t, (recipientList[i].PushToken), gotRecipient[i].PushToken)
	}

}

func TestGetAllRecipientsFilter(t *testing.T) {

	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	recipientList := []models.Recipient{
		{
			RecipientID: "1",
			Email:       "test1@test.com",
			PushToken:   "test1234",
		},
		{
			RecipientID: "2",
			Email:       "test2@test.com",
			WebToken:    "test1234",
		},
		{
			RecipientID: "3",
			WebToken:    "test5678",
			PushToken:   "test1234",
		},
	}

	err := SeedRecipients(&recipientList)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	pagination := serializers.Pagination{
		Limit:  1,
		Offset: 1,
	}

	filter := filter.Recipient{
		RecipientID:          "",
		Email:                0,
		PushToken:            0,
		WebToken:             1,
		PreferredChannelType: 0,
	}

	gotRecipient, gotError := recipients.GetAllRecipients(&pagination, &filter)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, len(gotRecipient), 1)

	assert.Equal(t, recipientList[2].ID, gotRecipient[0].ID)
	assert.Equal(t, (recipientList[2].RecipientID), gotRecipient[0].RecipientID)
	assert.Equal(t, (recipientList[2].Email), gotRecipient[0].Email)
	assert.Equal(t, (recipientList[2].WebToken), gotRecipient[0].WebToken)
	assert.Equal(t, (recipientList[2].PushToken), gotRecipient[0].PushToken)

}
