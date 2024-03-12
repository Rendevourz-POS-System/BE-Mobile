package errors

type ErrorWrapper struct {
	ErrorCode int    `json:"ErrorCode"`
	Field     string `json:"Field"`
	Message   string `json:"Message"`
}
