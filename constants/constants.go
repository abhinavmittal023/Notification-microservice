package constants

//Constants struct is for storing all the constants
type Constants struct {
	Regex Regex `json:"regex"`
}

//Regex struct stores all the required regex expressions
type Regex struct {
	Email    string `json:"email"`
}
