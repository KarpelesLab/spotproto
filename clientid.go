package spotproto

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/BottleFmt/gobottle"
)

// ClientId represents a client identifier in the protocol.
// The format is "Type.ServerId.Target" or "Type.Target" if ServerId is empty.
type ClientId struct {
	Type     byte   // Type is the identifier type ('k' for key-based, 'c' for connection, etc.)
	ServerId string // ServerId is the server identifier (can be empty for global IDs)
	Target   string // Target is the specific identifier value (key hash, connection name, etc.)
}

// String returns the client ID in its canonical string format.
func (c *ClientId) String() string {
	if c.ServerId == "" {
		return string([]byte{c.Type}) + "." + c.Target
	}
	return string([]byte{c.Type}) + "." + c.ServerId + "." + c.Target
}

// NewClientIdFromId creates a new key-based ClientId from a cryptutil IDCard.
// The target is derived from the SHA-256 hash of the IDCard's public key.
func NewClientIdFromId(id *gobottle.IDCard) *ClientId {
	h := sha256.Sum256(id.Self)
	return &ClientId{Type: 'k', Target: base64.RawURLEncoding.EncodeToString(h[:])}
}
