package presistence

const (
	timeFormat              = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	UserCollectionName      = "users"
	UserTokenCollectionName = "user_tokens"
)

//type (
//	SuccessResponse struct {
//		Message string      `json:"Message"`
//		Data    interface{} `json:"Data,omitempty"`
//	}
//	ErrorResponse struct {
//		Field   string      `json:"Field"`
//		Error   error `json:"Error"`
//	}
//)
