package middleware

import (
	"net/http"
)

func StripSlashes(fn http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 1 && path[len(path)-1] == '/' {
			path = path[:len(path)-1]
			r.URL.Path = path
		}
		fn.ServeHTTP(w, r)
	})
}
