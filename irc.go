package main

import "io"
import "net"
import "net/textproto"

type IRCConn struct {
	clientConn *net.TCPConn
}

func NewIRCConn(conn *net.TCPConn) *IRCConn {
	return &IRCConn {
		conn,
	}
}

func ReadLineIter(conn io.ReadWriteCloser) chan string {
 
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