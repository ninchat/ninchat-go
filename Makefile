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
	$(GO) get -u github.com/gopherjs/gopherjs

clean:
	rm -rf bin
	rm -rf pkg
	rm -rf src/bitbucket.org
	rm -rf src/code.google.com
	rm -rf src/github.com
	rm -rf src/gopkg.in

container-for-testing:
	$(DOCKER) build -t $(DOCKER_TAG) .

test-in-container:
	$(DOCKER) run -e DISPLAY=$(DISPLAY) -i --rm -t -v /tmp:/tmp -v $(PWD):/work $(DOCKER_TAG) $(DOCKER_BROWSER) file:///work/example/test.html

.PHONY: build clean container-for-testing test-in-container
