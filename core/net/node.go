package net

import (
	"distgraphia/core/svc"
	"log"
	"time"
)

type Role int

const (
	Coordinator Role = iota
	Worker
)

type Node struct {
	name     string
	services map[string]*svc.Serviceable // Services, by names
	count    int                         // incoming RPCs
	reqCh    chan reqMsg                 // requests to this particular node
}

func MakeNode(name string) *Node {
	return &Node{
		name:     name,
		services: make(map[string]*svc.Serviceable),
		count:    0,
		reqCh:    make(chan reqMsg),
	}
}

func (n *Node) Run(done chan struct{}) {
	for {
		select {
		case <-done:
			// Entire Network has been destroyed.
			return
		case req := <-n.reqCh:
			/*
				handle the request:
						1. decode and execute
						2. encode
				put encoded reply to the req.replyCh
			*/
		}
	}
}

func (n *Node) Dispatch(req reqMsg) ReplyMsg {
	/*
			1. Put the encoded request in the requestChan in the Node.
			2. The Goroutine running the node.Run() function
		       will handle the request and place the response
		       in the request.replyChan.
	*/
	n.reqCh <- req
	select {
	case reply := <-req.replyCh:
		return reply
	case <-time.After(time.Second * 5): // Timeout after 5 seconds
		log.Print("Node.Dispatch(): timeout waiting for reply")
		return ReplyMsg{false, nil}
	}
}

func (n *Node) GetRPCount() int {
	return n.count
}
