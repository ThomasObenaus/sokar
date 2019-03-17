package sokar

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/test/sokar"
)

func Test_Health(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, _ := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()
	sokar.Health(w, req, httprouter.Params{})

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
