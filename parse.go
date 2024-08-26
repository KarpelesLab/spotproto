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

	buf = buf[1:]

	switch pkt {
	case PingPong:
		return Ping(buf), nil
	case Handshake:
		if isClient {
			// this is a HandshakeStart
			return parseObjCbor[HandshakeRequest](buf)
		} else {
			// server side, so this is a HandshakeResponse
			return parseObjCbor[HandshakeResponse](buf)
		}
	case InstantMsg:
		msg := &Message{}
		return msg, msg.UnmarshalBinary(buf)
	default:
		return nil, fmt.Errorf("failed to parse message: unknown packet id %x", pkt)
	}
}

func parseObjCbor[T any](buf []byte) (*T, error) {
	var res *T
	err := cbor.Unmarshal(buf, &res)
	return res, err
}

// VersionAndPacket returns the version & packet ID for a given code byte
func VersionAndPacket(v byte) (byte, byte) {
	return v >> 4 & 0xf, v & 0xf
}
