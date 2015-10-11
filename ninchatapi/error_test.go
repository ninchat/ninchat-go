package ninchatapi_test

import (
	"testing"

	"."
)

func TestServerError(t *testing.T) {
	badChannelId := "1"

	action := ninchatapi.DescribeChannel{
		ChannelId: &badChannelId,
	}

	reply, err := action.Invoke(nil)
	if reply != nil {
		t.Error(reply)
	}

	if err == nil {
		t.Fatal("no err")
	}

	event, ok := err.(*ninchatapi.Error)
	if !ok {
		t.Fatal(err)
	}

	if event.String() != "error" {
		t.Fatal(event)
	}

	if event.ErrorType != "channel_not_found" {
		t.Error(event.ErrorType)
	}

	if event.ChannelId == nil {
		t.Error("no channel_id")
	} else if *event.ChannelId != badChannelId {
		t.Error(*event.ChannelId)
	}
}

func TestClientError(t *testing.T) {
	badAction := ninchatapi.DescribeChannel{}

	reply, err := badAction.Invoke(nil)
	if reply != nil {
		t.Fatal(reply)
	}
	if err == nil {
		t.Fail()
	}
	event, ok := err.(*ninchatapi.Error)
	if !ok {
		t.Fail()
	}
	if event.String() != "error" {
		t.Fatal(event)
	}
	if event.ErrorType != "request_malformed" {
		t.Fatal(event.ErrorType)
	}
	if event.ErrorReason == nil {
		t.Fail()
	}
	if *event.ErrorReason != "describe_channel action requires channel_id parameter" {
		t.Error(*event.ErrorReason)
	}
}
