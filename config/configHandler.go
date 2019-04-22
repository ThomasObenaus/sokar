package config

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

type ConfigHandler struct {
	Logger zerolog.Logger
	Config Config
}

// ConfigEndpoint represents the config end-point of sokar
func (ch *ConfigHandler) ConfigEndpoint(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := http.StatusOK
	//sk.logger.Info().Str("health", http.StatusText(code)).Msg("Health Check called.")

	w.WriteHeader(code)
	io.WriteString(w, "Sokar is Healthy")
}
