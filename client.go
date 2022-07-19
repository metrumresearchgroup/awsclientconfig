package awsclientconfig

import (
	"regexp"
	"time"

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
	AccessKey     string        `mapstructure:"access-key"`
	SecretKey     string        `mapstructure:"secret-key"`
	SessionToken  string        `mapstructure:"session-token"`
	Region        string        `mapstructure:"region"`
	ARN           string        `mapstructure:"arn"`
	MapId         string        `mapstructure:"map-id"`
	AppIdentifier string        `mapstructure:"application-identifier"`
	Refresh       time.Duration `mapstructure:"refresh-duration"`
}

// New creates a new ClientConfig, after verifying parameters are complete.
func New(credentials aws.Credentials, region, arn string, mapId, appIdentifier string, refreshInterval time.Duration) (*ClientConfig, error) {
	if err := ValidateCredentials(credentials); err != nil {
		return nil, err
	}

	if refreshInterval == 0 {
		refreshInterval = time.Minute * 5
	}

	return &ClientConfig{
		AccessKey:     credentials.AccessKeyID,
		SecretKey:     credentials.SecretAccessKey,
		SessionToken:  credentials.SessionToken,
		Region:        region,
		ARN:           arn,
		MapId:         mapId,
		AppIdentifier: appIdentifier,
		Refresh:       refreshInterval,
	}, nil
}

func (cc *ClientConfig) Credentials(source string) aws.Credentials {
	return aws.Credentials{
		AccessKeyID:     cc.AccessKey,
		SecretAccessKey: cc.SecretKey,
		SessionToken:    cc.SessionToken,
		Source:          source,
	}
}
