package failure

import (
	"errors"
	"fmt"
)

// common failures
var (
	ErrNotFound   = errors.New("not found")
	ErrBadRequest = errors.New("bad request")
	ErrDuplicate  = errors.New("duplicate")
	ErrCantModify = errors.New("can't modify")

	ErrUnauthorized = errors.New("not authorized")
	ErrForbidden    = errors.New("forbidden")

	ErrInternal = errors.New("internal error")
)

type Failure struct {
	Message string `json:"message"`
	Err     error  `json:"-"` // for internal logging (the underlying chained error messages)
}

func New(message string, err error) *Failure {
	return &Failure{
		Message: message,
		Err:     err,
	}
}

// Error method implements the standard go error interface which requires an Error() string method
// if there's a wrapped error (f.Err), it combines both error messages with a colon separator
func (f Failure) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("%s: %v", f.Message, f.Err)
	}
	return f.Message
}

func (f Failure) GetMessage() string {
	return f.Message
}

// Unwrap method is part of go's error wrapping mechanism. It allows the use of errors.Is, and errors.As
// to inspect wrapped errors. When we call errors.Is, go will repeatedly call Unwrap to check each
// wrapped error in the chain
func (f Failure) Unwrap() error {
	return f.Err
}

// ValidationFailure embeds the Failure and contains additional
// information about invalid fields
type ValidationFailure struct {
	Failure
	InvalidFields []InvalidField `json:"invalid_fields,omitempty"`
}

type InvalidField struct {
	Field    string `json:"field"`
	Message  string `json:"message"`
	Location string `json:"location"` // location specifies where the field is conmming from (path, query, body)
}

func NewValidation(message string, flds []InvalidField) *ValidationFailure {
	return &ValidationFailure{
		Failure: Failure{
			Message: message,
			Err:     ErrBadRequest,
		},
		InvalidFields: flds,
	}
}
