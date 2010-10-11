package main

import "net"
import "net/textproto"

type IRCConn struct {
	serverName   	string
	tcpConn   	*net.TCPConn
	textConn   	*textproto.Conn
	messageChannel 	chan *IRCMessage
}


func NewIRCConn(tcpConn *net.TCPConn, serverName string) *IRCConn {
	var newConn IRCConn

	newConn.messageChannel = make(chan *IRCMessage)
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

func (conn *IRCConn) MessageChannel() chan *IRCMessage {

	return conn.messageChannel
}

func (conn *IRCConn) Dispatch() {

	ircChannel := ReadLineIter(conn.tcpConn)

	for !closed(conn.messageChannel) && !closed(ircChannel) {
		select {

			//Message to be sent
			case incomingMessage := <-conn.messageChannel:
				if(incomingMessage == nil) { break }
				conn.textConn.PrintfLine(incomingMessage.raw)

			//Handle raw line from irc socket
			case line := <-ircChannel:
				if(len(line) == 0) { break }
				outgoingMessage := NewIRCMessage(line)
				conn.DispatchOrConsumeMessage(outgoingMessage)
		}
	}

	close(conn.messageChannel)
	close(ircChannel)
}

func (conn *IRCConn) DispatchOrConsumeMessage(message *IRCMessage) {
	logger.Logf("<- IRC: %s", message.raw)

	switch message.command {
		case "PING":
			conn.textConn.PrintfLine("PONG %s", message.args[0])
		default:
			conn.messageChannel <- message
	}
}
