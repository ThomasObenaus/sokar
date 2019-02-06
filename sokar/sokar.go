package sokar

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type Sokar struct {
	scaler Scaler

	logger zerolog.Logger
}

type Config struct {
	Logger zerolog.Logger
}

func (cfg *Config) New(scaler Scaler) (*Sokar, error) {
	if scaler == nil {
		return nil, fmt.Errorf("Given Scaler is nil")
	}

	return &Sokar{
		scaler: scaler,
		logger: cfg.Logger,
	}, nil
}

func (sk *Sokar) HandleScaler(w http.ResponseWriter, r *http.Request) {

}
