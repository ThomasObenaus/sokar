package nomadWorker

import "github.com/aws/aws-sdk-go/service/autoscaling"

// AutoScaling is the minimal interface needed to interact with aws autoscaling
type AutoScaling interface {
	DescribeAutoScalingGroups(input *autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
}
