package nomadWorker

import (
	"fmt"

	aws "github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/thomasobenaus/sokar/helper"
)

type autoScalingGroupQuery struct {

	// asgIF is the interface to aws used to query information about AutoScalingGroup's
	asgIF AutoScaling

	// tagKey is the name of the tag that should be used to find the ASG
	tagKey string
	// tagValue is the value of the tag that should be used to find the ASG
	tagValue string
}

// getScaleNumbers returns the numbers reflecting the scale of the AutoScalingGroup specified by the
// autoScalingGroupQuery.
func (asgQ *autoScalingGroupQuery) getScaleNumbers() (minCount uint, desiredCount uint, maxCount uint, err error) {

	asgs, err := getAutoScalingGroups(asgQ.asgIF)
	if err != nil {
		return 0, 0, 0, err
	}

	asg := filterAutoScalingGroupByTag(asgQ.tagKey, asgQ.tagValue, asgs)
	if asg == nil {
		return 0, 0, 0, fmt.Errorf("No ASG with %s=%s found", asgQ.tagKey, asgQ.tagValue)
	}

	minCount, err = helper.CastInt64ToUint(asg.MinSize)
	if err != nil {
		return 0, 0, 0, err
	}

	desiredCount, err = helper.CastInt64ToUint(asg.DesiredCapacity)
	if err != nil {
		return 0, 0, 0, err
	}

	maxCount, err = helper.CastInt64ToUint(asg.MaxSize)
	if err != nil {
		return 0, 0, 0, err
	}

	return minCount, desiredCount, maxCount, nil
}

// getTagValue returns the value of the TagDescription matching the given key.
// The first matching TagDescription will be taken. In case none of the TagDescriptions
// with the given key matches, an error is returned.
func getTagValue(key string, tags []*aws.TagDescription) (string, error) {

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

// filterAutoScalingGroupByTag filters the given AutoScalingGroups by the given tag (key-value pair).
// In case none of the AutoScalingGroups has the specified tag-key and tag-value nil is returned.
// If multiple AutoScalingGroups match the specified tag-key/ -value, only the first one is returned.
func filterAutoScalingGroupByTag(tagKey string, tagValue string, autoScalingGroups []*aws.Group) *aws.Group {

	for _, asg := range autoScalingGroups {
		if asg == nil {
			continue
		}

		tags := asg.Tags
		tagVal, err := getTagValue(tagKey, tags)

		if err != nil {
			continue
		}

		if tagValue == tagVal {
			return asg
		}
	}

	return nil
}

// getAutoScalingGroups obtains all AutoScalingGroup's
func getAutoScalingGroups(autoScaling AutoScaling) ([]*aws.Group, error) {
	input := aws.DescribeAutoScalingGroupsInput{}
	result, err := autoScaling.DescribeAutoScalingGroups(&input)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, fmt.Errorf("Result is nil")
	}

	return result.AutoScalingGroups, nil
}
