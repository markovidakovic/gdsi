package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// common errors
var (
	ErrNotFound        = errors.New("resource not found")
	ErrInvalidInput    = errors.New("invalid input")
	ErrBadRequest      = errors.New("bad request")
	ErrDuplicateRecord = errors.New("resource already exists")
	ErrUnauthorized    = errors.New("access unathorized")
	ErrForbidden       = errors.New("access forbidden")
	ErrInternal        = errors.New("internal server error")
)

// error status code mapping
var failureStatusCodes = map[error]int{
	ErrNotFound:        http.StatusNotFound,
	ErrInvalidInput:    http.StatusBadRequest,
	ErrBadRequest:      http.StatusBadRequest,
	ErrDuplicateRecord: http.StatusConflict,
	ErrUnauthorized:    http.StatusUnauthorized,
	ErrForbidden:       http.StatusForbidden,
	ErrInternal:        http.StatusInternalServerError,
}

// FailWriter interface defines methods that error types must implement
type FailureWriter interface {
	Error() string
	GetStatus() int
	GetMessage() string
}

// Fail represents a basic error response
type Failure struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error method implements the standard go error interface which requires an Error() string method
// if there's a wrapped error (be.Err), it combines both error messages with a colon separator
func (f Failure) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("%s: %v", f.Message, f.Err)
	}
	return f.Message
}

// Unwrap method is part of go's error wrapping mechanism. it allows the use of errors.Is() and errors.As()
// to inspect wrapped errors. when we call errors.Is(err, targetErr), go will repeatedly call Unwrap() to check
// each wrapped error in the chain
func (f Failure) Unwrap() error {
	return f.Err
}

func (f Failure) GetStatus() int {
	return f.Status
}

func (f Failure) GetMessage() string {
	return f.Message
}

// ValidationFailure embeds the Failure and contains additional information about invalid fields
type ValidationFailure struct {
	Failure
	InvalidFields []InvalidField `json:"invalid_fields,omitempty"`
}

type InvalidField struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Location string `json:"location"` // location specifies where the field is comming from (path, query, body)
}

func NewFailure(status int, message string, err error) Failure {
	return Failure{
		Status:  status,
		Message: message,
		Err:     err,
	}
}

func NewNotFoundFailure(message string) Failure {
	return Failure{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func NewBadRequestFailure(message string) Failure {
	return Failure{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func NewUnauthorizedFailure(message string) Failure {
	return Failure{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func NewForbiddenFailure(message string) Failure {
	return Failure{
		Status:  http.StatusForbidden,
		Message: message,
	}
}

func NewInternalFailure(err error) Failure {
	return Failure{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
		Err:     err,
	}
}

func NewValidationFailure(message string, fields []InvalidField) ValidationFailure {
	return ValidationFailure{
		Failure: Failure{
			Status:  http.StatusBadRequest,
			Message: message,
		},
		InvalidFields: fields,
	}
}

// helper funcs
func StatusCodeFromFailure(err error) int {
	if code, exists := failureStatusCodes[err]; exists {
		return code
	}
	return http.StatusInternalServerError
}

func WriteFailure(w http.ResponseWriter, fw FailureWriter) {
	// log server error
	if fw.GetStatus() >= 500 {
		log.Printf("Server error: %v", fw.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(fw.GetStatus())
	if err := json.NewEncoder(w).Encode(fw); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
