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
	name  string
	count int         // incoming RPCs
	reqCh chan reqMsg // requests to this particular node
	bc    *svc.BroadCaster
}

func MakeNode(name string) *Node {
	return &Node{
		name:  name,
		count: 0,
		reqCh: make(chan reqMsg),
		// svc is connected in Network
	}
}

// Run is simplified, Run() represents a goroutine,
// and also it is a Node in the system
// each Node can concurrently process the requests
// but in this implementation it is done sequentially
// simply add go handleRequest() and create chan of sub-replies
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
			switch req.to {
			case Coordinator:
				reps := n.handleCoordinator(req)
				// TODO sum the replies up, then put res into the req.replyCh
			case Worker:
				n.handleWorker(req)
			}
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

func (n *Node) handleCoordinator(req reqMsg) []ReplyMsg {
	/*
		1. Gather Quorum(Nodes)
		2. Send them the task
		3. Get the reply
	*/
	quorum := n.bc.GatherQuorum()
	// A channel to hold the responses from the Dispatch function
	replyCh := make(chan ReplyMsg)

	for _, node := range quorum {
		go func(node *Node) {
			task := reqMsg{
				clientName: node.GetName(),
				meth:       req.meth,
				to:         Worker,
				args:       req.args,
				replyCh:    replyCh,
			}
			node.Dispatch(task)
		}(node)
	}

	var replies []ReplyMsg
	for i := 0; i < len(quorum); i++ {
		replies = append(replies, <-replyCh)
	}

	return replies
}

func (n *Node) handleWorker(req reqMsg) {
	/*
		1. Process the task
		2. Send the result back
	*/
	methodName := req.meth
	/*
		pseudo: (reflectionMap[methodName]func()
		meth := reflectionMap[methodName]
		exec meth(req.args)
		send the result
	*/
}

func (n *Node) GetName() string {
	return n.name
}

// ConnBroadCaster connects all BroadCaster to the node
func (n *Node) ConnBroadCaster(bc *svc.BroadCaster) {
	n.bc = bc
}
