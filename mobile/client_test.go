package client

import (
	"encoding/base64"
	"testing"

	"github.com/dvsekhvalnov/jose2go"
)

func TestProps_Marshal(t *testing.T) {
	t.Run("ConvertPropsFull", func(t *testing.T) {
		strJson := `{"data":{"base":"test-base","currency":"EU","amount":99.87}}`
		stringArr := NewStrings()
		stringArr.Append("1")
		stringArr.Append("2")
		stringArr.Append("3")
		simplePros := &Props{
			m: map[string]interface{}{
				"sub-i": "ii",
				"sub-j": "jj",
			},
		}
		ninchatProps := Props{
			m: map[string]interface{}{},
		}
		ninchatProps.SetFloat("foo", 3.14159)
		ninchatProps.SetString("bar", "asdf")
		ninchatProps.SetStringArray("baz", stringArr)
		ninchatProps.SetInt("kaz", 1)
		ninchatProps.SetBool("taz", true)
		ninchatProps.SetJSON("uzz", NewJSON(strJson))
		ninchatProps.SetObject("qux", simplePros)

		bt, err := ninchatProps.MarshalJSON()
		if err != nil {
			t.Error("should be able to marshal ninchat json", err)
		}
		if bt == "" {
			t.Error("byte representation of the json will not be empty")
		}
		err = ninchatProps.UnmarshalJSON(bt)
		if err != nil {
			t.Error("should be able to unmarshal ninchat string", err)
		}
		nextBt, _ := ninchatProps.MarshalJSON()
		t.Logf("%+v\n", string(bt))
		t.Logf("%+v\n", string(nextBt))
		t.Logf("%+v\n", ninchatProps)
	})

	t.Run("ConvertPropsMalformed", func(t *testing.T) {
		strJson := `{"data":{"base":"test-base","currency":"EU","amount":}`
		stringArr := NewStrings()
		stringArr.Append("1")
		stringArr.Append("2")
		stringArr.Append("3")
		simplePros := &Props{
			m: map[string]interface{}{
				"sub-i": "ii",
				"sub-j": "jj",
			},
		}
		ninchatProps := Props{
			m: map[string]interface{}{},
		}
		ninchatProps.SetFloat("foo", 3.14159)
		ninchatProps.SetString("bar", "asdf")
		ninchatProps.SetStringArray("baz", stringArr)
		ninchatProps.SetInt("kaz", 1)
		ninchatProps.SetBool("taz", true)
		ninchatProps.SetJSON("uzz", NewJSON(strJson))
		ninchatProps.SetObject("qux", simplePros)

		bt, err := ninchatProps.MarshalJSON()
		if err == nil {
			t.Error("should failed to marshal malform json object", err)
		}
		if bt != "" {
			t.Error("byte representation of the json will empty")
		}
		err = ninchatProps.UnmarshalJSON(bt)
		if err == nil {
			t.Error("should failed to unmarshal with malformed json", err)
		}
	})

}

func TestProps_UnMarshalJSON(t *testing.T) {
	ninchatProps := Props{
		m: map[string]interface{}{},
	}
	propsString := `{"bar":"asdf","baz":["1","2","3"],"foo":3.14159,"kaz":1,"qux":{"sub-i":"ii","sub-j":"jj"},"taz":true,"uzz":{"data":{"amount":99.87,"base":"test-base","currency":"EU"}}}`
	err := ninchatProps.UnmarshalJSON(propsString)
	if err != nil {
		t.Error("should be able to unmarshal ninchat string", err)
	}
	if _, err := ninchatProps.GetObject("uzz"); err != nil {
		t.Error("should be able to get props object", err)
	}
	if _, err := ninchatProps.GetStringArray("baz"); err != nil {
		t.Error("should be able to get string array", err)
	}
	if _, err := ninchatProps.GetObject("qux"); err != nil {
		t.Error("should be able to get json object", err)
	}
}

func TestProps_EncryptJWT_InvalidParams(t *testing.T) {
	props := NewProps()

	token, err := props.EncryptToJWT("", "", -1.0)
	if err == nil {
		t.Error("expected to receive an error in caes of an invalid secret")
	}
	if token != "" {
		t.Error("expected to receive an empty token in caes of an invalid secret")
	}
}

func TestProps_EncryptJWT_InvalidSecret(t *testing.T) {
	props := NewProps()

	token, err := props.EncryptToJWT("", "ABC123", -1.0)
	if err == nil {
		t.Error("expected to receive an error in caes of an invalid secret")
	}
	if token != "" {
		t.Error("expected to receive an empty token in caes of an invalid secret")
	}
}

func TestProps_EncryptJWT(t *testing.T) {
	props := NewProps()
	err := props.UnmarshalJSON("{\"ninchat.com/metadata\":{\"key\":\"value\"}}")
	if err != nil {
		t.Error(err)
	}

	token, err := props.EncryptToJWT("123456789ABC", "VGhpcyBpcyBhIGZha2Uga2V5IGdlbmVyYXRlZCBmb3I=", 1)
	if err != nil {
		t.Error(err)
	}
	if len(token) == 0 {
		t.Error("expected to receive a valid token")
	}

	decodedSecret, err := base64.StdEncoding.DecodeString("VGhpcyBpcyBhIGZha2Uga2V5IGdlbmVyYXRlZCBmb3I=")
	if err != nil {
		t.Error(err)
	}

	decoded, _, err := jose.Decode(token, decodedSecret)
	if err != nil {
		t.Error(err)
	}
	if len(decoded) == 0 {
		t.Error("expected to decode successfully")
	}
	if decoded != "{\"exp\":1,\"ninchat.com/metadata\":{\"key\":\"value\"}}" {
		t.Error("expected to decode to the given props")
	}
}
