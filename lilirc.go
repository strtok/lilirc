package main

import "flag"
import "log"
import "os"
import "net"

var (
	SERVERNAME = "irc.lily.com"
	LILYADDRESS = "thales.strtok.net:7777"
)

var logger *log.Logger;

func main() {
	// read arguments
	listen_addr := flag.String("listen", "0.0.0.0:6667",
			"Port to listen for IRC connections on")
	log_file := flag.String("log-file", "/dev/stdout", "Where to write log files")
	flag.Parse()

	// Set up logging
	var outFile *os.File
	if *log_file == "/dev/stdout" {
		outFile = os.Stdout
	} else {
		var err os.Error
		outFile, err = os.Open(*log_file, os.O_WRONLY | os.O_CREAT, 0755)
		if err != nil {
			log.Crash("Oh noes!")
		}
	}
	logger = log.New(outFile, nil, "", log.Ldate | log.Ltime)

	listenAddress, err := net.ResolveTCPAddr(*listen_addr)

	if err != nil {
		logger.Log("failed creating listenAddress")
	}

	listener, err := net.ListenTCP("tcp", listenAddress)

	if err != nil {
		logger.Log("failed listening")
	}

	for {

		newConnection, err := listener.AcceptTCP()

		if err != nil {
			logger.Log("failed accepting TCP connection")
		}

		go StartGateway(newConnection)
	}
}


func StartGateway(conn *net.TCPConn) {

	logger.Logf("new client: %s->%s\n",
		conn.RemoteAddr(),
		conn.LocalAddr())

	ircConn := NewIRCConn(conn, SERVERNAME)
	lilyConn := NewLilyConn(LILYADDRESS)

	if ircConn == nil || lilyConn == nil {
		return
	}

	dis := NewDispatcher(ircConn, lilyConn)
	dis.Dispatch()

	logger.Log("ending session")
}
