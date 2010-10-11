package main

type Dispatcher struct {
	ircConn *IRCConn
	lilyConn *LilyConn

	ircNick string
	ircPass string
}

func NewDispatcher(ircConn *IRCConn, lilyConn *LilyConn) *Dispatcher {
	return &Dispatcher{ircConn: ircConn,
			   lilyConn: lilyConn }
}

func (dis *Dispatcher) Dispatch() {

	for !closed(dis.ircConn.MessageChannel()) &&
	    !closed(dis.lilyConn.MessageChannel()) {
		select {
			case message := <-dis.ircConn.MessageChannel():
				if(message == nil) { break }
				dis.DispatchIRC(message)
		}
	}

	close(dis.ircConn.MessageChannel())
	close(dis.lilyConn.MessageChannel()) 

	dis.ircConn.Close()
	dis.lilyConn.Close()
}

func (dis *Dispatcher) DispatchIRC(message *IRCMessage) {
	switch(message.command) {
		case "PASS":
			dis.ircPass = message.args[0]
		case "NICK":
			dis.ircNick = message.args[0]
		case "USER":
			//Send user and pass to lily
			dis.lilyConn.MessageChannel() <- &LilyMessage{raw: dis.ircNick}
			dis.lilyConn.MessageChannel() <- &LilyMessage{raw: dis.ircPass}
	}

}
