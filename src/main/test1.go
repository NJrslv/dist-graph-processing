package main

import (
	net2 "distgraphia/src/net"
)

func main() {
	net1 := net2.MakeNetwork("n1")
	defer net1.Cleanup()

	cl1 := net2.MakeClient("c1")
	cl1.ConnectTo(net1)

	reply := ""
	cl1.Call("n1", "CountNodes", "", &reply)

}
