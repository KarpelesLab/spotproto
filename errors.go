package spotproto

import "errors"

var (
	ErrEmptyBuf       = errors.New("empty buffer")
	ErrInvalidVersion = errors.New("invalid packet version")
)
