package spotproto

import (
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

func Parse(buf []byte, isClient bool) (Packet, error) {
	if len(buf) == 0 {
		return nil, ErrEmptyBuf
	}
	vers, pkt := VersionAndPacket(buf[0])
	if vers != 0 {
		return nil, ErrInvalidVersion
	}

	switch pkt {
	case PingPong:
		return Ping(buf[1:]), nil
	case Handshake:
		if isClient {
			// this is a HandshakeStart
			return parseObj[HandshakeStart](buf[1:])
		} else {
			// server side, so this is a HandshakeResponse
			return parseObj[HandshakeResponse](buf[1:])
		}
	case InstantMsg:
		msg := &Message{}
		return msg, msg.UnmarshalBinary(buf[1:])
	default:
		return nil, fmt.Errorf("failed to parse message: unknown packet id %x", pkt)
	}
}

func parseObj[T any](buf []byte) (*T, error) {
	var res *T
	err := cbor.Unmarshal(buf, &res)
	return res, err
}

// VersionAndPacket returns the version & packet ID for a given code byte
func VersionAndPacket(v byte) (byte, byte) {
	return v >> 4 & 0xf, v & 0xf
}
