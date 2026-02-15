package main

import (
	"log"
	"net/http"
	"social/internal/store"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Application struct {
	config Config
	store  store.Storage
}

type Config struct {
	addr    string
	db      dbConfig
	env     string
	version string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *Application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
	})

	r.Route("/v1/posts", func(r chi.Router) {
		r.Post("/", app.createPostsHandler)
		r.Route("/{postID}", func(r chi.Router) {
			r.Use(app.postsContextMiddleware)

			r.Get("/", app.getPostsHandler)
			r.Delete("/", app.deletePostHandler)
			r.Patch("/", app.updatePostHandler)

			// Comments
			r.Route("/comments", func(r chi.Router) {
				r.Post("/", app.createCommentHandler)
			})
		})
	})

	r.Route("/v1/users", func(r chi.Router) {
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(app.userContextMiddleware)

			r.Get("/", app.getUserHandler)
			r.Put("/follow", app.followUserHandler)
			r.Put("/unfollow", app.unfollowUserHandler)
		})
	})

	return r
}

func (app *Application) run(mux http.Handler) error {

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Server has started at %s", app.config.addr)
	return srv.ListenAndServe()
}
