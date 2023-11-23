package main

import (
	"distgraphia/src/net"
	"distgraphia/src/test"
	"strconv"
	"sync"
	"time"
)

func testCountNodes() string {
	net1 := net.MakeNetwork("n1")
	defer net1.Cleanup()

	cl1 := net.MakeClient("c1")
	cl1.ConnectTo(net1)

	start := time.Now()
	reply := ""
	cl1.Call("n1", "CountNodes", "", &reply)
	test.Duration("Count Nodes", start)
	return reply
}

func testCountNodesMultiClient(clientCount int) []string {
	net2 := net.MakeNetwork("n2")
	defer net2.Cleanup()

	clients := make([]net.Client, clientCount)
	replies := make([]string, clientCount)
	var wg sync.WaitGroup

	start := time.Now()
	for i := range clients {
		clients[i] = *net.MakeClient("cl" + strconv.Itoa(i))
		clients[i].ConnectTo(net2)

		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			reply := ""
			clients[i].Call("n2", "CountNodes", "", &reply)
			replies[i] = reply
		}(i)
	}

	wg.Wait()
	defer test.Duration("Count Nodes Multi Client", start)
	return replies
}

func testCountConnComponents(path string) string {
	test.CreateTestGraphs(path) // with 2 components in each node
	net3 := net.MakeNetwork("n3")
	defer net3.Cleanup()

	net.InitGraphs(net.GraphPath, net3.GetNodes())

	cl1 := net.MakeClient("cl1")
	cl1.ConnectTo(net3)

	start := time.Now()
	reply := ""
	cl1.Call("n3", "CountConnectedComponents", "", &reply)
	test.Duration("Count Connected Components", start)
	return reply
}

func main() {
	test.DisableLogs()
	// test1
	test.Assert("Test Count Nodes", testCountNodes(), strconv.Itoa(net.NumNodes))

	// test2
	clientCount := 1000
	expectedReplies := make([]string, clientCount)
	for i := range expectedReplies {
		expectedReplies[i] = strconv.Itoa(net.NumNodes)
	}
	multiCount := testCountNodesMultiClient(clientCount)
	test.Assert("Test Count Nodes on Multiple Client Calls", multiCount, expectedReplies)

	// test3
	dist := testCountConnComponents(net.GraphPath)
	seq := test.CountComponentsSequentially(net.GraphPath)
	test.Assert("Test Count Components", dist, seq)
}
