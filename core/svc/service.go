package svc

import "distgraphia/core/net"

type BroadCaster struct {
	NodeNetwork map[string]*net.Network // nodes <-> Networks
	nodeName    string                  // node that uses broadcaster
}

func MakeBroadCaster(n *net.Network) {
	bc := &BroadCaster{
		NodeNetwork: make(map[string]*net.Network),
	}
	for nodeName, node := range n.GetNodes() {
		bc.NodeNetwork[nodeName] = n
		bc.nodeName = nodeName
		node.ConnBroadCaster(bc)
	}
}

func (bc *BroadCaster) GatherQuorum() map[string]*net.Node {
	nodeNet := bc.NodeNetwork[bc.nodeName]
	return nodeNet.GetNodes()
}

type Storage struct {
}
