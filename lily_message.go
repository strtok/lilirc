package main

import "strings"
import "strconv"

type LilyMessage struct {
	raw string
	command string
	attributes map[string] string

	source string
	recip string
}

func NewLilyMessage(line string, userMap func(string) string) *LilyMessage {

	if !strings.HasPrefix(line, "%") {
		return &LilyMessage {raw: line }
	}

	tokens := strings.Split(line, " ", -1)
	command := strings.ToUpper(tokens[0][1:])

	lilyMessage := &LilyMessage { raw: line,
				      command: command,
				      attributes: make(map[string] string) }
	lilyMessage.ParseLilyMap()
	lilyMessage.MapNames(userMap)

	return lilyMessage
}

func (lilyMessage *LilyMessage) MapNames(userMap func(string) string) {
	if _, present := lilyMessage.attributes["SOURCE"] ; present {
		lilyMessage.source = userMap(lilyMessage.attributes["SOURCE"])
	}

	if _, present := lilyMessage.attributes["RECIPS"] ; present {
		lilyMessage.recip = userMap(lilyMessage.attributes["RECIPS"])
	}
}

func (lilyMessage *LilyMessage) ParseLilyMap() {

	mapString := lilyMessage.raw[strings.Index(lilyMessage.raw, " ") + 1:]

	const (
		KEY = iota
		VALUE
		VALUE_SIZE
		SEEK_NEXT_KEY
	)

	state := KEY

	var key string
	var value string

	keyStart := 0
	valueStart := 0
	valueSize := uint(0)

	for i,c := range mapString {
		switch state {
			case KEY:
				switch c {
					case ' ':
						keyStart = i + 1
					case '=':
						key = mapString[keyStart:i]

						valueStart = i + 1
						keyStart = 0
						state = VALUE
				}
			case VALUE:
				switch c {
					case '=':
						valueSize,_ = strconv.Atoui(mapString[valueStart:i])
						valueStart = i + 1
						state = VALUE_SIZE
					case ' ':
						value = mapString[valueStart:i]
						lilyMessage.attributes[key] = value
						valueStart = 0
						keyStart = 0
						state = SEEK_NEXT_KEY
				}
			case VALUE_SIZE:
				if uint(i - valueStart) + 1 == valueSize {
					value = mapString[valueStart:i + 1]
					lilyMessage.attributes[key] = value
					valueStart = 0
					keyStart = 0
					state = SEEK_NEXT_KEY
				}
			case SEEK_NEXT_KEY:
				if c != ' '  {
					keyStart = i
					state = KEY
				}
		}
	}

}
