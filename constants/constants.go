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

	// LoginPath stores the url of the login page of the front-end
	LoginPath = "http://localhost:4200/users/login"
)

// ChannelType is an function mapping type field of channel to its string counterpart.
// We are using a function as golang doesn't allow complex types to be constants
func ChannelType(index uint) string {
	return []string{"Email", "Push", "Web"}[int(index-1)]
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
	InternalError             string
	UnAuthorized              string
	CredentialsMismatch       string
	InvalidID                 string
	IDNotInRecords            string
	EmailAlreadyPresent       string
	RefreshTokenRequired      string
	EmailPasswordRequired     string
	EmailNotVerified          string
	EmailPasswordNameRequired string
	FileFormatError           string
	InvalidCSVFile            string
	EmailRoleRequired         string
	NewPasswordRequired       string
	OldPasswordRequired       string
	OldPasswordIncorrect      string
	InvalidRole               string
	InvalidEmail              string
	NameTypePriorityRequired  string
	InvalidType               string
	InvalidFilter             string
	InvalidPagination         string
}

// Errors is a function that returns all the error messages
func Errors() Error {
	return Error{
		InternalError:             "Internal Server Error",
		UnAuthorized:              "Unauthorized for the Route",
		CredentialsMismatch:       "Email and Password don't match",
		InvalidID:                 "ID should be a unsigned integer",
		IDNotInRecords:            "ID is not in the database",
		EmailAlreadyPresent:       "Email ID is Already Present",
		RefreshTokenRequired:      "Refresh Token is required",
		EmailPasswordRequired:     "Email and Password are required",
		EmailNotVerified:          "Email ID is not verified",
		EmailPasswordNameRequired: "Email, Password and First Name are required",
		FileFormatError:           "File Format Error",
		InvalidCSVFile:            "Invalid CSV File",
		EmailRoleRequired:         "Email and Role are required",
		OldPasswordIncorrect:      "Old Password is Incorrect",
		NewPasswordRequired:       "New Password is required",
		OldPasswordRequired:       "Old Password is required",
		InvalidRole:               "Invalid Role is provided",
		InvalidEmail:              "Email ID is invalid",
		NameTypePriorityRequired:  "Name, Type and Priority are required",
		InvalidType:               "Invalid Type Provided",
		InvalidFilter:             "Invalid Filter Parameters",
		InvalidPagination:         "Invalid Limit and Offset",
	}
}
