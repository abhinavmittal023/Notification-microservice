package configuration

import "time"

// Configuration struct is for storing all the config data
type Configuration struct {
	Server            Server            `json:"server"`
	Database          Database          `json:"database"`
	Token             Token             `json:"token"`
	EmailNotification EmailNotification `json:"email_notification"`
	PasswordHash      string            `json:"password_hash"`
}

// Server struct stores the server information
type Server struct {
	Port   string `json:"port"`
	Domain string `json:"domain"`
}

// Token struct stores the jwt configuration
type Token struct {
	SecretKey    string     `json:"secret_key"`
	HeaderPrefix string     `json:"header_prefix"`
	ExpiryTime   ExpiryTime `json:"expiry_time"`
}

// Database struct stores the database info
type Database struct {
	Dbstring     string `json:"dbstring"`
}

// EmailNotification struct stores the validation email id info
type EmailNotification struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	SMTPHost string `json:"smtp_host"`
	SMTPPort string `json:"smtp_port"`
}

// ExpiryTime struct stores the expiry times of different tokens
type ExpiryTime struct {
	ValidationToken time.Duration `json:"validation_token"`
	AccessToken     time.Duration `json:"access_token"`
	RefreshToken    time.Duration `json:"refresh_token"`
}
