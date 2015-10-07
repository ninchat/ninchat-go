package ninchatmodel

import (
	"sort"
	"strings"
	"sync/atomic"

	"github.com/ninchat/ninchat-go"
	api "github.com/ninchat/ninchat-go/ninchatapi"
)

const (
	missingLoadOverlap = 10
	missingLoadLength  = 1000
)

// messageGroup
type messageGroup interface {
	newLoadHistoryAction() *api.LoadHistory
}

// MessageWindow
type MessageWindow struct {
	Aux

	active       bool
	missing      bool
	missingSince string
	ids          []string
	minLoadedId  string
	maxLoadedId  string
	earliest     int32 // atomic boolean
}

func (w *MessageWindow) GetAux(key interface{}) interface{} {
	// for GopherJS
	return w.Aux.GetAux(key)
}

func (w *MessageWindow) SetAux(key, value interface{}) {
	// for GopherJS
	w.Aux.SetAux(key, value)
}

func (w *MessageWindow) IsActive() bool {
	return w.active
}

func (w *MessageWindow) activate(session *ninchat.Session, group messageGroup) {
	if !w.active {
		w.active = true

		if w.missing {
			w.loadMissing(session, group)
		}
	}
}

func (w *MessageWindow) Deactivate() {
	w.active = false
}

func (w *MessageWindow) GetLength() int {
	return len(w.ids)
}

func (w *MessageWindow) HasEarliest() bool {
	return atomic.LoadInt32(&w.earliest) != 0
}

func (w *MessageWindow) gotEarliest() {
	atomic.StoreInt32(&w.earliest, 1)
}

func (w *MessageWindow) getLatestId() (id string) {
	if len(w.ids) > 0 {
		id = w.ids[len(w.ids)-1]
	}
	return
}

func (w *MessageWindow) indexOf(id string) (i int, found bool) {
	i = sort.SearchStrings(w.ids, id)
	found = i < len(w.ids) && w.ids[i] == id
	return
}

func (w *MessageWindow) handleSecondarySession(session *ninchat.Session, group messageGroup) {
	if !w.missing {
		w.missing = true
		w.missingSince = w.maxLoadedId

		if len(w.ids) >= missingLoadOverlap {
			id := w.ids[len(w.ids)-missingLoadOverlap]
			if strings.Compare(id, w.missingSince) > 0 {
				w.missingSince = id
			}
		}
	}

	if w.active {
		w.loadMissing(session, group)
	}
}

func (w *MessageWindow) prepareReceive(e *api.MessageReceived) (ok bool) {
	if e.HistoryLength != nil {
		if w.minLoadedId == "" || strings.Compare(e.MessageId, w.minLoadedId) < 0 {
			w.minLoadedId = e.MessageId
		}

		if strings.Compare(e.MessageId, w.maxLoadedId) > 0 {
			w.maxLoadedId = e.MessageId
		}
	}

	if i, found := w.indexOf(e.MessageId); !found {
		w.ids = append(w.ids[:i], append([]string{e.MessageId}, w.ids[i:]...)...)
		ok = true
	}
	return
}

// loadEarlier fills in action.MessageId.  The caller must have set
// action.ChannelId or action.UserId.  The caller may set action.HistoryLength.
func (w *MessageWindow) loadEarlier(session *ninchat.Session, action *api.LoadHistory) <-chan error {
	loaded := make(chan error, 1)

	if w.HasEarliest() {
		close(loaded)
	} else {
		if w.minLoadedId != "" {
			id := w.minLoadedId
			action.MessageId = &id
		}

		go func() {
			var err error

			defer func() {
				if err != nil {
					loaded <- err
				}
				close(loaded)
			}()

			e, err := action.Invoke(session)
			if err != nil {
				panic(err)
			}
			if e == nil {
				return
			}

			if e.HistoryLength == 0 {
				w.gotEarliest()
			}
		}()
	}

	return loaded
}

func (w *MessageWindow) loadMissing(session *ninchat.Session, group messageGroup) {
	if w.missingSince == "" {
		w.loadEarlier(session, group.newLoadHistoryAction())
	} else {
		messageId := w.missingSince

		go func() {
			var (
				order  = 1
				length = missingLoadLength
			)

			for {
				action := group.newLoadHistoryAction()
				action.MessageId = &messageId
				action.HistoryOrder = &order
				action.HistoryLength = &length

				e, err := action.Invoke(session)
				if err != nil {
					// TODO: what?
					return
				}
				if e == nil {
					return
				}

				if e.HistoryLength == 0 {
					break
				}

				messageId = *e.MessageId
			}
		}()
	}

	w.missing = false
	w.missingSince = ""
}

// MessageState
type MessageState struct {
	OnReceive func(groupId string, event *api.MessageReceived)
	OnUpdate  func(groupId string, event *api.MessageUpdated)

	session   *ninchat.Session
	groupType string
}

func (state *MessageState) init(session *ninchat.Session, groupType string) {
	state.session = session
	state.groupType = groupType
}

func (state *MessageState) handleReceive(groupId string, w *MessageWindow, e *api.MessageReceived) {
	if w.prepareReceive(e) {
		state.log(groupId, e.MessageId, "received")

		state.OnReceive(groupId, e)
	} else {
		state.log(groupId, e.MessageId, "received again")
	}
}

func (state *MessageState) handleUpdate(groupId string, w *MessageWindow, e *api.MessageUpdated) {
	if _, found := w.indexOf(e.MessageId); found {
		state.log(groupId, e.MessageId, "updated with hidden set to", e.MessageHidden)

		state.OnUpdate(groupId, e)
	} else {
		state.log(groupId, e.MessageId, "referenced by", e.String(), "is unknown")
	}
}

func (state *MessageState) log(groupId string, fragments ...interface{}) {
	if state.session.OnLog != nil {
		fragments = append([]interface{}{state.groupType, groupId, "message:"}, fragments...)
		state.session.OnLog(fragments...)
	}
}
