package constants

import (
	"strings"
)

const (
	// EmailRegex is used to export regular expression for email
	EmailRegex = "^[a-zA-Z0-9_+&*-]+(?:\\.[a-zA-Z0-9_+&*-]+)*@(?:[a-zA-Z0-9-]+\\.)+[a-zA-Z]{2,7}$"

	// HostRegex is used to export regular expression for Host
	HostRegex = "[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}"

	// PortRegex is used to export regular expression for Port
	PortRegex = "[0-9]{1,5}"

	// Authorization is the header type for authorization token
	Authorization = "Authorization"

	// XAPIKey is the header type for API Key based open endpoint
	XAPIKey = "X-API-KEY"

	// SystemAdminRole defines the role value for system admin in the database
	SystemAdminRole = 2

	// AdminRole defines the role value for system admin in the database
	AdminRole = 1

	// MaxPriority Defines the maximum priority of channel
	MaxPriority = 3

	// APIKeyLength is the length of API Key
	APIKeyLength = 64

	// ResetTokenLength is the length of API Key
	ResetTokenLength = 64

	// Pending for notification status
	Pending = 1

	// Success for notification status
	Success = 2

	// Failure for notification status
	Failure = 3
	// Role defines the key for saving the role to the request
	Role = "role"

	// ID defines the key for saving the user_id to the request
	ID = "user_id"

	// MaxEmail specifies maximum length of the email
	MaxEmail = 320

	// MaxPushToken specifies the maximum length of the push token
	MaxPushToken = 255

	// MaxWebToken specifies the maximum length of the web token
	MaxWebToken = 255

	// LoginPath stores the url of the login page of the front-end
	LoginPath = "http://localhost:4200/users/login"

	// ResetPasswordPath stores the url of the reset password page of the front-end
	ResetPasswordPath = "http://localhost:4200/users/password/reset/"
)

// ChannelType is a function mapping type field of channel to its string counterpart.
// We are using a function as golang doesn't allow complex types to be constants
func ChannelType(index uint) string {
	return []string{"Email", "Push", "Web"}[int(index-1)]
}

// CSVHeaders is a function that returns the headers of our csv file
func CSVHeaders() []string {
	return []string{
		"ThirdPartyID", "Email", "PushToken", "WebToken", "PreferredChannelType",
	}
}

// ChannelTypeToInt converts channel type to its uint counterpart
func ChannelTypeToInt(channel string) uint {
	channel = strings.ToLower(channel)
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

// RoleType defines the role values of all possible role types
func RoleType() []int {
	return []int{1, 2}
}

// ChannelIntType defines the type int values of all possible channel types
func ChannelIntType() []int {
	return []int{1, 2, 3}
}

// TokenType defines the various token types
func TokenType() struct {
	Validation string
	Access     string
	Refresh    string
} {
	return struct {
		Validation string
		Access     string
		Refresh    string
	}{
		Validation: "validation",
		Access:     "access",
		Refresh:    "refresh",
	}
}

// Error is the struct used to store the error messages
type Error struct {
	InternalError         string
	UnAuthorized          string
	CredentialsMismatch   string
	InvalidID             string
	IDNotInRecords        string
	EmailAlreadyPresent   string
	RefreshTokenRequired  string
	EmailPasswordRequired string
	EmailNotVerified      string
	FileFormatError       string
	InvalidCSVFile        string
	OldPasswordRequired   string
	OldPasswordIncorrect  string
	InvalidRole           string
	InvalidEmail          string
	InvalidType           string
	InvalidFilter         string
	InvalidPagination     string
	InvalidPriority       string
	ChannelTypePresent    string
	RecipientIDIncorrect  string
	ChannelNamePresent    string
	NextPrevNonBool       string
	NextPrevBothProvided  string
	NoAPIKey              string
	InvalidToken          string
	InvalidHost           string
	InvalidPort           string
	InvalidJSON           string
	InvalidDataType       string
}

// Errors is a function that returns all the error messages
func Errors() Error {
	return Error{
		InternalError:         "Internal Server Error",
		UnAuthorized:          "Unauthorized for the Route",
		CredentialsMismatch:   "Email and Password don't match",
		InvalidID:             "ID should be a unsigned integer",
		IDNotInRecords:        "ID is not in the database",
		EmailAlreadyPresent:   "Email ID is Already Present",
		RefreshTokenRequired:  "Refresh Token is required",
		EmailPasswordRequired: "Email and Password are required",
		EmailNotVerified:      "Email ID is not verified",
		FileFormatError:       "File Format Error",
		InvalidCSVFile:        "Invalid CSV File",
		OldPasswordIncorrect:  "Old Password is Incorrect",
		OldPasswordRequired:   "Old Password is required",
		InvalidRole:           "Invalid Role is provided",
		InvalidEmail:          "Email ID is invalid",
		InvalidType:           "Invalid Type Provided",
		InvalidFilter:         "Invalid Filter Parameters",
		InvalidPagination:     "Invalid Limit and Offset",
		InvalidPriority:       "Invalid Priority Provided",
		ChannelTypePresent:    "Channel with Provided Type already exists",
		RecipientIDIncorrect:  "Recipient ID incorrect",
		ChannelNamePresent:    "Channel with Provided Name already exists",
		NextPrevNonBool:       "Non boolean next or prev provided",
		NextPrevBothProvided:  "Provide either next or previous",
		NoAPIKey:              "No API Key exists",
		InvalidToken:          "Provided token is invalid",
		InvalidHost:           "Host is invalid",
		InvalidPort:           "Port is invalid",
		InvalidJSON:           "Invalid JSON",
		InvalidDataType:       "Invalid Data Type Provided",
	}
}
