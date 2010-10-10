include $(GOROOT)/src/Make.inc

TARG=lilirc
GOFILES=\
	irc.go\
	irc_message.go\
	lilirc.go\
	lily.go\
	util.go\

include $(GOROOT)/src/Make.cmd
