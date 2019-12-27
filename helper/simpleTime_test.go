package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetMinutes(t *testing.T) {
	// GIVEN
	st, err := NewTime(0, 0)

	// WHEN
	minutes := st.Minutes()

	//THEN
	require.NoError(t, err)
	assert.Equal(t, uint(0), minutes)

	// GIVEN
	st, err = NewTime(23, 59)

	// WHEN
	minutes = st.Minutes()

	//THEN
	require.NoError(t, err)
	assert.Equal(t, uint(1439), minutes)

}

func Test_NewShouldSucceed(t *testing.T) {
	// GIVEN

	// WHEN
	st, err := NewTime(0, 0)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, uint(0), st.Hour)
	assert.Equal(t, uint(0), st.Minute)

	// GIVEN

	// WHEN
	st, err = NewTime(23, 59)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, uint(23), st.Hour)
	assert.Equal(t, uint(59), st.Minute)
}

func Test_NewShouldFailOnInvalidValues(t *testing.T) {
	// GIVEN

	// WHEN
	_, err := NewTime(24, 0)

	// THEN
	assert.Error(t, err)

	// GIVEN

	// WHEN
	_, err = NewTime(23, 60)

	// THEN
	assert.Error(t, err)
}
