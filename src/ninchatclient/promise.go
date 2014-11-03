package main

import (
	"github.com/gopherjs/gopherjs/js"
)

const (
	promiseSuccessInvocationName = namespace + " promise success callback"
	promiseFailureInvocationName = namespace + " promise failure callback"
)

// PromiseResolver
type PromiseResolver func(successful bool, args ...interface{})

// Promise
type Promise struct {
	settled    bool
	successful bool
	args       []interface{}

	successCallbacks []js.Object
	failureCallbacks []js.Object
}

// NewPromise
func NewPromise() (promise js.Object, resolve PromiseResolver) {
	var p Promise

	promise = NewObject()
	promise.Set("then", p.Then)

	resolve = p.Resolve

	return
}

// Then implements the Promise.then(function|null[, function|null]) JavaScript
// API.
func (p *Promise) Then(success, failure js.Object) {
	hasSuccess := !success.IsUndefined() && !success.IsNull()
	hasFailure := !failure.IsUndefined() && !failure.IsNull()

	if p.settled {
		if p.successful {
			if hasSuccess {
				jsInvoke(promiseSuccessInvocationName, success, p.args...)
			}
		} else {
			if hasFailure {
				jsInvoke(promiseFailureInvocationName, failure, p.args...)
			}
		}
	} else {
		if hasSuccess {
			p.successCallbacks = append(p.successCallbacks, success)
		}

		if hasFailure {
			p.failureCallbacks = append(p.failureCallbacks, failure)
		}
	}
}

// Resolve
func (p *Promise) Resolve(successful bool, args ...interface{}) {
	if p.settled {
		panic("promise has already been settled")
	}

	p.settled = true
	p.successful = successful
	p.args = args

	if successful {
		for _, callback := range p.successCallbacks {
			jsInvoke(promiseSuccessInvocationName, callback, p.args...)
		}
	} else {
		for _, callback := range p.failureCallbacks {
			jsInvoke(promiseFailureInvocationName, callback, p.args...)
		}
	}
}
