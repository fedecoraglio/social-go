package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type envelope struct {
	Error string `json:"error"`
}

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(response http.ResponseWriter, status int, data any) error {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	return json.NewEncoder(response).Encode(data)
}

func readJSON(response http.ResponseWriter, request *http.Request, data any) error {
	maxBytes := 1_048_578 // 1mb
	request.Body = http.MaxBytesReader(response, request.Body, int64(maxBytes))
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJSONError(response http.ResponseWriter, status int, message string) error {
	return writeJSON(response, status, envelope{Error: message})
}

func (app *Application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}
	return writeJSON(w, status, &envelope{Data: data})
}
