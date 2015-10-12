package ninchatapi_test

import (
	"testing"

	"github.com/tsavola/pointer"

	"."
)

func TestCall(t *testing.T) {
	action := ninchatapi.DescribeChannel{
		ChannelId: pointer.String("1p255nth008"),
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
