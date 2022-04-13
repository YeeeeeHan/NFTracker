package custerror

import "errors"

const (
	INVALID_SLUG = "INVALID_SLUG"
)

var InvalidSlugErr = errors.New(INVALID_SLUG)
