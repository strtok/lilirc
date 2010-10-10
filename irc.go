package main

import "net"
import "net/textproto"
import "strings"



type IRCMessage struct {
	raw string
	command string
	args []string
}

func NewIRCMessage(raw string) *IRCMessage {

	var newMessage IRCMessage
        newMessage.raw = raw

	tokens := strings.Split(raw, " ",  -1)
	newMessage.command = tokens[0]
	newMessage.args = tokens[1:]

	return &newMessage
}


type IRCConn struct {
	_serverName   string
	_clientConn   *net.TCPConn
	_clientText   *textproto.Conn
	_messageChannel chan *IRCMessage
}


func NewIRCConn(tcpConn *net.TCPConn, serverName string) *IRCConn {
	var newConn IRCConn

	newConn._messageChannel = make(chan *IRCMessage)
	newConn._serverName = serverName
	newConn._clientConn = tcpConn
	newConn._clientText = textproto.NewConn(tcpConn)

	go newConn.readMessages()

	return &newConn
}

func (conn IRCConn) sendCode(code int) {
	conn._clientText.PrintfLine(":%s %d :Error\r\n",
				    conn._serverName,
				    code) 
}

func (conn IRCConn) messageChannel() chan *IRCMessage {

	return conn._messageChannel
}

func (conn IRCConn) readMessages() {
	for line := range ReadLineIter(conn._clientConn) {
		message := NewIRCMessage(line)
		conn.dispatchOrConsumeMessage(message)
	}

	//Client closed connection
	close(conn._messageChannel)
}

func (conn IRCConn) dispatchOrConsumeMessage(message *IRCMessage) {
	conn._messageChannel <- message
}
