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
	AccessKey     string        `mapstructure:"access_key" json:"access_key"`
	SecretKey     string        `mapstructure:"secret_key" json:"secret_key"`
	SessionToken  string        `mapstructure:"session_token" json:"session_token"`
	Region        string        `mapstructure:"region" json:"region"`
	ARN           string        `mapstructure:"arn" json:"arn"`
	MapId         string        `mapstructure:"map_id" json:"map_id"`
	AppIdentifier string        `mapstructure:"application_identifier" json:"app_identifier"`
	Refresh       time.Duration `mapstructure:"refresh_duration" json:"refresh"`
}

// New creates a new ClientConfig, after verifying parameters are complete.
func New(credentials aws.Credentials, region, arn string, opts ...func(*ClientConfig)) (*ClientConfig, error) {
	if err := ValidateCredentials(credentials); err != nil {
		return nil, err
	}

	cc := &ClientConfig{
		AccessKey:    credentials.AccessKeyID,
		SecretKey:    credentials.SecretAccessKey,
		SessionToken: credentials.SessionToken,
		Region:       region,
		ARN:          arn,
	}

	for _, o := range opts {
		o(cc)
	}

	if cc.Refresh == 0 {
		cc.Refresh = time.Minute * 5
	}

	return cc, nil
}

func WithRefreshInterval(d time.Duration) func(*ClientConfig) {
	return func(cc *ClientConfig) {
		cc.Refresh = d
	}
}

func WithMapID(id string) func(*ClientConfig) {
	return func(cc *ClientConfig) {
		cc.MapId = id
	}
}

func WithAppIdentifier(id string) func(config *ClientConfig) {
	return func(cc *ClientConfig) {
		cc.AppIdentifier = id
	}
}
func (cc *ClientConfig) Credentials(source string) aws.Credentials {
	return aws.Credentials{
		AccessKeyID:     cc.AccessKey,
		SecretAccessKey: cc.SecretKey,
		SessionToken:    cc.SessionToken,
		Source:          source,
	}
}

func (cc *ClientConfig) Inherit(source *ClientConfig) {
	if len(cc.AccessKey) == 0 && len(source.AccessKey) > 0 {
		cc.AccessKey = source.AccessKey
	}

	if len(cc.SecretKey) == 0 && len(source.SecretKey) > 0 {
		cc.SecretKey = source.SecretKey
	}

	if len(cc.SessionToken) == 0 && len(source.SessionToken) > 0 {
		cc.SessionToken = source.SessionToken
	}

	if len(cc.Region) == 0 && len(source.Region) > 0 {
		cc.Region = source.Region
	}

	if len(cc.ARN) == 0 && len(source.ARN) > 0 {
		cc.ARN = source.ARN
	}

	if len(cc.MapId) == 0 && len(source.MapId) > 0 {
		cc.MapId = source.MapId
	}

	if len(cc.AppIdentifier) == 0 && len(source.AppIdentifier) > 0 {
		cc.AppIdentifier = source.AppIdentifier
	}

	if cc.Refresh == 0 && source.Refresh > 0 {
		cc.Refresh = source.Refresh
	}
}
