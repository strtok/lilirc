package main

import "fmt"
import "net"
import "net/textproto"

var (
   SERVERNAME = "irc.lily.com"
)

func newClient()
func sendIRCPreamble(conn *net.TCPConn)

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
		newConnection,err := listener.AcceptTCP()
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
   			  
	textConn := textproto.NewConn(conn)
    textConn.PrintfLine(":%s 001 :Welcome to lilirc!", SERVERNAME)
     
}

