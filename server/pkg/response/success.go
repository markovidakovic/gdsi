package response

import (
	"encoding/json"
	"net/http"
)

// type DeletedResource struct {
// 	Message string `json:"message"`
// }

func WriteSuccess(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
