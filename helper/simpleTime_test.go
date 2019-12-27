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

func Test_NewFromMinuteShouldSucceed(t *testing.T) {
	// GIVEN

	// WHEN
	time1, err1 := NewTimeFromMinute(0)
	time2, err2 := NewTimeFromMinute(1439)
	time3, err3 := NewTimeFromMinute(510)
	time4, err4 := NewTimeFromMinute(11)
	time5, err5 := NewTimeFromMinute(671)
	time6, err6 := NewTimeFromMinute(60)
	time7, err7 := NewTimeFromMinute(61)
	time8, err8 := NewTimeFromMinute(59)

	// THEN
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NoError(t, err3)
	assert.NoError(t, err4)
	assert.NoError(t, err5)
	assert.NoError(t, err6)
	assert.NoError(t, err7)
	assert.NoError(t, err8)
	assert.Equal(t, uint(0), time1.Hour)
	assert.Equal(t, uint(0), time1.Minute)
	assert.Equal(t, uint(23), time2.Hour)
	assert.Equal(t, uint(59), time2.Minute)
	assert.Equal(t, uint(8), time3.Hour)
	assert.Equal(t, uint(30), time3.Minute)
	assert.Equal(t, uint(0), time4.Hour)
	assert.Equal(t, uint(11), time4.Minute)
	assert.Equal(t, uint(11), time5.Hour)
	assert.Equal(t, uint(11), time5.Minute)
	assert.Equal(t, uint(1), time6.Hour)
	assert.Equal(t, uint(0), time6.Minute)
	assert.Equal(t, uint(1), time7.Hour)
	assert.Equal(t, uint(1), time7.Minute)
	assert.Equal(t, uint(0), time8.Hour)
	assert.Equal(t, uint(59), time8.Minute)
}
