package configuration

//Configuration struct is for storing all the config data
type Configuration struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
	Token    Token    `json:"token"`
}

//Server struct stores the server information
type Server struct {
	Port string `json:"port"`
}

//Token struct stores the jwt configuration
type Token struct {
	SecretKey string `json:"secret_key"`
}

//Database struct stores the database info
type Database struct {
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
	SSLMode  string `json:"sslmode"`
}
