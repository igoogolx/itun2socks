package routes

var (
	ErrBadRequest   = NewError("Body invalid")
	ErrUnauthorized = NewError("Unauthorized")
)

// HTTPError is custom HTTP error for API
type HTTPError struct {
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewError(msg string) *HTTPError {
	return &HTTPError{Message: msg}
}
