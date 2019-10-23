package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	aws "github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	iface "github.com/thomasobenaus/sokar/aws/iface"
	"github.com/thomasobenaus/sokar/helper"
)

// AutoScalingFactoryImpl implementation of a factory that can create an objects used for auto scaling AWS ASG's
type AutoScalingFactoryImpl struct {
}

// CreateAutoScaling  creates an object used for auto scaling AWS ASG's
func (asf *AutoScalingFactoryImpl) CreateAutoScaling(session *session.Session) iface.AutoScaling {

	if session == nil {
		return nil
	}

	return aws.New(session)
}

// AutoScalingGroupQuery structure to build queries to the AWS ASG API
type AutoScalingGroupQuery struct {

	// AsgIF is the interface to aws used to query information about AutoScalingGroup's
	AsgIF iface.AutoScaling

	// TagKey is the name of the tag that should be used to find the ASG
	TagKey string
	// TagValue is the value of the tag that should be used to find the ASG
	TagValue string
}

// GetAutoScalingGroupName returns the name of the AutoScalingGroup specified by the
// AutoScalingGroupQuery.
func (asgQ *AutoScalingGroupQuery) GetAutoScalingGroupName() (string, error) {
	asgs, err := GetAutoScalingGroups(asgQ.AsgIF)
	if err != nil {
		return "", err
	}

	asg := FilterAutoScalingGroupByTag(asgQ.TagKey, asgQ.TagValue, asgs)
	if asg == nil {
		return "", fmt.Errorf("No ASG with %s=%s found", asgQ.TagKey, asgQ.TagValue)
	}

	return *asg.AutoScalingGroupName, nil
}

// GetScaleNumbers returns the numbers reflecting the scale of the AutoScalingGroup specified by the
// AutoScalingGroupQuery.
func (asgQ *AutoScalingGroupQuery) GetScaleNumbers() (minCount uint, desiredCount uint, maxCount uint, err error) {

	asgs, err := GetAutoScalingGroups(asgQ.AsgIF)
	if err != nil {
		return 0, 0, 0, err
	}

	asg := FilterAutoScalingGroupByTag(asgQ.TagKey, asgQ.TagValue, asgs)
	if asg == nil {
		return 0, 0, 0, fmt.Errorf("No ASG with %s=%s found", asgQ.TagKey, asgQ.TagValue)
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

// FilterAutoScalingGroupByTag filters the given AutoScalingGroups by the given tag (key-value pair).
// In case none of the AutoScalingGroups has the specified tag-key and tag-value nil is returned.
// If multiple AutoScalingGroups match the specified tag-key/ -value, only the first one is returned.
func FilterAutoScalingGroupByTag(TagKey string, TagValue string, autoScalingGroups []*aws.Group) *aws.Group {

	for _, asg := range autoScalingGroups {
		if asg == nil {
			continue
		}

		tags := asg.Tags
		tagVal, err := getTagValue(TagKey, tags)

		if err != nil {
			continue
		}

		if TagValue == tagVal {
			return asg
		}
	}

	return nil
}

// GetAutoScalingGroups obtains all AutoScalingGroup's
func GetAutoScalingGroups(autoScaling iface.AutoScaling) ([]*aws.Group, error) {
	input := aws.DescribeAutoScalingGroupsInput{}
	result, err := autoScaling.DescribeAutoScalingGroups(&input)
	if err != nil {
		return nil, errors.WithMessage(err, "Failure while DescribeAutoScalingGroups call")
	}

	if result == nil {
		return nil, fmt.Errorf("Result of DescribeAutoScalingGroups is nil")
	}

	return result.AutoScalingGroups, nil
}

// TerminateInstanceInAsg removes the specified instance and decrements the desired capacity of the instance accordingly.
func TerminateInstanceInAsg(autoScaling iface.AutoScaling, instanceID string) (asgName string, activityID string, err error) {
	shouldDecDesiredCapa := true

	input := aws.TerminateInstanceInAutoScalingGroupInput{InstanceId: &instanceID, ShouldDecrementDesiredCapacity: &shouldDecDesiredCapa}
	if err := input.Validate(); err != nil {
		return "", "", errors.WithMessage(err, "Validation of TerminateInstanceInAutoScalingGroupInput failed")
	}

	// First create the request
	req, output := autoScaling.TerminateInstanceInAutoScalingGroupRequest(&input)

	if req == nil {
		return "", "", fmt.Errorf("Request from TerminateInstanceInAutoScalingGroupRequest is nil")
	}

	// Now send the request
	if err := req.Send(); err != nil {
		return "", "", errors.WithMessage(err, "Sending TerminateInstanceInAutoScalingGroupRequest failed")
	}

	if output == nil || output.Activity == nil || output.Activity.AutoScalingGroupName == nil || output.Activity.ActivityId == nil {
		return "", "", fmt.Errorf("Output from TerminateInstanceInAutoScalingGroupRequest is not valid")
	}

	return *output.Activity.AutoScalingGroupName, *output.Activity.ActivityId, nil
}
