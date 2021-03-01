package constants

import (
	"log"
	"strings"
)

const (
	// EmailRegex is used to export regular expression for email
	EmailRegex = "^[a-zA-Z0-9_+&*-]+(?:\\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,7}$"

	// Authorization is the header type for authorization token
	Authorization = "Authorization"

	// SystemAdminRole defines the role value for system admin in the database
	SystemAdminRole = 2

	// AdminRole defines the role value for system admin in the database
	AdminRole = 1

	// MaxType Defines the maximum types of notifications supported by the service
	MaxType = 3

	// MaxPriority Defines the maximum priority of channel
	MaxPriority = 3

	// APIKeyLength is the length of API Key
	APIKeyLength = 64

	// Pending for notification status
	Pending = 1

	// Success for notification status
	Success = 2

	// Failure for notification status
	Failure = 3
)

// ChannelType is a function mapping type field of channel to its string counterpart.
// We are using a function as golang doesn't allow complex types to be constants
func ChannelType(index uint) string {
	if index < 1 || index > MaxType {
		return ""
	}
	return []string{"Email", "Push", "Web"}[int(index-1)]
}

// CSVHeaders is a function that returns the headers of our csv file
func CSVHeaders() []string {
	return []string{
		"ID", "Email", "PushToken", "WebToken", "PreferredChannelType",
	}
}

// ChannelTypeToInt converts channel type to its uint counterpart
func ChannelTypeToInt(channel string) uint {
	channel = strings.ToLower(channel)
	log.Println(channel)
	if channel == "email" {
		return 1
	} else if channel == "push" {
		return 2
	} else if channel == "web" {
		return 3
	}
	return 0
}

// PriorityTypeToInt converts priority to respective int
func PriorityTypeToInt(priority string) int {
	priority = strings.ToLower(priority)
	if priority == "high" {
		return 1
	} else if priority == "medium" {
		return 2
	} else if priority == "low" {
		return 3
	}
	return 0
}
