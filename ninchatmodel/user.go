package ninchatmodel

import (
	"reflect"

	"github.com/ninchat/ninchat-go"
	api "github.com/ninchat/ninchat-go/ninchatapi"
)

// User
type User struct {
	Aux
	Id    string
	Attrs *api.UserAttrs
}

func newUser(id string, attrs *api.UserAttrs) *User {
	return &User{
		Id:    id,
		Attrs: attrs,
	}
}

func (user *User) update(attrs *api.UserAttrs) (c Change) {
	if !reflect.DeepEqual(user.Attrs, attrs) {
		user.Attrs = attrs
		c = Updated
	}

	return
}

// UserState
type UserState struct {
	OnChange func(Change, *User)
	Map      map[string]*User

	session *ninchat.Session
}

func (state *UserState) init(session *ninchat.Session) {
	state.session = session
	state.Map = make(map[string]*User)
}

func (state *UserState) handle(id string, attrs *api.UserAttrs, eventName string) {
	var c Change

	user := state.Map[id]
	if user != nil {
		c = user.update(attrs)
	} else {
		user = newUser(id, attrs)
		state.Map[id] = user
		c = Added
	}

	state.log(id, c.String(), "by", eventName)

	if c != unchanged {
		state.OnChange(c, user)
	}

	return
}

func (state *UserState) discover(id string, rediscover bool) {
	if rediscover || state.Map[id] == nil {
		api.Send(state.session, &api.DescribeUser{
			UserId: &id,
		})
	}
}

func (state *UserState) log(fragments ...interface{}) {
	log(state.session, "user:", fragments)
}
