package common

import "encoding/binary"

const (
	u16ByteSize = 2
	u32ByteSize = 4
)

func serializeString(value string, out []byte) []byte {
	b := []byte(value)
	out = serializeU16(uint16(len(b)), out)
	out = append(out, b...)
	return out
}

func serializeU8(value uint8, out []byte) []byte {
	return append(out, byte(value))
}

func serializeU16(value uint16, out []byte) []byte {
	tmp := make([]byte, u16ByteSize)
	binary.BigEndian.PutUint16(tmp, value)
	return append(out, tmp...)
}

func serializeU32(value uint32, out []byte) []byte {
	tmp := make([]byte, u32ByteSize)
	binary.BigEndian.PutUint32(tmp, value)
	return append(out, tmp...)
}
