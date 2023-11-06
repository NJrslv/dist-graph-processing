package net

import (
	"distgraphia/core/constants"
	"log"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
)

type reqMsg struct {
	clientName string // name of sending Client
	meth       string // e.g. "Print"
	to         Role   // client sends to Coordinator, Coordinator to Worker
	argsType   reflect.Type
	args       []byte
	replyCh    chan ReplyMsg
}

type ReplyMsg struct {
	ok    bool
	reply []byte
}

// TO DO: think about removing mutex from the Network, I simply use channels

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
	nodes, nodesByName := MakeNodes()
	net := &Network{
		name:        name,
		reliable:    true,
		clients:     make(map[string]*Client),
		nodes:       nodesByName,
		connections: make(map[string]string),
		lb:          MakeLoadBalancer(nodes),
		done:        make(chan struct{}),
		clientCh:    make(chan reqMsg),
		count:       0,
		bytes:       0,
	}

	// single goroutine to handle all Client Call()s
	go func() {
		for {
			select {
			case req := <-net.clientCh:
				atomic.AddInt32(&net.count, 1)
				atomic.AddInt64(&net.bytes, int64(len(req.args)))
				go net.processReq(req)
			case <-net.done:
				return
			}
		}
	}()

	return net
}

func (n *Network) processReq(req reqMsg) {
	client, coordinator := n.readClientInfo(req)

	if client != nil && coordinator != nil {
		// execute the request (call the RPC handler).
		ech := make(chan ReplyMsg)
		go func() {
			r := coordinator.Dispatch(req)
			ech <- r
		}()

		// wait for handler to return
		var reply ReplyMsg
		replyOK := false
		for replyOK == false {
			reply = <-ech
			replyOK = true
		}
		req.replyCh <- reply
	}
}

func (n *Network) readClientInfo(req reqMsg) (*Client, *Node) {
	client, clientOk := n.clients[req.clientName]
	nodeName, nodeOk := n.connections[req.clientName]

	if !clientOk || !nodeOk {
		log.Print("Network.readClientInfo(): There is no client or coordinator")
		return nil, nil
	}

	node := n.nodes[nodeName]
	return client, node
}

// MakeNodes creates 'NumNodes' workers
// with the names '1', '2', ... , 'NumNodes'
func MakeNodes() ([constants.NumNodes]*Node, map[string]*Node) {
	var nodes [constants.NumNodes]*Node
	nodeMap := make(map[string]*Node, constants.NumNodes)

	for i := 0; i < constants.NumNodes; i++ {
		nodeName := strconv.Itoa(i)
		node := MakeNode(nodeName)

		nodes[i] = node
		nodeMap[nodeName] = node
	}
	return nodes, nodeMap
}

// ConnectClient maps client and node
func (n *Network) ConnectClient(c *Client) {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, ok := n.clients[c.GetName()]
	if !ok {
		coordinator := n.lb.GetNextNode()
		n.connections[c.GetName()] = coordinator.name
	} else {
		log.Printf("Network.connect(): %s is already connected to the Network\n", c.GetName())
	}
}
