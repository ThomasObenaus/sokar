package nomadWorker

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/autoscaling"
)

// getTagValue returns the value of the TagDescription matching the given key.
// The first matching TagDescription will be taken. In case none of the TagDescriptions
// with the given key matches, an error is returned.
func getTagValue(key string, tags []*autoscaling.TagDescription) (string, error) {

	for _, tDesc := range tags {
		if tDesc == nil {
			continue
		}
		if *tDesc.Key == key {
			return *tDesc.Value, nil
		}
	}

	// not found
	return "", fmt.Errorf("Tag with key %s was not found", key)
}
