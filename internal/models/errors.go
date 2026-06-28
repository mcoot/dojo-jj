package models

type ErrorCode string

const (
	ErrJJNotOnPath     ErrorCode = "JJ_NOT_ON_PATH"
	ErrNotInJJRepo     ErrorCode = "NOT_IN_JJ_REPO"
	ErrJJGetRootFailed ErrorCode = "JJ_GET_ROOT_FAILED"
)

type DojoError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func NewDojoError(code ErrorCode, message string) *DojoError {
	return &DojoError{
		Code:    code,
		Message: message,
		Cause:   nil,
	}
}

func NewDojoErrorWithCause(code ErrorCode, message string, cause error) *DojoError {
	return &DojoError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

func (e *DojoError) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return e.Message + ": " + e.Cause.Error()
}

func (e *DojoError) Unwrap() error {
	return e.Cause
}
