package sokar

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	sokarIF "github.com/thomasobenaus/sokar/sokar/iface"
)

const (
	// ParamBy part of the path specifying the amount to be scaled
	ParamBy = "by"
	// PathScaleBy is the url path to the scale-by end-point
	PathScaleBy = "/scaler/scale/:" + ParamBy

	// PathHealth is the url path for health end-point
	PathHealth = "/health"
)

// Sokar component that can be used to scale jobs/instances
type Sokar struct {
	scaler            sokarIF.Scaler
	capacityPlanner   sokarIF.CapacityPlanner
	scaleEventEmitter sokarIF.ScaleEventEmitter

	// channel used to signal teardown/ stop
	stopChan chan struct{}

	logger zerolog.Logger
}

// Config cfg for sokar
type Config struct {
	Logger zerolog.Logger
}

// New creates a new instance of sokar
func (cfg *Config) New(scaleEventEmitter sokarIF.ScaleEventEmitter, capacityPlanner sokarIF.CapacityPlanner, scaler sokarIF.Scaler) (*Sokar, error) {
	if scaler == nil {
		return nil, fmt.Errorf("Given Scaler is nil")
	}

	if capacityPlanner == nil {
		return nil, fmt.Errorf("Given CapacityPlanner is nil")
	}

	if scaleEventEmitter == nil {
		return nil, fmt.Errorf("Given ScaleEventEmitter is nil")
	}

	return &Sokar{
		scaleEventEmitter: scaleEventEmitter,
		capacityPlanner:   capacityPlanner,
		scaler:            scaler,
		stopChan:          make(chan struct{}, 1),

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
	if scaResult.State == sokarIF.ScaleFailed {
		code = http.StatusInternalServerError
	}
	sk.logger.Info().Msgf("Scale %s: %s", scaResult.State, scaResult.StateDescription)

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(scaResult)
}
