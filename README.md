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
	creds := aws.Credentials{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		SessionToken:    os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Source:          "environment",
	}
	region := os.Getenv("AWS_REGION")
	arn := os.Getenv("AWS_TARGET_ARN")
	token := os.Getenv("COGNITO_TOKEN")

	config, err := awsclientconfig.NewClientConfig(creds, region, arn)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Configuration error: %v", err)
		os.Exit(1)
	}

	cognitoAwsConfig, err := config.Login(context.Background(), "test-cognito-permissions")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Login error: %v", err)
		os.Exit(1)
	}

	svc, err := cognitoidentityprovider.NewFromConfig(cognitoAwsConfig)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Spawn cognito client: %v", err)
		os.Exit(1)
	}

	user, err := svc.GetUser(context.Background(), &cognitoidentityprovider.GetUserInput{AccessToken: token})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Spawn cognito client: %v", err)
		os.Exit(1)
	}

	fmt.Println(user)
}
```