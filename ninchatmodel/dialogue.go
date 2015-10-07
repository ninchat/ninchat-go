package ninchatmodel

import (
	"reflect"

	"github.com/ninchat/ninchat-go"
	api "github.com/ninchat/ninchat-go/ninchatapi"
)

// Dialogue
type Dialogue struct {
	PeerId           string
	Status           string
	SelfMemberAttrs  *api.DialogueMemberAttrs
	PeerMemberAttrs  *api.DialogueMemberAttrs
	AudienceMetadata map[string]interface{}
	Window           MessageWindow
}

func (d *Dialogue) update(status string, selfMemberAttrs, peerMemberAttrs *api.DialogueMemberAttrs, audienceMetadata map[string]interface{}) (c Change) {
	if d.Status != status {
		d.Status = status
		c = Updated
	}

	if !reflect.DeepEqual(d.SelfMemberAttrs, selfMemberAttrs) {
		d.SelfMemberAttrs = selfMemberAttrs
		c = Updated
	}

	if !reflect.DeepEqual(d.PeerMemberAttrs, peerMemberAttrs) {
		d.PeerMemberAttrs = peerMemberAttrs
		c = Updated
	}

	if !reflect.DeepEqual(d.AudienceMetadata, audienceMetadata) {
		d.AudienceMetadata = audienceMetadata
		c = Updated
	}

	return
}

func (d *Dialogue) updateStatus(status string) (c Change) {
	if d.Status != status {
		d.Status = status
		c = Updated
	}
	return
}

func (d *Dialogue) updateStatusIfHigher(status string) (c Change) {
	if compareStatus(status, d.Status) > 0 {
		d.Status = status
		c = Updated
	}
	return
}

func (d *Dialogue) updateStatusIfLower(status string) (c Change) {
	if compareStatus(status, d.Status) < 0 {
		d.Status = status
		c = Updated
	}
	return
}

func (d *Dialogue) updateStatusIfRead(messageId string) (c Change) {
	if messageId >= d.Window.getLatestId() {
		c = d.updateStatusIfLower("visible")
	}
	return
}

func (d *Dialogue) newLoadHistoryAction() *api.LoadHistory {
	return &api.LoadHistory{
		UserId: &d.PeerId,
	}
}

// DialogueState
type DialogueState struct {
	Messages MessageState
	OnChange func(Change, *Dialogue)
	Map      map[string]*Dialogue

	session *ninchat.Session
}

func (state *DialogueState) init(session *ninchat.Session) {
	state.session = session
	state.Messages.init(session, "dialogue")
	state.Map = make(map[string]*Dialogue)
}

func (state *DialogueState) handleSessionStatus(e *api.SessionStatusUpdated) {
	if d := state.Map[*e.UserId]; d != nil {
		c := d.updateStatusIfRead(e.MessageId)

		state.log(d.PeerId, c.String(), "with status", d.Status, "by", e.String())

		if c != unchanged {
			state.OnChange(c, d)
		}
	} else {
		state.log(d.PeerId, "referenced by", e.String(), "is unknown")
	}
}

func (state *DialogueState) handleUser(selfId string, userDialogues map[string]*api.UserDialogue, eventName string) {
	var discard []*Dialogue

	for peerId, d := range state.Map {
		if userDialogues[peerId] == nil {
			discard = append(discard, d)
		}
	}

	for _, d := range discard {
		delete(state.Map, d.PeerId)

		state.log(d.PeerId, "removed by", eventName)

		state.OnChange(Removed, d)
	}

	for peerId, ud := range userDialogues {
		var c Change

		d := state.Map[peerId]
		if d != nil {
			c = d.update(*ud.DialogueStatus, ud.DialogueMembers[selfId], ud.DialogueMembers[peerId], ud.AudienceMetadata)
		} else {
			d = &Dialogue{
				PeerId:           peerId,
				Status:           *ud.DialogueStatus,
				SelfMemberAttrs:  ud.DialogueMembers[selfId],
				PeerMemberAttrs:  ud.DialogueMembers[peerId],
				AudienceMetadata: ud.AudienceMetadata,
			}
			state.Map[peerId] = d
			c = Added
		}

		state.log(peerId, c.String(), "with status", d.Status, "by", eventName)

		if c != unchanged {
			state.OnChange(c, d)
		}
	}
}

func (state *DialogueState) handleDialogue(selfId string, e *api.DialogueUpdated) {
	d := state.Map[e.UserId]

	if e.DialogueStatus != nil {
		var c Change

		if d != nil {
			c = d.update(*e.DialogueStatus, e.DialogueMembers[selfId], e.DialogueMembers[e.UserId], e.AudienceMetadata)
		} else {
			d = &Dialogue{
				PeerId:           e.UserId,
				Status:           *e.DialogueStatus,
				SelfMemberAttrs:  e.DialogueMembers[selfId],
				PeerMemberAttrs:  e.DialogueMembers[e.UserId],
				AudienceMetadata: e.AudienceMetadata,
			}
			state.Map[e.UserId] = d
			c = Added
		}

		state.log(e.UserId, c.String(), "with status", d.Status, "by", e.String())

		if c != unchanged {
			state.OnChange(c, d)
		}
	} else if d != nil {
		delete(state.Map, e.UserId)

		state.log(e.UserId, "removed by", e.String())

		state.OnChange(Removed, d)
	}
}

func (state *DialogueState) handleReceive(selfId string, e *api.MessageReceived) {
	status := "highlight"

	if e.MessageUserId != nil && *e.MessageUserId == selfId {
		status = "visible"
	}

	var c Change

	d := state.Map[*e.UserId]
	if d != nil {
		c = d.updateStatusIfHigher(status)
	} else {
		d = &Dialogue{
			PeerId:          *e.UserId,
			Status:          status,
			SelfMemberAttrs: new(api.DialogueMemberAttrs),
			PeerMemberAttrs: new(api.DialogueMemberAttrs),
		}
		state.Map[*e.UserId] = d
		c = Added
	}

	state.log(d.PeerId, c.String(), "with status", d.Status, "by", e.String())

	if c != unchanged {
		state.OnChange(c, d)
	}

	state.Messages.handleReceive(*e.UserId, &d.Window, e)
}

func (state *DialogueState) LoadEarlier(peerId string) <-chan error {
	status := "visible"

	var c Change

	d := state.Map[peerId]
	if d != nil {
		c = d.updateStatusIfHigher(status)
	} else {
		d = &Dialogue{
			PeerId:          peerId,
			Status:          status,
			SelfMemberAttrs: new(api.DialogueMemberAttrs),
			PeerMemberAttrs: new(api.DialogueMemberAttrs),
		}
		state.Map[peerId] = d
		c = Added
	}

	state.log(peerId, "status", c.String(), "with value", status)

	if c != unchanged {
		api.Send(state.session, &api.UpdateDialogue{
			UserId:         &peerId,
			DialogueStatus: &status,
		})

		state.OnChange(c, d)
	}

	return d.Window.loadEarlier(state.session, &api.LoadHistory{
		UserId: &peerId,
	})
}

func (state *DialogueState) UpdateStatus(d *Dialogue, status string) {
	c := d.updateStatus(status)

	state.log(d.PeerId, "status", c.String(), "with value", status)

	if c != unchanged {
		api.Send(state.session, &api.UpdateDialogue{
			UserId:         &d.PeerId,
			DialogueStatus: &status,
		})

		state.OnChange(c, d)
	}
}

func (state *DialogueState) Activate(d *Dialogue) {
	d.Window.activate(state.session, d)
}

func (state *DialogueState) Discard(d *Dialogue) {
	delete(state.Map, d.PeerId)

	if messageId := d.Window.getLatestId(); messageId != "" {
		api.Send(state.session, &api.DiscardHistory{
			UserId:    &d.PeerId,
			MessageId: &messageId,
		})
	}

	state.log(d.PeerId, "removed")

	state.OnChange(Removed, d)
}

func (state *DialogueState) log(fragments ...interface{}) {
	log(state.session, "dialogue:", fragments)
}
