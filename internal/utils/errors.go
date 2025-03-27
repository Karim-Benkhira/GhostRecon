package utils

import "fmt"


type ErrorType int

const (
	DNSError ErrorType = iota
	APIError
	NetworkError
	TimeoutError
	ConfigError
)


type ReconError struct {
	Type    ErrorType
	Message string
	Cause   error
}

func (e *ReconError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}


func NewError(errType ErrorType, message string, cause error) *ReconError {
	return &ReconError{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}
