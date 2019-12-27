package scaleschedule

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ShouldSortInAscendingOrder(t *testing.T) {

	// GIVEN
	entries := make([]*entry, 0)
	e := &entry{startMinute: 3}
	entries = append(entries, e)
	e = &entry{startMinute: 2}
	entries = append(entries, e)
	e = &entry{startMinute: 1}
	entries = append(entries, e)
	e = &entry{startMinute: 1}
	entries = append(entries, e)

	// WHEN
	sort.Sort(byStartMinute(entries))

	//THEN
	assert.Equal(t, uint(1), entries[0].startMinute)
	assert.Equal(t, uint(1), entries[1].startMinute)
	assert.Equal(t, uint(2), entries[2].startMinute)
	assert.Equal(t, uint(3), entries[3].startMinute)
}
