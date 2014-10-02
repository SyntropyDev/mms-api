package mware

import (
	"encoding/json"
	"net/http"

	"github.com/SyntropyDev/httperr"
)

func ConstantHandler(src interface{}) httperr.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return json.NewEncoder(w).Encode(src)
	}
}

func ServeFile(name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,x-requested-with")
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, name)
	})
}
