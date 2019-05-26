package nomadWorker

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thomasobenaus/sokar/test/nomadWorker"
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
		fnCreateSessionFromProfile: newAWSSessionFromProfile,
	}
	sess, err = connector.createSession()
	assert.NoError(t, err)
	assert.NotNil(t, sess)

	// success, no profile
	connector = Connector{
		fnCreateSession: newAWSSession,
	}
	sess, err = connector.createSession()
	assert.NoError(t, err)
	assert.NotNil(t, sess)
}

func TestSetJobCount(t *testing.T) {

	cfg := Config{}
	connector, err := cfg.New(0)

	require.NotNil(t, connector)
	require.NoError(t, err)

	err = connector.SetJobCount("public-services", 10)
	assert.NoError(t, err)
}

func TestGetJobCount(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgFactory := mock_nomadWorker.NewMockAutoScalingFactory(mockCtrl)
	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	cfg := Config{}
	connector, err := cfg.New(0)
	require.NotNil(t, connector)
	require.NoError(t, err)

	connector.autoScalingFactory = asgFactory

	// error, no numbers
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, nil)
	count, err := connector.GetJobCount("invalid")
	assert.Error(t, err)
	assert.Equal(t, uint(0), count)

	// no error
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	minCount := int64(1)
	desiredCount := int64(123)
	maxCount := int64(3)
	autoScalingGroups := make([]*autoscaling.Group, 0)
	key := "datacenter"
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
	count, err = connector.GetJobCount(tagVal)
	assert.NoError(t, err)
	assert.Equal(t, uint(123), count)
}

func Test_IsJobDead(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgFactory := mock_nomadWorker.NewMockAutoScalingFactory(mockCtrl)
	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	cfg := Config{}
	connector, err := cfg.New(0)
	require.NotNil(t, connector)
	require.NoError(t, err)

	// no error, dead
	dead, err := connector.IsJobDead("invalid")
	assert.NoError(t, err)
	assert.True(t, dead)

	connector.autoScalingFactory = asgFactory

	// error, no asgs
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(nil, nil)
	dead, err = connector.IsJobDead("public-services")
	assert.Error(t, err)
	assert.True(t, dead)

	// no error, not found, dead
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	asgOut := make([]*autoscaling.Group, 0)
	group := autoscaling.Group{}
	asgOut = append(asgOut, &group)
	output := &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: asgOut}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	dead, err = connector.IsJobDead("public-services")
	assert.NoError(t, err)
	assert.True(t, dead)

	// no error, found, NOT dead
	asgFactory.EXPECT().CreateAutoScaling(gomock.Any()).Return(asgIF)
	autoScalingGroups := make([]*autoscaling.Group, 0)
	key := "datacenter"
	tagVal := "private-services"
	asgName := "my-asg"
	var tags []*autoscaling.TagDescription
	td := autoscaling.TagDescription{Key: &key, Value: &tagVal}
	tags = append(tags, &td)
	asgIn := autoscaling.Group{Tags: tags, AutoScalingGroupName: &asgName}
	autoScalingGroups = append(autoScalingGroups, &asgIn)
	output = &autoscaling.DescribeAutoScalingGroupsOutput{AutoScalingGroups: autoScalingGroups}
	asgIF.EXPECT().DescribeAutoScalingGroups(gomock.Any()).Return(output, nil)
	dead, err = connector.IsJobDead("private-services")
	assert.NoError(t, err)
	assert.False(t, dead)
}
