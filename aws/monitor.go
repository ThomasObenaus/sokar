package aws

import (
	"fmt"
	"time"

	aws "github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	iface "github.com/thomasobenaus/sokar/aws/iface"
)

var monitorAWSStateBackoff time.Duration = time.Millisecond * 500

// MonitorInstanceScaling will block until the instance is scaled up/ down
// The function returns the number of iterations that where needed to monitor the scaling of the instance
func MonitorInstanceScaling(autoScaling iface.AutoScaling, autoScalingGroupName string, activityID string, timeout time.Duration, log zerolog.Logger) (uint, error) {
	start := time.Now()
	iterations := uint(0)
	for {
		iterations++
		state, err := getCurrentScalingState(autoScaling, autoScalingGroupName, activityID, log)
		if err != nil {
			return iterations, err
		}

		if state.progress >= 100 {
			// scaling completed
			return iterations, nil
		}

		time.Sleep(monitorAWSStateBackoff)

		if time.Since(start) >= timeout {
			return iterations, fmt.Errorf("MonitorInstanceScaling timed out after %v (%d iterations)", timeout, iterations)
		}
	}
}

type scalingState struct {
	status   string
	progress int64
}

func getCurrentScalingState(autoScaling iface.AutoScaling, autoScalingGroupName string, activityID string, log zerolog.Logger) (*scalingState, error) {

	activityIDs := make([]*string, 0)
	activityIDs = append(activityIDs, &activityID)
	input := aws.DescribeScalingActivitiesInput{AutoScalingGroupName: &autoScalingGroupName, ActivityIds: activityIDs}
	if err := input.Validate(); err != nil {
		return nil, errors.Wrap(err, "Validation of input failed")
	}

	// First create the request
	req, output := autoScaling.DescribeScalingActivitiesRequest(&input)
	if req == nil {
		return nil, fmt.Errorf("Request generated by DescribeScalingActivitiesInput is nil")
	}

	// Now send the request
	if err := req.Send(); err != nil {
		return nil, errors.Wrap(err, "Failed sending request to describe the scaling activities")
	}

	if output == nil {
		return nil, fmt.Errorf("DescribeScalingActivitiesOutput is invalid")
	}

	if len(output.Activities) == 0 || output.Activities[0].StatusCode == nil || output.Activities[0].Progress == nil {
		return nil, fmt.Errorf("DescribeScalingActivitiesOutput contains no valid activities")
	}

	log.Debug().Str("activityID", activityID).Str("autoScalingGroupName", autoScalingGroupName).Msgf("Output: %v\n", output)
	state := &scalingState{status: *output.Activities[0].StatusCode, progress: *output.Activities[0].Progress}
	return state, nil
}
