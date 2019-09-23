package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_nomadWorker "github.com/thomasobenaus/sokar/test/aws"
)

func Test_CreateAutoScaling(t *testing.T) {

	asgF := AutoScalingFactoryImpl{}

	// nil, no session
	as := asgF.CreateAutoScaling(nil)
	assert.Nil(t, as)

	//  no session
	sess, _ := NewAWSSession("eu-central-1")
	as = asgF.CreateAutoScaling(sess)
	assert.NotNil(t, as)
}

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

func TestFilterAutoScalingGroupByTag(t *testing.T) {

	var autoScalingGroups []*autoscaling.Group

	// none, empty
	asg := FilterAutoScalingGroupByTag("key", "value", autoScalingGroups)
	assert.Nil(t, asg)

	// none, no match
	asgIn := autoscaling.Group{}
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	asg = FilterAutoScalingGroupByTag("key", "value", autoScalingGroups)
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
	asg = FilterAutoScalingGroupByTag(key, tagVal, autoScalingGroups)
	require.NotNil(t, asg)
	assert.Equal(t, asgName, *asg.AutoScalingGroupName)

	// not found, no match
	autoScalingGroups = make([]*autoscaling.Group, 0)
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	asg = FilterAutoScalingGroupByTag(key, "tagVal", autoScalingGroups)
	assert.Nil(t, asg)

	// robust against nil
	autoScalingGroups = make([]*autoscaling.Group, 0)
	autoScalingGroups = append(autoScalingGroups, nil)
	asg = FilterAutoScalingGroupByTag(key, "tagVal", autoScalingGroups)
	assert.Nil(t, asg)
}

func Test_GetAutoScalingGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	// no result, error
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, fmt.Errorf("ERR"))
	asgList, err := GetAutoScalingGroups(asgIF)
	assert.Error(t, err)
	assert.Empty(t, asgList)

	// no result, result is nil
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, nil)
	asgList, err = GetAutoScalingGroups(asgIF)
	assert.Error(t, err)
	assert.Empty(t, asgList)

	// result
	asgOut := make([]*autoscaling.Group, 0)
	group := autoscaling.Group{}
	asgOut = append(asgOut, &group)
	output := &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: asgOut}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	asgList, err = GetAutoScalingGroups(asgIF)
	assert.NoError(t, err)
	assert.NotEmpty(t, asgList)
}

func Test_GetScaleNumbers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	asgQ := AutoScalingGroupQuery{
		AsgIF:    asgIF,
		TagKey:   "key",
		TagValue: "value",
	}

	// no result, error
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, fmt.Errorf("ERR"))
	min, desired, max, err := asgQ.GetScaleNumbers()
	assert.Error(t, err)
	assert.Equal(t, uint(0), min)
	assert.Equal(t, uint(0), desired)
	assert.Equal(t, uint(0), max)

	minCount := int64(1)
	desiredCount := int64(2)
	maxCount := int64(3)
	tagKey := "datacenter"
	tagValue := "private-services"
	tagDesc := autoscaling.TagDescription{Key: &tagKey, Value: &tagValue}
	tags := make([]*autoscaling.TagDescription, 0)
	tags = append(tags, &tagDesc)
	asgOut := make([]*autoscaling.Group, 0)
	group := autoscaling.Group{
		Tags:            tags,
		MinSize:         &minCount,
		MaxSize:         &maxCount,
		DesiredCapacity: &desiredCount,
	}
	asgOut = append(asgOut, &group)
	output := &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: asgOut}

	// no result, tag not match
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	min, desired, max, err = asgQ.GetScaleNumbers()
	assert.Error(t, err)
	assert.Equal(t, uint(0), min)
	assert.Equal(t, uint(0), desired)
	assert.Equal(t, uint(0), max)

	// result
	asgQ = AutoScalingGroupQuery{
		AsgIF:    asgIF,
		TagKey:   "datacenter",
		TagValue: "private-services",
	}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	min, desired, max, err = asgQ.GetScaleNumbers()
	assert.NoError(t, err)
	assert.Equal(t, uint(1), min)
	assert.Equal(t, uint(2), desired)
	assert.Equal(t, uint(3), max)
}

func Test_GetAutoScalingGroupName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	asgQ := AutoScalingGroupQuery{
		AsgIF:    asgIF,
		TagKey:   "key",
		TagValue: "value",
	}

	// no result, error
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, fmt.Errorf("ERR"))
	name, err := asgQ.GetAutoScalingGroupName()
	assert.Error(t, err)
	assert.Empty(t, name)

	asgName := "myASG"
	tagKey := "datacenter"
	tagValue := "private-services"
	tagDesc := autoscaling.TagDescription{Key: &tagKey, Value: &tagValue}
	tags := make([]*autoscaling.TagDescription, 0)
	tags = append(tags, &tagDesc)
	asgOut := make([]*autoscaling.Group, 0)
	group := autoscaling.Group{
		Tags:                 tags,
		AutoScalingGroupName: &asgName,
	}
	asgOut = append(asgOut, &group)
	output := &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: asgOut}

	// no result, tag not match
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	name, err = asgQ.GetAutoScalingGroupName()
	assert.Error(t, err)
	assert.Empty(t, name)

	// result
	asgQ = AutoScalingGroupQuery{
		AsgIF:    asgIF,
		TagKey:   "datacenter",
		TagValue: "private-services",
	}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	name, err = asgQ.GetAutoScalingGroupName()
	assert.NoError(t, err)
	assert.Equal(t, "myASG", name)
}

func Test_TerminateInstanceInAsg(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	instanceID := "1234"
	shouldDecDesiredCapa := true
	input := autoscaling.TerminateInstanceInAutoScalingGroupInput{InstanceId: &instanceID, ShouldDecrementDesiredCapacity: &shouldDecDesiredCapa}

	// invalid request returned
	asgIF.EXPECT().TerminateInstanceInAutoScalingGroupRequest(&input).Return(nil, nil)
	_, _, err := TerminateInstanceInAsg(asgIF, instanceID)
	assert.Error(t, err)

	// invalid output returned
	req := request.Request{}
	asgIF.EXPECT().TerminateInstanceInAutoScalingGroupRequest(&input).Return(&req, nil)
	_, _, err = TerminateInstanceInAsg(asgIF, instanceID)
	assert.Error(t, err)

	// success
	activityID := "ActivityId"
	asgName := "AsgName"
	activity := autoscaling.Activity{ActivityId: &activityID, AutoScalingGroupName: &asgName}
	output := autoscaling.TerminateInstanceInAutoScalingGroupOutput{Activity: &activity}
	asgIF.EXPECT().TerminateInstanceInAutoScalingGroupRequest(&input).Return(&req, &output)
	asgNameResult, activityIDResult, err := TerminateInstanceInAsg(asgIF, instanceID)
	assert.NoError(t, err)
	assert.Equal(t, activityID, activityIDResult)
	assert.Equal(t, asgName, asgNameResult)
}
