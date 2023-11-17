package net

import (
	"encoding/binary"
)

func CountNodesMap(interface{}) interface{} {
	return "1"
}

func CountNodesReduce(replies interface{}) interface{} {
	count := 0
	for _, reply := range replies.([]ReplyMsg) {
		count += int(binary.LittleEndian.Uint64(reply.Reply))
	}
	return ReplyMsg{Ok: true, Reply: intToBytes(count)}
}

func intToBytes(i int) []byte {
	b := make([]byte, 4) // Assuming int is 4 bytes
	binary.LittleEndian.PutUint32(b, uint32(i))
	return b
}
