package sokar

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_sokar "github.com/thomasobenaus/sokar/test/sokar"
)

func Test_PlanScaleByPercentage(t *testing.T) {

	assert.Equal(t, uint(0), planScaleByPercentage(0.1, 0))
	assert.Equal(t, uint(2), planScaleByPercentage(0.1, 1))
	assert.Equal(t, uint(0), planScaleByPercentage(-0.1, 1))
	assert.Equal(t, uint(110), planScaleByPercentage(0.1, 100))
	assert.Equal(t, uint(90), planScaleByPercentage(-0.1, 100))
	assert.Equal(t, uint(0), planScaleByPercentage(-1, 100))
	assert.Equal(t, uint(200), planScaleByPercentage(1, 100))
	assert.Equal(t, uint(300), planScaleByPercentage(2, 100))
	assert.Equal(t, uint(0), planScaleByPercentage(-1.1, 100))

	assert.Equal(t, uint(33), planScaleByPercentage(10, 3))
}
func Test_PlanScaleByValue(t *testing.T) {

	assert.Equal(t, uint(100), planScaleByValue(0.1, 100))
	assert.Equal(t, uint(101), planScaleByValue(1, 100))
	assert.Equal(t, uint(101), planScaleByValue(0.9, 100))
	assert.Equal(t, uint(200), planScaleByValue(100, 100))
	assert.Equal(t, uint(100), planScaleByValue(-0.1, 100))
	assert.Equal(t, uint(99), planScaleByValue(-1, 100))
	assert.Equal(t, uint(99), planScaleByValue(-0.9, 100))
	assert.Equal(t, uint(0), planScaleByValue(-100, 100))
	assert.Equal(t, uint(0), planScaleByValue(-101, 100))
}

func Test_ScaleByPercentage_HTTPHandler_InvalidParam(t *testing.T) {
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

	req := httptest.NewRequest("PUT", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	// no params -> BadRequest
	sokar.ScaleByPercentage(w, req, httprouter.Params{})
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// invalid param -> BadRequest
	params := []httprouter.Param{httprouter.Param{Key: PathPartValue, Value: "invalid"}}
	sokar.ScaleByPercentage(w, req, params)
	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func Test_ScaleByPercentage_HTTPHandler_OK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	// valid param -> Ok
	currentScale := uint(1)
	scaleTo := uint(2)
	//scaleFactor := float32(0.1)
	params := []httprouter.Param{httprouter.Param{Key: PathPartValue, Value: "10"}}
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		scalerIF.EXPECT().GetTimeOfLastScaleAction().Return(time.Now()),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), false).Return(false),
		scalerIF.EXPECT().ScaleTo(scaleTo, true).Return(nil),
	)
	metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale))
	metricMocks.plannedJobCount.EXPECT().Set(float64(scaleTo))
	sokar.dryRunMode = true
	sokar.ScaleByPercentage(w, req, params)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_ScaleByPercentage_HTTPHandler_IntError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	// valid param -> Ok
	currentScale := uint(1)
	scaleTo := uint(2)
	//scaleFactor := float32(0.1)
	params := []httprouter.Param{httprouter.Param{Key: PathPartValue, Value: "10"}}
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		scalerIF.EXPECT().GetTimeOfLastScaleAction().Return(time.Now()),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), false).Return(false),
		scalerIF.EXPECT().ScaleTo(scaleTo, true).Return(fmt.Errorf("Failed to scale")),
		metricMocks.failedScalingTotal.EXPECT().Inc(),
	)
	metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale))
	metricMocks.plannedJobCount.EXPECT().Set(float64(scaleTo))
	sokar.dryRunMode = true
	sokar.ScaleByPercentage(w, req, params)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func Test_ScaleByValue_HTTPHandler_InvalidParam(t *testing.T) {
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

	req := httptest.NewRequest("PUT", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	// no params -> BadRequest
	sokar.ScaleByValue(w, req, httprouter.Params{})
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// invalid param -> BadRequest
	params := []httprouter.Param{httprouter.Param{Key: PathPartValue, Value: "invalid"}}
	sokar.dryRunMode = true
	sokar.ScaleByValue(w, req, params)
	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func Test_ScaleByValue_HTTPHandler_OK(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	// valid param -> Ok
	currentScale := uint(1)
	scaleTo := uint(11)
	//scaleFactor := float32(0.1)
	params := []httprouter.Param{httprouter.Param{Key: PathPartValue, Value: "10"}}
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		scalerIF.EXPECT().GetTimeOfLastScaleAction().Return(time.Now()),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), false).Return(false),
		scalerIF.EXPECT().ScaleTo(scaleTo, true).Return(nil),
	)
	metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale))
	metricMocks.plannedJobCount.EXPECT().Set(float64(scaleTo))
	sokar.dryRunMode = true
	sokar.ScaleByValue(w, req, params)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func Test_ScaleByValue_HTTPHandler_IntError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	evEmitterIF := mock_sokar.NewMockScaleEventEmitter(mockCtrl)
	scalerIF := mock_sokar.NewMockScaler(mockCtrl)
	capaPlannerIF := mock_sokar.NewMockCapacityPlanner(mockCtrl)
	metrics, metricMocks := NewMockedMetrics(mockCtrl)

	cfg := Config{}
	sokar, err := cfg.New(evEmitterIF, capaPlannerIF, scalerIF, metrics)
	require.NotNil(t, sokar)
	require.NoError(t, err)

	req := httptest.NewRequest("PUT", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	// valid param -> Ok
	currentScale := uint(1)
	scaleTo := uint(11)
	//scaleFactor := float32(0.1)
	params := []httprouter.Param{httprouter.Param{Key: PathPartValue, Value: "10"}}
	gomock.InOrder(
		scalerIF.EXPECT().GetCount().Return(currentScale, nil),
		scalerIF.EXPECT().GetTimeOfLastScaleAction().Return(time.Now()),
		capaPlannerIF.EXPECT().IsCoolingDown(gomock.Any(), false).Return(false),
		scalerIF.EXPECT().ScaleTo(scaleTo, true).Return(fmt.Errorf("Failed to scale")),
		metricMocks.failedScalingTotal.EXPECT().Inc(),
	)
	metricMocks.preScaleJobCount.EXPECT().Set(float64(currentScale))
	metricMocks.plannedJobCount.EXPECT().Set(float64(scaleTo))
	sokar.dryRunMode = true
	sokar.ScaleByValue(w, req, params)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
