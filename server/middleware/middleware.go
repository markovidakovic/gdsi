package middleware

import (
	"fmt"
	"net/http"
)

// New will create a new middleware handler from a http.Handler
func New(h http.Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
		})
	}
}

// contextKey is a value for use with context.WithValue. It's used as pointer
// so it fits in an interface{} without allocation. This technique for defining
// context keys was copied from go's 1.17 new use of context in net/http
type contextKey struct {
	name string
}

func (ck *contextKey) String() string {
	return fmt.Sprintf("middleware context value: %q", ck.name)
}
