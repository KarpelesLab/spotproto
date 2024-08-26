package spotproto

type Ping []byte

func (p Ping) Bytes() []byte {
	return append([]byte{PingPong}, p...)
}
