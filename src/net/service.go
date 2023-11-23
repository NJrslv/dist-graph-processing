package net

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

//
// --- BroadCaster start ---
//

// BroadCaster broadcasts the message across the network,
// each node has its own BroadCaster service
type BroadCaster struct {
	Net  *Network // Network
	node *Node    // node that uses broadcaster
}

func MakeBroadCaster(n *Network, node *Node) *BroadCaster {
	return &BroadCaster{
		Net:  n,
		node: node,
	}
}

func (bc *BroadCaster) GatherQuorum() map[string]*Node {
	return bc.Net.GetNodes()
}

//
// --- BroadCaster end ---
//

//
// --- Method Invoker Start ---
//

type anyFunc func(*Node, interface{}) interface{}

type MethodInvoker struct {
	reflectionMap map[string]anyFunc // func name <-> func
	node          *Node
}

func MakeMethodInvoker(methods []string, node *Node) *MethodInvoker {
	mi := &MethodInvoker{
		reflectionMap: make(map[string]anyFunc),
		node:          node,
	}
	for _, methodName := range methods {
		mi.reflectionMap[methodName+"Map"] = mi.getFuncByName(methodName + "Map")
		mi.reflectionMap[methodName+"Reduce"] = mi.getFuncByName(methodName + "Reduce")
	}
	return mi
}

func (mi *MethodInvoker) InvokeMethod(methodName string, args interface{}) interface{} {
	if method, ok := mi.reflectionMap[methodName]; ok {
		a := method(mi.node, args)
		return a
	}
	log.Fatalf("Service.InvokeMethod(): Method '%s' not found", methodName)
	return ""
}

func (mi *MethodInvoker) getFuncByName(name string) anyFunc {
	switch name {
	case "CountNodesMap":
		return CountNodesMap
	case "CountNodesReduce":
		return CountNodesReduce
	case "CountConnectedComponentsMap":
		return CountConnectedComponentsMap
	case "CountConnectedComponentsReduce":
		return CountConnectedComponentsReduce
	default:
		log.Fatalf("Service:getFuncByName(): No func named %s", name)
		return nil
	}
}

//
// --- Method Invoker End ---
//

//
// --- Storage Start ---
//

/*
	Graph format:
	graph.txt:							For example Node1:
		nodeName1						c <- a <-> b -> g
		a:bc
		b:ag
		nodeName2
		a:c
		b:a
*/

type (
	Vertex  rune
	AdjList map[Vertex][]Vertex
	Graph   AdjList
)

// InitGraphs reads a file containing graphs
// and distributes them among the nodes in the system
func InitGraphs(path string, nodes map[string]*Node) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("File %s not found", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nodeName := ""
	g := make(Graph)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 1 {
			// flush previous graph
			if nodeName != "" {
				nodes[nodeName].g = &g
				g = make(Graph)
				nodeName = ""
			}
			// node name
			nodeName = string(line[0])
		} else {
			// Otherwise, it's an adjacency list
			vs := make([]Vertex, len(line)-2)
			var vertex Vertex
			isFirstV := true
			for i, c := range line {
				if isFirstV && c != ':' {
					vertex = Vertex(c)
					isFirstV = false
				} else if i >= 2 && c != ':' {
					vs[i-2] = Vertex(c)
				}
			}
			g[vertex] = vs
		}
	}
	// the last one
	if nodeName == "" {
		log.Fatalf("%s is empty", GraphPath)
	} else {
		nodes[nodeName].g = &g
	}
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

	for i := 0; i < NumNodes; i++ {
		fmt.Fprintf(writer, strconv.Itoa(i)+"\n")
		graph := Graph{'a': {'a'}}
		writeGraphToFile(writer, graph)
	}
}

func writeGraphToFile(writer *bufio.Writer, graph Graph) {
	for vertex, neighbors := range graph {
		fmt.Fprint(writer, string(vertex)+":"+verticesToString(neighbors)+"\n")
	}
}

func verticesToString(vertices []Vertex) string {
	result := ""
	for _, v := range vertices {
		result += string(v)
	}
	return result
}

//
// --- Storage End ---
//
