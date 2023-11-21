package main

import (
	"distgraphia/src/net"
	"fmt"
)

func main() {
	net1 := net.MakeNetwork("n1")
	defer net1.Cleanup()

	cl1 := net.MakeClient("c1")
	cl1.ConnectTo(net1)

	reply := ""
	cl1.Call("n1", "CountNodes", "", &reply)
	defer fmt.Printf("Reply for %s is %s", cl1.GetName(), reply)
}
