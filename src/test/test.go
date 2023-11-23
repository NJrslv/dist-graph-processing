package test

import (
	"bufio"
	"distgraphia/src/net"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

const (
	ColorBlue      = "\033[0;34m"
	ColorReset     = "\033[0m"
	ColorRed       = "\033[0;31m"
	ColorGreen     = "\033[0;32m"
	AssertFormat   = "%s%s %s\n"
	DurationFormat = ColorReset + "Execution time of %s: " + ColorBlue + "%dmcs" + "\n"
)

func Duration(function string, start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf(DurationFormat, function, elapsed.Microseconds())
}

func Assert(testName string, result, expected interface{}) {
	var color, assert string
	if reflect.DeepEqual(result, expected) {
		color = ColorGreen
		assert = "success"
	} else {
		color = ColorRed
		assert = "fail"
	}
	fmt.Printf(AssertFormat, color, testName, assert)
}

func DisableLogs() {
	log.SetFlags(0)
	log.SetOutput(io.Discard)
}

func CountComponentsSequentially(path string) string {
	// fictive chan, in order to create a node we need a chan
	done := make(chan struct{})
	nodesByName := make(map[string]*net.Node) // node name <-> *node
	nodes := make([]*net.Node, net.NumNodes)
	for i := range nodes {
		nodes[i] = net.MakeNode(strconv.Itoa(i), done)
		nodesByName[strconv.Itoa(i)] = nodes[i]
	}
	net.InitGraphs(path, nodesByName)

	count := 0
	for _, node := range nodesByName {
		// This function takes (*Node, request.arguments)
		// For this case we can assume arguments are ""
		// because we do not need them
		components, _ := strconv.Atoi(net.CountConnectedComponentsMap(node, "").(string))
		count += components
	}

	return strconv.Itoa(count)
}

func CreateTestGraphs(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("error creating graph file")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < net.NumNodes; i++ {
		fmt.Fprintf(writer, strconv.Itoa(i)+"\n")
		graph := net.Graph{'a': {'b'}, 'c': {'d'}}
		writeGraphToFile(writer, graph)
	}
}

func writeGraphToFile(writer *bufio.Writer, graph net.Graph) {
	for vertex, neighbors := range graph {
		fmt.Fprint(writer, string(vertex)+":"+verticesToString(neighbors)+"\n")
	}
}

func verticesToString(vertices []net.Vertex) string {
	result := ""
	for _, v := range vertices {
		result += string(v)
	}
	return result
}
