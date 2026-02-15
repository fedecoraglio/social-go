package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"social/internal/store"
	"strconv"
)

func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponseError(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := app.store.Users.GetByID(ctx, userID)

	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.badRequestResponseError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}

}
