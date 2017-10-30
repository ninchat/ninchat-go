package main

/*
#include <stdlib.h>
#include <string.h>

#define NINCHAT_NO_PROTOTYPES
#include "../include/ninchat.h"

static inline void do_session_event_callback(
		uintptr_t cb,
		uintptr_t ctx,
		const uint8_t *params,
		size_t paramslen)
{
	((ninchat_session_event_callback) cb)((void *) ctx, (const char *) params, paramslen);
}

static inline void do_event_callback(
		uintptr_t cb,
		uintptr_t ctx,
		const uint8_t *params,
		size_t paramslen,
		uintptr_t payload,
		unsigned int payloadlen,
		bool lastreply)
{
	((ninchat_event_callback) cb)((void *) ctx, (const char *) params, paramslen, (ninchat_frame *) payload, payloadlen, lastreply);
}

static inline void do_close_callback(
		uintptr_t cb,
		uintptr_t ctx)
{
	((ninchat_close_callback) cb)((void *) ctx);
}

static inline void do_conn_state_callback(
		uintptr_t cb,
		uintptr_t ctx,
		const char *state)
{
	((ninchat_conn_state_callback) cb)((void *) ctx, state);
}

static inline void do_conn_active_callback(
		uintptr_t cb,
		uintptr_t ctx)
{
	((ninchat_conn_active_callback) cb)((void *) ctx);
}

static inline void do_log_callback(
		uintptr_t cb,
		uintptr_t ctx,
		const char *message,
		size_t message_len)
{
	((ninchat_log_callback) cb)((void *) ctx, message, message_len);
}

static inline size_t payload_part_size(
		uintptr_t vector,
		int index)
{
	return ((ninchat_frame *) vector)[index].size;
}

static inline void copy_payload_part(
		uintptr_t target_buffer,
		uintptr_t source_vector,
		int index)
{
	ninchat_frame *source_frame = &((ninchat_frame *) source_vector)[index];
	memcpy((void *) target_buffer, source_frame->data, source_frame->size);
}

*/
import "C"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"unsafe"

	ninchat "github.com/ninchat/ninchat-go"
)

var (
	sessionRefs uintptr
	sessions    map[uintptr]*ninchat.Session
)

//export ninchat_session_new
func ninchat_session_new() (ref uintptr) {
	if sessions == nil {
		sessions = make(map[uintptr]*ninchat.Session)
	}

	sessionRefs++
	ref = sessionRefs

	sessions[ref] = new(ninchat.Session)
	return
}

//export ninchat_session_delete
func ninchat_session_delete(ref uintptr) {
	delete(sessions, ref)
}

//export ninchat_session_on_session_event
func ninchat_session_on_session_event(ref uintptr, cb uintptr, ctx uintptr) {
	sessions[ref].OnSessionEvent = func(e *ninchat.Event) {
		params, err := json.Marshal(e.Params)
		if err != nil {
			panic(err)
		}
		C.do_session_event_callback(C.uintptr_t(cb), C.uintptr_t(ctx), (*C.uint8_t)(&params[0]), C.size_t(len(params)))
		runtime.KeepAlive(params)
	}
}

//export ninchat_session_on_event
func ninchat_session_on_event(ref uintptr, cb uintptr, ctx uintptr) {
	sessions[ref].OnEvent = func(e *ninchat.Event) {
		params, err := json.Marshal(e.Params)
		if err != nil {
			panic(err)
		}
		payload := (*reflect.SliceHeader)(unsafe.Pointer(&e.Payload))
		C.do_event_callback(C.uintptr_t(cb), C.uintptr_t(ctx), (*C.uint8_t)(&params[0]), C.size_t(len(params)), C.uintptr_t(payload.Data), C.uint(payload.Len), C.bool(e.LastReply))
		runtime.KeepAlive(params)
		runtime.KeepAlive(e.Payload)
	}
}

//export ninchat_session_on_close
func ninchat_session_on_close(ref uintptr, cb uintptr, ctx uintptr) {
	sessions[ref].OnClose = func() {
		C.do_close_callback(C.uintptr_t(cb), C.uintptr_t(ctx))
	}
}

//export ninchat_session_on_conn_state
func ninchat_session_on_conn_state(ref uintptr, cb uintptr, ctx uintptr) {
	sessions[ref].OnConnState = func(state string) {
		cState := C.CString(state)
		defer C.free(unsafe.Pointer(cState))

		C.do_conn_state_callback(C.uintptr_t(cb), C.uintptr_t(ctx), cState)
	}
}

//export ninchat_session_on_conn_active
func ninchat_session_on_conn_active(ref uintptr, cb uintptr, ctx uintptr) {
	sessions[ref].OnConnActive = func() {
		C.do_conn_active_callback(C.uintptr_t(cb), C.uintptr_t(ctx))
	}
}

//export ninchat_session_on_log
func ninchat_session_on_log(ref uintptr, cb uintptr, ctx uintptr) {
	sessions[ref].OnLog = func(fragments ...interface{}) {
		var buf bytes.Buffer

		for i, x := range fragments {
			fmt.Fprint(&buf, x)
			if i < len(fragments)-1 {
				buf.WriteString(" ")
			}
		}

		cMsg := C.CString(buf.String())
		defer C.free(unsafe.Pointer(cMsg))

		C.do_log_callback(C.uintptr_t(cb), C.uintptr_t(ctx), cMsg, C.size_t(buf.Len()))
	}
}

//export ninchat_session_set_params
func ninchat_session_set_params(ref uintptr, cParams unsafe.Pointer, cParamsLen C.size_t) {
	var params map[string]interface{}

	if err := json.Unmarshal(C.GoBytes(cParams, C.int(cParamsLen)), &params); err != nil {
		panic(err)
	}

	sessions[ref].SetParams(params)
}

//export ninchat_session_open
func ninchat_session_open(ref uintptr) {
	sessions[ref].Open()
}

//export ninchat_session_close
func ninchat_session_close(ref uintptr) {
	sessions[ref].Close()
}

//export ninchat_session_send
func ninchat_session_send(ref uintptr, cParams unsafe.Pointer, cParamsLen C.size_t, cPayload uintptr, cPayloadLen C.uint, actionIdPtr *C.int64_t) (cErr *C.char) {
	action := new(ninchat.Action)

	if err := json.Unmarshal(C.GoBytes(cParams, C.int(cParamsLen)), &action.Params); err != nil {
		panic(err)
	}

	if actionIdPtr == nil {
		action.Params["action_id"] = nil
	}

	if cPayloadLen > 0 {
		action.Payload = make([]ninchat.Frame, cPayloadLen)

		for i := range action.Payload {
			size := C.payload_part_size(C.uintptr_t(cPayload), C.int(i))
			action.Payload[i] = make([]byte, size)
			part := (*reflect.SliceHeader)(unsafe.Pointer(&action.Payload[i]))
			C.copy_payload_part(C.uintptr_t(part.Data), C.uintptr_t(cPayload), C.int(i))
		}
	}

	if err := sessions[ref].Send(action); err != nil {
		cErr = C.CString(err.Error())
		return
	}

	if actionIdPtr != nil {
		if x, found := action.Params["action_id"]; found {
			*actionIdPtr = C.int64_t(x.(int64))
		}
	}
	return
}

func main() {}
