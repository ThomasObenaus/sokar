package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
