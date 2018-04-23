GO ?= go
GOPATH := $(CURDIR)
GOBIN := bin
CROSS := GOOS=linux

all:
	@$(GO) install -v bin/
	@mv bin/bin bin/server

linux:
	@$(CROSS) $(GO) install -v bin/
	@mv bin/linux_amd64/bin bin/linux_amd64/server

gopath:
	@echo $(GOPATH)
