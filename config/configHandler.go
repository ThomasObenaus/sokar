package config

import (
	"encoding/json"
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
	ch.Logger.Info().Msg("Config end-point called.")
	code := http.StatusOK

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	if err := enc.Encode(ch.Config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		ch.Logger.Error().Err(err).Msg("Error encoding config.")
		return
	}
}
