package sokar

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// ScaleByPercentage is the end-point for receiving scale-by events. These are events for a relative
// scaling of the scaling-object. In this case the scaling is made basend on the given percentage value
func (sk *Sokar) ScaleByPercentage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	percentageStr := ps.ByName(PathPartValue)

	sk.logger.Info().Msgf("ScaleBy Percentage Endpoint with '%s %%' called.", percentageStr)
	percentage, err := strconv.ParseInt(percentageStr, 10, 64)
	if err != nil {
		sk.logger.Error().Err(err).Msg("Percentage parameter is invalid.")
		http.Error(w, fmt.Sprintf("Percentage parameter is invalid: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if !sk.dryRunMode {
		sk.logger.Error().Msg("The scale by endpoint is only supported if sokar is running in dry-run mode.")
		http.Error(w, "The scale by endpoint is only supported if sokar is running in dry-run mode.", http.StatusBadRequest)
		return
	}

	percentageFract := float32(percentage) / 100.00
	// this is used in manual (override) mode --> force has to be true
	err = sk.triggerScale(true, percentageFract, planScaleByPercentage)
	if err != nil {
		sk.logger.Error().Err(err).Msg("Unable to trigger scale")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Scaling triggered")
}

// ScaleByValue is the end-point for receiving scale-by events. These are events for a relative
// scaling of the scaling-object. In this case the scaling is made basend on the given value.
func (sk *Sokar) ScaleByValue(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	valueStr := ps.ByName(PathPartValue)

	sk.logger.Info().Msgf("ScaleBy Value Endpoint with '%s %%' called.", valueStr)
	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		sk.logger.Error().Err(err).Msg("Percentage parameter is invalid.")
		http.Error(w, fmt.Sprintf("Value parameter is invalid: %s", err.Error()), http.StatusBadRequest)
		return
	}

	if !sk.dryRunMode {
		sk.logger.Error().Msg("The scale by endpoint is only supported if sokar is running in dry-run mode.")
		http.Error(w, "The scale by endpoint is only supported if sokar is running in dry-run mode.", http.StatusBadRequest)
		return
	}

	// this is used in manual (override) mode --> force has to be true
	err = sk.triggerScale(true, float32(value), planScaleByValue)
	if err != nil {
		sk.logger.Error().Err(err).Msg("Unable to trigger scale")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Scaling triggered")
}

// planScaleByPercentage plans the new scale based on the current scale and the given percentage.
// The percentage has to be expressed by fractions (100% is 1.0, 10% is 0.1).
// Negative values will plan a down scale, positive ones an upscale.
func planScaleByPercentage(percentOfChange float32, currentScale uint) uint {

	deltaRaw := float64(percentOfChange * float32(currentScale))

	// round up
	var delta float64
	if deltaRaw > 0 {
		delta = math.Ceil(deltaRaw)
	} else {
		delta = math.Floor(deltaRaw)
	}

	var result uint
	result = uint(float64(currentScale) + delta)
	if float64(currentScale) < math.Abs(delta) && delta < 0 {
		result = 0
	}

	return result
}

func planScaleByValue(scaleBy float32, currentScale uint) uint {

	scaleByRounded := math.Round(float64(scaleBy))

	var result uint
	result = uint(float64(currentScale) + scaleByRounded)

	if scaleByRounded < 0 && math.Abs(scaleByRounded) > float64(currentScale) {
		result = 0
	}
	return result
}
