package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type appError struct {
	Error   error  `json:"error,omitempty"  bson:"error"`
	Message string `json:"message,omitempty"  bson:"message"`
	Code    int    `json:"code,omitempty"  bson:"code"`
}

type errorHandler func(w http.ResponseWriter, r *http.Request) *appError

func (fn errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Printf("%v", e.Error)
		w.WriteHeader(e.Code)
		json.NewEncoder(w).Encode(e)
	}
}
