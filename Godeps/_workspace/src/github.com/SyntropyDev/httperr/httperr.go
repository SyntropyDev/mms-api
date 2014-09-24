package httperr

import (
	"encoding/json"
	"net/http"
)

const (
	internalMessage   = "There was a problem with the system.  If the problem persists contact the administrator."
	jsonEncodingError = `{"statusCode":500,"message":"There was a problem with the system.  If the problem persists contact the administrator.","error":"json response encoding error"}`
)

type Error interface {
	StatusCode() int
	Message() string
	Error() string
}

type apiError struct {
	Code int    `json:"statusCode"`
	M    string `json:"message"`
	Err  string `json:"error"`
}

func New(code int, message string, err error) Error {
	return &apiError{
		Code: code,
		M:    message,
		Err:  err.Error(),
	}
}

func NewInternal(err error) Error {
	if err == nil {
		return nil
	}
	return New(http.StatusInternalServerError, internalMessage, err)
}

func (e *apiError) StatusCode() int {
	return e.Code
}

func (e *apiError) Message() string {
	return e.M
}

func (e *apiError) Error() string {
	return e.Err
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cors(w, r)

	// CORS Preflighted OPTIONS request
	if r.Method == "OPTIONS" {
		return
	}

	err := h(w, r)
	if err != nil {
		writeError(w, r, err)
	}
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	e, ok := err.(Error)
	if !ok {
		e = NewInternal(err)
	}

	w.WriteHeader(e.StatusCode())

	if err := json.NewEncoder(w).Encode(e); err != nil {
		http.Error(w, jsonEncodingError, http.StatusInternalServerError)
	}
}

func cors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,HEAD,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,x-requested-with")
	w.Header().Set("Content-Type", "application/json")
}
