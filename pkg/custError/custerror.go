package custError

import "errors"

const (
	InvalidSlug = "INVALID_SLUG"
)

var InvalidSlugErr = errors.New(InvalidSlug)
