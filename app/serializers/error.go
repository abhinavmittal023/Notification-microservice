package serializers

// ErrorInfo struct is a serializer for the error
type ErrorInfo struct {
	Error map[int][]string `json:"error"`
}
