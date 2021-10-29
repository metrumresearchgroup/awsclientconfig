package awsclientconfig

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// DefaultRegion is a universal default in case you have a different
// value you'd like to assign globally.
var DefaultRegion = "us-east-1"

// ClientConfig contains everything to login and perform
// STS transition to another ARN, if set.
type ClientConfig struct {
	Credentials aws.Credentials
	Region      string
	ARN         string
}

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

var (
	AccessKeyPattern = regexp.MustCompile(`^[\w]+$`)

	ErrAKIDRequired    = errors.New("access key id required")
	ErrAKIDMinLen      = errors.New("access key id is less than 16 characters long")
	ErrAKIDMaxLen      = errors.New("access key id greater than 128 characters long")
	ErrAKIDBadPrefix   = errors.New("access key id must begin with either AKIA or ASIA")
	ErrAKIDInvalidChar = fmt.Errorf("access key id must match the pattern %s", AccessKeyPattern)
	ErrSAKRequired     = errors.New("secret access key required")
	ErrSTRequired      = errors.New("session token required for ASIA access key id")
)

func ValidateCredentials(c aws.Credentials) error {
	if c.AccessKeyID == "" {
		return ErrAKIDRequired
	}
	if len(c.AccessKeyID) < 16 {
		return ErrAKIDMinLen
	}
	if len(c.AccessKeyID) > 128 {
		return ErrAKIDMaxLen
	}
	if !(strings.HasPrefix(c.AccessKeyID, "AKIA") || strings.HasPrefix(c.AccessKeyID, "ASIA")) {
		return ErrAKIDBadPrefix
	}
	if !AccessKeyPattern.MatchString(c.AccessKeyID) {
		return ErrAKIDInvalidChar
	}
	if c.SecretAccessKey == "" {
		return ErrSAKRequired
	}
	if strings.HasPrefix(c.AccessKeyID, "ASIA") && c.SessionToken == "" {
		return ErrSTRequired
	}

	return nil
}

type opts struct {
	fns []func(loadOptions *config.LoadOptions) error
}

func (o *opts) add(opts ...func(loadOptions *config.LoadOptions) error) {
	o.fns = append(o.fns, opts...)
}

func (cc ClientConfig) login(ctx context.Context, optFns ...func(loadOptions *config.LoadOptions) error) (aws.Config, error) {
	var o opts

	copy(o.fns, optFns)

	o.add(config.WithDefaultRegion(DefaultRegion))

	if cc.Region != "" {
		o.add(config.WithRegion(cc.Region))
	}

	o.add(config.WithCredentialsProvider(
		credentials.StaticCredentialsProvider{
			Value: cc.Credentials,
		},
	))

	return config.LoadDefaultConfig(ctx, o.fns...)
}

func (cc ClientConfig) Login(ctx context.Context, sessionName string, optFns ...func(loadOptions *config.LoadOptions) error) (aws.Config, error) {
	awsConfig, err := cc.login(ctx, optFns...)
	if err != nil {
		return aws.Config{}, err
	}

	// If the ARN is unset, there is nothing more to do.
	if cc.ARN == "" {
		return awsConfig, nil
	}

	client := sts.NewFromConfig(awsConfig)

	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return aws.Config{}, fmt.Errorf("getting caller identity: %w", err)
	}

	// If the ARN is set, and matches the user, return the original.
	if aws.ToString(identity.Arn) == cc.ARN {
		return awsConfig, nil
	}

	// Use current role/user credentials to assume another role
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(cc.ARN),
		RoleSessionName: aws.String(sessionName),
	}

	aro, err := client.AssumeRole(ctx, input)
	if err != nil {
		return aws.Config{}, fmt.Errorf("assuming role %s: %w", cc.ARN, err)
	}

	canExpire := aws.ToTime(aro.Credentials.Expiration) != time.Time{}

	ncc, err := New(aws.Credentials{
		AccessKeyID:     aws.ToString(aro.Credentials.AccessKeyId),
		SecretAccessKey: aws.ToString(aro.Credentials.SecretAccessKey),
		SessionToken:    aws.ToString(aro.Credentials.SessionToken),
		Source:          cc.Credentials.Source,
		CanExpire:       canExpire,
		Expires:         aws.ToTime(aro.Credentials.Expiration),
	}, cc.Region, cc.ARN)
	if err != nil {
		return aws.Config{}, err
	}

	return ncc.login(ctx, optFns...)
}
