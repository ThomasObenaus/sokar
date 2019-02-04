package logging

import (
	"testing"

	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStack(t *testing.T) {

	stack := NewStack()
	require.NotNil(t, stack)

}
