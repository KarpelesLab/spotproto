package spotproto

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/KarpelesLab/cryptutil"
)

type ClientId struct {
	Type     byte   // 'k', 'c', etc
	ServerId string // can be empty
	Target   string // key, etc
}

func (c *ClientId) String() string {
	if c.ServerId == "" {
		return string([]byte{c.Type}) + "." + c.Target
	}
	return string([]byte{c.Type}) + "." + c.ServerId + "." + c.Target
}

func NewClientIdFromId(id *cryptutil.IDCard) *ClientId {
	h := sha256.Sum256(id.Self)
	return &ClientId{Type: 'k', Target: base64.RawURLEncoding.EncodeToString(h[:])}
}
