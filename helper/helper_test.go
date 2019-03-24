package helper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SubUint(t *testing.T) {
	r := SubUint(1, 1)
	assert.Equal(t, 0, r)

	r = SubUint(0, 1)
	assert.Equal(t, -1, r)

	r = SubUint(0, maxUint)
	assert.Equal(t, minInt, r)

	r = SubUint(uint(maxInt), 0)
	assert.Equal(t, maxInt, r)
}

func Test_IncUint(t *testing.T) {

	r := IncUint(0, 1)
	assert.Equal(t, uint(1), r)

	r = IncUint(0, -1)
	assert.Equal(t, uint(0), r)

	r = IncUint(2, -1)
	assert.Equal(t, uint(1), r)

	r = IncUint(2, 2)
	assert.Equal(t, uint(4), r)

	r = IncUint(2, maxInt)
	assert.Equal(t, uint(maxInt)+uint(2), r)

	r = IncUint(maxUint, 1)
	assert.Equal(t, maxUint, r)

	r = IncUint(maxUint, -1)
	assert.Equal(t, maxUint-1, r)
}

func Test_Hash(t *testing.T) {

	hashed1 := Hash("HALLO")
	hashed2 := Hash("HALLO")
	assert.Equal(t, hashed1, hashed2)

	hashed1 = Hash("HALLO")
	hashed2 = Hash("HALLO.")
	assert.NotEqual(t, hashed1, hashed2)
}

func willFail(fail bool) (string, error) {
	if fail {
		return "Fail", fmt.Errorf("Failed")
	}
	return "No Fail", nil
}

func Test_Must(t *testing.T) {
	result := Must(willFail(false))

	assert.NotEmpty(t, result)
	assert.Panics(t, func() { Must(willFail(true)) })
}
