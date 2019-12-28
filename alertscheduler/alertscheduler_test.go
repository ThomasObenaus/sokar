package alertscheduler

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_alertscheduler "github.com/thomasobenaus/sokar/test/alertscheduler"
)

func Test_NewShouldCreateInstance(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	schedule := mock_alertscheduler.NewMockAlertSchedule(mockCtrl)
	alertscheduler := New(schedule)
	assert.NotNil(t, alertscheduler)
}

func Test_WithLogger(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	schedule := mock_alertscheduler.NewMockAlertSchedule(mockCtrl)

	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel)
	am := New(schedule, WithLogger(logger))
	require.NotNil(t, am)
	assert.Equal(t, zerolog.DebugLevel, logger.GetLevel())
}
