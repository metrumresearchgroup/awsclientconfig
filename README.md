# awsclientconfig

This package allows us to create an STS-targeted user for creation of a session on the fly within AWS. It's smart enough
to skip steps when the answers are obvious re: whether we need to do STS assignment or not. This allows us to focus the
"blast-radius" of a user in aws with too many inherent permissions to a set of roles in AWS that are accessible to a
service user. (least privilege for the task)

## Usage

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/metrumresearchgroup/awsclientconfig"
)

func main() {
	var err error

	// We need to set up the basic configuration structure.
	// Don't worry about how we handle some of these things,
	// It's all verified before using, so if you put in a bad
	// value or have an insufficient set of credentials, we
	// catch it before sending on to AWS here.
	var (
		clientConfig awsclientconfig.ClientConfig
	)
	{
		creds := aws.Credentials{
			AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
			SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
			SessionToken:    os.Getenv("AWS_SESSION_TOKEN"),
			Source:          "environment",
		}
		region := os.Getenv("AWS_REGION")
		arn := os.Getenv("AWS_TARGET_ARN")

		if clientConfig, err = awsclientconfig.New(creds, region, arn); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Configuration error: %v", err)
			os.Exit(1)
		}
	}

	// We then log in using this config, and the ARN will 
	// automatically become an STS "sudo" to the role, if present.
	var (
		cognitoAwsConfig aws.Config
	)
	{
		if cognitoAwsConfig, err = clientConfig.Login(context.Background(), "test-cognito-permissions"); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Login error: %v", err)
			os.Exit(1)
		}
	}

	// Since what we get back is an opaque AWS config, we can
	// then use that to start new services from this config.
	var (
		cognitoProvider cognitoidentityprovider.Client
	)
	{
		if cognitoProvider, err = cognitoidentityprovider.NewFromConfig(cognitoAwsConfig); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Spawn cognito client: %v", err)
			os.Exit(1)
		}
	}

	// Then we use the service.
	var (
		user cognitoidentityprovider.GetUserOutput
	)
	{
		token := os.Getenv("COGNITO_TOKEN")

		if user, err = cognitoProvider.GetUser(context.Background(), &cognitoidentityprovider.GetUserInput{AccessToken: token}); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Spawn cognito client: %v", err)
			os.Exit(1)
		}
	}

	fmt.Println(user)
}
```