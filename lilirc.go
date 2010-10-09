package main

import "fmt"
import "net"
import "net/textproto"
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
		newConnection,err := listener.AcceptTCP()
		if err != nil {
			fmt.Print("failed accepting TCP connection")
		}
		go NewClient(newConnection)
   	}
}


func ReadLine(conn *net.TCPConn) chan string {
 
   ch := make(chan string)
   textConn := textproto.NewConn(conn)
   
   go func() {
      for {
         line, err := textConn.ReadLine()
         
         if err != nil {
            break
         }
         ch <- line
      }
      close(ch)
   }()
   
   return ch
}

func NewClient(conn *net.TCPConn) {
   fmt.Printf("new client: %s->%s\n", 
   			  conn.RemoteAddr(), 
   			  conn.LocalAddr())
   			  
   for line := range ReadLine(conn) {
      fmt.Printf("read: %s\n", line)
   } 
}