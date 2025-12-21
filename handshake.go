package spotproto

import (
	"crypto"
	"crypto/rand"
	"crypto/x509"

	"github.com/BottleFmt/gobottle"
	"github.com/fxamacker/cbor/v2"
)

// HandshakeRequest is sent from server to client to initiate the handshake.
// It contains server identification and a nonce for authentication.
type HandshakeRequest struct {
	Ready      bool     `json:"rdy,omitempty"` // Ready indicates handshake completion when true
	ServerCode string   `json:"srv"`           // ServerCode is the short name of the server
	ClientId   string   `json:"cid"`           // ClientId is the assigned connection identifier
	Nonce      []byte   `json:"rnd"`           // Nonce is a random blob for authentication
	Groups     [][]byte `json:"grp"`           // Groups the client belongs to
	raw        []byte
}

// Bytes serializes the handshake request to CBOR format.
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
	sig, err := gobottle.Sign(rand.Reader, s, rawBuf, crypto.SHA256)
	if err != nil {
		return nil, err
	}
	res := &HandshakeResponse{
		Key: pubBin,
		Sig: sig,
	}
	return res, nil
}

// HandshakeResponse is sent from client to server in response to a HandshakeRequest.
// It contains the client's public key and a signature proving key ownership.
type HandshakeResponse struct {
	ID  []byte `json:"id"`  // ID is an optional client identifier
	Key []byte `json:"key"` // Key is the client's PKIX-encoded public key
	Sig []byte `json:"sig"` // Sig is the signature over the request nonce
}

// Bytes serializes the handshake response to CBOR format.
func (p *HandshakeResponse) Bytes() []byte {
	buf, err := cbor.Marshal(p)
	if err != nil {
		return nil
	}
	return buf
}
