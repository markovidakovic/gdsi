package response

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	// ErrBadRequest   = errors.New("bad request")
	ErrDuplicateRecord = errors.New("resource already exists")
	ErrUnauthorized    = errors.New("access unathorized")
	ErrInternal        = errors.New("internal server error")
)

type ErrorWriter interface {
	GetStatus() int
	GetMessage() string
}

type Error struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (er Error) GetStatus() int {
	return er.Status
}

func (er Error) GetMessage() string {
	return er.Message
}

type ValidationError struct {
	Error
	InvalidFields []InvalidField `json:"invalid_fields,omitempty"`
}

func (ver ValidationError) GetStatus() int {
	return ver.Status
}

func (ver ValidationError) GetMessage() string {
	return ver.Message
}

type InvalidField struct {
	Field string `json:"field"`
	Error string `json:"error"`
	// Location  string `json:"location"` // location would specify where is the invalid param located (path, query, body)
}

func WriteError(w http.ResponseWriter, err ErrorWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.GetStatus())
	_ = json.NewEncoder(w).Encode(err)
}
