package network

import (
	"distgraphia/src/node"
	"reflect"
	"sync"
)

type reqMsg struct {
	client   interface{} // name of sending Client
	meth     string      // e.g. "Print"
	argsType reflect.Type
	args     []byte
	replyCh  chan replyMsg
}

type replyMsg struct {
	ok    bool
	reply []byte
}

type Network struct {
	mu          sync.Mutex
	reliable    bool
	ends        map[interface{}]*Client     // clients, by name
	servers     map[interface{}]*node.Node  // nodes, by name
	connections map[interface{}]interface{} // endname -> servername
	done        chan struct{}               // closed when Network is cleaned up
	clientCh    chan reqMsg                 // chan with requests from clients
	count       int32                       // total RPC count, for statistics
	bytes       int64                       // total bytes send, for statistics
}
