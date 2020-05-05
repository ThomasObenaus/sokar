package helper

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CastInt64ToUint(t *testing.T) {
	result, err := CastInt64ToUint(nil)
	assert.Error(t, err)
	assert.Equal(t, uint(0), result)

	in := int64(10)
	result, err = CastInt64ToUint(&in)
	assert.NoError(t, err)
	assert.Equal(t, uint(10), result)

	in = int64(-10)
	result, err = CastInt64ToUint(&in)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), result)
}

func Test_SubUint2(t *testing.T) {
	r := SubUint2(1, 1)
	assert.Equal(t, uint(0), r)

	r = SubUint2(0, 1)
	assert.Equal(t, uint(0), r)

	r = SubUint2(0, MaxUint)
	assert.Equal(t, uint(0), r)

	r = SubUint2(MaxUint, 0)
	assert.Equal(t, MaxUint, r)
}
func Test_SubUint(t *testing.T) {
	r := SubUint(1, 1)
	assert.Equal(t, 0, r)

	r = SubUint(0, 1)
	assert.Equal(t, -1, r)

	r = SubUint(0, MaxUint)
	assert.Equal(t, MinInt, r)

	r = SubUint(uint(MaxInt), 0)
	assert.Equal(t, MaxInt, r)
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

	r = IncUint(2, MaxInt)
	assert.Equal(t, uint(MaxInt)+uint(2), r)

	r = IncUint(MaxUint, 1)
	assert.Equal(t, MaxUint, r)

	r = IncUint(MaxUint, -1)
	assert.Equal(t, MaxUint-1, r)
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

func Test_CastToStringMapSlice(t *testing.T) {
	mapStrings, err := CastToStringMapSliceE(nil)
	assert.Empty(t, mapStrings)
	assert.Error(t, err)

	mapStrings, err = CastToStringMapSliceE("invalid input")
	assert.Empty(t, mapStrings)
	assert.Error(t, err)

	m1 := make(map[string]string)
	m1["a"] = "A"
	mapStringList := make([]map[string]string, 0)
	mapStringList = append(mapStringList, m1)
	mapStrings, err = CastToStringMapSliceE(mapStringList)
	assert.NoError(t, err)
	assert.NotEmpty(t, mapStrings)
	assert.Len(t, mapStrings, 1)
	assert.Equal(t, m1, mapStrings[0])
}
