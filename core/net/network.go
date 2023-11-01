package net

import (
	"log"
	"reflect"
	"sync"
)

type reqMsg struct {
	client   interface{} // name of sending Client
	meth     string      // e.g. "Print"
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
	ends        map[string]*Client // clients, by name
	nodes       map[string]*Runner // nodes, by name
	connections map[string]string  // endname -> servername
	done        chan struct{}      // closed when Network is cleaned up
	clientCh    chan reqMsg        // chan with requests from clients
	count       int32              // total RPC count, for statistics
	bytes       int64              // total bytes send, for statistics
}

// ConnectClient maps client and node
func (n *Network) ConnectClient(c *Client) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.ends[c.GetName()] = c
	// find the most free node to connect and make him coordinator to this user

}

// maps client in the Network and node(coordinator)
func (n *Network) connect(clientName string, nodeName string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, exists := n.ends[clientName]
	if exists {
		n.connections[clientName] = nodeName
	} else {
		log.Printf("Network.connect(): %s is not connected to the Network\n", clientName)
	}
}
