package scaleAlertAggregator

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPool(t *testing.T) {
	scap := NewScaleAlertPool(time.Second * 60)
	assert.NotNil(t, scap.entries)
}

func Test_Cleanup(t *testing.T) {
	scap := NewScaleAlertPool(time.Second * 60)

	now := time.Now()
	expired := now.Add(time.Minute * -1)
	entry1 := ScaleAlertPoolEntry{
		expiresAt: expired,
	}
	scap.entries[1234] = entry1
	scap.cleanup()
	assert.Equal(t, 0, len(scap.entries))

	notExpired := now.Add(time.Minute)
	entry2 := ScaleAlertPoolEntry{
		expiresAt: notExpired,
	}
	scap.entries[5678] = entry2
	scap.cleanup()
	assert.Equal(t, 1, len(scap.entries))
}

func Test_Update(t *testing.T) {
	scap := NewScaleAlertPool(time.Second * 60)

	var scaleAlerts []ScaleAlert

	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert1", Firing: true})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert2", Firing: true})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert3", Firing: true})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert4", Firing: false})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "", Firing: true})

	weightMap := make(ScaleAlertWeightMap, 0)
	weightMap["Alert1"] = 1
	weightMap["Alert2"] = -1
	receiver := "AM"
	scap.update(receiver, scaleAlerts, weightMap)
	keyAlert1 := toID(receiver, "Alert1")
	keyAlert2 := toID(receiver, "Alert2")
	keyAlert3 := toID(receiver, "Alert3")
	keyAlert4 := toID(receiver, "Alert4")

	scap.cleanup()

	assert.Equal(t, 3, len(scap.entries))

	entry, ok := scap.entries[keyAlert1]
	require.True(t, ok)
	assert.Equal(t, float32(1), entry.weight)

	entry, ok = scap.entries[keyAlert2]
	require.True(t, ok)
	assert.Equal(t, float32(-1), entry.weight)

	entry, ok = scap.entries[keyAlert3]
	require.True(t, ok)
	assert.Equal(t, float32(0), entry.weight)

	_, ok = scap.entries[keyAlert4]
	assert.False(t, ok)
}

func Test_Sync(t *testing.T) {
	scap := NewScaleAlertPool(time.Second * 60)

	weightMap := make(ScaleAlertWeightMap, 0)
	weightMap["Alert1"] = 1
	weightMap["Alert2"] = -1
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
			scap.update("alertmanager", scaleAlerts, weightMap)
			if stop {
				break
			}
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		for {
			scap.update("cloudwatch", scaleAlerts, weightMap)
			if stop {
				break
			}
		}
	}()

	time.Sleep(time.Second * 2)
	stop = true
	wg.Wait()

	assert.Len(t, scap.entries, 4)
}

func Test_Iterate(t *testing.T) {
	scap := NewScaleAlertPool(time.Second * 60)

	var scaleAlerts []ScaleAlert

	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert1", Firing: true})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert2", Firing: true})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert3", Firing: false})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "Alert4", Firing: false})
	scaleAlerts = append(scaleAlerts, ScaleAlert{Name: "", Firing: true})

	weightMap := make(ScaleAlertWeightMap, 0)
	weightMap["Alert1"] = 1
	weightMap["Alert2"] = -1
	receiver := "AM"
	scap.update(receiver, scaleAlerts, weightMap)
	keyAlert1 := toID(receiver, "Alert1")
	keyAlert2 := toID(receiver, "Alert2")
	keyAlert3 := toID(receiver, "Alert3")
	keyAlert4 := toID(receiver, "Alert4")

	var keys []uint32
	scap.iterate(func(key uint32, entry ScaleAlertPoolEntry) {
		keys = append(keys, key)
	})

	assert.Equal(t, 2, scap.size())
	assert.Contains(t, keys, keyAlert1)
	assert.Contains(t, keys, keyAlert2)
	assert.NotContains(t, keys, keyAlert3)
	assert.NotContains(t, keys, keyAlert4)
}
