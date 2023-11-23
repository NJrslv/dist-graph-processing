package net

import (
	"bufio"
	"log"
	"os"
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
	Nodes are named as chars(Runes)
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
	for scanner.Scan() {
		line := scanner.Text()
		if isNodeName(line) {
			nodeName = line
			nodes[nodeName].g = make(Graph)
		} else {
			v, vs := parseLine(line)
			nodes[nodeName].g[v] = vs
		}

	}
}

// parseLine parses line with format: vertex1:vertex2vertex3...
// vertex#i is a rune(char)
func parseLine(line string) (Vertex, []Vertex) {
	var v Vertex
	vs := make([]Vertex, len(line)-2)
	for i, c := range line {
		// line[1] is ':'
		if i == 0 {
			v = Vertex(c)
		} else if i >= 2 {
			vs[i-2] = Vertex(c)
		}
	}
	return v, vs
}

func isNodeName(line string) bool {
	for _, c := range line {
		if c == ':' {
			return false
		}
	}
	return true
}

//
// --- Storage End ---
//
