package custommiddleware

import (
	"fmt"
	"net/http"
)

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("log from authentication middleware")

		next.ServeHTTP(w, r)
	})
}
