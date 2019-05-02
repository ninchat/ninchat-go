GO		?= go
GOMOBILE	?= gomobile

LOCALGO		:= GOPATH=$(PWD) $(GO)

build-client-go:
	$(LOCALGO) get github.com/gorilla/websocket
	$(LOCALGO) vet .
	$(LOCALGO) build .

all: build-client-go lib/libninchat.a lib/ninchat-client.aar lib/NinchatClient.framework

check: all check-client check-api check-message check-c

check-client: check-client-go check-client-js

check-client-go: build-client-go
	$(LOCALGO) test -v .

check-client-js: bin/gopherjs
	bin/gopherjs build

check-api:
	$(LOCALGO) get github.com/tsavola/pointer
	$(LOCALGO) vet ./ninchatapi
	$(LOCALGO) test -v ./ninchatapi

check-message:
	$(LOCALGO) vet ./ninchatmessage
	$(LOCALGO) test -v ./ninchatmessage

bin/gopherjs: Makefile
	$(LOCALGO) get github.com/gopherjs/gopherjs
	$(LOCALGO) build -o bin/gopherjs github.com/gopherjs/gopherjs

lib/libninchat.a: $(wildcard *.go c/*.go include/*.h) Makefile
	@ mkdir -p lib
	$(LOCALGO) get github.com/gorilla/websocket
	$(LOCALGO) build -buildmode=c-archive -o $@ c/library.go

bin/c-test: $(wildcard c/test/*.c include/*.h) lib/libninchat.a Makefile
	@ mkdir -p bin
	$(CC) $(CPPFLAGS) $(CFLAGS) -pthread -Wall -Wextra -Wno-unused-parameter -Wno-missing-field-initializers -g -o $@ c/test/test.c lib/libninchat.a

check-c: bin/c-test
	bin/c-test

lib/ninchat-client.aar: $(wildcard *.go mobile/*.go) Makefile
	@ mkdir -p lib
	$(GOMOBILE) bind -target=android -javapkg=com.ninchat -o $@ github.com/ninchat/ninchat-go/mobile

lib/NinchatClient.framework: $(wildcard *.go mobile/*.go) Makefile
	@ mkdir -p lib
	$(GOMOBILE) bind -target=ios -prefix=Ninchat -o $@ github.com/ninchat/ninchat-go/mobile

clean:
	rm -rf bin lib pkg

.PHONY: all build-client-go check check-client check-client-go check-client-js check-api check-message check-c clean
