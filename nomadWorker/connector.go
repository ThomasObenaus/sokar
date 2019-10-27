package nomadWorker

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	nomadApi "github.com/hashicorp/nomad/api"
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

	// Interface that is used to interact with nomad nodes
	nodesIF Nodes

	// nodeDrainDeadline the maximum amount of time nomad will wait before the containers will be forced to be moved
	nodeDrainDeadline time.Duration

	// instanceTerminationTimeout is the timeout used to monitor the scale of an aws instance at maximum
	instanceTerminationTimeout time.Duration
}

// Option represents an option for the nomadWorker Connector
type Option func(c *Connector)

// WithLogger adds a configured Logger to the nomadWorker Connector
func WithLogger(logger zerolog.Logger) Option {
	return func(c *Connector) {
		c.log = logger
	}
}

// WithAwsRegion sets the aws region in which the resource to be scaled can be found
func WithAwsRegion(region string) Option {
	return func(c *Connector) {
		c.awsRegion = region
	}
}

// TimeoutForInstanceTermination sets the maximum time the instance termination will be monitored before assuming that this action failed.
func TimeoutForInstanceTermination(timeout time.Duration) Option {
	return func(c *Connector) {
		c.instanceTerminationTimeout = timeout
	}
}

// New creates a new nomad worker connector.
// The profile represents the name of the aws profile that shall be used to access the resources to scale the aws AutoScalingGroup.
// This parameter is optional. If the profile is NOT set the instance where sokar runs on has to have enough permissions to access the
// resources (ASG) for scaling (e.g. granted by a AWS Instance Profile). In this case the region parameter has to be specified instead (via WithAwsRegion()).
func New(nomadServerAddress, awsProfile string, options ...Option) (*Connector, error) {
	if len(nomadServerAddress) == 0 {
		return nil, fmt.Errorf("required configuration 'nomadServerAddress' is missing/ empty")
	}

	nomadConn := Connector{
		tagKey:                     "datacenter",
		autoScalingFactory:         &aws.AutoScalingFactoryImpl{},
		fnCreateSession:            aws.NewAWSSession,
		fnCreateSessionFromProfile: aws.NewAWSSessionFromProfile,
		nodeDrainDeadline:          time.Second * 60,
		instanceTerminationTimeout: time.Minute * 10,
		awsProfile:                 awsProfile,
	}

	// config needed to set up a nomad api client
	config := nomadApi.DefaultConfig()
	config.Address = nomadServerAddress
	//config.SecretID = token
	//config.TLSConfig.TLSServerName = tls_server_name

	client, err := nomadApi.NewClient(config)
	if err != nil {
		return nil, err
	}
	nomadConn.nodesIF = client.Nodes()

	// apply the options
	for _, opt := range options {
		opt(&nomadConn)
	}

	nomadConn.log.Info().Str("srvAddr", nomadServerAddress).Str("awsProfile", nomadConn.awsProfile).Str("awsRegion", nomadConn.awsRegion).Msg("Setting up nomad worker connector ...")

	if err := validate(nomadConn); err != nil {
		return nil, err
	}

	nomadConn.log.Info().Msg("Setting up nomad worker connector ... done")
	return &nomadConn, nil
}

func (c *Connector) String() string {
	return "Nomad-DC (Nomad DataCenter, on AWS)"
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
