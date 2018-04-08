GO ?= go
GOPATH := $(CURDIR)
GOBIN := bin

all:
	@$(GO) install bin/
	@mv bin/bin bin/server

gopath:
	@echo $(GOPATH)
