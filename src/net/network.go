package net

import (
	"log"
	"runtime"
	"strconv"
	"sync"
)

type Network struct {
	name        string
	reliable    bool
	clients     map[string]*Client // clients, by name
	nodes       map[string]*Node   // nodes, by name
	connections map[string]string  // client name -> servername
	lb          LoadBalancer       // circular queue of nodes
	done        chan struct{}      // closed when Network is cleaned up
	clientCh    chan reqMsg        // chan with requests from clients
	mu          sync.Mutex         // protects Network data(concurrent clients)
}

func MakeNetwork(name string) *Network {
	done := make(chan struct{})
	nodes, nodesByName := MakeNodes(done)
	net := &Network{
		name:        name,
		reliable:    true,
		clients:     make(map[string]*Client),
		nodes:       nodesByName,
		connections: make(map[string]string),
		lb:          MakeLoadBalancer(nodes),
		done:        done,
		clientCh:    make(chan reqMsg, 1),
	}

	net.connectServices()

	// single goroutine to handle all Client Call()s
	go func() {
		for {
			select {
			case req := <-net.clientCh:
				go net.processReq(req)
			case <-net.done:
				log.Printf("Network %s is done...", net.name)
				return
			}
		}
	}()

	return net
}

func (n *Network) processReq(req reqMsg) {
	log.Printf(" %d : network.processReq()", runtime.NumGoroutine())
	client, coordinator := n.readClientInfo(req)

	if client != nil && coordinator != nil {
		reply := <-coordinator.Dispatch(req)
		req.replyCh <- reply
	} else {
		log.Fatalf("network.processReq(): client == nil or coordinator == nil")
	}
}

func (n *Network) readClientInfo(req reqMsg) (*Client, *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()

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
func MakeNodes(done chan struct{}) ([NumNodes]*Node, map[string]*Node) {
	var nodes [NumNodes]*Node
	nodeMap := make(map[string]*Node, NumNodes)

	for i := 0; i < NumNodes; i++ {
		nodeName := strconv.Itoa(i)
		node := MakeNode(nodeName, done)

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

func (n *Network) connectServices() {
	methods := []string{"CountNodes"}
	methInv := MakeMethodInvoker(methods)
	bc := MakeBroadCaster(n)

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

func (n *Network) GetRPCount() int32 {
	n.mu.Lock()
	defer n.mu.Unlock()

	rpcCount := int32(0)
	for _, node := range n.nodes {
		rpcCount += node.count
	}
	return rpcCount
}
