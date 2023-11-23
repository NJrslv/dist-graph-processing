package net

import (
	"strconv"
)

func CountNodesMap(_ *Node, _ interface{}) interface{} {
	return "1"
}

func CountNodesReduce(_ *Node, replies interface{}) interface{} {
	count := 0
	for _, reply := range replies.([]ReplyMsg) {
		repInt, _ := strconv.Atoi(string(reply.Reply))
		count += repInt
	}
	return count
}

func CountConnectedComponentsMap(node *Node, _ interface{}) interface{} {
	graph := node.g
	visited := make(map[Vertex]bool)
	count := 0

	for vertex := range graph {
		if !visited[vertex] {
			dfs(graph, vertex, visited)
			count++
		}
	}

	return strconv.Itoa(count)
}

func dfs(graph Graph, vertex Vertex, visited map[Vertex]bool) {
	visited[vertex] = true

	for _, neighbor := range (graph)[vertex] {
		if !visited[neighbor] {
			dfs(graph, neighbor, visited)
		}
	}
}

func CountConnectedComponentsReduce(_ *Node, replies interface{}) interface{} {
	count := 0
	for _, reply := range replies.([]ReplyMsg) {
		repInt, _ := strconv.Atoi(string(reply.Reply))
		count += repInt
	}
	return count
}

func intToStrBytes(i int) []byte {
	s := strconv.Itoa(i)
	return []byte(s)
}
