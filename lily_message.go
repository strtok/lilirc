package main

import "strings"

type LilyMessage struct {
	raw string
	command string
}

func NewLilyMessage(line string) *LilyMessage {

	if !strings.HasPrefix(line, "%") {
		return &LilyMessage {raw: line }
	}

	tokens := strings.Split(line, " ", -1)
	command := strings.ToUpper(tokens[0][1:])

	return &LilyMessage {
			raw: line,
			command: command,
		}
}
