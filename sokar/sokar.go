package sokar

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
)

const (
	ParamBy     = "by"
	PathScaleBy = "/scaler/scale/:" + ParamBy
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

// ScaleBy is the http end-point for scale actions
func (sk *Sokar) ScaleBy(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	byStr := ps.ByName(ParamBy)
	sk.logger.Debug().Msgf("%s called with param='%s'.", PathScaleBy, byStr)

	by, err := strconv.ParseInt(byStr, 10, 64)
	if err != nil {
		errMsg := fmt.Sprintf("Unable to parse parameter by=%s: %s.", byStr, err.Error())
		http.Error(w, errMsg, http.StatusBadRequest)
		sk.logger.Error().Msg(errMsg)
		return
	}
	err = sk.scaler.ScaleBy(int(by))

	if err != nil {
		errMsg := fmt.Sprintf("Failed to scale: %s.", err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
		sk.logger.Error().Msg(errMsg)
	}
}
