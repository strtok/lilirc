package main

import "fmt"
import "net"

var (
	SERVERNAME = "irc.lily.com"
)

func newClient()

func main() {

	listenAddress, err := net.ResolveTCPAddr("0.0.0.0:6667")

	if err != nil {
		fmt.Print("failed creating listenAddress")
	}

	listener, err := net.ListenTCP("tcp", listenAddress)

	if err != nil {
		fmt.Print("failed listening")
	}

	for {

		newConnection, err := listener.AcceptTCP()

		if err != nil {
			fmt.Print("failed accepting TCP connection")
		}

		go NewClient(newConnection)
	}
}


func NewClient(conn *net.TCPConn) {

	fmt.Printf("new client: %s->%s\n",
		conn.RemoteAddr(),
		conn.LocalAddr())

	ircConn := NewIRCConn(conn, SERVERNAME)

	for ev := range ircConn.eventChannel() {
		fmt.Printf("CMD: %s %s\n", ev.command, ev.args)
		ircConn.sendCode(464)
	}

	fmt.Printf("client %s closed connection\n", conn.RemoteAddr())

}
