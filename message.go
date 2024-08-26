package spotproto

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
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
