package nomadWorker

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/test/nomadWorker"
)

func Test_CreateAutoScaling(t *testing.T) {

	asgF := autoScalingFactoryImpl{}

	// nil, no session
	as := asgF.CreateAutoScaling(nil)
	assert.Nil(t, as)

	//  no session
	sess, _ := newAWSSession("eu-central-1")
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
	asg := filterAutoScalingGroupByTag("key", "value", autoScalingGroups)
	assert.Nil(t, asg)

	// none, no match
	asgIn := autoscaling.Group{}
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	asg = filterAutoScalingGroupByTag("key", "value", autoScalingGroups)
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
	asg = filterAutoScalingGroupByTag(key, tagVal, autoScalingGroups)
	require.NotNil(t, asg)
	assert.Equal(t, asgName, *asg.AutoScalingGroupName)

	// not found, no match
	autoScalingGroups = make([]*autoscaling.Group, 0)
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	asg = filterAutoScalingGroupByTag(key, "tagVal", autoScalingGroups)
	assert.Nil(t, asg)

	// robust against nil
	autoScalingGroups = make([]*autoscaling.Group, 0)
	autoScalingGroups = append(autoScalingGroups, nil)
	asg = filterAutoScalingGroupByTag(key, "tagVal", autoScalingGroups)
	assert.Nil(t, asg)
}

func Test_GetAutoScalingGroups(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	// no result, error
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, fmt.Errorf("ERR"))
	asgList, err := getAutoScalingGroups(asgIF)
	assert.Error(t, err)
	assert.Empty(t, asgList)

	// no result, result is nil
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, nil)
	asgList, err = getAutoScalingGroups(asgIF)
	assert.Error(t, err)
	assert.Empty(t, asgList)

	// result
	asgOut := make([]*autoscaling.Group, 0)
	group := autoscaling.Group{}
	asgOut = append(asgOut, &group)
	output := &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: asgOut}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	asgList, err = getAutoScalingGroups(asgIF)
	assert.NoError(t, err)
	assert.NotEmpty(t, asgList)
}

func Test_GetScaleNumbers(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	asgQ := autoScalingGroupQuery{
		asgIF:    asgIF,
		tagKey:   "key",
		tagValue: "value",
	}

	// no result, error
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, fmt.Errorf("ERR"))
	min, desired, max, err := asgQ.getScaleNumbers()
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
	min, desired, max, err = asgQ.getScaleNumbers()
	assert.Error(t, err)
	assert.Equal(t, uint(0), min)
	assert.Equal(t, uint(0), desired)
	assert.Equal(t, uint(0), max)

	// result
	asgQ = autoScalingGroupQuery{
		asgIF:    asgIF,
		tagKey:   "datacenter",
		tagValue: "private-services",
	}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	min, desired, max, err = asgQ.getScaleNumbers()
	assert.NoError(t, err)
	assert.Equal(t, uint(1), min)
	assert.Equal(t, uint(2), desired)
	assert.Equal(t, uint(3), max)
}

func Test_GetAutoScalingGroupName(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	asgQ := autoScalingGroupQuery{
		asgIF:    asgIF,
		tagKey:   "key",
		tagValue: "value",
	}

	// no result, error
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, fmt.Errorf("ERR"))
	name, err := asgQ.getAutoScalingGroupName()
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
	name, err = asgQ.getAutoScalingGroupName()
	assert.Error(t, err)
	assert.Empty(t, name)

	// result
	asgQ = autoScalingGroupQuery{
		asgIF:    asgIF,
		tagKey:   "datacenter",
		tagValue: "private-services",
	}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	name, err = asgQ.getAutoScalingGroupName()
	assert.NoError(t, err)
	assert.Equal(t, "myASG", name)
}
