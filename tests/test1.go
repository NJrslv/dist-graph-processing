package tests

import "distgraphia/core/net"

func main() {
	net1 := net.MakeNetwork("n1")
	defer net1.Cleanup()

	cl1 := net.MakeClient("c1")
	cl1.ConnectTo(net1)

	reply := ""
	cl1.Call("n1", "countNodes", "1", &reply)

}
