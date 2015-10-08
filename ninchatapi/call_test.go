package ninchatapi_test

import (
	"testing"

	"."
)

func TestCall(t *testing.T) {
	channelId := "1p255nth008"

	action := ninchatapi.DescribeChannel{
		ChannelId: &channelId,
	}

	reply, err := action.Invoke(nil)
	if err != nil {
		t.Fatal(err)
	}

	if reply.RealmId != nil {
		t.Log(reply, reply.ChannelId, "realm", *reply.RealmId)
	} else {
		t.Log(reply, reply.ChannelId)
	}
}
