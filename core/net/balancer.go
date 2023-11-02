package net

import (
	"log"
	"sync"
)

/*
	We use Round Robin algorithm because we only have RPC counts
	to assess node usage, and they can change quickly as nodes interact.
	Using these counts isn't practical. A cyclic queue is the simplest
	way to distribute the load evenly.

	Since the number of nodes is statically defined,
	we don't need to create a traditional queue.
*/

type LoadBalancer struct {
	nodes   []*Node
	current int
	mu      sync.Mutex
}

func MakeLoadBalancer(nodes []*Node) LoadBalancer {
	return LoadBalancer{
		nodes:   nodes,
		current: 0,
	}
}

func (lb *LoadBalancer) GetNextNode() *Node {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.nodes) == 0 {
		log.Print("balancer.GetNextNode(): there are 0 nodes for balancing\n")
		return nil
	}

	node := lb.nodes[lb.current]
	lb.current = (lb.current + 1) % len(lb.nodes)
	return node
}
