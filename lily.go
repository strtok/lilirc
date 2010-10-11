package main

import "net"
import "net/textproto"

type LilyConn struct {
	tcpConn	*net.TCPConn
	textConn *textproto.Conn
	messageChannel chan *LilyMessage
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
	newLilyConn.messageChannel = make(chan *LilyMessage)

	newLilyConn.SendOptions()
	go newLilyConn.Dispatch()

	return &newLilyConn
}

func (conn *LilyConn) MessageChannel() chan *LilyMessage {
	return conn.messageChannel
}

func (conn *LilyConn) SendOptions() {
	//Send options before all else
	conn.textConn.PrintfLine("#$# options +version +prompt +prompt2 +leaf-notify +leaf-cmd +connected")
}

func (conn *LilyConn) Dispatch() {

	lilyChannel := ReadLineIter(conn.tcpConn)

	for !closed(conn.messageChannel) && !closed(lilyChannel) {
		select {

			//Message to be sent
			case incomingMessage := <-conn.messageChannel:
				if(incomingMessage == nil) { break }
				conn.textConn.PrintfLine(incomingMessage.raw)

			//Handle raw line from lily socket
			case line := <-lilyChannel:
				if(len(line) == 0) { break }
				outgoingMessage := NewLilyMessage(line)
				conn.DispatchOrConsumeMessage(outgoingMessage)
		}
	}

	//Server closed connection
	close(conn.messageChannel)
}

func (conn *LilyConn) DispatchOrConsumeMessage(message *LilyMessage) {
	logger.Logf("<- LILY: %s", message.raw)

	switch(message.command) {
		case "PROMPT":
			conn.textConn.PrintfLine("");
	}
}
