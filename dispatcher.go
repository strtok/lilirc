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
	for !closed(dis.ircConn.incomingChannel) &&
	    !closed(dis.lilyConn.incomingChannel) {
		select {
			case message := <-dis.ircConn.incomingChannel:
				if message == nil { break }
				dis.DispatchIRC(message)

			case message := <-dis.lilyConn.incomingChannel:
				if message == nil { break }
				dis.DispatchLily(message)
		}
	}

	//Closing the connections below should
	//cause a ripple effect, but we close the channels
	//anyway for good measure
	close(dis.ircConn.incomingChannel)
	close(dis.lilyConn.incomingChannel)

	dis.ircConn.Close()
	dis.lilyConn.Close()
}

func (dis *Dispatcher) DispatchIRC(message *IRCMessage) {
	switch message.command {
		case "PASS":
			dis.ircPass = message.args[0]
		case "NICK":
			dis.ircNick = message.args[0]
		case "USER":
			//Send user and pass to lily
			dis.lilyConn.Send(dis.ircNick)
			dis.lilyConn.Send(dis.ircPass)
		case "PRIVMSG":
			dis.lilyConn.Send(message.target + ";" + message.text)
	}

}

func (dis *Dispatcher) DispatchLily(message *LilyMessage) {
	switch message.command {
		//TODO: We may receive this even after connected!
		case "CONNECTED":
			dis.ircConn.outgoingChannel <- &IRCMessage { raw: ":" + SERVERNAME + " 001 " + dis.ircNick + " :Login successful!" }
		case "NOTIFY":
			if event, present := message.attributes["EVENT"] ; present {
				switch event {
					case "private":
						dis.ircConn.outgoingChannel <- &IRCMessage { raw: ":" + message.source + " PRIVMSG " + dis.ircNick + " :" + message.attributes["VALUE"] }
				}
			}
	}
}
