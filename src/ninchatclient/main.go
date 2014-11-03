package main

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	namespace = "NinchatClient"
)

var (
	module = NewObject()
)

func main() {
	module.Set("call", Call)
	module.Set("newSession", NewSession)
	module.Set("stringifyFrame", StringifyFrame)

	js.Global.Set(namespace, module)
}
