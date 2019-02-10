package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type Api struct {
	Router *httprouter.Router

	port     int
	logger   zerolog.Logger
	srv      *http.Server
	stopChan chan struct{}
}

func New(port int, logger zerolog.Logger) Api {
	return Api{
		Router:   httprouter.New(),
		port:     port,
		logger:   logger,
		stopChan: make(chan struct{}, 1),
	}
}

func (api *Api) Stop() {

	// context: wait for 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := api.srv.Shutdown(ctx)
	if err != nil {
		panic(err)
	}
}

func (api *Api) Run() {
	api.srv = &http.Server{Addr: ":" + strconv.Itoa(api.port)}

	// Run listening for messages in background
	go func() {
		api.logger.Info().Msgf("Start listening at %d.", api.port)
		err := api.srv.ListenAndServe()

		if err != nil && err == http.ErrServerClosed {
			api.logger.Info().Msg("Api Srv shut down gracefully")
		} else {
			api.logger.Fatal().Err(err).Msg("Failed serving.")
		}

		// send the stop message
		api.stopChan <- struct{}{}
	}()
}

func (api *Api) Join() {
	<-api.stopChan
}
