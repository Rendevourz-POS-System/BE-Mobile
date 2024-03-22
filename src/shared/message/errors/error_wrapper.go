package errors

type ErrorWrapper struct {
	ErrorCode int    `json:"ErrorCode,omitempty"`
	Field     string `json:"Field,omitempty"`
	Message   string `json:"Message,omitempty"`
}
