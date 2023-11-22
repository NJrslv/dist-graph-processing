package net

import (
	"strconv"
)

func CountNodesMap(interface{}) interface{} {
	return "1"
}

func CountNodesReduce(replies interface{}) interface{} {
	count := 0
	for _, reply := range replies.([]ReplyMsg) {
		repInt, _ := strconv.Atoi(string(reply.Reply))
		count += repInt
	}
	return count
}

func CountConnectedComponentsMap(graph interface{}) interface{} {
	visited := make(map[Vertex]bool)
	count := 0

	for vertex := range graph.(Graph) {
		if !visited[vertex] {
			dfs(graph.(Graph), vertex, visited)
			count++
		}
	}

	return count
}

func dfs(graph Graph, vertex Vertex, visited map[Vertex]bool) {
	visited[vertex] = true

	for _, neighbor := range graph[vertex] {
		if !visited[neighbor] {
			dfs(graph, neighbor, visited)
		}
	}
}

func CountConnectedComponentsReduce(replies interface{}) interface{} {
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
