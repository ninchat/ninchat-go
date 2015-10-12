GO	:= go
GOFMT	:= gofmt
GOPATH	:= $(PWD)

PYTHON	:= python

export GOPATH

check: check-client check-api check-message

check-client: check-client-go check-client-js

check-client-go:
	$(GO) get github.com/gorilla/websocket
	$(GOFMT) -d -s *.go
	$(GO) vet .
	$(GO) test -v .

check-client-js: bin/gopherjs
	bin/gopherjs build

check-api:
	$(GO) get github.com/tsavola/pointer
	$(GOFMT) -d -s ninchatapi/*.go
	$(GO) vet ./ninchatapi
	$(GO) test -v ./ninchatapi

check-message:
	$(GOFMT) -d -s ninchatmessage/*.go
	$(GO) vet ./ninchatmessage
	$(GO) test -v ./ninchatmessage

bin/gopherjs:
	$(GO) get github.com/gopherjs/gopherjs
	$(GO) build -o bin/gopherjs github.com/gopherjs/gopherjs

api:
	PYTHONPATH=src/github.com/ninchat/ninchat-python $(PYTHON) -B ninchatapi/generate.py
	$(GO) fmt ninchatapi/*_gen.go
	$(GO) vet ninchatapi/*_gen.go

clean:
	rm -rf bin pkg

.PHONY: check check-client check-client-go check-client-js check-api check-message api clean
