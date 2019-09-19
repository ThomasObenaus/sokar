package awsEc2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func newAWSSessionFromProfile(profile string) (*session.Session, error) {
	verboseCredErrors := true

	cfg := aws.Config{CredentialsChainVerboseErrors: &verboseCredErrors}
	sessionOpts := session.Options{Profile: profile, Config: cfg, SharedConfigState: session.SharedConfigEnable}

	return session.NewSessionWithOptions(sessionOpts)
}

func newAWSSession(region string) (*session.Session, error) {
	if len(region) == 0 {
		return nil, fmt.Errorf("Required region parameter is empty")
	}

	verboseCredErrors := true

	cfg := aws.Config{CredentialsChainVerboseErrors: &verboseCredErrors, Region: &region}
	sessionOpts := session.Options{Config: cfg, SharedConfigState: session.SharedConfigEnable}

	return session.NewSessionWithOptions(sessionOpts)
}
