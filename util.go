package main

import "io"
import "net/textproto"
import "time"

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

func Timer(ms int64) chan bool {
	ch := make(chan bool, 1)
	go func() {
		time.Sleep(ms * 1e6)
		ch <- true
	}()
	return ch
}

func IRCToLily(name string) string {
	if name[0] == '#' {
		return "-" + name[1:]
	}
	return name
}

func LilyToIRC(name string) string {
	if name[0] == '-' {
		return "#" + name[1:]
	}
	return name
}
