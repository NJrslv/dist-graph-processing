package net

import (
	"log"
	"runtime"
	"time"
)

/*
	- Nodes are communicating using their channels (node.replyCh / node.reqCh)
	- Nodes are communicating with clients using request channel (req.ReplyCh)
*/

// TODO check why 1110 when 1000 nodes, check for races, run the system with 1e3 nodes

type Role int

const (
	Coordinator Role = iota
	Worker
)

type Node struct {
	name    string
	count   int            // incoming RPCs
	reqCh   chan reqMsg    // requests to this node
	bc      *BroadCaster   // see service.go/BroadCaster
	methInv *MethodInvoker // see service.go/MethodInvoker
	done    chan struct{}  // closed when Network is cleaned up
}

func MakeNode(name string, done chan struct{}) *Node {
	return &Node{
		name:  name,
		count: 0,
		reqCh: make(chan reqMsg, 1),
		done:  done,
	}
}

// Run represents a goroutine,
// and also it is a node in the system
func (n *Node) Run(replyCh chan ReplyMsg) {
	log.Printf(" %d : node.Run()", runtime.NumGoroutine())

	for {
		select {
		case <-n.done:
			// Entire Network has been destroyed.
			return
		case req := <-n.reqCh:
			/*
				1. handle a request
				2. send a reply
			*/
			var reply ReplyMsg
			switch req.to {
			case Coordinator:
				reply = <-n.handleCoordinator(req)
			case Worker:
				reply = <-n.handleWorker(req)
			}
			replyCh <- reply
		}
	}
}

func (n *Node) Dispatch(req reqMsg) <-chan ReplyMsg {
	log.Printf(" %d : node.Dispatch()", runtime.NumGoroutine())

	reply := make(chan ReplyMsg)
	go func() {
		// send the reply to the node.requestChan
		n.reqCh <- req

		repl := make(chan ReplyMsg)
		go n.Run(repl)

		select {
		case r := <-repl:
			reply <- r
			log.Printf(" %d : node.Dispatch() Reply with %s", runtime.NumGoroutine(), string(r.Reply))
		case <-time.After(time.Second * 2): // Timeout after 2 seconds
			log.Printf(" %d : node.Dispatch(): timeout waiting for Reply", runtime.NumGoroutine())
			reply <- ReplyMsg{false, nil}
		}
	}()
	return reply
}

func (n *Node) GetRPCount() int {
	return n.count
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
		replyCh := make(chan ReplyMsg, len(quorum)+1)

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
		res := n.methInv.InvokeMethod(methodName+"Map", string(req.args))

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
