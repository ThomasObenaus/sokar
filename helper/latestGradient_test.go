package helper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateAndGet(t *testing.T) {

	lg := LatestGradient{}
	now := time.Now()
	ago60Secs := now.Add(time.Second * -60)
	ago50Secs := now.Add(time.Second * -50)
	ago40Secs := now.Add(time.Second * -40)
	ago30Secs := now.Add(time.Second * -30)

	assert.Equal(t, float32(0), lg.UpdateAndGet(0, ago60Secs))
	assert.Equal(t, float32(1), lg.UpdateAndGet(10, ago50Secs))
	assert.Equal(t, float32(0), lg.UpdateAndGet(0, ago50Secs))
	assert.Equal(t, float32(0), lg.UpdateAndGet(10, ago60Secs))
	assert.Equal(t, float32(-1), lg.UpdateAndGet(0, ago50Secs))
	assert.Equal(t, float32(-1), lg.UpdateAndGet(-10, ago40Secs))
	assert.Equal(t, float32(1), lg.UpdateAndGet(0, ago30Secs))
}

func Test_Update(t *testing.T) {

	lg := LatestGradient{}
	now := time.Now()
	ago60Secs := now.Add(time.Second * -60)

	lg.Update(0, ago60Secs)
	assert.Equal(t, float32(0), lg.Value)
	assert.Equal(t, ago60Secs, lg.Timestamp)

	lg.Update(10, ago60Secs)
	assert.Equal(t, float32(10), lg.Value)
	assert.Equal(t, ago60Secs, lg.Timestamp)
}

func Test_Get(t *testing.T) {

	lg := LatestGradient{}
	now := time.Now()
	ago60Secs := now.Add(time.Second * -60)
	ago50Secs := now.Add(time.Second * -50)
	ago40Secs := now.Add(time.Second * -40)

	lg.Update(0, ago50Secs)
	assert.Equal(t, float32(1), lg.Get(10, ago40Secs))
	assert.Equal(t, float32(0), lg.Get(0, ago40Secs))
	assert.Equal(t, float32(0), lg.Get(10, ago60Secs))
}
