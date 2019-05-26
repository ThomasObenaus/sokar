package nomadWorker

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/rs/zerolog"
	iface "github.com/thomasobenaus/sokar/nomadWorker/iface"
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
	fnCreateSession func() (*session.Session, error)

	// fnCreateSessionFromProfile is the function that should be used to create the aws session
	// which is needed to access the aws resources.
	// Here a given aws profile name is regarded.
	fnCreateSessionFromProfile func(profile string) (*session.Session, error)

	// awsProfile is used to specify which shared credentials shall be used in order to
	// gain permission to access the needed AWS resources.
	// If empty the default credentials will be used.
	awsProfile string
}

// Config contains the main configuration for the nomad worker connector
type Config struct {
	Logger     zerolog.Logger
	AWSProfile string
}

// New creates a new nomad connector
func (cfg *Config) New() (*Connector, error) {

	nc := &Connector{
		log:                        cfg.Logger,
		tagKey:                     "datacenter",
		fnCreateSession:            newAWSSession,
		fnCreateSessionFromProfile: newAWSSessionFromProfile,
		awsProfile:                 cfg.AWSProfile,
		autoScalingFactory:         &autoScalingFactoryImpl{},
	}

	cfg.Logger.Info().Msg("Setting up nomad worker connector ... done")
	return nc, nil
}
