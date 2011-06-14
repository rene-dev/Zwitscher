include $(GOROOT)/src/Make.inc

TARG=zwitscher
GOFILES=controller.go gui.go zwitscher.go

all:
	gomake -C gotter

install: all
	gomake -C gotter install

include $(GOROOT)/src/Make.cmd