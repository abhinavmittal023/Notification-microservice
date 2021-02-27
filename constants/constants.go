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

	// Role defines the key for saving the role to the request
	Role = "role"

	// ID defines the key for saving the user_id to the request
	ID = "user_id"
)

// Error is the struct used to store the error messages
type Error struct {
	InternalError       string
	UnAuthorized        string
	CredentialsMismatch string
	InvalidID           string
	IDNotInRecords      string
	EmailAlreadyPresent string
}

// Errors is a function that returns all the error messages
func Errors() Error {
	return Error{
		InternalError:       "Internal Server Error",
		UnAuthorized:        "Unauthorized for the Route",
		CredentialsMismatch: "Email and Password don't match",
		InvalidID:           "ID should be a unsigned integer",
		IDNotInRecords:      "ID is not in the database",
		EmailAlreadyPresent: "Email ID is Already Present",
	}
}
