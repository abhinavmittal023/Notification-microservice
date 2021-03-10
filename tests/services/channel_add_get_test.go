package services

import (
	"log"
	"testing"

	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/serializers/filter"
	"code.jtg.tools/ayush.singhal/notifications-microservice/app/services/channels"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
	"github.com/go-playground/assert/v2"
	"github.com/jinzhu/gorm"
)

func TestCreateChannel(t *testing.T) {
	if err := RefreshAllTables(); err != nil {
		t.Fail()
	}

	channel := models.Channel{
		Name:     "email",
		Type:     1,
		Priority: 1,
	}

	gotError := channels.AddChannel(&channel)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, channel.ID, uint(1))
	assert.Equal(t, channel.Name, "email")
	assert.Equal(t, channel.Type, 1)
	assert.Equal(t, channel.Priority, 1)
}

func TestGetChannelWithID(t *testing.T) {

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

	gotChannel, gotError := channels.GetChannelWithID(channel.ID)

	assert.Equal(t, gotError, nil)
	assert.Equal(t, channel.ID, gotChannel.ID)
	assert.Equal(t, channel.Name, gotChannel.Name)
	assert.Equal(t, channel.Type, gotChannel.Type)
	assert.Equal(t, channel.Priority, gotChannel.Priority)

	_, gotError = channels.GetChannelWithID(channel.ID + 1)
	assert.Equal(t, gotError, gorm.ErrRecordNotFound)
}

func TestGetChannelWithType(t *testing.T) {

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

	gotChannel, gotError := channels.GetChannelWithType(uint(channelsList[0].Type))

	assert.Equal(t, gotError, nil)
	assert.Equal(t, channelsList[0].ID, gotChannel.ID)
	assert.Equal(t, channelsList[0].Name, gotChannel.Name)
	assert.Equal(t, channelsList[0].Type, gotChannel.Type)
	assert.Equal(t, channelsList[0].Priority, gotChannel.Priority)

	_, gotError = channels.GetChannelWithType(uint(channelsList[2].Type) + 1)
	assert.Equal(t, gotError, gorm.ErrRecordNotFound)
}

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

	pagination := serializers.Pagination{
		Limit:  3,
		Offset: 0,
	}

	filter := filter.Channel{
		Name:     "",
		Type:     0,
		Priority: 0,
	}

	gotChannel, gotError := channels.GetAllChannels(&pagination, &filter)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, len(gotChannel), 3)

	for i := range channelsList {
		assert.Equal(t, channelsList[i].ID, gotChannel[i].ID)
		assert.Equal(t, channelsList[i].Name, gotChannel[i].Name)
		assert.Equal(t, (channelsList[i].Type), gotChannel[i].Type)
		assert.Equal(t, (channelsList[i].Priority), gotChannel[i].Priority)
	}
}

func TestGetAllChannelsFilters(t *testing.T) {

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

	pagination := serializers.Pagination{
		Limit:  3,
		Offset: 0,
	}

	filter := filter.Channel{
		Name:     "",
		Type:     1,
		Priority: 0,
	}

	gotChannel, gotError := channels.GetAllChannels(&pagination, &filter)
	assert.Equal(t, gotError, nil)
	assert.Equal(t, len(gotChannel), 1)

	assert.Equal(t, channelsList[0].ID, gotChannel[0].ID)
	assert.Equal(t, channelsList[0].Name, gotChannel[0].Name)
	assert.Equal(t, channelsList[0].Type, gotChannel[0].Type)
	assert.Equal(t, channelsList[0].Priority, gotChannel[0].Priority)
}
