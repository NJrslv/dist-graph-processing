package net

import (
	"log"
	"runtime"
	"time"
)

type Role int

const (
	Coordinator Role = iota
	Worker
)

type Node struct {
	name    string
	count   int         // incoming RPCs
	reqCh   chan reqMsg // requests to this particular node
	bc      *BroadCaster
	methInv *MethodInvoker
	done    chan struct{} // closed when Network is cleaned up
}

func MakeNode(name string, done chan struct{}) *Node {
	return &Node{
		name:  name,
		count: 0,
		reqCh: make(chan reqMsg),
		done:  done,
		// svc is connected in Network
	}
}

// Run represents a goroutine,
// and also it is a node in the system
// each Node can concurrently process the requests
func (n *Node) Run() {
	log.Printf(" %d : node.Run()", runtime.NumGoroutine())
	for {
		select {
		case <-n.done:
			// Entire Network has been destroyed.
			return
		case req := <-n.reqCh:
			/*
				handle the request:
						1. decode and execute
						2. encode
				put encoded Reply to the req.replyCh
			*/
			switch req.to {
			case Coordinator:
				// TODO go handleCoordinator + handleWorker
				reply := n.handleCoordinator(req)
				req.replyCh <- reply
			case Worker:
				n.handleWorker(req)
			}
		}
	}
}

// TODO my workers cannot pass the replies to the coordinator => we sleep, check logs

func (n *Node) Dispatch(req reqMsg) ReplyMsg {
	log.Printf(" %d : node.Dispatch()", runtime.NumGoroutine())
	/*
			1. Put the encoded request in the requestChan in the Node.
			2. The Goroutine running the node.Run() function
		       will handle the request and place the response
		       in the request.replyChan.
	*/
	go n.Run()
	n.reqCh <- req
	select {
	case reply := <-req.replyCh:
		return reply
	case <-time.After(time.Second * 5): // Timeout after 5 seconds
		log.Print("Node.Dispatch(): timeout waiting for Reply")
		return ReplyMsg{false, nil}
	}
}

func (n *Node) GetRPCount() int {
	return n.count
}

func (n *Node) handleCoordinator(req reqMsg) ReplyMsg {
	log.Printf(" %d : node.handleCoordinator()", runtime.NumGoroutine())
	/*
		1. Gather Quorum(Nodes)
		2. Send them the task
		3. Get the Reply
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

	aggReply := n.methInv.InvokeMethod(req.meth+"Reduce", replies)
	return aggReply.(ReplyMsg)
}

func (n *Node) handleWorker(req reqMsg) {
	log.Printf(" %d : node.handleWorker()", runtime.NumGoroutine())
	/*
		1. Process the task
		2. Send the result back
	*/
	methodName := req.meth
	res := n.methInv.InvokeMethod(methodName+"Map", string(req.args))

	var repl ReplyMsg
	if len(res.(string)) == 0 {
		repl.Ok = false
	}

	repl.Reply = []byte(res.(string))
	req.replyCh <- repl
}

func (n *Node) GetName() string {
	return n.name
}

func (n *Node) ConnServices(bc *BroadCaster, methInv *MethodInvoker) {
	n.bc = bc
	n.methInv = methInv
}
