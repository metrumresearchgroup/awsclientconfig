package awsclientconfig

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// ValidateCredentials validates that the correct combinations of settings
// are available in a combination that AWS would check, without needing
// to log in to confirm. Its returning of nil is not a predictor of whether
// keys are correct, only that they aren't incorrect in format.
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
