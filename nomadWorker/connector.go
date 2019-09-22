package nomadWorker

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

// Config contains the main configuration for the nomad worker connector
type Config struct {
	Logger zerolog.Logger

	// AWSProfile represents the name of the aws profile that shall be used to access the resources to scale the data-center.
	// This parameter is optional. If it is empty the instance where sokar runs on has to have enough permissions to access the
	// resources (ASG) for scaling. In this case the AWSRegion parameter has to be specified as well.
	AWSProfile string

	// AWSRegion is an optional parameter and is regarded only if the parameter AWSProfile is empty.
	// The AWSRegion has to specify the region in which the data-center to be scaled resides in.
	AWSRegion string
}

// New creates a new nomad connector
func (cfg *Config) New() (*Connector, error) {

	cfg.Logger.Info().Msg("Setting up nomad worker connector ...")
	if len(cfg.AWSProfile) == 0 && len(cfg.AWSRegion) == 0 {
		return nil, fmt.Errorf("The parameters AWSRegion and AWSProfile are empty")
	}

	nc := &Connector{
		log:                        cfg.Logger,
		tagKey:                     "datacenter",
		autoScalingFactory:         &aws.AutoScalingFactoryImpl{},
		fnCreateSession:            aws.NewAWSSession,
		fnCreateSessionFromProfile: aws.NewAWSSessionFromProfile,
		awsProfile:                 cfg.AWSProfile,
		awsRegion:                  cfg.AWSRegion,
	}

	cfg.Logger.Info().Msg("Setting up nomad worker connector ... done")
	return nc, nil
}

func (c *Connector) String() string {
	return "Nomad-DC (Nomad DataCenter, on AWS)"
}
