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
	go newLilyConn.ReadMessages()

	return &newLilyConn
}

func (conn LilyConn) Send(line string) {
	logger.Logf("-> LILY: %s", line)
	conn.textConn.PrintfLine(line)
}

func (conn LilyConn) SendOptions() {
	//Send options before all else
	conn.Send("#$# options +version +prompt +prompt2 +leaf-notify +leaf-cmd +connected")
}

func (conn LilyConn) ReadMessages() {
	for line := range ReadLineIter(conn.tcpConn) {
		message := NewLilyMessage(line)
		conn.DispatchOrConsumeMessage(message)
	}

	//Server closed connection
	close(conn.messageChannel)
}

func (conn LilyConn) DispatchOrConsumeMessage(message *LilyMessage) {
	logger.Logf("<- LILY: %s", message.raw)
}
