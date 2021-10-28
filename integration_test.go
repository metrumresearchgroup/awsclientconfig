//go:build integration
// +build integration

package awsclientconfig_test

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/metrumresearchgroup/wrapt"

	. "github.com/metrumresearchgroup/awsclientconfig"
)

func fromEnv() aws.Credentials {
	akid := os.Getenv("AWS_ACCESS_KEY_ID")
	sak := os.Getenv("AWS_SECRET_ACCESS_KEY")
	st := os.Getenv("AWS_SESSION_TOKEN_KEY")

	return aws.Credentials{
		AccessKeyID:     akid,
		SecretAccessKey: sak,
		SessionToken:    st,
		Source:          "environment",
	}
}

func Test_StsLogin(tt *testing.T) {
	t := wrapt.WrapT(tt)

	creds := fromEnv()
	region := "us-east-1"
	arn := "arn:aws:iam::654062675744:role/test-cognito-permissions"

	cc, err := NewClientConfig(creds, region, arn)
	t.R.NoError(err)

	_, err = cc.Login(context.Background(), "new-session")
	t.R.NoError(err)
}
