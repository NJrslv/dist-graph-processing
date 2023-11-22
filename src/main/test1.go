package main

import (
	"distgraphia/src/net"
)

// TODO  runtime error: index out of range [-1] distgraphia/src/net.InitGraphs({0x1025475bf, 0xd}, 0x102564b80?)

func main() {
	//net.CreateTestGraphs("src/graph.txt")
	done := make(chan struct{})
	m := make(map[string]*net.Node)
	m["1"] = net.MakeNode("1", done)
	m["2"] = net.MakeNode("2", done)
	net.InitGraphs(net.GraphPath, m)
}
