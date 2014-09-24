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
