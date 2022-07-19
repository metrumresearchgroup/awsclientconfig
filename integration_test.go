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

func fromEnv() (aws.Credentials, string) {
	akid := os.Getenv("AWS_ACCESS_KEY_ID")
	sak := os.Getenv("AWS_SECRET_ACCESS_KEY")
	st := os.Getenv("AWS_SESSION_TOKEN")
	arn := os.Getenv("AWS_TARGET_ARN")
	return aws.Credentials{
		AccessKeyID:     akid,
		SecretAccessKey: sak,
		SessionToken:    st,
		Source:          "environment",
	}, arn
}

func Test_StsLogin_SuccessSTS(tt *testing.T) {
	t := wrapt.WrapT(tt)

	creds, arn := fromEnv()
	region := "us-east-1"

	cc, err := New(creds, region, arn, "", "", 0)
	t.R.NoError(err)

	cli, err := cc.Login(context.Background(), "new-session")
	t.R.NoError(err)
	t.R.NotEqual(aws.Config{}, cli)
}

func Test_StsLogin_SuccessSkipARN(tt *testing.T) {
	t := wrapt.WrapT(tt)

	creds, arn := fromEnv()
	region := "us-east-1"

	// clearing ARN, going to be a normal login as user.
	arn = ""

	cc, err := New(creds, region, arn, "", "", 0)
	t.R.NoError(err)

	cli, err := cc.Login(context.Background(), "new-session")
	t.R.NoError(err)
	t.R.NotEqual(aws.Config{}, cli)
}

func Test_StsLogin_FailureBadKey(tt *testing.T) {
	t := wrapt.WrapT(tt)

	creds, arn := fromEnv()

	fuzz := []byte(creds.SecretAccessKey)
	// just swap two characters in the middle of the key
	a := fuzz[len(fuzz)/2]
	b := fuzz[len(fuzz)/2+2]
	fuzz[len(fuzz)/2] = b
	fuzz[len(fuzz)/2+2] = a

	creds.SecretAccessKey = string(fuzz)
	region := "us-east-1"

	cc, err := New(creds, region, arn, "", "", 0)
	t.R.NoError(err)

	cli, err := cc.Login(context.Background(), "new-session")
	t.R.Error(err)
	t.R.Equal(aws.Config{}, cli)
}

func Test_StsLogin_FailureBadArn(tt *testing.T) {
	t := wrapt.WrapT(tt)

	creds, arn := fromEnv()

	fuzz := []byte(creds.SecretAccessKey)
	// just swap two characters in the middle of the key
	a := fuzz[len(fuzz)-1]
	b := fuzz[len(fuzz)-2]
	fuzz[len(fuzz)-1] = b
	fuzz[len(fuzz)-2] = a

	creds.SecretAccessKey = string(fuzz)
	region := "us-east-1"

	cc, err := New(creds, region, arn, "", "", 0)
	t.R.NoError(err)

	cli, err := cc.Login(context.Background(), "new-session")
	t.R.Error(err)
	t.R.Equal(aws.Config{}, cli)
}
