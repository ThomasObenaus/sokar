package scaleschedule

import (
	"fmt"

	"github.com/thomasobenaus/sokar/helper"
)

// this type implements the needed functions to satisfy the sorting of entries in
// ascending order based on the entries startMinute
type byStartMinute []*entry

type entry struct {
	startMinute uint
	endMinute   uint

	minScale uint
	maxScale uint
}

func (entries byStartMinute) Len() int {
	return len(entries)
}

func (entries byStartMinute) Less(i, j int) bool {
	return entries[i].startMinute < entries[j].startMinute
}

func (entries byStartMinute) Swap(i, j int) {
	entries[i], entries[j] = entries[j], entries[i]
}

func (e entry) String() string {
	start, _ := helper.NewTimeFromMinute(e.startMinute)
	end, _ := helper.NewTimeFromMinute(e.endMinute)
	return fmt.Sprintf("%s-%s->%d-%d", start, end, e.minScale, e.maxScale)
}
