package spotproto

import "github.com/fxamacker/cbor/v2"

type HandshakeStart struct { // pid=1(S→C)
	ServerCode string `json:"srv"`
	ClientId   string `json:"cid"`
	Rand       []byte `json:"rnd"`
}

func (p *HandshakeStart) Bytes() []byte {
	res := []byte{Handshake} // v=0 p=1
	buf, err := cbor.Marshal(p)
	if err != nil {
		return nil
	}
	return append(res, buf...)
}

type HandshakeResponse struct {
	ID  []byte `json:"id"`
	Key []byte `json:"key"`
	Sig []byte `json:"sig"`
}

func (p *HandshakeResponse) Bytes() []byte {
	res := []byte{Handshake} // v=0 p=1
	buf, err := cbor.Marshal(p)
	if err != nil {
		return nil
	}
	return append(res, buf...)
}
