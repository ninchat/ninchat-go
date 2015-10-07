GO	:= go
GOFMT	:= gofmt
GOPATH	:= $(PWD)

PYTHON	:= python

export GOPATH

check: check-go check-js

check-go:
	$(GO) get github.com/gorilla/websocket
	$(GOFMT) -d -s *.go ninchat*/*.go
	$(GO) vet . ./ninchatmessage ./ninchatapi
	$(GO) test -v . ./ninchatmessage ./ninchatapi

check-js: bin/gopherjs
	bin/gopherjs build

bin/gopherjs:
	$(GO) get github.com/gopherjs/gopherjs
	$(GO) build -o bin/gopherjs github.com/gopherjs/gopherjs

api:
	PYTHONPATH=src/github.com/ninchat/ninchat-python $(PYTHON) -B ninchatapi/generate.py
	$(GO) fmt ninchatapi/*_gen.go
	$(GO) vet ninchatapi/*_gen.go

clean:
	rm -rf bin pkg

.PHONY: check check-go check-js api clean
