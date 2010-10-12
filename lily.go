package main

import "net"
import "net/textproto"

type LilyConn struct {
	tcpConn	*net.TCPConn
	textConn *textproto.Conn
	incomingChannel chan *LilyMessage
	outgoingChannel chan *LilyMessage
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
	newLilyConn.SendOptions()
	go newLilyConn.Dispatch()

	return &newLilyConn
}

func (conn *LilyConn) Close() {
	conn.tcpConn.Close()
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
				incomingMessage := NewLilyMessage(line)
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
			conn.textConn.PrintfLine("");
		case "CONNECTED":
			conn.incomingChannel <- message
	}
}
