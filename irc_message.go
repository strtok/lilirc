package main

import "strings"

type IRCMessage struct {
	raw string
	command string
	args []string
}

func NewIRCMessage(raw string) *IRCMessage {

	var newMessage IRCMessage
        newMessage.raw = raw

	tokens := strings.Split(raw, " ",  -1)
	newMessage.command = tokens[0]
	newMessage.args = tokens[1:]

	return &newMessage
}



