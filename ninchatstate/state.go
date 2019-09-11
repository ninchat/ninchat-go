package ninchatstate

import (
	"log"

	"github.com/ninchat/ninchat-go"
)

type State struct {
	session            *ninchat.Session
	userOnSessionEvent func(*ninchat.Event)
	userOnEvent        func(*ninchat.Event)
	userOnClose        func()

	created  bool
	channels map[string]*channel
}

// New takes ownership of the session object.  Don't touch it after this call.
func New(session *ninchat.Session) *State {
	s := &State{
		session:            session,
		userOnSessionEvent: session.OnSessionEvent,
		userOnEvent:        session.OnEvent,
		userOnClose:        session.OnClose,
		channels:           make(map[string]*channel),
	}
	session.OnSessionEvent = s.stateOnSessionEvent
	session.OnEvent = s.stateOnEvent
	session.OnClose = s.stateOnClose
	return s
}

func (s *State) SetParams(params map[string]interface{}) {
	s.session.SetParams(params)
}

func (s *State) Open() {
	s.session.Open()
}

func (s *State) Close() {
	s.session.Close()
}

func (s *State) Send(action *ninchat.Action) (err error) {
	err = s.session.Send(action)
	if err != nil {
		return
	}

	return
}

func (s *State) stateOnSessionEvent(e *ninchat.Event) {
	log.Printf("stateOnSessionEvent: %s", e)

	if e.String() != "session_created" {
		s.userOnSessionEvent(e)
		return
	}

	var joined []string
	var parted []string

	userChannels, _ := e.Map("user_channels")
	for id, info := range userChannels {
		if c, known := s.channels[id]; known {
			c.refresh(id, info.(map[string]interface{}), s)
		} else {
			c = new(channel)
			s.channels[id] = c
			joined = append(joined, id)
		}
	}
	for id := range s.channels {
		if _, exist := userChannels[id]; !exist {
			parted = append(parted, id)
		}
	}
	e.Params["user_channels"] = map[string]interface{}{}

	if !s.created {
		s.created = true
		s.userOnSessionEvent(e)
	}

	for _, id := range parted {
		c := s.channels[id]
		delete(s.channels, id)
		c.fakeParted(id, s)
	}
	for _, id := range joined {
		s.channels[id].fakeJoined(id, s)
	}
}

func (s *State) stateOnEvent(e *ninchat.Event) {
	log.Printf("stateOnEvent: %s", e)

	switch e.String() {
	case "channel_joined":
		id, _ := e.Str("channel_id")
		if _, exist := s.channels[id]; !exist {
			c := new(channel)
			s.channels[id] = c
			c.realJoined(id, s)
		}

	case "channel_parted":
		id, _ := e.Str("channel_id")
		delete(s.channels, id)

	case "message_received":
		id, _ := e.Str("channel_id")
		if c, exist := s.channels[id]; exist {
			c.messageReceived(e)
		}
	}

	s.userOnEvent(e)
}

func (s *State) stateOnClose() {
	log.Printf("stateOnClose")

	if s.userOnClose != nil {
		s.userOnClose()
	}
}

type channel struct {
	lastMessageID string
}

func (c *channel) realJoined(id string, s *State) {
	s.session.Send(&ninchat.Action{Params: map[string]interface{}{
		"action":     "load_history",
		"channel_id": id,
	}})
}

func (c *channel) fakeJoined(id string, s *State) {
	s.session.Send(&ninchat.Action{
		Params: map[string]interface{}{
			"action":     "describe_channel",
			"channel_id": id,
		},
		OnReply: func(e *ninchat.Event) {
			if e != nil && e.String() != "error" {
				params := map[string]interface{}{
					"event": "channel_joined",
				}
				for _, key := range []string{
					"channel_id",
					"channel_attrs",
					"channel_members",
					"message_time",
					"realm_id",
					"audience_metadata",
					"channel_metadata",
				} {
					if value, found := e.Params[key]; found {
						params[key] = value
					}
				}
				s.userOnEvent(&ninchat.Event{Params: params})

				c.realJoined(id, s)
			}
		},
	})
}

func (c *channel) fakeParted(id string, s *State) {
	s.userOnEvent(&ninchat.Event{
		Params: map[string]interface{}{
			"event":      "channel_parted",
			"channel_id": id,
		},
	})
}

func (c *channel) refresh(id string, info map[string]interface{}, s *State) {
	s.userOnEvent(&ninchat.Event{
		Params: map[string]interface{}{
			"event":         "channel_updated",
			"channel_id":    id,
			"channel_attrs": info["channel_attrs"],
			"realm_id":      info["realm_id"],
		},
	})

	c.refreshContinue(id, c.lastMessageID, s)
}

func (c *channel) refreshContinue(channelID, sinceMessageID string, s *State) {
	params := map[string]interface{}{
		"action":        "load_history",
		"channel_id":    channelID,
		"history_order": 1,
	}
	if sinceMessageID != "" {
		params["message_id"] = sinceMessageID
	}

	s.session.Send(&ninchat.Action{
		Params: params,
		OnReply: func(e *ninchat.Event) {
			if e != nil && e.LastReply && e.String() == "message_received" {
				messageID, _ := e.Str("message_id")
				c.refreshContinue(channelID, messageID, s)
			}
		},
	})
}

func (c *channel) messageReceived(e *ninchat.Event) {
	if _, action := e.Str("action_id"); !action {
		messageID, _ := e.Str("message_id")
		if messageID > c.lastMessageID {
			c.lastMessageID = messageID
		}
	}
}
