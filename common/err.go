package common

import "errors"

var (
	ErrInvalidFormat = errors.New("invalid format")
	ErrNoTracklist   = errors.New("no tracklist found")
)
