package sokar

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ScaleByPercentage is the end-point for receiving scale-by events. These are events for a relative
// scaling of the job. In this case the scaling is made basend on the given percentage value
func (sk *Sokar) ScaleByPercentage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sk.logger.Info().Msgf("SCALE-BY PERCENTAGE: %s", ps.ByName(PathPartValue))
}

// ScaleByValue is the end-point for receiving scale-by events. These are events for a relative
// scaling of the job. In this case the scaling is made basend on the given value.
func (sk *Sokar) ScaleByValue(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sk.logger.Info().Msgf("SCALE-BY VALUE: %s", ps.ByName(PathPartValue))
}
