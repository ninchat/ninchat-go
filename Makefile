GO		?= go
GOMOBILE	?= gomobile

build-client-go:
	$(GO) vet .
	$(GO) build .

all: build-client-go lib/libninchat.a lib/ninchat-client.aar lib/NinchatLowLevelClient.framework/NinchatLowLevelClient

check: all check-client check-api check-message check-c

check-client: check-client-go check-client-js

check-client-go: build-client-go
	$(GO) test -v .

check-client-js: bin/gopherjs
	bin/gopherjs build

check-api:
	$(GO) vet ./ninchatapi
	$(GO) test -v ./ninchatapi

check-message:
	$(GO) vet ./ninchatmessage
	$(GO) test -v ./ninchatmessage

bin/gopherjs: Makefile
	$(GO) get -d github.com/gopherjs/gopherjs
	$(GO) build -o bin/gopherjs github.com/gopherjs/gopherjs

lib/libninchat.a: $(wildcard *.go c/*.go include/*.h) Makefile
	@ mkdir -p lib
	$(GO) build -buildmode=c-archive -o $@ c/library.go

bin/c-test: $(wildcard c/test/*.c include/*.h) lib/libninchat.a Makefile
	@ mkdir -p bin
	$(CC) $(CPPFLAGS) $(CFLAGS) -pthread -Wall -Wextra -Wno-unused-parameter -Wno-missing-field-initializers -g -o $@ c/test/test.c lib/libninchat.a

check-c: bin/c-test
	bin/c-test

lib/ninchat-client.aar: $(wildcard *.go mobile/*.go) Makefile
	@ mkdir -p lib
	$(GOMOBILE) bind -target=android -javapkg=com.ninchat -o $@ github.com/ninchat/ninchat-go/mobile

lib/NinchatLowLevelClient.framework/NinchatLowLevelClient: $(wildcard *.go mobile/*.go) Makefile
	@ mkdir -p lib
	$(GOMOBILE) bind -target=ios -prefix=NINLowLevel -o $(patsubst %/,%,$(dir $@)) github.com/ninchat/ninchat-go/mobile	
	mv $(dir $@)/Versions/Current/Client $(dir $@)/Versions/Current/NinchatLowLevelClient
	rm $(dir $@)/Client
	ln -sf Versions/Current/NinchatLowLevelClient $@

clean:
	rm -rf bin lib pkg

.PHONY: all build-client-go check check-client check-client-go check-client-js check-api check-message check-c clean
