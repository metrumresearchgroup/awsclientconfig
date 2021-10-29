package awsclientconfig

import (
	"errors"
	"fmt"
)

var (
	ErrAKIDRequired    = errors.New("access key id required")
	ErrAKIDMinLen      = errors.New("access key id is less than 16 characters long")
	ErrAKIDMaxLen      = errors.New("access key id greater than 128 characters long")
	ErrAKIDBadPrefix   = errors.New("access key id must begin with either AKIA or ASIA")
	ErrAKIDInvalidChar = fmt.Errorf("access key id must match the pattern %s", AccessKeyPattern)
	ErrSAKRequired     = errors.New("secret access key required")
	ErrSTRequired      = errors.New("session token required for ASIA access key id")
)
