package spotproto

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"

	"github.com/fxamacker/cbor/v2"
)

type HandshakeStart struct { // pid=1(S→C)
	ServerCode string `json:"srv"`
	ClientId   string `json:"cid"`
	Rand       []byte `json:"rnd"`
}

func (p *HandshakeStart) Bytes() []byte {
	buf, err := cbor.Marshal(p)
	if err != nil {
		return nil
	}
	return buf
}

// Respond generates a response to the current handshake start
func (p *HandshakeStart) Respond(rawBuf []byte, s crypto.Signer) (*HandshakeResponse, error) {
	pubBin, err := x509.MarshalPKIXPublicKey(s.Public())
	if err != nil {
		return nil, err
	}
	if rawBuf == nil {
		// guess what the rawBuf was, but really we should have the original buffer
		rawBuf = append([]byte{Handshake}, p.Bytes()...)
	}
	sum := sha256.Sum256(rawBuf)
	sig, err := s.Sign(rand.Reader, sum[:], crypto.SHA256)
	if err != nil {
		return nil, err
	}
	res := &HandshakeResponse{
		Key: pubBin,
		Sig: sig,
	}
	return res, nil
}

type HandshakeResponse struct { // pid=1(C→S)
	ID  []byte `json:"id"`
	Key []byte `json:"key"`
	Sig []byte `json:"sig"`
}

func (p *HandshakeResponse) Bytes() []byte {
	buf, err := cbor.Marshal(p)
	if err != nil {
		return nil
	}
	return buf
}
