include $(GOROOT)/src/Make.inc

TARG=lilirc
GOFILES=\
	irc.go\
	lilirc.go\
	util.go\

include $(GOROOT)/src/Make.cmd
