package helper

import (
	"fmt"
	"log"
)

// PrintCheckPoint simple helper for printing check-points.
func PrintCheckPoint(testCase, message string, a ...interface{}) {
	msg := fmt.Sprintf("[TestCase=%s] %s", testCase, message)

	if len(a) > 0 {
		log.Printf(msg, a)
	} else {
		log.Printf(msg)
	}
}
