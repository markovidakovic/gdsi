package response

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// don't write body for 204
	if status == http.StatusNoContent {
		return
	}
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding success response: %v", err)
	}
}
