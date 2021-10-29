package awsclientconfig_test

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"

	. "github.com/metrumresearchgroup/awsclientconfig"
)

func TestValidateCredentials(t *testing.T) {
	tests := []struct {
		name        string
		credentials aws.Credentials
		wantErr     bool
		expErr      error
	}{
		{
			name: "no errors ASIA",
			credentials: aws.Credentials{
				AccessKeyID:     "ASIAASDFASDFASDF",
				SecretAccessKey: "ASDFASDFASDFASDFASDFASDFASDFASDFASDFASDF",
				SessionToken:    "FwoGZXIvYXdzEC0aDNTmPAer6krz7/fnGCLnAZbl8nF7f6nPc3Jg4CM3JEsPsuEf6Zlxl7x8kBKtf9OC+b6GBscsRoQK4c0u8yeP8LnGjnK9KRKZa7YQFNz+6TSkYjkI44bZ2iPjXa/UGmj0dUtl19u58BJPM1u77w3W8mC2Luz4yApWSmn17bggpHCB/oS8hiY+Y/Iry5GhnBaVGqh3AeclcRYYqpWcACzYqws6VwdHg3uWixGMQqSBPM9YCyiacnyOUP/0YJQuSqX5pCWKyP9saNAZnnrFDHf0hcbpmm1/BG2IQdpJ2gLMOjxXro0PJiJA2xOA0l1MQ0/LhoykznGL8Sj3+t+BBjIrd8x0f3jauroAHvbxcY/ol0TP3MFnZBLtMsS/cMn8IAyQ7P6p8DVAstiCfw==",
			},
			wantErr: false,
		},
		{
			name: "no errors AKIA",
			credentials: aws.Credentials{
				AccessKeyID:     "AKIAASDFASDFASDF",
				SecretAccessKey: "ASDFASDFASDFASDFASDFASDFASDFASDFASDFASDF",
				SessionToken:    "",
			},
			wantErr: false,
		},
		{
			name: "short AKID",
			credentials: aws.Credentials{
				AccessKeyID:     "ASIAASDFASDFASD",
				SecretAccessKey: "ASDFASDFASDFASDFASDFASDFASDFASDFASDFASDF",
				SessionToken:    "FwoGZXIvYXdzEC0aDNTmPAer6krz7/fnGCLnAZbl8nF7f6nPc3Jg4CM3JEsPsuEf6Zlxl7x8kBKtf9OC+b6GBscsRoQK4c0u8yeP8LnGjnK9KRKZa7YQFNz+6TSkYjkI44bZ2iPjXa/UGmj0dUtl19u58BJPM1u77w3W8mC2Luz4yApWSmn17bggpHCB/oS8hiY+Y/Iry5GhnBaVGqh3AeclcRYYqpWcACzYqws6VwdHg3uWixGMQqSBPM9YCyiacnyOUP/0YJQuSqX5pCWKyP9saNAZnnrFDHf0hcbpmm1/BG2IQdpJ2gLMOjxXro0PJiJA2xOA0l1MQ0/LhoykznGL8Sj3+t+BBjIrd8x0f3jauroAHvbxcY/ol0TP3MFnZBLtMsS/cMn8IAyQ7P6p8DVAstiCfw==",
			},
			wantErr: true,
			expErr:  ErrAKIDMinLen,
		},
		{
			name: "long AKID",
			credentials: aws.Credentials{
				AccessKeyID:     "FwoGZXIvYXdzEC0aDNTmPAer6krz7/fnGCLnAZbl8nF7f6nPc3Jg4CM3JEsPsuEf6Zlxl7x8kBKtf9OC+b6GBscsRoQK4c0u8yeP8LnGjnK9KRKZa7YQFNz+6TSkYjkI44bZ2iPjXa/UGmj0dUtl19u58BJPM1u77w3W8mC2Luz4yApWSmn17bggpHCB/oS8hiY+Y/Iry5GhnBaVGqh3AeclcRYYqpWcACzYqws6VwdHg3uWixGMQqSBPM9YCyiacnyOUP/0YJQuSqX5pCWKyP9saNAZnnrFDHf0hcbpmm1/BG2IQdpJ2gLMOjxXro0PJiJA2xOA0l1MQ0/LhoykznGL8Sj3+t+BBjIrd8x0f3jauroAHvbxcY/ol0TP3MFnZBLtMsS/cMn8IAyQ7P6p8DVAstiCfw==",
				SecretAccessKey: "ASDFASDFASDFASDFASDFASDFASDFASDFASDFASDF",
				SessionToken:    "FwoGZXIvYXdzEC0aDNTmPAer6krz7/fnGCLnAZbl8nF7f6nPc3Jg4CM3JEsPsuEf6Zlxl7x8kBKtf9OC+b6GBscsRoQK4c0u8yeP8LnGjnK9KRKZa7YQFNz+6TSkYjkI44bZ2iPjXa/UGmj0dUtl19u58BJPM1u77w3W8mC2Luz4yApWSmn17bggpHCB/oS8hiY+Y/Iry5GhnBaVGqh3AeclcRYYqpWcACzYqws6VwdHg3uWixGMQqSBPM9YCyiacnyOUP/0YJQuSqX5pCWKyP9saNAZnnrFDHf0hcbpmm1/BG2IQdpJ2gLMOjxXro0PJiJA2xOA0l1MQ0/LhoykznGL8Sj3+t+BBjIrd8x0f3jauroAHvbxcY/ol0TP3MFnZBLtMsS/cMn8IAyQ7P6p8DVAstiCfw==",
			},
			wantErr: true,
			expErr:  ErrAKIDMaxLen,
		},
		{
			name: "empty AKID",
			credentials: aws.Credentials{
				AccessKeyID:     "",
				SecretAccessKey: "ASDFASDFASDFASDFASDFASDFASDFASDFASDFASDF",
				SessionToken:    "FwoGZXIvYXdzEC0aDNTmPAer6krz7/fnGCLnAZbl8nF7f6nPc3Jg4CM3JEsPsuEf6Zlxl7x8kBKtf9OC+b6GBscsRoQK4c0u8yeP8LnGjnK9KRKZa7YQFNz+6TSkYjkI44bZ2iPjXa/UGmj0dUtl19u58BJPM1u77w3W8mC2Luz4yApWSmn17bggpHCB/oS8hiY+Y/Iry5GhnBaVGqh3AeclcRYYqpWcACzYqws6VwdHg3uWixGMQqSBPM9YCyiacnyOUP/0YJQuSqX5pCWKyP9saNAZnnrFDHf0hcbpmm1/BG2IQdpJ2gLMOjxXro0PJiJA2xOA0l1MQ0/LhoykznGL8Sj3+t+BBjIrd8x0f3jauroAHvbxcY/ol0TP3MFnZBLtMsS/cMn8IAyQ7P6p8DVAstiCfw==",
			},
			wantErr: true,
			expErr:  ErrAKIDRequired,
		},
		{
			name: "AKID does not start with ASIA or AKIA",
			credentials: aws.Credentials{
				AccessKeyID:     "AFIAASDFASDFASDF",
				SecretAccessKey: "ASDFASDFASDFASDFASDFASDFASDFASDFASDFASDF",
				SessionToken:    "FwoGZXIvYXdzEC0aDNTmPAer6krz7/fnGCLnAZbl8nF7f6nPc3Jg4CM3JEsPsuEf6Zlxl7x8kBKtf9OC+b6GBscsRoQK4c0u8yeP8LnGjnK9KRKZa7YQFNz+6TSkYjkI44bZ2iPjXa/UGmj0dUtl19u58BJPM1u77w3W8mC2Luz4yApWSmn17bggpHCB/oS8hiY+Y/Iry5GhnBaVGqh3AeclcRYYqpWcACzYqws6VwdHg3uWixGMQqSBPM9YCyiacnyOUP/0YJQuSqX5pCWKyP9saNAZnnrFDHf0hcbpmm1/BG2IQdpJ2gLMOjxXro0PJiJA2xOA0l1MQ0/LhoykznGL8Sj3+t+BBjIrd8x0f3jauroAHvbxcY/ol0TP3MFnZBLtMsS/cMn8IAyQ7P6p8DVAstiCfw==",
			},
			wantErr: true,
			expErr:  ErrAKIDBadPrefix,
		},
		{
			name: "empty SAK",
			credentials: aws.Credentials{
				AccessKeyID:     "ASIAASDFASDFASDF",
				SecretAccessKey: "",
				SessionToken:    "FwoGZXIvYXdzEC0aDNTmPAer6krz7/fnGCLnAZbl8nF7f6nPc3Jg4CM3JEsPsuEf6Zlxl7x8kBKtf9OC+b6GBscsRoQK4c0u8yeP8LnGjnK9KRKZa7YQFNz+6TSkYjkI44bZ2iPjXa/UGmj0dUtl19u58BJPM1u77w3W8mC2Luz4yApWSmn17bggpHCB/oS8hiY+Y/Iry5GhnBaVGqh3AeclcRYYqpWcACzYqws6VwdHg3uWixGMQqSBPM9YCyiacnyOUP/0YJQuSqX5pCWKyP9saNAZnnrFDHf0hcbpmm1/BG2IQdpJ2gLMOjxXro0PJiJA2xOA0l1MQ0/LhoykznGL8Sj3+t+BBjIrd8x0f3jauroAHvbxcY/ol0TP3MFnZBLtMsS/cMn8IAyQ7P6p8DVAstiCfw==",
			},
			wantErr: true,
			expErr:  ErrSAKRequired,
		},
		{
			name: "ASIA missing ST",
			credentials: aws.Credentials{
				AccessKeyID:     "ASIAASDFASDFASDF",
				SecretAccessKey: "ASDFASDFASDFASDFASDFASDFASDFASDFASDFASDF",
				SessionToken:    "",
			},
			wantErr: true,
			expErr:  ErrSTRequired,
		},
		{
			name: "invalid AKID",
			credentials: aws.Credentials{
				AccessKeyID:     "ASIAASDFASDFASD=",
				SecretAccessKey: "ASDFASDFASDFASDFASDFASDFASDFASDFASDFASDF",
				SessionToken:    "FwoGZXIvYXdzEC0aDNTmPAer6krz7/fnGCLnAZbl8nF7f6nPc3Jg4CM3JEsPsuEf6Zlxl7x8kBKtf9OC+b6GBscsRoQK4c0u8yeP8LnGjnK9KRKZa7YQFNz+6TSkYjkI44bZ2iPjXa/UGmj0dUtl19u58BJPM1u77w3W8mC2Luz4yApWSmn17bggpHCB/oS8hiY+Y/Iry5GhnBaVGqh3AeclcRYYqpWcACzYqws6VwdHg3uWixGMQqSBPM9YCyiacnyOUP/0YJQuSqX5pCWKyP9saNAZnnrFDHf0hcbpmm1/BG2IQdpJ2gLMOjxXro0PJiJA2xOA0l1MQ0/LhoykznGL8Sj3+t+BBjIrd8x0f3jauroAHvbxcY/ol0TP3MFnZBLtMsS/cMn8IAyQ7P6p8DVAstiCfw==",
			},
			wantErr: true,
			expErr:  ErrAKIDInvalidChar,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCredentials(tt.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !errors.Is(err, tt.expErr) {
				t.Errorf("Validate() error = %v, expErr %v", err, tt.expErr)
			}
		})
	}
}
