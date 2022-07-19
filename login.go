package awsclientconfig

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go-v2/service/sts/types"
)

// Login signs in to the aws service and retrieves an aws.Config for passing
// to aws service clients.
func (cc *ClientConfig) Login(ctx context.Context, sessionName string, optFns ...func(loadOptions *config.LoadOptions) error) (aws.Config, error) {
	var (
		err       error
		awsConfig aws.Config
	)
	{
		if awsConfig, err = cc.loadConfig(ctx, optFns...); err != nil {
			return aws.Config{}, err
		}
	}

	// If the ARN is unset, there is nothing more to do.
	if cc.ARN == "" {
		return awsConfig, nil
	}

	stsClient := sts.NewFromConfig(awsConfig)

	// Retrieve the current user's STS identity, for checking its ARN.
	var (
		identity *sts.GetCallerIdentityOutput
	)
	{
		if identity, err = stsClient.GetCallerIdentity(
			ctx, &sts.GetCallerIdentityInput{},
		); err != nil {
			return aws.Config{}, fmt.Errorf("getting caller identity: %w", err)
		}
	}

	// If the ARN is set, and matches the user, return the original.
	if aws.ToString(identity.Arn) == cc.ARN {
		return awsConfig, nil
	}

	// Retrieve the credentials from STS to assume the role.
	var (
		creds *types.Credentials
	)
	{
		var assumeRoleOutput *sts.AssumeRoleOutput

		if assumeRoleOutput, err = stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
			RoleArn:         aws.String(cc.ARN),
			RoleSessionName: aws.String(sessionName),
		}); err != nil {
			return aws.Config{}, fmt.Errorf("assuming role %s: %w", cc.ARN, err)
		}

		creds = assumeRoleOutput.Credentials
	}

	// Create a new client configuration for the STS data.
	var (
		newClientConfig *ClientConfig
	)
	{
		if newClientConfig, err = New(
			aws.Credentials{
				AccessKeyID:     aws.ToString(creds.AccessKeyId),
				SecretAccessKey: aws.ToString(creds.SecretAccessKey),
				SessionToken:    aws.ToString(creds.SessionToken),
				Source:          "awsclientconfig",
				CanExpire:       aws.ToTime(creds.Expiration) != time.Time{},
				Expires:         aws.ToTime(creds.Expiration),
			},
			cc.Region,
			cc.ARN,
			WithMapID(cc.MapId),
			WithAppIdentifier(cc.AppIdentifier),
			WithRefreshInterval(cc.Refresh),
		); err != nil {
			return aws.Config{}, err
		}
	}

	// Load the client configuration in AWS, returning the aws.Config.
	return newClientConfig.loadConfig(ctx, optFns...)
}

func (cc *ClientConfig) loadConfig(ctx context.Context, optFns ...func(loadOptions *config.LoadOptions) error) (aws.Config, error) {
	var o opts

	copy(o.fns, optFns)

	o.add(config.WithDefaultRegion(DefaultRegion))

	if cc.Region != "" {
		o.add(config.WithRegion(cc.Region))
	}

	o.add(config.WithCredentialsProvider(
		credentials.StaticCredentialsProvider{
			Value: cc.Credentials("awsclientconfig"),
		},
	))

	return config.LoadDefaultConfig(ctx, o.fns...)
}

// opts internally tracks aws's opts functions cleanly.
type opts struct {
	fns []func(loadOptions *config.LoadOptions) error
}

func (o *opts) add(opts ...func(loadOptions *config.LoadOptions) error) {
	o.fns = append(o.fns, opts...)
}
