package nomadWorker

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGetAutoScalingGoupByKey(t *testing.T) {

	var autoScalingGroups []*autoscaling.Group

	// none, empty
	asg := filterAutoScalingGroupByKey("key", "value", autoScalingGroups)
	assert.Nil(t, asg)

	// none, no match
	asgIn := autoscaling.Group{}
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	asg = filterAutoScalingGroupByKey("key", "value", autoScalingGroups)
	assert.Nil(t, asg)

	// found, match
	autoScalingGroups = make([]*autoscaling.Group, 0)
	key := "datacenter"
	tagVal := "private-services"
	asgName := "my-asg"
	var tags []*autoscaling.TagDescription
	td := autoscaling.TagDescription{Key: &key, Value: &tagVal}
	tags = append(tags, &td)
	asgIn = autoscaling.Group{Tags: tags, AutoScalingGroupName: &asgName}
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	asg = filterAutoScalingGroupByKey(key, tagVal, autoScalingGroups)
	require.NotNil(t, asg)
	assert.Equal(t, asgName, *asg.AutoScalingGroupName)

	// not found, no match
	autoScalingGroups = make([]*autoscaling.Group, 0)
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	asg = filterAutoScalingGroupByKey(key, "tagVal", autoScalingGroups)
	assert.Nil(t, asg)

	// robust against nil
	autoScalingGroups = make([]*autoscaling.Group, 0)
	autoScalingGroups = append(autoScalingGroups, nil)
	asg = filterAutoScalingGroupByKey(key, "tagVal", autoScalingGroups)
	assert.Nil(t, asg)
}
