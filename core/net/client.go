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

func (c *Client) connectTo(net *Network) {
	net.ConnectClient(c)
}
