package main

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/kenzosashida090/metro-math-go/docs"
	"github.com/kenzosashida090/metro-math-go/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	config  config
	store   store.Storage
	env     string
	version string
}

type config struct {
	addr        string
	frontEndUrl string
	db          dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(":3030/swagger/doc.json"),
		))
		r.Route("/songs", func(r chi.Router) {
			r.Get("/", app.GetAllSongsHandler)
			r.Post("/", app.CreateSongHandler)
			r.Route("/{songId}", func(r chi.Router) {
				r.Get("/", app.GetSongByIdHandler)
				r.Delete("/", app.DeleteSongByIdHandler)
				r.Patch("/", app.UpdateSongByIdHandler)
			})
		})
	})
	return r
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = "0.0.1"
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	shutdown := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdown <- srv.Shutdown(ctx)
	}()
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
