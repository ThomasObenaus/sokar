package sokar

import (
	"encoding/json"
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
	scaler          Scaler
	capacityPlanner CapacityPlanner

	logger zerolog.Logger
}

type Config struct {
	Logger zerolog.Logger
}

func (cfg *Config) New(scaler Scaler, capacityPlanner CapacityPlanner) (*Sokar, error) {
	if scaler == nil {
		return nil, fmt.Errorf("Given Scaler is nil")
	}

	if capacityPlanner == nil {
		return nil, fmt.Errorf("Given CapacityPlanner is nil")
	}

	return &Sokar{
		scaler:          scaler,
		capacityPlanner: capacityPlanner,

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
	scaResult := sk.scaler.ScaleBy(int(by))

	code := http.StatusOK
	if scaResult.State == ScaleFailed {
		code = http.StatusInternalServerError
	}
	sk.logger.Info().Msgf("Scale %s: %s", scaResult.State, scaResult.StateDescription)

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(scaResult)
}
