GO	:= go
GOFMT	:= gofmt
GOPATH	:= $(PWD)

export GOPATH

check: check-go check-js

check-go:
	$(GO) get github.com/gorilla/websocket
	$(GOFMT) -d -s *.go ninchat*/*.go
	$(GO) vet . ./ninchatmessage
	$(GO) test -v . ./ninchatmessage

check-js: bin/gopherjs
	bin/gopherjs build

bin/gopherjs:
	$(GO) get github.com/gopherjs/gopherjs
	$(GO) build -o bin/gopherjs github.com/gopherjs/gopherjs

clean:
	rm -rf bin pkg

.PHONY: check check-go check-js clean
