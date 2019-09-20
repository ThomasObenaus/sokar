package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewAWSSessionFromProfile creates a new session needed for interaction with the aws api. Here a profile can be specified.
func NewAWSSessionFromProfile(profile string) (*session.Session, error) {
	verboseCredErrors := true

	cfg := aws.Config{CredentialsChainVerboseErrors: &verboseCredErrors}
	sessionOpts := session.Options{Profile: profile, Config: cfg, SharedConfigState: session.SharedConfigEnable}

	return session.NewSessionWithOptions(sessionOpts)
}

// NewAWSSession creates a new session needed for interaction with the aws api
func NewAWSSession(region string) (*session.Session, error) {
	if len(region) == 0 {
		return nil, fmt.Errorf("Required region parameter is empty")
	}

	verboseCredErrors := true

	cfg := aws.Config{CredentialsChainVerboseErrors: &verboseCredErrors, Region: &region}
	sessionOpts := session.Options{Config: cfg, SharedConfigState: session.SharedConfigEnable}

	return session.NewSessionWithOptions(sessionOpts)
}
