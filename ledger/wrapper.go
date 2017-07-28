package ledger

import (
	"encoding/binary"
)

var codec = binary.BigEndian

// WrapCommandAPDU turns the command into a sequence of 64 byte packets
func WrapCommandAPDU(channel uint16, command []byte, packetSize int, ble bool) []byte {
	if packetSize < 3 {
		panic("packet size must be at least 3")
	}

	var sequenceIdx uint16
	var offset, extraHeaderSize, blockSize int
	var result = make([]byte, 64)
	var buf = result

	if !ble {
		codec.PutUint16(buf, channel)
		extraHeaderSize = 2
		buf = buf[2:]
	}

	buf[0] = 0x05
	codec.PutUint16(buf[1:], sequenceIdx)
	codec.PutUint16(buf[3:], uint16(len(command)))
	sequenceIdx++
	buf = buf[5:]

	blockSize = packetSize - 5 - extraHeaderSize
	copy(buf, command)
	offset += blockSize

	for offset < len(command) {
		// TODO: optimize this
		end := len(result)
		result = append(result, make([]byte, 64)...)
		buf = result[end:]
		if !ble {
			codec.PutUint16(buf, channel)
			buf = buf[2:]
		}
		buf[0] = 0x05
		codec.PutUint16(buf[1:], sequenceIdx)
		sequenceIdx++
		buf = buf[3:]

		blockSize = packetSize - 3 - extraHeaderSize
		copy(buf, command[offset:])
		offset += blockSize
	}

	return result
}
