package scaleEventAggregator

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

// ScaleAlertPool is a structure for organizing ScaleAlerts.
// Adding, removing, updating and obtaining them.
// Based on the configured TTL the ScaleAlerts will be removed automatically if they were not updated.
// Access is thread-safe.
type ScaleAlertPool struct {
	// entries is a map that provides fast access to a ScaleAlertPoolEntry using it's name
	// key: ScaleAlertPoolEntry.scaleAlert.Name, value: ScaleAlertPoolEntry
	entries map[string]ScaleAlertPoolEntry
	ttl     time.Duration

	// ensures thread safety for the entries of the pool
	lock sync.RWMutex
}

// ScaleAlertPoolEntry represents a ScaleAlert with an expiration time.
// This is needed in order to be able to clean up alerts which are not updated/ fired for a long time.
type ScaleAlertPoolEntry struct {
	scaleAlert ScaleAlert
	expiresAt  time.Time
}

// NewScaleAlertPool creates a new empty pool
func NewScaleAlertPool() ScaleAlertPool {
	return ScaleAlertPool{
		ttl:     time.Second * 60,
		entries: make(map[string]ScaleAlertPoolEntry),
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
func (sp *ScaleAlertPool) update(scaleAlerts ScaleAlertList) {

	expiresAt := time.Now().Add(sp.ttl)

	sp.lock.Lock()
	defer sp.lock.Unlock()
	for _, alert := range scaleAlerts {
		// ignore alerts without name
		if len(alert.Name) == 0 {
			continue
		}

		// remove resolved alert
		if !alert.Firing {
			delete(sp.entries, alert.Name)
			continue
		}

		// add entry, even override it if it already exists
		// for now there is no information to keep
		sp.entries[alert.Name] = ScaleAlertPoolEntry{expiresAt: expiresAt, scaleAlert: alert}
	}
}

// String returns the content of the pool in a string representation
func (sp *ScaleAlertPool) String() string {
	var buf bytes.Buffer

	sp.lock.RLock()
	defer sp.lock.RUnlock()

	buf.WriteString(fmt.Sprintf("%d entries (ttl: %s)\n", len(sp.entries), sp.ttl))
	for key, entry := range sp.entries {
		buf.WriteString(fmt.Sprintf("\t%s %t %s\n", key, entry.scaleAlert.Firing, entry.expiresAt))
	}
	return buf.String()
}

// iterate ensures thread-safe iteration over the ScalingAlerts inside the pool
func (sp *ScaleAlertPool) iterate(fn func(scaleAlert ScaleAlert)) {
	sp.lock.RLock()
	defer sp.lock.RUnlock()

	for _, entry := range sp.entries {
		fn(entry.scaleAlert)
	}
}

// size returns the number of entries in the pool
func (sp *ScaleAlertPool) size() int {
	sp.lock.RLock()
	defer sp.lock.RUnlock()
	return len(sp.entries)
}
