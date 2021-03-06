package main

import "net"
import "net/textproto"

type LilyHandle struct {
	name string
}

type LilyConn struct {
	tcpConn	*net.TCPConn
	textConn *textproto.Conn
	incomingChannel chan *LilyMessage
	outgoingChannel chan *LilyMessage

	//Map of user id (e.g. #105) to LilyUser.
	//This map is kept up to date from %USER messages
	handleMap map[string] *LilyHandle
}

func NewLilyConn(address string) *LilyConn {
	localAddress, _ := net.ResolveTCPAddr("0.0.0.0:0")
	lilyAddress, err := net.ResolveTCPAddr(address)

	if(err != nil) {
		logger.Logf("failed resolving %s: %s", address, err)
		return nil
	}

	tcpConn, err := net.DialTCP("tcp", localAddress, lilyAddress)

	if(err != nil) {
		logger.Logf("failed connecting to %s: %s", address, err)
		return nil
	}

	var newLilyConn LilyConn

	newLilyConn.tcpConn = tcpConn
	newLilyConn.textConn = textproto.NewConn(tcpConn)
	newLilyConn.incomingChannel = make(chan *LilyMessage, 100)
	newLilyConn.outgoingChannel = make(chan *LilyMessage, 100)
	newLilyConn.handleMap = make(map[string] *LilyHandle)

	newLilyConn.SendOptions()

	go newLilyConn.Dispatch()

	return &newLilyConn
}

func (conn *LilyConn) HandleMap(name string) string {
	handle, ok := conn.handleMap[name]
	if(ok) {
		return handle.name
	}
	return ""
}

func (conn *LilyConn) Close() {
	conn.tcpConn.Close()
}

func (conn *LilyConn) Send(raw string) {
	conn.outgoingChannel <- &LilyMessage{raw: raw}
}

func (conn *LilyConn) SendOptions() {
	//Send options before all else
	conn.textConn.PrintfLine("#$# options +version +prompt +prompt2 +leaf-notify +leaf-cmd +connected")
}

func (conn *LilyConn) Dispatch() {

	tcpChannel := ReadLineIter(conn.tcpConn)

	for !closed(conn.incomingChannel) && !closed(conn.outgoingChannel) && !closed(tcpChannel) {
		select {

			//Message to be sent
			case outgoingMessage := <-conn.outgoingChannel:
				if(outgoingMessage == nil) { break }
				logger.Logf("-> LILY: %s", outgoingMessage.raw)
				conn.textConn.PrintfLine(outgoingMessage.raw)

			//Handle raw line from lily socket
			case line := <-tcpChannel:
				if(len(line) == 0) { break }
				incomingMessage := NewLilyMessage(line, func(name string) string { return conn.HandleMap(name) } )
				conn.DispatchOrConsumeMessage(incomingMessage)
		}
	}

	close(conn.incomingChannel)
	close(conn.outgoingChannel)
	close(tcpChannel)
}

func (conn *LilyConn) DispatchOrConsumeMessage(message *LilyMessage) {
	logger.Logf("<- LILY: %s", message.raw)

	switch(message.command) {
		case "PROMPT":
			conn.textConn.PrintfLine("")
		case "CONNECTED":
			conn.textConn.PrintfLine("/where")
			conn.incomingChannel <- message
		case "USER", "DISC":
			conn.DispatchHandleUpdate(message)
		case "NOTIFY":
			conn.DispatchNotify(message)
	}
}

func (conn *LilyConn) DispatchHandleUpdate(message *LilyMessage) {
	//Example:
	//%USER HANDLE=#100 NAME=14=System Manager BLURB=0= LOGIN=1286726924 INPUT=1286747362 STATE=detach ATTRIB=0= PRONOUN=5=their

	handle := message.attributes["HANDLE"]
	name := message.attributes["NAME"]

	conn.handleMap[handle] = &LilyHandle{ name: name }
}

func (conn *LilyConn) DispatchNotify(message *LilyMessage) {
	conn.incomingChannel <- message
}
