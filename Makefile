GO		:= go
GOFMT		:= gofmt
GOPATH		:= $(PWD)
GOPHERJS	:= bin/gopherjs

DOCKER		:= docker
DOCKER_TAG	:= ninchat-js
DOCKER_BROWSER	:= chromium-browser --disable-setuid-sandbox

export GOPATH

build: gen/ninchatclient.js gen/ninchatclient.min.js

gen/ninchatclient.js gen/ninchatclient.min.js: $(wildcard src/ninchatclient/*.go) $(GOPHERJS)
	@ mkdir -p gen
	$(GOPHERJS) build -o gen/ninchatclient.js ninchatclient
	$(GOPHERJS) build -m -o gen/ninchatclient.min.js ninchatclient
	$(GOFMT) -d -s src/ninchatclient
	$(GO) vet ninchatclient

$(GOPHERJS):
	$(GO) get bitbucket.org/kardianos/osext
	$(GO) get github.com/neelance/sourcemap
	$(GO) get github.com/spf13/cobra
	$(GO) get github.com/spf13/pflag
	$(GO) get golang.org/x/crypto/ssh/terminal
	$(GO) get golang.org/x/tools/go/exact
	$(GO) get golang.org/x/tools/go/gcimporter
	$(GO) get golang.org/x/tools/go/types
	$(GO) get gopkg.in/fsnotify.v1
	$(GO) build -o $@ github.com/gopherjs/gopherjs

clean:
	rm -rf bin
	rm -rf pkg
	rm -rf src/bitbucket.org
	rm -rf src/github.com/neelance
	rm -rf src/github.com/spf13
	rm -rf src/golang.org
	rm -rf src/gopkg.in

container-for-testing:
	$(DOCKER) build -t $(DOCKER_TAG) .

test-in-container:
	$(DOCKER) run -e DISPLAY=$(DISPLAY) -i --rm -t -v /tmp:/tmp -v $(PWD):/work $(DOCKER_TAG) $(DOCKER_BROWSER) file:///work/example/test.html

.PHONY: build clean container-for-testing test-in-container
