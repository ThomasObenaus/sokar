package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	mock_nomadWorker "github.com/thomasobenaus/sokar/test/aws"
)

func Test_MonitorInstanceScaling(t *testing.T) {
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
	require.NotNil(t, state.nextToken)
	assert.Equal(t, nextToken, *state.nextToken)
}
