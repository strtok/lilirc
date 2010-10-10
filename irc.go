package main

import "net"
import "net/textproto"
import "strings"

type IRCEvent struct {
	raw string
	command string
	args []string
}

func NewIRCEvent(raw string) *IRCEvent {

	var newEvent IRCEvent
        newEvent.raw = raw

	tokens := strings.Split(raw, " ",  -1)
	newEvent.command = tokens[0]
	newEvent.args = tokens[1:]

	return &newEvent
}


type IRCConn struct {
	_serverName   string
	_clientConn   *net.TCPConn
	_clientText   *textproto.Conn
	_eventChannel chan *IRCEvent
}


func NewIRCConn(conn *net.TCPConn, serverName string) *IRCConn {
	return &IRCConn{
		_serverName: 	serverName,
		_clientConn:  	conn,
		_clientText:  	textproto.NewConn(conn),
		_eventChannel: 	nil,
	}
}

func (conn IRCConn) sendCode(code int) {
	conn._clientText.PrintfLine(":%s %d :Error\r\n", 
				    conn._serverName,
				    code) 
}

func (conn IRCConn) eventChannel() chan *IRCEvent {

	if conn._eventChannel != nil {
		return conn._eventChannel
	}

	conn._eventChannel = make(chan *IRCEvent)
	go conn.readEvents()
	return conn._eventChannel
}

func (conn IRCConn) readEvents() {
	for line := range ReadLineIter(conn._clientConn) {
		event := NewIRCEvent(line)
		conn.dispatchOrConsumeEvent(event)
	}

	//Client closed connection
	close(conn._eventChannel)
}

func (conn IRCConn) dispatchOrConsumeEvent(event *IRCEvent) {
	conn._eventChannel <- event
}
