package spotproto

// Ping represents a ping/pong packet used for connection keep-alive.
// The payload is echoed back unchanged.
type Ping []byte

// Bytes returns the raw ping payload.
func (p Ping) Bytes() []byte {
	return []byte(p)
}
