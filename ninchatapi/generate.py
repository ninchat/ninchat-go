#!/usr/bin/env python

import cStringIO as stringio
import os
import subprocess
import sys

sys.dont_write_bytecode = True
sys.path.insert(0, "../ninchat-python")

import ninchat.api

GO_COMMAND = "go"

requiredtypes = {
    "bool":         "bool",
    "float":        "float64",
    "int":          "int",
    "object":       "map[string]interface{}",
    "string array": "[]string",
    "string":       "string",
    "time":         "int",
    None:           "interface{}",
}

optionaltypes = {
    "bool":         "bool",
    "float":        "*float64",
    "int":          "*int",
    "object":       "map[string]interface{}",
    "string array": "[]string",
    "string":       "*string",
    "time":         "*int",
    None:           "interface{}",
}

checks = {
    "bool":         "x",
    "float":        "x != nil",
    "int":          "x != nil",
    "object":       "x != nil",
    "string":       "x != nil",
    "string array": "x != nil",
    "time":         "x != nil",
}

evals = {
    "bool":         "x",
    "float":        "*x",
    "int":          "*x",
    "object":       "x",
    "string":       "*x",
    "string array": "x",
    "time":         "*x",
}

payloadactions = {
    "send_message",
    "update_user",
}

payloadevents = {
    "message_received",
}

attrobjects = [
    (ninchat.api.attrs.channel,        "Channel",        "channel"),
    (ninchat.api.attrs.channelmember,  "ChannelMember",  "channel membership"),
    (ninchat.api.attrs.dialoguemember, "DialogueMember", "dialogue membership"),
    (ninchat.api.attrs.identity,       "Identity",       "identity"),
    (ninchat.api.attrs.puppet,         "Puppet",         "puppet"),
    (ninchat.api.attrs.queue,          "Queue",          "queue"),
    (ninchat.api.attrs.queuemember,    "QueueMember",    "queue membership"),
    (ninchat.api.attrs.realm,          "Realm",          "realm"),
    (ninchat.api.attrs.realmmember,    "RealmMember",    "realm membership"),
    (ninchat.api.attrs.tag,            "Tag",            "tag"),
    (ninchat.api.attrs.user,           "User",           "user"),
]

unaryreplies = {
    "accept_audience": "dialogue_updated",
    "add_member": "member_joined",          # pseudo-event implemented manually
    "close_session": None,
    "create_access": "access_created",
    "create_channel": "channel_joined",
    "create_identity": "identity_created",
    "create_master_key": "master_key_created",
    "create_queue": "queue_created",
    "create_realm": "realm_joined",
    "create_tag": "tag_created",
    "create_user": "user_created",
    "delete_identity": "identity_deleted",
    "delete_master_key": "master_key_deleted",
    "delete_queue": "queue_deleted",
    "delete_queue_transcripts": "queue_transcripts_deleted",
    "delete_realm": "realm_deleted",
    "delete_tag": "tag_deleted",
    "delete_transcript": "transcript_deleted",
    "delete_user": "user_deleted",
    "describe_access": "access_found",
    "describe_channel": "channel_found",
    "describe_file": "file_found",
    "describe_identity": "identity_found",
    "describe_master_keys": "master_keys_found",
    "describe_queue": "queue_found",
    "describe_queue_transcripts": "queue_transcripts_found",
    "describe_realm": "realm_found",
    "describe_realm_queues": "realm_queues_found",
    "describe_tag": "tag_found",
    "describe_tags": "tags_found",
    "describe_user": "user_found",
    "discard_history": "history_discarded",
    "follow_channel": "channel_found",
    "get_queue_stats": "queue_stats_contents",
    "get_transcript": "transcript_contents",
    "join_channel": "channel_joined",
    "load_history": "history_results",      # message_received events ignored
    "part_channel": "channel_parted",
    "ping": "pong",
    "remove_member": "member_parted",       # pseudo-event implemented manually
    "request_audience": "audience_enqueued",
    "request_identity_auth_reset_access": "access_created",
    "request_identity_verify_access": "access_created",
    "reset_identity_auth": "identity_updated",
    "resume_session": None,
    "search": None,
    "send_access": "access_found",
    "send_file": "message_received",
    "send_message": "message_received",
    "track": "ack",
    "update_channel": "channel_updated",
    "update_dialogue": "dialogue_updated",
    "update_identity": "identity_updated",
    "update_identity_auth": "identity_updated",
    "update_member": "member_updated",      # pseudo-event implemented manually
    "update_message": "message_updated",
    "update_queue": "queue_updated",
    "update_realm": "realm_updated",
    "update_session": None,
    "update_tag": "tag_updated",
    "update_user": "user_updated",
    "update_user_messages": None,
    "verify_identity": "identity_updated",
}


class Output(object):

    def __init__(self, name):
        self.filename_tmp = os.path.join("ninchatapi", name + "_gen_tmp.go")
        self.filename = os.path.join("ninchatapi", name + "_gen.go")

    def __enter__(self):
        self.buf = stringio.StringIO()

        self.orig_stdout = sys.stdout
        sys.stdout = self.buf

    def __exit__(self, exc_type, exc_value, traceback):
        assert sys.stdout is self.buf
        sys.stdout = self.orig_stdout

        if not exc_type:
            with open(self.filename_tmp, "w") as f:
                f.write(self.buf.getvalue())

            subprocess.check_call([GO_COMMAND, "fmt", self.filename_tmp])
            subprocess.check_call([GO_COMMAND, "vet", self.filename_tmp])
            os.rename(self.filename_tmp, self.filename)


def main():
    with Output("actions"):
        print_header()

        print
        print 'import ('
        print '  "github.com/ninchat/ninchat-go"'
        print ')'

        for action in sorted(ninchat.api.actions.values(), key=lambda a: a.name):
            if action.name != "create_session":
                print_action(action)

    with Output("events"):
        print_header()

        print
        print 'import ('
        print '  "github.com/ninchat/ninchat-go"'
        print ')'

        for event in sorted(ninchat.api.events.values(), key=lambda e: e.name):
            print_event(event)

    with Output("eventfactory"):
        print_header()
        print_eventfactory()

    with Output("attrs"):
        print_header()

        for attrs, objectname, commentname in attrobjects:
            print_attrs(attrs, objectname, commentname)

    with Output("params"):
        print_header()

        for obj in sorted(ninchat.api.objecttypes.values(), key=lambda o: o.name):
            if obj.name != "master_key" and not obj.name.endswith("_2"):
                print_object(obj)


def print_header():
    print 'package ninchatapi'
    print
    print '// THIS FILE IS AUTO-GENERATED BY generate.py - DO NOT EDIT BY HAND!'


def print_action(action):
    action_id = action.params.get("action_id")
    payload = action.name in payloadactions

    print
    print '// {} action.  https://ninchat.com/api/v2#{}'.format(title(action.name), action.name)
    print 'type {} struct {{'.format(title(action.name))

    if action_id and not action_id.required:
        print '  ActionIdDisabled bool `json:"-"`'

    for _, p in sorted(action.params.items()):
        if p.name != "action_id":
            if p.name == "member_attrs":
                if action.name.endswith("_member"):
                    typ = "MemberAttrs"
                else:
                    typ = "*" + title(action.name.split("_")[-1]) + title(p.name)
            elif p.type == "object" and p.name in ninchat.api.objecttypes:
                obj = ninchat.api.objecttypes[p.name]
                typ = "map[string]{}".format(obj.value)
            elif p.type == "object" and p.name.endswith("_attrs"):
                typ = "*" + title(p.name)
            else:
                typ = optionaltypes[p.type]

            if p.required:
                tag = ""
            else:
                tag = ",omitempty"

            print '  {} {} `json:"{}{}"`'.format(title(p.name), typ, p.name, tag)

    if payload:
        print '  Payload []ninchat.Frame'

    print '}'
    print
    print '// String returns "{}".'.format(action.name)
    print 'func (*{}) String() string {{'.format(title(action.name))
    print '  return "{}"'.format(action.name)
    print '}'
    print

    print 'func (action *{}) newClientAction() (clientAction *ninchat.Action, err error) {{'.format(title(action.name))
    print '  clientAction = &ninchat.Action{'
    print '    Params: map[string]interface{}{'
    print '      "action": "{}",'.format(action.name)

    if not action_id:
        print '      "action_id": nil,'

    print '    },'

    if payload:
        print '    Payload: action.Payload,'

    print '  }'

    if action_id and not action_id.required:
        print
        print '  if action.ActionIdDisabled {'
        print '    clientAction.Params["action_id"] = nil'
        print '  }'

    for _, p in sorted(action.params.items()):
        if p.name != "action_id":
            print
            print '  if x := action.{}; {} {{'.format(title(p.name), checks[p.type])
            print '    clientAction.Params["{}"] = {}'.format(p.name, evals[p.type])

            if p.required:
                print '  } else {'
                print '    err = newRequestMalformedError("{} action requires {} parameter")'.format(action.name, p.name)
                print '    return'

            print '  }'

    print
    print '  return'
    print '}'

    reply = unaryreplies[action.name]
    if reply:
        print
        print '// Invoke the action synchronously.'.format(action.name)
        print 'func (action *{}) Invoke(sender Sender) (reply *{}, err error) {{'.format(title(action.name), title(reply))
        print '  var buf {}'.format(title(reply))
        print
        print '  ok, err := unaryCall(sender, action, &buf)'
        print '  if err != nil {'
        print '    return nil, err'
        print '  }'
        print
        print '  if ok {'
        print '    return &buf, nil'
        print '  }'
        print
        print '  return nil, nil'
        print '}'


def print_event(event):
    payload = event.name in payloadevents

    params = event.params.copy()
    params["event_id"] = ninchat.api.Parameter("event_id", False)
    try:
        del params["action_id"]
    except KeyError:
        pass

    print
    print '// {} event.  https://ninchat.com/api/v2#{}'.format(title(event.name), event.name)
    print 'type {} struct {{'.format(title(event.name))

    for _, p in sorted(params.items()):
        if p.name == "member_attrs":
            typ = "*" + title(event.name.split("_", 1)[0]) + title(p.name)
        elif p.type == "object" and p.name.endswith("_attrs"):
            typ = "*" + title(p.name)
        elif p.type == "object" and p.name in ninchat.api.objecttypes:
            obj = ninchat.api.objecttypes[p.name]
            if obj.value:
                val = obj.value
                if val not in requiredtypes:
                    obj2 = ninchat.api.objecttypes.get(val)
                    if obj2 and obj2.value:
                        val2 = obj2.value
                        if val2 == "master_key":
                            val2 = "struct{}"
                        else:
                            val2 = "*" + title(val2)
                        val = "map[string]" + val2
                    else:
                        val = "*" + title(val)
                typ = "map[string]{}".format(val)
            elif obj.item:
                typ = "[]*{}".format(title(obj.item))
            else:
                typ = "*" + title(p.name)
        elif p.type == "object array":
            obj = ninchat.api.objecttypes[p.name]
            typ = "[]*" + title(obj.item)
        elif p.required or p.name == "event_id":
            typ = requiredtypes[p.type]
        else:
            typ = optionaltypes[p.type]

        if p.required:
            tag = ""
        else:
            tag = ",omitempty"

        print '  {} {} `json:"{}{}"`'.format(title(p.name), typ, p.name, tag)

    if payload:
        print
        print '  payload []ninchat.Frame'

    print '}'

    if payload:
        initwhat = "parameters and payload"
    else:
        initwhat = "parameters"

    print
    print '// New{} creates an event object with the {}'.format(title(event.name), initwhat)
    print '// specified by the clientEvent.'
    print '// Its type must be "{}".'.format(event.name)
    print 'func New{0}(clientEvent *ninchat.Event) (event *{0}) {{'.format(title(event.name))
    print '  if clientEvent != nil {'
    print '    e := new({})'.format(title(event.name))
    print '    if err := e.Init(clientEvent); err != nil {'
    print '      panic(err)'
    print '    }'
    print '    event = e'
    print '  }'
    print '  return'
    print '}'
    print
    print '// Init fills in the {}'.format(initwhat)
    print '// specified by the clientEvent (other fields are not touched).'
    print '// An UnexpectedEventError is returned if its type is not'
    print '// "{}".'.format(event.name)
    print 'func (target *{}) Init(clientEvent *ninchat.Event) error {{'.format(title(event.name))
    print '  if clientEvent.String() != "{}" {{'.format(event.name)
    print '    return &UnexpectedEventError{clientEvent}'
    print '  }'

    if params:
        print
        print '  source := clientEvent.Params'

        for _, p in sorted(params.items()):
            print
            print 'if x := source["{}"]; x != nil {{'.format(p.name)

            if p.type == "bool":
                print '  target.{} = true'.format(title(p.name))
            elif p.type == "string array":
                print '  if y, ok := x.([]interface{}); ok {'
                print '    target.{} = AppendStrings(nil, y)'.format(title(p.name))
                print '  }'
            elif p.type == "object array":
                print '  if y, ok := x.([]interface{}); ok {'
                print '    target.{0} = Append{0}(nil, y)'.format(title(p.name))
                print '  }'
            elif p.type != "object":
                if p.type == "int":
                    typ = "float64"
                    if p.required or p.name == "event_id":
                        expr = "int(y)"
                    else:
                        expr = "intPointer(y)"
                else:
                    typ = requiredtypes[p.type]
                    if p.required:
                        expr = "y"
                    else:
                        expr = "&y"

                print '  if y, ok := x.({}); ok {{'.format(typ)
                print '    target.{} = {}'.format(title(p.name), expr)
                print '  }'
            else:
                print '  if y, ok := x.(map[string]interface{}); ok {'

                pobj = ninchat.api.objecttypes.get(p.name)
                if pobj or p.name.endswith("_attrs"):
                    if p.name == "member_attrs":
                        typ = title(event.name.split("_", 1)[0]) + title(p.name)
                    else:
                        typ = title(p.name)

                    if pobj and pobj.value:
                        print '    target.{} = Make{}(y)'.format(title(p.name), typ)
                    elif pobj and pobj.item:
                        print '    target.{} = Append{}(nil, y)'.format(title(p.name), typ)
                    else:
                        print '    target.{} = New{}(y)'.format(title(p.name), typ)
                else:
                    print '    target.{} = y'.format(title(p.name))

                print '  }'

            print '}'

    if payload:
        print
        print '  target.InitPayload(clientEvent.Payload)'

    print
    print '  return nil'
    print '}'

    if payload:
        print
        print '// InitPayload sets the payload (other fields are not touched).'
        print 'func (event *{}) InitPayload(payload []ninchat.Frame) {{'.format(title(event.name))
        print '  event.payload = payload'
        print '}'
        print
        print '// Payload of the event.'
        print 'func (event *{}) Payload() []ninchat.Frame {{'.format(title(event.name))
        print '  return event.payload'
        print '}'

    print
    print '// Id returns the EventId parameter.'
    print 'func (event *{}) Id() int {{'.format(title(event.name))
    print '  return event.EventId'
    print '}'
    print
    print '// String returns "{}".'.format(event.name)
    print 'func (*{}) String() string {{'.format(title(event.name))
    print '  return "{}"'.format(event.name)
    print '}'


def print_attrs(attrs, objectname, commentname):
    if attrs:
        print
        print '// {}Attrs.  https://ninchat.com/api/v2#{}'.format(objectname, commentname.replace(" ", "-"))
    else:
        print
        print '// {}Attrs.'.format(objectname)

    print 'type {}Attrs struct {{'.format(objectname)

    for _, a in sorted(attrs.items()):
        if a.type == "object":
            t = "*{}{}Attr".format(objectname, title(a.name))
        else:
            t = optionaltypes[a.type]

        print '  {} {} `json:"{},omitempty"`'.format(title(a.name), t, a.name)

    print '}'
    print
    print '// New{}Attrs creates an object with the attributes specified by the source.'.format(objectname)
    print 'func New{0}Attrs(source map[string]interface{{}}) (target *{0}Attrs) {{'.format(objectname)
    print '  target = new({}Attrs)'.format(objectname)
    print '  target.Init(source)'
    print '  return'
    print '}'
    print
    print '// Init fills in the attributes specified by the source'
    print '// (other fields are not touched).'
    print 'func (target *{}Attrs) Init(source map[string]interface{{}}) {{'.format(objectname)

    for i, (_, p) in enumerate(sorted(attrs.items())):
        if i > 0:
            print

        print '  if x := source["{}"]; x != nil {{'.format(p.name)

        if p.type == "bool":
            print '    target.{} = true'.format(title(p.name))
        elif p.type == "string array":
            print '    if y, ok := x.([]interface{}); ok {'
            print '      target.{} = AppendStrings(nil, y)'.format(title(p.name))
            print '    }'
        elif p.type != "object":
            print '    if y, ok := x.({}); ok {{'.format(requiredtypes[p.type], p.name)
            print '      target.{} = &y'.format(title(p.name))
            print '    }'
        else:
            print '    if y, ok := x.(map[string]interface{}); ok {'
            print '      target.{} = New{}{}Attr(y)'.format(title(p.name), objectname, title(p.name))
            print '    }'

        print '  }'

    print '}'


def print_object(obj):
    if obj.value:
        obj2 = None
        val2 = None
        val = obj.value
        if val not in requiredtypes:
            obj2 = ninchat.api.objecttypes.get(val)
            if obj2 and obj2.value:
                val2 = obj2.value
                if val2 == "master_key":
                    val2 = "struct{}"
                else:
                    val2 = "*" + title(val2)
                val = "map[string]" + val2
            else:
                val = "*" + title(val)

        print
        print '// Make{} duplicates the map while unwrapping the values.'.format(title(obj.name))
        print 'func Make{}(source map[string]interface{{}}) (target map[string]{}) {{'.format(title(obj.name), val)
        print '  target = make(map[string]{})'.format(val)
        print
        print '  for key, x := range source {'

        if val2:
            print '    if y, ok := x.(map[string]interface{}); ok {'
            print '      t := make({})'.format(val)
            print
            print '      for key2, x2 := range y {'

            if val2 == "struct{}":
                print '        if _, ok2 := x2.(map[string]interface{}); ok2 {'
                print '          t[key2] = struct{}{}'
                print '        }'
            else:
                print '        if y2, ok2 := x2.(map[string]interface{}); ok2 {'
                print '          t[key2] = New{}(y2)'.format(title(obj2.value))
                print '        }'

            print '      }'
            print
            print '      target[key] = t'.format(val)
            print '    }'
        elif obj.value in requiredtypes or obj.name == "member_attrs" or obj.name.endswith("_settings"):
            print '    if y, ok := x.({}); ok {{'.format(requiredtypes[obj.value])
            print '      target[key] = y'
            print '    }'
        else:
            print '    if y, ok := x.(map[string]interface{}); ok {'
            print '      target[key] = New{}(y)'.format(title(obj.value))
            print '    }'

        print '  }'
        print
        print '  return'
        print '}'
    elif obj.item:
        print
        print '// Append{} duplicates the source slice while unwrapping the elements.'.format(title(obj.name))
        print 'func Append{0}(target []*{1}, source []interface{{}}) []*{1} {{'.format(title(obj.name), title(obj.item))
        print '  if source != nil {'
        print '    if target == nil || cap(target) < len(target)+len(source) {'
        print '      t := make([]*{}, len(target), len(target)+len(source))'.format(title(obj.item))
        print '      copy(t, target)'
        print '      target = t'
        print '    }'
        print
        print '    for _, x := range source {'
        print '      var z *{}'.format(title(obj.item))
        print '      if y, ok := x.(map[string]interface{}); ok {'
        print '        z = New{}(y)'.format(title(obj.item))
        print '      }'
        print '      target = append(target, z)'
        print '    }'
        print '  }'
        print
        print '  return target'
        print '}'
    else:
        print
        print '// {} event parameter type.'.format(title(obj.name))
        print 'type {} struct {{'.format(title(obj.name))

        for _, p in sorted(obj.params.items()):
            if p.required:
                types = requiredtypes
                tag = ""
            else:
                types = optionaltypes
                tag = ",omitempty"

            pobj = ninchat.api.objecttypes.get(p.type)
            if pobj and pobj.value:
                pval = pobj.value
                if pval not in requiredtypes:
                    pval = "*" + title(pval)

                t = 'map[string]{}'.format(pval)
            elif pobj and pobj.item:
                t = '[]*{}'.format(title(pobj.item))
            elif pobj or p.name.endswith("_attrs"):
                t = '*{}'.format(title(p.type))
            else:
                t = types[p.type]

            print '{} {} `json:"{}{}"`'.format(title(p.name), t, p.name, tag)

        print '}'
        print
        print '// New{} creates an object with the parameters specified by the source.'.format(title(obj.name))
        print 'func New{0}(source map[string]interface{{}}) (target *{0}) {{'.format(title(obj.name))
        print '  target = new({})'.format(title(obj.name))
        print '  target.Init(source)'
        print '  return'
        print '}'
        print
        print '// Init fills in the parameters specified by the source'
        print '// (other fields are not touched).'
        print 'func (target *{0}) Init(source map[string]interface{{}}) {{'.format(title(obj.name))

        for i, (_, p) in enumerate(sorted(obj.params.items())):
            if i > 0:
                print

            print '  if x := source["{}"]; x != nil {{'.format(p.name)

            if p.type == "bool":
                print '    target.{} = true'.format(title(p.name))
            elif p.type == "string array":
                print '    if y, ok := x.([]interface{}); ok {'
                print '      target.{} = AppendStrings(nil, y)'.format(title(p.name))
                print '    }'
            elif p.type in requiredtypes:
                if p.required or p.name.endswith("_metadata"):
                    expr = "y"
                else:
                    expr = "&y"

                print '    if y, ok := x.({}); ok {{'.format(requiredtypes[p.type])
                print '      target.{} = {}'.format(title(p.name), expr)
                print '    }'
            else:
                pobj = ninchat.api.objecttypes.get(p.type)
                if pobj and pobj.value:
                    print '    if y, ok := x.(map[string]interface{}); ok {'
                    print '      target.{} = Make{}(y)'.format(title(p.name), title(pobj.name))
                    print '    }'
                elif pobj and pobj.item:
                    print '    if y, ok := x.([]interface{}); ok {'
                    print '      target.{} = Append{}(nil, y)'.format(title(p.name), title(pobj.name))
                    print '    }'
                elif pobj or p.name.endswith("_attrs"):
                    typ = title(p.type)

                    print '    if y, ok := x.(map[string]interface{}); ok {'
                    print '      target.{} = New{}(y)'.format(title(p.name), typ)
                    print '    }'

            print '  }'

        print '}'


def print_eventfactory():
    print
    print '// EventFactories contains default constructors for all known event types.'
    print 'var EventFactories = map[string]func() Event{'

    for event in sorted(ninchat.api.events.values(), key=lambda e: e.name):
        print '  "{}": func() Event {{ return new({}) }},'.format(event.name, title(event.name))

    print '}'


def title(s):
    return s.title().replace("_", "")


if __name__ == "__main__":
    main()
