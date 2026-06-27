package models

type ErrorCode string

const (
	ErrJJNotOnPath ErrorCode = "JJ_NOT_ON_PATH"
)

type DojoError struct {
	Code    ErrorCode
	Message string
}

func (e *DojoError) Error() string {
	return e.Message
}
