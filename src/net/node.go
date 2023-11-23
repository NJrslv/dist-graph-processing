package net

import (
	"log"
	"runtime"
	"sync/atomic"
	"time"
)

/*
	- Nodes are communicating using connections
		(channels that are created when needed to send/get the request/reply)
	- Nodes are communicating with clients using request channel (req.ReplyCh)
*/

type Role int

const (
	Coordinator Role = iota
	Worker
)

type Node struct {
	name string
	// reqCh   chan reqMsg    // requests to this node
	bc      *BroadCaster   // see service.go/BroadCaster
	methInv *MethodInvoker // see service.go/MethodInvoker
	done    chan struct{}  // closed when Network is cleaned up
	count   int32          // total RPC count, for statistics
	g       *Graph
}

func MakeNode(name string, done chan struct{}) *Node {
	return &Node{
		name: name,
		//reqCh: make(chan reqMsg, 1),
		done:  done,
		count: 0,
	}
}

// Run represents a goroutine,
// and also it is a node in the system
func (n *Node) Run(req reqMsg, replyCh chan ReplyMsg) {
	log.Printf(" %d : node.Run()", runtime.NumGoroutine())

	var reply ReplyMsg
	switch req.to {
	case Coordinator:
		reply = <-n.handleCoordinator(req)
	case Worker:
		reply = <-n.handleWorker(req)
	}
	replyCh <- reply
}

func (n *Node) Dispatch(req reqMsg) <-chan ReplyMsg {
	log.Printf(" %d : node.Dispatch()", runtime.NumGoroutine())

	// stat
	atomic.AddInt32(&n.count, 1)

	reply := make(chan ReplyMsg)
	go func() {
		// send the reply to the node.requestChan
		repl := make(chan ReplyMsg)
		go n.Run(req, repl)

		// get the reply
		select {
		case r := <-repl:
			reply <- r
			log.Printf(" %d : node.Dispatch() Reply with %s", runtime.NumGoroutine(), string(r.Reply))
		case <-time.After(time.Second * 3): // Timeout after 2 seconds
			log.Printf(" %d : node.Dispatch(): timeout waiting for Reply", runtime.NumGoroutine())
			reply <- ReplyMsg{false, nil}
		}
	}()
	return reply
}

func (n *Node) GetRPCount() int32 {
	return atomic.LoadInt32(&n.count)
}

func (n *Node) handleCoordinator(req reqMsg) <-chan ReplyMsg {
	log.Printf(" %d : node.handleCoordinator()", runtime.NumGoroutine())
	/*
		1. Gather Quorum(Nodes)
		2. Send them the task
		3. Get the Reply
	*/
	reply := make(chan ReplyMsg)
	go func() {
		quorum := n.bc.GatherQuorum()
		replyCh := make(chan ReplyMsg, len(quorum))

		// Scatter
		for _, node := range quorum {
			go func(node *Node) {
				task := reqMsg{
					clientName: node.GetName(),
					meth:       req.meth,
					to:         Worker,
					args:       req.args,
					replyCh:    req.replyCh,
				}
				rep := <-node.Dispatch(task)
				replyCh <- rep
			}(node)
		}

		// Gather
		var replies []ReplyMsg
		for i := 0; i < len(quorum); i++ {
			replies = append(replies, <-replyCh)
		}

		agg := n.methInv.InvokeMethod(req.meth+"Reduce", replies).(int)
		reply <- ReplyMsg{true, intToStrBytes(agg)}
	}()
	return reply
}

func (n *Node) handleWorker(req reqMsg) <-chan ReplyMsg {
	log.Printf(" %d : node.handleWorker()", runtime.NumGoroutine())
	/*
		1. Process the task
		2. Send the result back
	*/
	reply := make(chan ReplyMsg)
	go func() {
		methodName := req.meth
		res := n.methInv.InvokeMethod(methodName+"Map", req.args)

		var repl ReplyMsg
		if len(res.(string)) == 0 {
			repl.Ok = false
		}

		repl.Reply = []byte(res.(string))
		reply <- repl
	}()
	return reply
}

func (n *Node) GetName() string {
	return n.name
}

func (n *Node) ConnServices(bc *BroadCaster, methInv *MethodInvoker) {
	n.bc = bc
	n.methInv = methInv
}
