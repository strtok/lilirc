package main

import "strings"

type IRCMessage struct {
	raw	string
	command string
	target 	string
	args []	string
	text	string
}

func NewIRCMessage(raw string) *IRCMessage {

	var newMessage IRCMessage

        newMessage.raw = raw

	msgPos := strings.Index(raw, ":")

	if msgPos != -1 {

		//Split out args and text
		textSplit := strings.Split(raw, ":", 2)

		//IRCMessage.text is everything after the first :
		newMessage.text = textSplit[1]

		//Split everything before : into args
		tokens := strings.Split(textSplit[0], " ", -1)
		newMessage.command = tokens[0]
		newMessage.args = tokens[1:]
	} else {
		tokens := strings.Split(newMessage.raw, " ",  -1)
		newMessage.command = tokens[0]
		newMessage.args = tokens[1:]
	}

	return &newMessage
}



