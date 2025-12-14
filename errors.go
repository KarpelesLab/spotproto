// Package spotproto implements the spot protocol for real-time messaging
// between clients and servers.
package spotproto

import "errors"

// Sentinel errors returned by the protocol parser.
var (
	// ErrEmptyBuf is returned when attempting to parse an empty buffer.
	ErrEmptyBuf = errors.New("empty buffer")
	// ErrInvalidVersion is returned when a packet has an unsupported protocol version.
	ErrInvalidVersion = errors.New("invalid packet version")
)
