package serializers

// Logs serializer to bind logs data
type Logs struct {
	Level string `json:"level"`
	Msg   string `json:"msg"`
	Time  string `json:"time"`
}

// LogsListResponse serializer for logs list response
type LogsListResponse struct {
	RecordsAffected int64  `json:"records_count"`
	LogInfo         []Logs `json:"logs"`
}
