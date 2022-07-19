package awsclientconfig_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/metrumresearchgroup/wrapt"

	. "github.com/metrumresearchgroup/awsclientconfig"
)

func TestNew(tt *testing.T) {
	type args struct {
		credentials aws.Credentials
		region      string
		arn         string
	}
	tests := []struct {
		name    string
		args    args
		want    *ClientConfig
		wantErr bool
	}{
		{
			name: "successful AKIA with defaults",
			args: args{
				credentials: aws.Credentials{
					AccessKeyID:     "AKIAASDFASDFASDF",
					SecretAccessKey: "Secretkey",
				},
			},
			want: &ClientConfig{
				AccessKey: "AKIAASDFASDFASDF",
				SecretKey: "Secretkey",
				Refresh:   300000000000,
			},
		},
		{
			name: "successful ASIA with defaults",
			args: args{
				credentials: aws.Credentials{
					AccessKeyID:     "ASIAASDFASDFASDF",
					SecretAccessKey: "Secretkey",
					SessionToken:    "Sessiontoken",
				},
				region: "us-east-2",
			},
			want: &ClientConfig{
				AccessKey:    "ASIAASDFASDFASDF",
				SecretKey:    "Secretkey",
				SessionToken: "Sessiontoken",
				Region:       "us-east-2",
				Refresh:      300000000000,
			},
		},
		{
			name:    "failure to prove validation runs",
			args:    args{},
			want:    nil,
			wantErr: true,
		},
	}
	for _, test := range tests {
		tt.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(tt)
			got, err := New(test.args.credentials, test.args.region, test.args.arn, "", "", 0)

			t.R.WantError(test.wantErr, err)
			t.R.Equal(got, test.want)
		})
	}
}
