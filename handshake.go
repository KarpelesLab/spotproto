package spotproto

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"

	"github.com/KarpelesLab/cryptutil"
	"github.com/fxamacker/cbor/v2"
)

type HandshakeRequest struct { // pid=1(S→C)
	Ready      bool     `json:"rdy,omitempty"` // ready state. If true, handshake is complete
	ServerCode string   `json:"srv"`           // short name of server
	ClientId   string   `json:"cid"`           // name of connection
	Nonce      []byte   `json:"rnd"`           // random blob
	Groups     [][]byte `json:"grp"`
	raw        []byte
}

func (p *HandshakeRequest) Bytes() []byte {
	buf, err := cbor.Marshal(p)
	if err != nil {
		return nil
	}
	p.raw = buf
	return buf
}

// Respond generates a response to the current handshake start
func (p *HandshakeRequest) Respond(rawBuf []byte, s crypto.Signer) (*HandshakeResponse, error) {
	pubBin, err := x509.MarshalPKIXPublicKey(s.Public())
	if err != nil {
		return nil, err
	}
	if rawBuf == nil {
		if p.raw != nil {
			rawBuf = p.raw
		} else {
			// guess what the rawBuf was, but really we should have the original buffer
			rawBuf = p.Bytes()
		}
	}
	sig, err := cryptutil.Sign(rand.Reader, s, rawBuf, crypto.SHA256)
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
