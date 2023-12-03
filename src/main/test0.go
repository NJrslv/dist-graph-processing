package main

import (
	"distgraphia/src/net"
	"distgraphia/src/test"
	"fmt"
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

func testCountConnComponentsSequentially(path string) string {
	// fictive chan, in order to create a node we need a chan
	done := make(chan struct{})
	nodesByName := make(map[string]*net.Node) // node name <-> *node
	nodes := make([]*net.Node, net.NumNodes)
	for i := range nodes {
		nodes[i] = net.MakeNode(strconv.Itoa(i), done)
		nodesByName[strconv.Itoa(i)] = nodes[i]
	}
	net.InitGraphs(path, nodesByName)

	start := time.Now()
	count := 0
	for _, node := range nodesByName {
		// This function takes (*Node, request.arguments)
		// For this case we can assume arguments are ""
		// because we do not need them
		components, _ := strconv.Atoi(net.CountConnectedComponentsMap(node, "").(string))
		count += components
	}
	test.Duration("Count Connected Components Sequentially", start)
	return strconv.Itoa(count)
}

func main() {
	test.DisableLogs()

	fmt.Printf(test.ColorRed+"Number of nodes in the system: %d\n", net.NumNodes)

	fmt.Print(test.ColorBlue + "--- Basic Tests ---\n")
	// test1
	test.AssertEq("Test Count Nodes", testCountNodes(), strconv.Itoa(net.NumNodes))
	// test2
	clientCount := 1000
	expectedReplies := make([]string, clientCount)
	for i := range expectedReplies {
		expectedReplies[i] = strconv.Itoa(net.NumNodes)
	}
	multiCount := testCountNodesMultiClient(clientCount)
	test.AssertEq("Test Count Nodes on Multiple Client Calls", multiCount, expectedReplies)

	fmt.Print(test.ColorBlue + "--- Benchmark && Tests ---\n")
	test.CreateTestGraphs(net.GraphPath)
	dist := testCountConnComponents(net.GraphPath)
	seq := testCountConnComponentsSequentially(net.GraphPath)
	test.AssertEq("Test Count Components", dist, seq)
}
