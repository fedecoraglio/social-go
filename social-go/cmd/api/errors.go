package main

import (
	"log"
	"net/http"
)

func (app *Application) internalServerError(w http.ResponseWriter, r *http.Request, err error) error {
	log.Printf("Internal server error %s path: %s error: %s", r.Method, r.URL.Path, err)
	return writeJSONError(w, http.StatusInternalServerError, err.Error())
}

func (app *Application) badRequestResponseError(w http.ResponseWriter, r *http.Request, err error) error {
	log.Printf("bad request error %s path: %s error: %s", r.Method, r.URL.Path, err)
	return writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *Application) conflictResponseError(w http.ResponseWriter, r *http.Request, err error) error {
	log.Printf("conflict error %s path: %s error: %s", r.Method, r.URL.Path, err)
	return writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *Application) notFoundResponseError(w http.ResponseWriter, r *http.Request, err error) error {
	log.Printf("not found error %s path: %s error: %s", r.Method, r.URL.Path, err)
	return writeJSONError(w, http.StatusNotFound, err.Error())
}
