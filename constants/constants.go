package constants

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

	// APIKeyLength is the length of API Key
	APIKeyLength = 64
)

// ChannelType is a function mapping type field of channel to its string counterpart.
// We are using a function as golang doesn't allow complex types to be constants
func ChannelType(index uint) string {
	if index < 1 || index > MaxType {
		return ""
	}
	return []string{"Email", "Push", "Web"}[int(index-1)]
}

// ChannelTypeToInt converts channel type to its uint counterpart
func ChannelTypeToInt(channel string) uint {
	if channel == "Email" {
		return 1
	} else if channel == "Push" {
		return 2
	} else if channel == "Web" {
		return 3
	}
	return 0
}
