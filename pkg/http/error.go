package http

var _ error = &Error{}

type Error struct {
	Cause   error  `json:"-"`
	Message string `json:"message"`
}

func (e *Error) Unwrap() error {
	return e.Cause
}

func (e *Error) Error() string {
	return e.Message
}

func NewError(cause error, message string) *Error {
	return &Error{
		Cause:   cause,
		Message: message,
	}
}
