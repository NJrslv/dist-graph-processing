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

func AssertEq(testName string, result, expected interface{}) {
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

func CreateTestGraphs(path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Fatal("error creating graph file")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	graph := makeGraph()
	for i := 0; i < net.NumNodes; i++ {
		fmt.Fprintf(writer, strconv.Itoa(i)+"\n")
		writeGraphToFile(writer, graph)
	}
}

func makeGraph() net.Graph {
	g := net.Graph{}
	for r := 'a'; r <= 's'; r++ {
		g[net.Vertex(r)] = make([]net.Vertex, int('s')-int('a')+1)
		for rr := 'a'; rr <= 's'; rr++ {
			g[net.Vertex(r)][int(rr)-int('a')] = net.Vertex(rr)
		}
	}
	return g
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
