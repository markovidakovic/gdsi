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
var errorStatusCodes = map[error]int{
	ErrNotFound:        http.StatusNotFound,
	ErrInvalidInput:    http.StatusBadRequest,
	ErrBadRequest:      http.StatusBadRequest,
	ErrDuplicateRecord: http.StatusConflict,
	ErrUnauthorized:    http.StatusUnauthorized,
	ErrForbidden:       http.StatusForbidden,
	ErrInternal:        http.StatusInternalServerError,
}

// ErrorWriter interface defines methods that error types must implement
type ErrorWriter interface {
	Error() string
	GetStatus() int
	GetMessage() string
}

// BaseError represents a basic error response
type BaseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error method implements the standard go error interface which requires an Error() string method
// if there's a wrapped error (be.Err), it combines both error messages with a colon separator
func (be BaseError) Error() string {
	if be.Err != nil {
		return fmt.Sprintf("%s: %v", be.Message, be.Err)
	}
	return be.Message
}

// Unwrap method is part of go's error wrapping mechanism. it allows the use of errors.Is() and errors.As()
// to inspect wrapped errors. when we call errors.Is(err, targetErr), go will repeatedly call Unwrap() to check
// each wrapped error in the chain
func (be BaseError) Unwrap() error {
	return be.Err
}

func (be BaseError) GetStatus() int {
	return be.Status
}

func (e BaseError) GetMessage() string {
	return e.Message
}

// ValidationError embeds the BaseError and contains additional information about invalid fields
type ValidationError struct {
	BaseError
	InvalidFields []InvalidField `json:"invalid_fields,omitempty"`
}

type InvalidField struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Location string `json:"location"` // location specifies where the field is comming from (path, query, body)
}

func NewError(status int, message string, err error) BaseError {
	return BaseError{
		Status:  status,
		Message: message,
		Err:     err,
	}
}

func NewNotFoundError(message string) BaseError {
	return BaseError{
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func NewBadRequestError(message string) BaseError {
	return BaseError{
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func NewUnauthorizedError(message string) BaseError {
	return BaseError{
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func NewInternalError(err error) BaseError {
	return BaseError{
		Status:  http.StatusInternalServerError,
		Message: "internal server error",
		Err:     err,
	}
}

func NewValidationError(message string, fields []InvalidField) ValidationError {
	return ValidationError{
		BaseError: BaseError{
			Status:  http.StatusBadRequest,
			Message: message,
		},
		InvalidFields: fields,
	}
}

// helper funcs
func StatusCodeFromError(err error) int {
	if code, exists := errorStatusCodes[err]; exists {
		return code
	}
	return http.StatusInternalServerError
}

func WriteError(w http.ResponseWriter, ew ErrorWriter) {
	// log server error
	if ew.GetStatus() >= 500 {
		log.Printf("Server error: %v", ew.Error())
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ew.GetStatus())
	if err := json.NewEncoder(w).Encode(ew); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}
