package helper

import (
	"fmt"
	"hash/fnv"
	"log"
	"reflect"

	"github.com/spf13/cast"
)

const (
	MaxUint = ^uint(0)
	MinUint = 0
	MaxInt  = int(MaxUint >> 1)
	MinInt  = -MaxInt - 1
)

// IncUint increments the given value (unsigned) by the given amount (signed)
// and returns the result as unsigned int (safe from under-/overflow)
func IncUint(value uint, by int) uint {

	// upper thresh of uint reached
	// avoid overflow
	if value >= MaxUint && by >= 0 {
		return MaxUint
	}

	// positive by --> just add
	if by >= 0 {
		return value + uint(by)
	}

	// negative by --> convert to uint and substract
	byUint := uint(by * -1)

	// avoid underflow
	if byUint >= value {
		return 0
	}

	return value - byUint
}

// SubUint substracts two uint's and returns the signed difference
// Avoids over-/underflow.
func SubUint(a uint, by uint) int {
	r := float64(a) - float64(by)

	if r >= float64(MaxInt) {
		return MaxInt
	}

	if r <= float64(MinInt) {
		return MinInt
	}

	return int(r)
}

// SubUint2 substracts two uint's and returns the unsigned difference
// Avoids over-/underflow.
func SubUint2(a uint, by uint) uint {

	if a <= by {
		return 0
	}

	return a - by
}

// Hash creates a 32-bit FNV-1a hash out of the given string.
func Hash(s string) uint32 {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		log.Fatalf("Error creating hash: %s", err.Error())
	}

	return h.Sum32()
}

// Must is a helper that checks if a error returned by a function is nil
// in this case Must will end the program with a fatal printing out the error.
// If the error is nil the result of the function will be returned.
func Must(v interface{}, err error) interface{} {
	if err != nil {
		panic(err.Error())
	}
	return v
}

// CastToStringMapSliceE casts the given input to a slice of string:string map's
func CastToStringMapSliceE(iface interface{}) ([]map[string]string, error) {
	if iface == nil {
		return nil, fmt.Errorf("Given input is nil")
	}
	result := make([]map[string]string, 0)

	switch typedIface := iface.(type) {
	case []map[string]string:
		for _, element := range typedIface {
			result = append(result, cast.ToStringMapString(element))
		}
	case []interface{}:
		for _, element := range typedIface {
			result = append(result, cast.ToStringMapString(element))
		}
	default:
		return nil, fmt.Errorf("Invalid type (%v)", reflect.TypeOf(iface))
	}

	return result, nil
}

// CastToStringMapSlice casts the given input to a slice of string:string map's
func CastToStringMapSlice(iface interface{}) []map[string]string {
	result, _ := CastToStringMapSliceE(iface)
	return result
}

// CastInt64ToUint safely casts the given int64 pointer value to uint value.
// If the input is nil then an error will be returned.
// If the input is < 0 then 0 is returned.
func CastInt64ToUint(in *int64) (uint, error) {
	if in == nil {
		return 0, fmt.Errorf("In is nil")
	}

	// avoid underflow
	if *in <= 0 {
		return 0, nil
	}

	return uint(*in), nil
}
