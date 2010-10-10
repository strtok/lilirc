package main

import "net"
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
	_clientConn   *net.TCPConn
	_eventChannel chan *IRCEvent
}


func NewIRCConn(conn *net.TCPConn) *IRCConn {
	return &IRCConn{
		_clientConn:   conn,
		_eventChannel: nil,
	}
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
