package main

import "net"
import "net/textproto"

type IRCConn struct {
	serverName	string
	tcpConn		*net.TCPConn
	textConn	*textproto.Conn
	incomingChannel	chan *IRCMessage
	outgoingChannel	chan *IRCMessage
}


func NewIRCConn(tcpConn *net.TCPConn, serverName string) *IRCConn {
	var newConn IRCConn

	newConn.incomingChannel = make(chan *IRCMessage, 100)
	newConn.outgoingChannel = make(chan *IRCMessage, 100)

	newConn.serverName = serverName
	newConn.tcpConn = tcpConn
	newConn.textConn = textproto.NewConn(tcpConn)

	go newConn.Dispatch()

	return &newConn
}

func (conn *IRCConn) Close() {
	conn.tcpConn.Close()
}

func (conn *IRCConn) SendCode(code int) {
	conn.textConn.PrintfLine(":%s %d :Error\r\n",
				    conn.serverName,
				    code)
}


func (conn *IRCConn) Dispatch() {

	tcpChannel := ReadLineIter(conn.tcpConn)

	for !closed(conn.incomingChannel) && !closed(conn.outgoingChannel) && !closed(tcpChannel) {
		select {

			//Message to be sent
			case outgoingMessage := <-conn.outgoingChannel:
				if(outgoingMessage == nil) { break }
				logger.Logf("-> IRC: %s", outgoingMessage.raw)
				conn.textConn.PrintfLine(outgoingMessage.raw)

			//Handle raw line from irc socket
			case line := <-tcpChannel:
				if(len(line) == 0) { break }
				incomingMessage := NewIRCMessage(line)
				conn.DispatchOrConsumeMessage(incomingMessage)
		}
	}

	close(conn.incomingChannel)
	close(conn.outgoingChannel)
	close(tcpChannel)
}

func (conn *IRCConn) DispatchOrConsumeMessage(message *IRCMessage) {
	logger.Logf("<- IRC: %s", message.raw)

	switch message.command {
		case "PING":
			conn.textConn.PrintfLine("PONG %s", message.args[0])
		default:
			conn.incomingChannel <- message
	}
}
