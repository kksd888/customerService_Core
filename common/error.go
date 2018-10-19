package common

import "errors"

var (
	AUTHORIZATION_REQUIRED = errors.New("API Authorization required")
	AUTHORIZATION_FAILED   = errors.New("API Authorization failed")
)
