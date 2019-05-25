package nomadWorker

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stretchr/testify/assert"
)

func TestGetTagValue(t *testing.T) {

	// not found, empty
	var tags []*autoscaling.TagDescription
	value, err := getTagValue("key", tags)
	assert.Error(t, err)
	assert.Empty(t, value)

	key := "datacenter"
	tagVal := "private-services"

	// not found, no match
	td := autoscaling.TagDescription{Key: &key, Value: &tagVal}
	tags = append(tags, &td)
	value, err = getTagValue("key", tags)
	assert.Error(t, err)
	assert.Empty(t, value)

	// found, match
	value, err = getTagValue(key, tags)
	assert.NoError(t, err)
	assert.NotEmpty(t, value)
	assert.Equal(t, tagVal, value)

	// found, first match
	key = "name"
	tagVal = "something"
	td = autoscaling.TagDescription{Key: &key, Value: &tagVal}
	tags = append(tags, &td)
	value, err = getTagValue(key, tags)
	assert.NoError(t, err)
	assert.NotEmpty(t, value)
	assert.Equal(t, tagVal, value)

	// robust against nil
	tags = append(tags, nil)
	value, err = getTagValue("key", tags)
	assert.Error(t, err)
	assert.Empty(t, value)
}
