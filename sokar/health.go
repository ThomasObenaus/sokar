package sokar

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Health represents the health end-point of sokar
func (sk *Sokar) Health(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	code := http.StatusOK
	sk.logger.Info().Str("health", http.StatusText(code)).Msg("Health Check called.")

	w.WriteHeader(code)
	io.WriteString(w, "Sokar is Healthy")
}
