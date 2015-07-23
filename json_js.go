// +build js

package ninchat

import (
	"errors"

	"github.com/gopherjs/gopherjs/js"
)

func jsonMarshal(v map[string]interface{}) (data *js.Object, err error) {
	defer func() {
		err = jsError(recover())
	}()

	data = js.Global.Get("JSON").Call("stringify", v)
	return
}

func jsonUnmarshalArray(data *js.Object, v *[]interface{}) (err error) {
	object, err := jsonParse(data)
	if err != nil {
		return
	}

	x, ok := object.([]interface{})
	if !ok {
		err = errors.New("json: cannot unmarshal value into Go value of type *[]interface{}")
		return
	}

	*v = x
	return
}

func jsonUnmarshalObject(data *js.Object, v *map[string]interface{}) (err error) {
	object, err := jsonParse(data)
	if err != nil {
		return
	}

	x, ok := object.(map[string]interface{})
	if !ok {
		err = errors.New("json: cannot unmarshal value into Go value of type *map[string]interface{}")
		return
	}

	*v = x
	return
}

func jsonParse(data *js.Object) (x interface{}, err error) {
	defer func() {
		err = jsError(recover())
	}()

	x = js.Global.Get("JSON").Call("parse", data).Interface()
	return
}
