package main

type LilyMessage struct {
	raw string
}

func NewLilyMessage(line string) *LilyMessage {
	return &LilyMessage {raw: line}
}
