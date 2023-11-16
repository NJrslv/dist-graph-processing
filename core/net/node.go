package net

import (
	"encoding/binary"
	"log"
	"strconv"
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
				put encoded Reply to the req.replyCh
			*/
			switch req.to {
			case Coordinator:
				reply := n.handleCoordinator(req)
				req.replyCh <- reply
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
		log.Print("Node.Dispatch(): timeout waiting for Reply")
		return ReplyMsg{false, nil}
	}
}

func (n *Node) GetRPCount() int {
	return n.count
}

func (n *Node) handleCoordinator(req reqMsg) ReplyMsg {
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

	aggReply := n.aggregateReplies(req.meth, replies)
	return aggReply
}

func (n *Node) aggregateReplies(method string, replies []ReplyMsg) ReplyMsg {
	// TODO after testing move methods to the algorithms + check Mutexes
	switch method {
	case "SUM":
		sum := 0
		for _, reply := range replies {
			sum += int(binary.LittleEndian.Uint64(reply.Reply))
		}
		return ReplyMsg{Ok: true, Reply: []byte(strconv.Itoa(sum))}
	default:
		log.Print("Node.aggregateReplies(): Unknown method")
		return ReplyMsg{Ok: false, Reply: []byte("Unknown method")}
	}
}

func (n *Node) handleWorker(req reqMsg) {
	/*
		1. Process the task
		2. Send the result back
	*/
	methodName := req.meth
	res := n.methInv.InvokeMethod(methodName, string(req.args))
	/*
		meth := req.meth
		res := n.methInv.InvokeMethod(methodName + ".Map()", string(req.args))
	*/
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
