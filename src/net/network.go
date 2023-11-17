package net

import (
	"log"
	"strconv"
	"sync"
	"sync/atomic"
)

type reqMsg struct {
	clientName string // name of sending Client
	meth       string // e.g. "Print"
	to         Role   // client sends to Coordinator, Coordinator to Worker
	args       []byte
	replyCh    chan ReplyMsg
}

type ReplyMsg struct {
	Ok    bool
	Reply []byte
}

// TODO: think about removing mutex from the Network, I simply use channels
// Structure: Client.Call() --> [network <- request] --> network.processRequest() -->

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

	methods := []string{"CountNodes"}
	methInv := MakeMethodInvoker(methods)
	bc := MakeBroadCaster(net)
	net.connectServices(bc, methInv)

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
		log.Fatalf("Network.readClientInfo(): There is no client or coordinator")
		return nil, nil
	}

	node := n.nodes[nodeName]
	return client, node
}

// MakeNodes creates 'NumNodes' workers
// with the names '1', '2', ... , 'NumNodes'
func MakeNodes() ([NumNodes]*Node, map[string]*Node) {
	var nodes [NumNodes]*Node
	nodeMap := make(map[string]*Node, NumNodes)

	for i := 0; i < NumNodes; i++ {
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

		n.clients[c.GetName()] = c
		n.connections[c.GetName()] = coordinator.name
		c.connections[n.name] = n
	} else {
		log.Printf("Network.connect(): %s is already connected to the Network\n", c.GetName())
	}
}

func (n *Network) Cleanup() {
	close(n.done)
}

func (n *Network) connectServices(bc *BroadCaster, methInv *MethodInvoker) {
	for nodeName, node := range n.nodes {
		bc.Net = n
		bc.NodeName = nodeName
		methInv.NodeName = nodeName

		node.ConnServices(bc, methInv)
	}
}

func (n *Network) GetNodes() map[string]*Node {
	return n.nodes
}
