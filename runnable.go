package main

import (
	"os"

	"github.com/rs/zerolog"
)

type Runnable interface {
	Run()
	Join()
	Stop()
	GetName() string
}

// Run calls Run() on all Runnables in the list as they are ordered there.
func Run(orderedRunnables []Runnable, logger zerolog.Logger) {
	for _, runnable := range orderedRunnables {
		name := runnable.GetName()
		logger.Debug().Msgf("Starting %s ...", name)
		runnable.Run()
		logger.Info().Msgf("%s running.", name)
	}
}

// Join calls Join() on all Runnables in the list as they are ordered there.
func Join(orderedRunnables []Runnable, logger zerolog.Logger) {
	for _, runnable := range orderedRunnables {
		name := runnable.GetName()
		logger.Debug().Msgf("Join %s ...", name)
		runnable.Join()
		logger.Debug().Msgf("Join %s ... done.", name)
	}
}

// Stop calls Stop() on all Runnables in the list in reverse order.
func Stop(orderedRunnables []Runnable, logger zerolog.Logger) {
	for i := len(orderedRunnables) - 1; i >= 0; i-- {
		runnable := orderedRunnables[i]
		name := runnable.GetName()
		logger.Debug().Msgf("Stopping %s ...", name)
		runnable.Stop()
		logger.Info().Msgf("%s stopped.", name)
	}
}

// shutdownHandler handler that shuts down the running components in case
// a signal was sent on the given channel
func shutdownHandler(shutdownChan <-chan os.Signal, orderedRunnables []Runnable, logger zerolog.Logger) {
	s := <-shutdownChan
	logger.Info().Msgf("Received %v. Shutting down...", s)

	// Stop all components
	Stop(orderedRunnables, logger)
}
