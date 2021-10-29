package awsclientconfig

import (
	"regexp"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// DefaultRegion is a universal default in case you have a different
// value you'd like to assign globally.
var DefaultRegion = "us-east-1"

// AccessKeyPattern represents that access keys don't go beyond
// alphanumerics.
var AccessKeyPattern = regexp.MustCompile(`^[\w]+$`)

// ClientConfig contains everything to login and perform
// STS transition to another ARN, if set.
type ClientConfig struct {
	Credentials aws.Credentials
	Region      string
	ARN         string
}

// New creates a new ClientConfig, after verifying parameters are complete.
func New(credentials aws.Credentials, region, arn string) (ClientConfig, error) {
	if err := ValidateCredentials(credentials); err != nil {
		return ClientConfig{}, err
	}

	return ClientConfig{
		Credentials: credentials,
		Region:      region,
		ARN:         arn,
	}, nil
}
