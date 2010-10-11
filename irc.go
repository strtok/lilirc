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

	go newConn.ReadMessages()

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

func (conn *IRCConn) ReadMessages() {
	for line := range ReadLineIter(conn.tcpConn) {
		message := NewIRCMessage(line)
		conn.DispatchOrConsumeMessage(message)
	}

	//Client closed connection
	close(conn.messageChannel)
}

func (conn *IRCConn) DispatchOrConsumeMessage(message *IRCMessage) {
	logger.Logf("<- IRC: %s", message.raw)
	conn.messageChannel <- message
}
