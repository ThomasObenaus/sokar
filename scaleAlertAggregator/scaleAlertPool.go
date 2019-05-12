package scaleAlertAggregator

import (
	"bytes"
	"fmt"
	"sync"
	"time"

	"github.com/thomasobenaus/sokar/helper"
)

// ScaleAlertPool is a structure for organizing ScaleAlerts.
// Adding, removing, updating and obtaining them.
// Based on the configured TTL the ScaleAlerts will be removed automatically if they were not updated.
// Access is thread-safe.
type ScaleAlertPool struct {
	// entries is a map that provides fast access to a ScaleAlertPoolEntry using it's id
	// key: ScaleAlertPoolEntry.id, value: ScaleAlertPoolEntry
	entries map[uint32]ScaleAlertPoolEntry
	ttl     time.Duration

	// ensures thread safety for the entries of the pool
	lock sync.RWMutex
}

// ScaleAlertPoolEntry represents a ScaleAlert with an expiration time.
// This is needed in order to be able to clean up alerts which are not updated/ fired for a long time.
type ScaleAlertPoolEntry struct {
	id uint32

	scaleAlert ScaleAlert
	expiresAt  time.Time

	// Name of the receiver of the alert
	// This shows where the alert actually came from.
	receiver string

	// weight the weight of the ScaleAlert according to the configured
	// ScaleAlertWeightMap
	weight float32
}

// NewScaleAlertPool creates a new empty pool
// The parameter alertExpirationTime defines after which time an alert will be pruned if he did not
// get updated again by the ScaleAlertEmitter, assuming that the alert is not relevant any more.
func NewScaleAlertPool(alertExpirationTime time.Duration) ScaleAlertPool {
	return ScaleAlertPool{
		ttl:     alertExpirationTime,
		entries: make(map[uint32]ScaleAlertPoolEntry),
	}
}

// cleanup removes all expired entries from the pool
func (sp *ScaleAlertPool) cleanup() {

	now := time.Now()
	sp.lock.Lock()
	defer sp.lock.Unlock()
	for key, entry := range sp.entries {
		// expired --> remove it
		if now.After(entry.expiresAt) {
			delete(sp.entries, key)
		}
	}
}

// update adds firing and non expired ScaleAlerts to the pool
func (sp *ScaleAlertPool) update(receiver string, scaleAlerts []ScaleAlert, weightMap ScaleAlertWeightMap) {

	expiresAt := time.Now().Add(sp.ttl)

	sp.lock.Lock()
	defer sp.lock.Unlock()
	for _, alert := range scaleAlerts {

		// ignore alerts without name
		if len(alert.Name) == 0 {
			continue
		}

		// generate id
		id := toID(receiver, alert.Name)

		// remove resolved alert
		if !alert.Firing {
			delete(sp.entries, id)
			continue
		}

		weight := getWeight(alert.Name, weightMap)
		// add entry, even override it if it already exists
		// for now there is no information to keep
		sp.entries[id] = ScaleAlertPoolEntry{id: id, expiresAt: expiresAt, scaleAlert: alert, receiver: receiver, weight: weight}
	}
}

// String returns the content of the pool in a string representation
func (sp *ScaleAlertPool) String() string {
	var buf bytes.Buffer

	sp.lock.RLock()
	defer sp.lock.RUnlock()

	buf.WriteString(fmt.Sprintf("%d entries (ttl: %s)\n", len(sp.entries), sp.ttl))
	for key, entry := range sp.entries {
		buf.WriteString(fmt.Sprintf("\t%d %t %s\n", key, entry.scaleAlert.Firing, entry.expiresAt))
	}
	return buf.String()
}

// iterate ensures thread-safe iteration over the ScaleAlertPoolEntry inside the pool
func (sp *ScaleAlertPool) iterate(fn func(key uint32, entry ScaleAlertPoolEntry)) {
	sp.lock.RLock()
	defer sp.lock.RUnlock()

	for key, entry := range sp.entries {
		fn(key, entry)
	}
}

// size returns the number of entries in the pool
func (sp *ScaleAlertPool) size() int {
	sp.lock.RLock()
	defer sp.lock.RUnlock()
	return len(sp.entries)
}

// toID creates a unique id based on the given receiver and alertname
func toID(receiver string, alertName string) uint32 {
	concat := receiver + alertName

	return helper.Hash(concat)
}
