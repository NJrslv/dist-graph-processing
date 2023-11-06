package net

import (
	"bytes"
	"encoding/gob"
	"log"
	"reflect"
)

type Client struct {
	name        string
	connections map[string]*Network // Networks, by name
}

func MakeClient(name string) *Client {
	return &Client{
		name:        name,
		connections: make(map[string]*Network),
	}
}

func (c *Client) GetName() string {
	return c.name
}

func (c *Client) ConnectTo(net *Network) {
	net.ConnectClient(c)
}

// Call sends an RPC, waits for the reply.
// The return value indicates success, false means that
// no reply was received from the server.
func (c *Client) Call(netName string, meth string, args interface{}, reply interface{}) bool {
	req := reqMsg{
		clientName: c.GetName(),
		meth:       meth,
		to:         Coordinator,
		argsType:   reflect.TypeOf(args),
		replyCh:    make(chan ReplyMsg),
	}

	// Encode args
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(args); err != nil {
		log.Fatalf("ClientEnd.Call(): gob encode args: %v\n", err)
	}
	req.args = buf.Bytes()

	// Send the requests
	net := c.connections[netName]
	select {
	case net.clientCh <- req:
		// The request has been sent.
	case <-net.done:
		// Entire Network has been destroyed.
		return false
	}

	// Wait for the reply.
	rep := <-req.replyCh
	if rep.ok {
		// Decode reply
		dec := gob.NewDecoder(bytes.NewBuffer(rep.reply))
		if err := dec.Decode(reply); err != nil {
			log.Fatalf("ClientEnd.Call(): gob decode reply: %v\n", err)
		}
		return true
	} else {
		return false
	}
}
