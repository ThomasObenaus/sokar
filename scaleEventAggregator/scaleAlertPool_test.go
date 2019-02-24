package scaleEventAggregator

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPool(t *testing.T) {
	scap := NewScaleAlertPool()
	assert.NotNil(t, scap.entries)
}

func Test_Cleanup(t *testing.T) {
	scap := NewScaleAlertPool()

	now := time.Now()
	expired := now.Add(time.Minute * -1)
	entry1 := ScaleAlertPoolEntry{
		expiresAt: expired,
	}
	scap.entries["AlertA"] = entry1
	scap.cleanup()
	assert.Equal(t, 0, len(scap.entries))

	notExpired := now.Add(time.Minute)
	entry2 := ScaleAlertPoolEntry{
		expiresAt: notExpired,
	}
	scap.entries["AlertB"] = entry2
	scap.cleanup()
	assert.Equal(t, 1, len(scap.entries))
}

func Test_Update(t *testing.T) {
	scap := NewScaleAlertPool()

	var scaleAlerts []ScaleAlert

	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert1", Firing: true})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert2", Firing: true})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert3", Firing: false})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert4", Firing: false})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "", Firing: true})

	scap.update("AM", scaleAlerts)
	scap.cleanup()

	assert.Equal(t, 2, len(scap.entries))

	_, ok := scap.entries["Alert1"]
	assert.True(t, ok)
	_, ok = scap.entries["Alert2"]
	assert.True(t, ok)
	_, ok = scap.entries["Alert3"]
	assert.False(t, ok)
	_, ok = scap.entries["Alert4"]
	assert.False(t, ok)
}

func Test_Sync(t *testing.T) {
	scap := NewScaleAlertPool()

	var scaleAlerts []ScaleAlert

	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert1", Firing: true})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert2", Firing: true})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert3", Firing: false})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert4", Firing: false})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "", Firing: true})

	var wg sync.WaitGroup
	stop := false
	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			scap.update("alertmanager", scaleAlerts)
			assert.Equal(t, 2, len(scap.entries))

			if stop {
				break
			}
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			scap.cleanup()
			scap.update("cloudwatch", scaleAlerts)
			assert.Equal(t, 2, len(scap.entries))
			if stop {
				break
			}
		}
	}()

	time.Sleep(time.Second * 2)
	stop = true
	wg.Wait()

	assert.Equal(t, 2, len(scap.entries))
}
