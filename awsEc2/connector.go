package awsEc2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rs/zerolog"
	"github.com/thomasobenaus/sokar/aws"
	iface "github.com/thomasobenaus/sokar/aws/iface"
)

// Connector is a object that allows to interact with nomad worker
type Connector struct {
	log zerolog.Logger

	// tagKey is the name of the tag that is used to find the instances/ autoscalinggroup/ node
	// of the nomad worker that should be scaled.
	tagKey string

	// autoScalingFactory factory used to create the component to access
	// the AWS AutoScaling resources
	autoScalingFactory iface.AutoScalingFactory

	// fnCreateSession is the function that should be used to create the aws session
	// which is needed to access the aws resources.
	fnCreateSession func(region string) (*session.Session, error)

	// fnCreateSessionFromProfile is the function that should be used to create the aws session
	// which is needed to access the aws resources.
	// Here a given aws profile name is regarded.
	fnCreateSessionFromProfile func(profile string) (*session.Session, error)

	// awsProfile is used to specify which shared credentials shall be used in order to
	// gain permission to access the needed AWS resources.
	// If empty the default credentials will be used.
	awsProfile string

	// awsRegion is the region where the datacenter to be scaled is located in.
	awsRegion string
}

// Option represents an option for the awsEc2 Connector
type Option func(c *Connector)

// WithLogger adds a configured Logger to the awsEc2 Connector
func WithLogger(logger zerolog.Logger) Option {
	return func(c *Connector) {
		c.log = logger
	}
}

// WithAwsProfile sets the aws profile to be used.
// The profile represents the name of the aws profile that shall be used to access the resources to scale the aws AutoScalingGroup.
// This parameter is optional. If the profile is NOT set the instance where sokar runs on has to have enough permissions to access the
// resources (ASG) for scaling (e.g. granted by a AWS Instance Profile). In this case the region parameter has to be specified instead (via WithAwsRegion()).
func WithAwsProfile(profile string) Option {
	return func(c *Connector) {
		c.awsProfile = profile
	}
}

// WithAwsRegion sets the aws region in which the resource to be scaled can be found
func WithAwsRegion(region string) Option {
	return func(c *Connector) {
		c.awsRegion = region
	}
}

// New creates a new nomad connector
func New(asgTagKey string, options ...Option) (*Connector, error) {
	awsEc2Conn := Connector{
		tagKey:                     asgTagKey,
		autoScalingFactory:         &aws.AutoScalingFactoryImpl{},
		fnCreateSession:            aws.NewAWSSession,
		fnCreateSessionFromProfile: aws.NewAWSSessionFromProfile,
	}

	// apply the options
	for _, opt := range options {
		opt(&awsEc2Conn)
	}

	if err := validate(awsEc2Conn); err != nil {
		return nil, err
	}

	return &awsEc2Conn, nil
}

func (c *Connector) String() string {
	return "AWS-EC2 (AWS AutoScalingGroup, EC2 instances)"
}

func validate(c Connector) error {

	if len(c.tagKey) == 0 {
		return fmt.Errorf("the tagkey to identify the AutoScalingGroup that should be scaled is not specified")
	}
	if len(c.awsProfile) == 0 && len(c.awsRegion) == 0 {
		return fmt.Errorf("aws profile and region are not specified")
	}
	return nil
}
