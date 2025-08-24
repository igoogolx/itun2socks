package routes

var (
	ErrBadRequest   = NewError("Body invalid")
	ErrUnauthorized = NewError("Unauthorized")
)

// HTTPError is custom HTTP error for API
type HTTPError struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code,omitempty"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewError(msg string) *HTTPError {
	return &HTTPError{Message: msg}
}

func NewErrorWithCode(msg string, code ErrorCode) *HTTPError {
	return &HTTPError{Message: msg, Code: code}
}

type ErrorCode int

const (
	DefaultErrorCode     ErrorCode = 0
	NotElevatedErrorCode ErrorCode = 10000
)
