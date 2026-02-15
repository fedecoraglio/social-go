package main

import (
	"net/http"
	"social/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=500"`
}

func (app *Application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var post = app.getPostFromContext(r)
	var commentPayload CreateCommentPayload

	if err := readJSON(w, r, &commentPayload); err != nil {
		app.badRequestResponseError(w, r, err)
		return
	}

	if validationErr := Validate.Struct(commentPayload); validationErr != nil {
		app.badRequestResponseError(w, r, validationErr)
		return
	}

	var comment = &store.Comment{
		PostID:  post.ID,
		Content: commentPayload.Content,
		UserID:  220,
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
