package main

import "fmt"
import "net"

var (
	SERVERNAME = "irc.lily.com"
)

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

		go StartGateway(newConnection)
	}
}


func StartGateway(conn *net.TCPConn) {

	fmt.Printf("new client: %s->%s\n",
		conn.RemoteAddr(),
		conn.LocalAddr())

	ircConn := NewIRCConn(conn, SERVERNAME)

	for !closed(ircConn.eventChannel()) {
		select {
			case ircEvent := <-ircConn.eventChannel():
				if(ircEvent == nil) { break }
				fmt.Printf("CMD: %s %s\n", ircEvent.command, ircEvent.args)
		}
	}

	fmt.Printf("client %s closed connection\n", conn.RemoteAddr())

}
