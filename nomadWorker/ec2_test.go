package nomadWorker

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/aws"
	mock_aws "github.com/thomasobenaus/sokar/test/aws"
)

func Test_CreateSession(t *testing.T) {

	// error, session fun nil
	connector := Connector{}
	sess, err := connector.createSession()
	assert.Error(t, err)
	assert.Nil(t, sess)

	// error, session fun nil
	connector = Connector{awsProfile: "xyz"}
	sess, err = connector.createSession()
	assert.Error(t, err)
	assert.Nil(t, sess)

	// success, profile
	connector = Connector{
		awsProfile:                 "xyz",
		fnCreateSessionFromProfile: aws.NewAWSSessionFromProfile,
	}
	sess, err = connector.createSession()
	assert.NoError(t, err)
	assert.NotNil(t, sess)

	// success, no profile
	connector = Connector{
		awsRegion:       "xyz",
		fnCreateSession: aws.NewAWSSession,
	}
	sess, err = connector.createSession()
	assert.NoError(t, err)
	assert.NotNil(t, sess)
}

func TestAdjustScalingObjectCount_Error(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgFactory := mock_aws.NewMockAutoScalingFactory(mockCtrl)
	asgIF := mock_aws.NewMockAutoScaling(mockCtrl)

	connector, err := New("http://nomad.io", "profile")
	require.NotNil(t, connector)
	require.NoError(t, err)

	connector.autoScalingFactory = asgFactory

	// error, no numbers
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, nil)
	err = connector.AdjustScalingObjectCount("invalid", 2, 10, 4, 5)
	assert.Error(t, err)
}
func TestAdjustScalingObjectCount_Upscale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgFactory := mock_aws.NewMockAutoScalingFactory(mockCtrl)
	asgIF := mock_aws.NewMockAutoScaling(mockCtrl)

	key := "datacenter"
	connector, err := New("http://nomad.io", "profile")
	require.NotNil(t, connector)
	require.NoError(t, err)

	connector.autoScalingFactory = asgFactory

	// no error - UpScale
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	minCount := int64(1)
	desiredCount := int64(123)
	maxCount := int64(3)
	autoScalingGroups := make([]*autoscaling.Group, 0)
	tagVal := "private-services"
	asgName := "my-asg"
	var tags []*autoscaling.TagDescription
	td := autoscaling.TagDescription{Key: &key, Value: &tagVal}
	tags = append(tags, &td)
	asgIn := autoscaling.Group{
		Tags:                 tags,
		AutoScalingGroupName: &asgName,
		MinSize:              &minCount,
		MaxSize:              &maxCount,
		DesiredCapacity:      &desiredCount,
	}
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	output := &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: autoScalingGroups}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)

	asgIF.EXPECT().UpdateAutoScalingGroup(gomock.Any())
	err = connector.AdjustScalingObjectCount(tagVal, 2, 10, 4, 5)
	assert.NoError(t, err)
}

func TestAdjustScalingObjectCount_NoScale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgFactory := mock_aws.NewMockAutoScalingFactory(mockCtrl)

	connector, err := New("http://nomad.io", "profile")
	require.NotNil(t, connector)
	require.NoError(t, err)

	connector.autoScalingFactory = asgFactory

	// no error - DownScale
	tagVal := "private-services"
	err = connector.AdjustScalingObjectCount(tagVal, 2, 10, 4, 4)
	assert.NoError(t, err)
}

func TestGetScalingObjectCount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgFactory := mock_aws.NewMockAutoScalingFactory(mockCtrl)
	asgIF := mock_aws.NewMockAutoScaling(mockCtrl)

	key := "datacenter"
	connector, err := New("http://nomad.io", "profile")
	require.NotNil(t, connector)
	require.NoError(t, err)

	connector.autoScalingFactory = asgFactory

	// error, no numbers
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, nil)
	count, err := connector.GetScalingObjectCount("invalid")
	assert.Error(t, err)
	assert.Equal(t, uint(0), count)

	// no error
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	minCount := int64(1)
	desiredCount := int64(123)
	maxCount := int64(3)
	autoScalingGroups := make([]*autoscaling.Group, 0)
	tagVal := "private-services"
	asgName := "my-asg"
	var tags []*autoscaling.TagDescription
	td := autoscaling.TagDescription{Key: &key, Value: &tagVal}
	tags = append(tags, &td)
	asgIn := autoscaling.Group{
		Tags:                 tags,
		AutoScalingGroupName: &asgName,
		MinSize:              &minCount,
		MaxSize:              &maxCount,
		DesiredCapacity:      &desiredCount,
	}
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	output := &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: autoScalingGroups}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	count, err = connector.GetScalingObjectCount(tagVal)
	assert.NoError(t, err)
	assert.Equal(t, uint(123), count)
}

func Test_IsScalingObjectDead(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgFactory := mock_aws.NewMockAutoScalingFactory(mockCtrl)
	asgIF := mock_aws.NewMockAutoScaling(mockCtrl)

	key := "datacenter"
	connector, err := New("http://nomad.io", "profile")
	require.NotNil(t, connector)
	require.NoError(t, err)

	connector.autoScalingFactory = asgFactory

	// error, no asgs
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, nil)
	dead, err := connector.IsScalingObjectDead("public-services")
	assert.Error(t, err)
	assert.True(t, dead)

	// no error, not found, dead
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	asgOut := make([]*autoscaling.Group, 0)
	group := autoscaling.Group{}
	asgOut = append(asgOut, &group)
	output := &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: asgOut}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	dead, err = connector.IsScalingObjectDead("public-services")
	assert.NoError(t, err)
	assert.True(t, dead)

	// no error, found, NOT dead
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	autoScalingGroups := make([]*autoscaling.Group, 0)
	tagVal := "private-services"
	asgName := "my-asg"
	var tags []*autoscaling.TagDescription
	td := autoscaling.TagDescription{Key: &key, Value: &tagVal}
	tags = append(tags, &td)
	asgIn := autoscaling.Group{Tags: tags, AutoScalingGroupName: &asgName}
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	output = &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: autoScalingGroups}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	dead, err = connector.IsScalingObjectDead("private-services")
	assert.NoError(t, err)
	assert.False(t, dead)
}

func TestAdjustScalingObjectCount_Downscale(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgFactory := mock_aws.NewMockAutoScalingFactory(mockCtrl)

	connector, err := New("http://nomad.io", "profile")
	require.NotNil(t, connector)
	require.NoError(t, err)

	connector.autoScalingFactory = asgFactory

	// no error - DownScale
	tagVal := "private-services"
	err = connector.AdjustScalingObjectCount(tagVal, 2, 10, 5, 4)
	assert.Error(t, err)
}
