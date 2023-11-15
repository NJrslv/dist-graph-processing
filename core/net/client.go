package net

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
func (c *Client) Call(netName string, meth string, args string, reply *string) bool {
	req := reqMsg{
		clientName: c.GetName(),
		meth:       meth,
		to:         Coordinator,
		args:       []byte(args),
		replyCh:    make(chan ReplyMsg),
	}

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
		*reply = string(rep.reply)
		return true
	} else {
		return false
	}
}
