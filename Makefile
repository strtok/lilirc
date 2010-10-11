include $(GOROOT)/src/Make.inc

TARG=lilirc
GOFILES=\
	dispatcher.go\
	irc.go\
	irc_message.go\
	lilirc.go\
	lily.go\
	lily_message.go\
	util.go\

include $(GOROOT)/src/Make.cmd
