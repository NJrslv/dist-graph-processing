package svc

import "distgraphia/core/net"

// BroadCaster broadcasts the message across the network,
// each node has its own BroadCaster service
type BroadCaster struct {
	Net      *net.Network // Network
	NodeName string       // node that uses broadcaster
}

func MakeBroadCaster(n *net.Network) *BroadCaster {
	return &BroadCaster{
		Net:      n,
		NodeName: "",
	}
}

func (bc *BroadCaster) GatherQuorum() map[string]*net.Node {
	return bc.Net.GetNodes()
}

type Storage struct {
	// for each node different Storage instance
}
