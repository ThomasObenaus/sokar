package nomadWorker

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/thomasobenaus/sokar/aws"
)

func (c *Connector) upscale(datacenter string, min uint, max uint, desiredCount uint) error {
	sess, err := c.createSession()
	if err != nil {
		return err
	}
	autoScalingIF := c.autoScalingFactory.CreateAutoScaling(sess)

	asgQuery := aws.AutoScalingGroupQuery{
		AsgIF:    autoScalingIF,
		TagKey:   c.tagKey,
		TagValue: datacenter,
	}

	asgName, err := asgQuery.GetAutoScalingGroupName()
	if err != nil {
		return err
	}

	size := int64(desiredCount)
	minSize := int64(min)
	maxSize := int64(max)

	input := &autoscaling.UpdateAutoScalingGroupInput{
		AutoScalingGroupName: &asgName,
		MinSize:              &minSize,
		MaxSize:              &maxSize,
		DesiredCapacity:      &size,
	}

	_, err = autoScalingIF.UpdateAutoScalingGroup(input)
	if err != nil {
		return err
	}

	c.log.Info().Msgf("Upscaled min=max=desiredCapacity of %s to %d.", asgName, size)
	return nil
}
