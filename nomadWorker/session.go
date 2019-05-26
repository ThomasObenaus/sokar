package nomadWorker

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func newAWSSessionFromProfile(profile string) (*session.Session, error) {
	verboseCredErrors := true

	cfg := aws.Config{CredentialsChainVerboseErrors: &verboseCredErrors}
	sessionOpts := session.Options{Profile: profile, Config: cfg, SharedConfigState: session.SharedConfigEnable}

	return session.NewSessionWithOptions(sessionOpts)
}

func newAWSSession() (*session.Session, error) {
	verboseCredErrors := true

	cfg := aws.Config{CredentialsChainVerboseErrors: &verboseCredErrors}
	sessionOpts := session.Options{Config: cfg, SharedConfigState: session.SharedConfigEnable}

	return session.NewSessionWithOptions(sessionOpts)
}
