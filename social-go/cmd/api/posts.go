package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"social/internal/store"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type postKey = string

const postKeyContext postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"tags"`
}

type UpdatePostPayload struct {
	Title   string `json:"title" validate:"omitempty,max=100"`
	Content string `json:"content" validate:"omitempty,max=1000"`
}

func (app *Application) createPostsHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponseError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponseError(w, r, err)
		return
	}

	// Create a new post for the store layer
	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		UserID:  1,
	}

	ctx := r.Context()

	err := app.store.Posts.Create(ctx, post)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	errWrite := app.jsonResponse(w, http.StatusCreated, post)
	if errWrite != nil {
		app.internalServerError(w, r, errWrite)
		return
	}
}

func (app *Application) getPostsHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)
	if post == nil {
		app.notFoundResponseError(w, r, store.ErrNotFound)
		return
	}

	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonResponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *Application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.store.Posts.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponseError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *Application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromContext(r)

	var payload UpdatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponseError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponseError(w, r, err)
		return
	}

	if payload.Content != "" {
		post.Content = payload.Content
	}
	if payload.Title != "" {
		post.Title = payload.Title
	}

	if err := app.store.Posts.Update(r.Context(), post); err != nil {
		app.badRequestResponseError(w, r, err)
		return
	}

	errWrite := app.jsonResponse(w, http.StatusCreated, post)
	if errWrite != nil {
		app.internalServerError(w, r, errWrite)
		return
	}
}

func (app *Application) postsContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParam := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParam, 10, 64)
		log.Printf("idParam=%q", idParam)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		ctx := r.Context()

		post, err := app.store.Posts.GetByID(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponseError(w, r, err)
				return
			default:
				app.internalServerError(w, r, err)
				return
			}
		}
		ctx = context.WithValue(ctx, postKeyContext, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getPostFromContext(r *http.Request) *store.Post {
	post, ok := r.Context().Value(postKeyContext).(*store.Post)
	if !ok || post == nil {
		return nil
	}
	return post
}
