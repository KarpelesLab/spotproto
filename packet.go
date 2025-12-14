package spotproto

// Packet is the interface implemented by all protocol packet types.
type Packet interface {
	// Bytes serializes the packet to its wire format.
	Bytes() []byte
}

// Packet type identifiers encoded in the lower 4 bits of the first byte.
const (
	PingPong   = 0x0 // Ping/pong keep-alive packet
	Handshake  = 0x1 // Handshake request or response
	InstantMsg = 0x2 // Instant message packet
)
