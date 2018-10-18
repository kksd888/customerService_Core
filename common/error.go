package common

import "errors"

var (
	HEADER_TYPE_REQUIRED   = errors.New("request header content-type must be application/json")
	AUTHORIZATION_REQUIRED = errors.New("API Authorization required")
	AUTHORIZATION_FAILED   = errors.New("API Authorization failed")
)
