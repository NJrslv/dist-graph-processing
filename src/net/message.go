package net

type reqMsg struct {
	clientName string // name of sending Client
	meth       string // e.g. "Print"
	to         Role   // client sends to Coordinator, Coordinator to Worker
	args       []byte
	replyCh    chan ReplyMsg
}

type ReplyMsg struct {
	Ok    bool
	Reply []byte
}
