package spotproto

type Ping []byte

func (p Ping) Bytes() []byte {
	return []byte(p)
}
