package client

import (
	"testing"
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
