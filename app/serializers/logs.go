package serializers

// Logs serializer to bind logs data
type Logs struct {
	Level	string	`json:"level"`
	Msg		string	`json:"msg"`
	Time	string	`json:"time"`
}