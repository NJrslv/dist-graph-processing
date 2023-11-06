package tests

import "distgraphia/core/net"

func main() {
	cl1 := net.MakeClient("c1")
	net1 := net.MakeNetwork("n1")
	cl1.ConnectTo(net1)

}
