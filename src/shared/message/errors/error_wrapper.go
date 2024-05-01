package errors

type ErrorWrapper struct {
	ErrorCode int         `json:"ErrorCode,omitempty"`
	Field     string      `json:"Field,omitempty"`
	Error     string      `json:"Error,omitempty"`
	ErrorS    []string    `json:"ErrorS,omitempty"`
	Errors    error       `json:"Errors,omitempty"`
	Data      interface{} `json:"Data,omitempty"`
	Message   string      `json:"Message,omitempty"`
	Messages  []string    `json:"Messages,omitempty"`
}

type SuccessWrapper struct {
	Message string      `json:"Message,omitempty"`
	Data    interface{} `json:"Data,omitempty"`
}
