package main

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

func jsError(x interface{}) (err error) {
	if x == nil {
		return
	}

	if jsErr, ok := x.(*js.Error); ok {
		msg := jsErr.Get("message").Str()
		if msg == "" {
			msg = "error"
		}

		err = errors.New(msg)
		return
	}

	err = x.(error)
	return
}

func jsInvoke(name string, function js.Object, args ...interface{}) (ok bool) {
	defer func() {
		if err := jsError(recover()); err != nil {
			println(name + " invocation error: " + err.Error())
		}
	}()

	function.Invoke(args...)

	ok = true
	return
}

func Atob(string js.Object) (binary js.Object, err error) {
	defer func() {
		err = jsError(recover())
	}()

	binary = js.Global.Call("atob", string)
	return
}

func ParseDataURI(string js.Object) (base64 js.Object, err error) {
	defer func() {
		err = jsError(recover())
	}()

	base64 = string.Call("split", ",").Index(1)
	return
}

func NewArray() js.Object {
	return js.Global.Get("Array").New()
}

func NewArrayBuffer(length int) js.Object {
	return js.Global.Get("ArrayBuffer").New(length)
}

func NewUint8Array(arrayBuffer js.Object) js.Object {
	return js.Global.Get("Uint8Array").New(arrayBuffer)
}

func NewObject() js.Object {
	return js.Global.Get("Object").New()
}

func EncodeURIComponent(s string) string {
	return js.Global.Call("encodeURIComponent", s).Str()
}

func ParseJSON(json string) (object js.Object, err error) {
	defer func() {
		err = jsError(recover())
	}()

	object = js.Global.Get("JSON").Call("parse", json)
	return
}

func StringifyJSON(object interface{}) (json string, err error) {
	defer func() {
		err = jsError(recover())
	}()

	json = js.Global.Get("JSON").Call("stringify", object).Str()
	return
}

func Random() float64 {
	return js.Global.Get("Math").Call("random").Float()
}

func SetTimeout(callback func(), timeout Duration) (id js.Object) {
	return js.Global.Call("setTimeout", callback, timeout)
}

func ClearTimeout(id js.Object) {
	js.Global.Call("clearTimeout", id)
}
