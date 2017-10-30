GO	:= go
GOPATH	:= $(PWD)

export GOPATH

check: check-client check-api check-message check-c

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

lib/libninchat.a: $(wildcard c/*.go include/*.h)
	@ mkdir -p lib
	$(GO) get github.com/gorilla/websocket
	$(GO) build -buildmode=c-archive -o $@ c/library.go

bin/c-test: $(wildcard c/test/*.c include/*.h) lib/libninchat.a
	@ mkdir -p bin
	$(CC) $(CPPFLAGS) $(CFLAGS) -pthread -Wall -Wextra -Wno-unused-parameter -Wno-missing-field-initializers -g -o $@ c/test/test.c lib/libninchat.a

check-c: bin/c-test
	bin/c-test

clean:
	rm -rf bin lib pkg

.PHONY: check check-client check-client-go check-client-js check-api check-message check-c clean
