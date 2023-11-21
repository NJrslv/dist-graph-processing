package net

import (
	"strconv"
)

func CountNodesMap(interface{}) interface{} {
	return "1"
}

func CountNodesReduce(replies interface{}) interface{} {
	count := 0
	for _, reply := range replies.([]ReplyMsg) {
		repInt, _ := strconv.Atoi(string(reply.Reply))
		count += repInt
	}
	return count
}

func intToStrBytes(i int) []byte {
	s := strconv.Itoa(i)
	return []byte(s)
}
