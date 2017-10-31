GO	:= go
GOPATH	:= $(PWD)

export GOPATH

check: check-client check-api check-message

check-client: check-client-go check-client-js

check-client-go:
	$(GO) get github.com/gorilla/websocket
	$(GO) vet .
	$(GO) test -v .

check-client-js: bin/gopherjs
	bin/gopherjs build

check-api:
	$(GO) get github.com/tsavola/pointer
	$(GO) vet ./ninchatapi
	$(GO) test -v ./ninchatapi

check-message:
	$(GO) vet ./ninchatmessage
	$(GO) test -v ./ninchatmessage

bin/gopherjs:
	$(GO) get github.com/gopherjs/gopherjs
	$(GO) build -o bin/gopherjs github.com/gopherjs/gopherjs

clean:
	rm -rf bin pkg

.PHONY: check check-client check-client-go check-client-js check-api check-message clean
