package response

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/markovidakovic/gdsi/server/failure"
)

type FailureWriter interface {
	Error() string
	GetMessage() string
}

var failureStatusCodes = map[error]int{
	failure.ErrNotFound:     http.StatusNotFound,
	failure.ErrBadRequest:   http.StatusBadRequest,
	failure.ErrDuplicate:    http.StatusConflict,
	failure.ErrCantModify:   http.StatusConflict,
	failure.ErrUnauthorized: http.StatusUnauthorized,
	failure.ErrForbidden:    http.StatusForbidden,
	failure.ErrInternal:     http.StatusInternalServerError,
}

func statusCodeFromFailure(err error) int {
	for k, v := range failureStatusCodes {
		if errors.Is(err, k) {
			return v
		}
	}
	return http.StatusInternalServerError
}

func WriteFailure(w http.ResponseWriter, fw FailureWriter) {
	log.Printf("%v", fw.Error())
	statusCode := statusCodeFromFailure(fw)

	// if statusCode >= 500 {
	// 	log.Printf("%v", fw.Error())
	// }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(fw); err != nil {
		log.Printf("error encoding the response: %v", err)
	}
}
