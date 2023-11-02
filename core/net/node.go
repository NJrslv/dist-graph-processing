package net

import (
	"distgraphia/core/svc"
	"sync"
)

type Role int

const (
	Coordinator Role = iota
	Worker
)

type Node struct {
	mu       sync.Mutex
	name     string
	services map[string]*svc.Serviceable // Services, by names
	role     Role                        // coordinator or worker
	count    int                         // incoming RPCs
}

func (n *Node) Run() {
	n.mu.Lock()
	defer n.mu.Unlock()

	switch n.role {
	case Coordinator:
		// Coordinator logic
	case Worker:
		// Worker logic
	}
}

func (n *Node) GetRPCount() int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.count
}

func (n *Node) SetRole(role Role) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.role = role
}
