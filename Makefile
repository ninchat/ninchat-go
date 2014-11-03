GO		:= go
GOFMT		:= gofmt
GOPATH		:= $(PWD)
GOPHERJS	:= bin/gopherjs

export GOPATH

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

.PHONY: clean
