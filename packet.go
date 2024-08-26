package spotproto

type Packet interface {
	Bytes() []byte
}

const (
	PingPong   = 0x0
	Handshake  = 0x1
	InstantMsg = 0x2
)
