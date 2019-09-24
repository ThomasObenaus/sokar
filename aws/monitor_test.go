package aws

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_nomadWorker "github.com/thomasobenaus/sokar/test/aws"
)

func Test_GetCurrentScalingState(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	// nil request
	asgName := "asgName"
	activityID := "activityID"
	asgIF.EXPECT().DescribeScalingActivitiesRequest(gomock.Any()).Return(nil, nil)
	state, err := getCurrentScalingState(asgIF, asgName, activityID)
	assert.Error(t, err)
	assert.Nil(t, state)

	// nil output
	req := request.Request{}
	asgIF.EXPECT().DescribeScalingActivitiesRequest(gomock.Any()).Return(&req, nil)
	state, err = getCurrentScalingState(asgIF, asgName, activityID)
	assert.Error(t, err)
	assert.Nil(t, state)

	// success
	progress := int64(50)
	statusCode := "InProgress"
	activity := autoscaling.Activity{ActivityId: &activityID, AutoScalingGroupName: &asgName, Progress: &progress, StatusCode: &statusCode}
	activities := make([]*autoscaling.Activity, 0)
	activities = append(activities, &activity)
	nextToken := "next"
	output := autoscaling.DescribeScalingActivitiesOutput{Activities: activities, NextToken: &nextToken}
	asgIF.EXPECT().DescribeScalingActivitiesRequest(gomock.Any()).Return(&req, &output)
	state, err = getCurrentScalingState(asgIF, asgName, activityID)
	assert.NoError(t, err)
	require.NotNil(t, state)
	assert.Equal(t, progress, state.progress)
	assert.Equal(t, statusCode, state.status)
}

func Test_MonitorInstanceScaling(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	asgIF := mock_nomadWorker.NewMockAutoScaling(mockCtrl)

	// err no loop
	asgName := "asgName"
	activityID := "activityID"
	asgIF.EXPECT().DescribeScalingActivitiesRequest(gomock.Any()).Return(nil, nil)
	err := MonitorInstanceScaling(asgIF, asgName, activityID, time.Second*10)
	assert.Error(t, err)

	// success one loop
	progress := int64(100)
	statusCode := "InProgress"
	activity := autoscaling.Activity{ActivityId: &activityID, AutoScalingGroupName: &asgName, Progress: &progress, StatusCode: &statusCode}
	activities := make([]*autoscaling.Activity, 0)
	activities = append(activities, &activity)
	output := autoscaling.DescribeScalingActivitiesOutput{Activities: activities}
	req := request.Request{}
	asgIF.EXPECT().DescribeScalingActivitiesRequest(gomock.Any()).Return(&req, &output)
	err = MonitorInstanceScaling(asgIF, asgName, activityID, time.Second*10)
	assert.NoError(t, err)

	// timeout
	progress = int64(50)
	statusCode = "InProgress"
	activity = autoscaling.Activity{ActivityId: &activityID, AutoScalingGroupName: &asgName, Progress: &progress, StatusCode: &statusCode}
	activities = make([]*autoscaling.Activity, 0)
	activities = append(activities, &activity)
	output = autoscaling.DescribeScalingActivitiesOutput{Activities: activities}
	asgIF.EXPECT().DescribeScalingActivitiesRequest(gomock.Any()).Return(&req, &output).AnyTimes()
	err = MonitorInstanceScaling(asgIF, asgName, activityID, time.Second*1)
	assert.Error(t, err)
}
