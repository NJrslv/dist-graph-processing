package net

import (
	"log"
	"reflect"
	"sync"
)

type reqMsg struct {
	client   string // name of sending Client
	meth     string // e.g. "Print"
	argsType reflect.Type
	args     []byte
	replyCh  chan replyMsg
}

type replyMsg struct {
	ok    bool
	reply []byte
}

type Network struct {
	mu          sync.Mutex
	name        string
	reliable    bool
	clients     map[string]*Client // clients, by name
	nodes       map[string]*Node   // nodes, by name
	connections map[string]string  // client name -> servername
	lb          LoadBalancer       // circular queue of nodes
	done        chan struct{}      // closed when Network is cleaned up
	clientCh    chan reqMsg        // chan with requests from clients
	count       int32              // total RPC count, for statistics
	bytes       int64              // total bytes send, for statistics
}

func MakeNetwork(name string) *Network {
	return &Network{
		name:        name,
		reliable:    true,
		clients:     make(map[string]*Client),
		nodes:       make(map[string]*Node),
		connections: make(map[string]string),
		lb:          MakeLoadBalancer(nil),
		done:        make(chan struct{}),
		clientCh:    make(chan reqMsg),
		count:       0,
		bytes:       0,
	}
}

// ConnectClient maps client and node
func (n *Network) ConnectClient(c *Client) {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, exists := n.clients[c.name]
	if exists {
		coordinator := n.lb.GetNextNode()
		n.connections[c.name] = coordinator.name
	} else {
		log.Printf("Network.connect(): %s is not connected to the Network\n", c.name)
	}
}
