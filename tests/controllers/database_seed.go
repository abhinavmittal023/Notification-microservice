package controllers

import (
	"log"

	"code.jtg.tools/ayush.singhal/notifications-microservice/db"
	"code.jtg.tools/ayush.singhal/notifications-microservice/db/models"
)

// RefreshAllTables function clears the records in the database
func RefreshAllTables() error {

	dbG := db.Get()

	err := dbG.DropTableIfExists(&models.User{}, &models.Recipient{}, &models.Notification{}, &models.Organisation{}, &models.Channel{}, &models.RecipientNotifications{}).Error
	if err != nil {
		return err
	}

	err = dbG.AutoMigrate(&models.User{}, &models.Recipient{}, &models.Notification{}, &models.Organisation{}, &models.Channel{}, &models.RecipientNotifications{}).Error
	if err != nil {
		return err
	}

	err = dbG.Model(&models.User{}).AddUniqueIndex("email_date", "email", "deleted_at").Error
	if err != nil {
		return err
	}

	err = dbG.Model(&models.Channel{}).AddUniqueIndex("type_date", "type", "deleted_at").Error
	if err != nil {
		return err
	}

	log.Printf("Successfully refreshed table")
	return nil

}

// SeedOneUser function adds one record to the users table
func SeedOneUser(user *models.User) error {
	return db.Get().Model(&models.User{}).Create(user).Error
}

// SeedUsers function adds multiple users to the user table
func SeedUsers(users *[]models.User) error {

	for i := range *users {
		err := SeedOneUser(&((*users)[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

// SeedOneChannel function adds one channel to the channel table
func SeedOneChannel(channel *models.Channel) error {
	return db.Get().Model(&models.Channel{}).Create(channel).Error
}

// SeedChannels function adds multiple records to the channel table
func SeedChannels(channels *[]models.Channel) error {

	for i := range *channels {
		err := SeedOneChannel(&((*channels)[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

// SeedOneRecipient function adds single record to the recipient table
func SeedOneRecipient(recipient *models.Recipient) error {
	return db.Get().Model(&models.Recipient{}).Create(recipient).Error
}

// SeedRecipients function adds multiple record to the recipients table
func SeedRecipients(recipients *[]models.Recipient) error {

	for i := range *recipients {
		err := SeedOneRecipient(&((*recipients)[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

// SeedOneNotification function adds single record to the notifications table
func SeedOneNotification(notification *models.Notification) error {
	return db.Get().Model(&models.Notification{}).Create(notification).Error
}

// SeedNotifications function adds multiple records to the notification table
func SeedNotifications(notifications *[]models.Notification) error {

	for i := range *notifications {
		err := SeedOneNotification(&((*notifications)[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

// SeedOneRecipientNotification function adds single record to the recipient_notification table
func SeedOneRecipientNotification(recipientNotification *models.RecipientNotifications) error {
	return db.Get().Model(&models.RecipientNotifications{}).Create(recipientNotification).Error
}

// SeedRecipientNotifications function adds multiple records to the recipient_notifications table
func SeedRecipientNotifications(recipientNotifications *[]models.RecipientNotifications) error {

	for i := range *recipientNotifications {
		err := SeedOneRecipientNotification(&((*recipientNotifications)[i]))
		if err != nil {
			return err
		}
	}

	return nil
}
