package nomadWorker

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

// AutoScaling is the minimal interface needed to interact with aws autoscaling
type AutoScaling interface {
	DescribeAutoScalingGroups(input *autoscaling.DescribeAutoScalingGroupsInput) (*autoscaling.DescribeAutoScalingGroupsOutput, error)
}

// AutoScalingFactory is an interface to create the AutoScaling type based on the given session.
type AutoScalingFactory interface {
	CreateAutoScaling(session *session.Session) AutoScaling
}
