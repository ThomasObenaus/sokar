package config

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ServeConfig(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	cfg := Config{}

	cfgHandler := ConfigHandler{
		Config: cfg,
	}
	require.NotNil(t, cfgHandler)

	req := httptest.NewRequest("PUT", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	cfgHandler.ConfigEndpoint(w, req, httprouter.Params{})
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	require.NotNil(t, resp.Header["Content-Type"])
	assert.Contains(t, resp.Header["Content-Type"], "application/json")
}
