package spotproto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	MsgFlagResponse  = 1 << iota // message is a response (and must not cause any further response to be generated)
	MsgFlagError                 // message body is an error string
	MsgFlagNotBottle             // message body is not a bottle
)

type Message struct {
	MessageID [16]byte
	Flags     uint64
	Recipient string
	Sender    string
	Body      []byte
}

func (msg *Message) Bytes() []byte {
	buf := msg.MessageID[:]
	buf = binary.AppendUvarint(buf, msg.Flags)
	buf = binary.AppendUvarint(buf, uint64(len(msg.Recipient)))
	buf = append(buf, msg.Recipient...)
	buf = binary.AppendUvarint(buf, uint64(len(msg.Sender)))
	buf = append(buf, msg.Sender...)
	buf = append(buf, msg.Body...)
	return buf
}

func (msg *Message) UnmarshalBinary(b []byte) error {
	return msg.ReadFrom(bytes.NewReader(b))
}

// IsEncrypted returns if the message must be encrypted to be sent. This method can be used in handlers to
// ensure only encrypted messages are being handled.
func (msg *Message) IsEncrypted() bool {
	return msg.Flags&MsgFlagNotBottle == 0
}

func (msg *Message) ReadFrom(r io.Reader) error {
	buf := bufio.NewReader(r)
	_, err := io.ReadFull(buf, msg.MessageID[:])
	if err != nil {
		return err
	}
	msg.Flags, err = binary.ReadUvarint(buf)
	if err != nil {
		return err
	}
	msg.Recipient, err = readName(buf)
	if err != nil {
		return err
	}
	msg.Sender, err = readName(buf)
	if err != nil {
		return err
	}
	msg.Body, err = io.ReadAll(buf)
	if err != nil {
		return err
	}
	return nil
}

func readName(r *bufio.Reader) (string, error) {
	ln, err := binary.ReadUvarint(r)
	if err != nil {
		return "", err
	}
	if ln > 256 {
		return "", errors.New("cannot read name from packet: too long")
	}
	buf := make([]byte, ln)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
