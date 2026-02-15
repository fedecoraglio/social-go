package main

import "net/http"

func (app *Application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	posts, err := app.store.Posts.GetUserFeed(ctx, int64(311))

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, posts); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
