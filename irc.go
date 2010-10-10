package main

import "net"

type IRCEvent struct {
	raw string
}

type IRCConn struct {
	_clientConn *net.TCPConn;
	_eventChannel chan *IRCEvent
}


func NewIRCConn(conn *net.TCPConn) *IRCConn {
	return &IRCConn {
		_clientConn: conn,
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
      conn._eventChannel <- &IRCEvent{ 
      									raw: line, 
      								 }
   }
   
   //Client closed connection
   close(conn._eventChannel)
}
