package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

// API represents the implementation of the HTTP api of sokar
type API struct {
	Router *httprouter.Router

	port     int
	logger   zerolog.Logger
	srv      *http.Server
	stopChan chan struct{}
}

// Option represents an option for the api
type Option func(api *API)

// WithLogger adds a configured Logger to the api
func WithLogger(logger zerolog.Logger) Option {
	return func(api *API) {
		api.logger = logger
	}
}

// New creates a new runnable api server
func New(port int, options ...Option) *API {
	api := API{
		Router:   httprouter.New(),
		port:     port,
		stopChan: make(chan struct{}, 1),
	}

	// apply the options
	for _, opt := range options {
		opt(&api)
	}
	return &api
}

// GetName returns the name of this component
func (api *API) GetName() string {
	return "api"
}

// Stop stops/ tears down the api server
func (api *API) Stop() error {

	// context: wait for 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := api.srv.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Run starts the api server for sokar
func (api *API) Run() {
	api.srv = &http.Server{Addr: ":" + strconv.Itoa(api.port), Handler: api.Router}

	// Run listening for messages in background
	go func() {
		api.logger.Info().Msgf("Start listening at %d.", api.port)
		err := api.srv.ListenAndServe()

		if err != nil && err == http.ErrServerClosed {
			api.logger.Info().Msg("API Srv shut down gracefully")
		} else {
			api.logger.Fatal().Err(err).Msg("Failed serving.")
		}

		// send the stop message
		api.stopChan <- struct{}{}
	}()
}

// Join waits until the api server has been teared down
func (api *API) Join() {
	<-api.stopChan
}
