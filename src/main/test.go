package main

import (
	"distgraphia/src/net"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"sync"
)

const (
	colorRed   = "\033[0;31m"
	colorGreen = "\033[0;32m"
	format     = "%s%s %s\n"
)

func assert(testName string, result, expected interface{}) {
	var color, assert string
	if reflect.DeepEqual(result, expected) {
		color = colorGreen
		assert = "success"
	} else {
		color = colorRed
		assert = "fail"
	}
	fmt.Printf(format, color, testName, assert)
}

func disableLogs() {
	log.SetFlags(0)
	log.SetOutput(io.Discard)
}

func testCountNodes() string {
	net1 := net.MakeNetwork("n1")
	defer net1.Cleanup()

	cl1 := net.MakeClient("c1")
	cl1.ConnectTo(net1)

	reply := ""
	cl1.Call("n1", "CountNodes", "", &reply)
	return reply
}

func testCountNodesMultClient(clientCount int) []string {
	net2 := net.MakeNetwork("n2")
	defer net2.Cleanup()

	clients := make([]net.Client, clientCount)
	replies := make([]string, clientCount)
	var wg sync.WaitGroup

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
	return replies
}

// TODO multiple clients fail, protect node's data

func main() {
	disableLogs()
	// test1
	assert("Test CountNodes", testCountNodes(), strconv.Itoa(net.NumNodes))

	// test2
	clientCount := 10000
	expectedReplies := make([]string, clientCount)
	for i := range expectedReplies {
		expectedReplies[i] = strconv.Itoa(net.NumNodes)
	}
	res := testCountNodesMultClient(clientCount)
	for i := range res {
		if res[i] != strconv.Itoa(net.NumNodes) {
			fmt.Print(res[i])
		}
	}
	//assert("Test CountNodes on Multiple Client Calls", testCountNodesMultClient(clientCount), expectedReplies)
}
