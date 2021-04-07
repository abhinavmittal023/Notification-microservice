package services

import (
	"log"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/gorm"
)

func TestPatchChannel(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	channel := models.Channel{
		Name:     "email",
		Type:     1,
		Priority: 1,
	}

	err := SeedOneChannel(&channel)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	channel.Name = "email_channel"
	channel.Priority = 2

	gotError := channels.PatchChannel(&channel)
	assert.Equal(t, gotError, nil)

	newChannel := models.Channel{}
	err = db.Get().First(&newChannel).Error
	if err != nil {
		log.Println(err)
		t.Fail()
	}

	assert.Equal(t, newChannel.ID, uint(1))
	assert.Equal(t, newChannel.Name, "email_channel")
	assert.Equal(t, newChannel.Priority, 2)
	assert.NotEqual(t, newChannel.CreatedAt, newChannel.UpdatedAt)
}

func TestSoftDeleteChannel(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	channel := models.Channel{
		Name:     "email",
		Type:     1,
		Priority: 1,
	}

	err := SeedOneChannel(&channel)
	if err != nil {
		log.Println(err.Error())
		t.Fail()
	}

	gotError := channels.DeleteChannel(&channel)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, db.Get().First(models.Channel{}, channel.ID).Error, gorm.ErrRecordNotFound)
}
