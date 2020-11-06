package client

import (
	"testing"
)

func TestStrings_MarshalJSON(t *testing.T) {
	t.Run("Test with non empty string array", func(t *testing.T) {
		stringArr := NewStrings()
		stringArr.Append("1")
		stringArr.Append("2")
		stringArr.Append("3")
		bt, err := stringArr.MarshalJSON()
		if err != nil {
			t.Error("should be able to marshal ninchat string", err)
		}
		if bt == nil {
			t.Error("byte representation of the ninchat string will not be empty")
		}

		err = stringArr.UnMarshalJSON(bt)
		if err != nil {
			t.Error("should be able to unmarshal ninchat string", err)
		}
		if stringArr.a[0] != "1" && stringArr.a[1] != "2" && stringArr.a[2] != "3" {
			t.Error("should match the original string")
		}
		t.Logf("%+v\n", string(bt))
		t.Logf("%+v\n", stringArr)
	})
	t.Run("Test with empty string array", func(t *testing.T) {
		emptyArr := NewStrings()
		bt, err := emptyArr.MarshalJSON()
		if err != nil {
			t.Error("should be able to marshal ninchat string", err)
		}
		if bt == nil {
			t.Error("byte representation of the ninchat string will not be empty")
		}

		err = emptyArr.UnMarshalJSON(bt)
		if err != nil {
			t.Error("should be able to unmarshal ninchat string", err)
		}
		t.Logf("%+v\n", string(bt))
		t.Logf("%+v\n", emptyArr)
	})
}

func TestJSON_MarshalJSON(t *testing.T) {
	t.Run("Test with non empty json", func(t *testing.T) {
		jsonAsString := `{"data":{"base":"test-base","currency":"EU","amount":99.87}}`
		ninchatJson := NewJSON(jsonAsString)
		bt, err := ninchatJson.MarshalJSON()
		if err != nil {
			t.Error("should be able to marshal ninchat json")
		}
		if bt == nil {
			t.Error("byte representation of the json will not be empty")
		}
		err = ninchatJson.UnMarshalJSON(bt)
		if err != nil {
			t.Error("should be able to unmarshal ninchat string", err)
		}

		t.Logf("%+v\n", string(bt))
		t.Logf("%+v\n", ninchatJson)
	})

	t.Run("Test with empty json", func(t *testing.T) {
		jsonAsString := `{}`
		ninchatJson := NewJSON(jsonAsString)
		bt, err := ninchatJson.MarshalJSON()
		if err != nil {
			t.Error("should be able to marshal ninchat json", err)
		}
		if bt == nil {
			t.Error("byte representation of the json will not be empty")
		}
		err = ninchatJson.UnMarshalJSON(bt)
		if err != nil {
			t.Error("should be able to unmarshal ninchat string", err)
		}
		t.Logf("%+v\n", string(bt))
		t.Logf("%+v\n", ninchatJson)
	})
	t.Run("Test with invalid json", func(t *testing.T) {
		jsonAsString := `{`
		ninchatJson := NewJSON(jsonAsString)
		bt, err := ninchatJson.MarshalJSON()
		if err == nil {
			t.Error("should failed with invalid json", err)
		}
		if bt != nil {
			t.Error("byte representation of the invalid json will be empty")
		}
		err = ninchatJson.UnMarshalJSON(bt)
		if err == nil {
			t.Error("should be able to unmarshal ninchat string", err)
		}
		t.Logf("%+v\n", string(bt))
		t.Logf("%+v\n", ninchatJson)
	})

}

func TestProps_Marshal(t *testing.T) {
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
	val := map[string]interface{}{
		"foo":  3.14159,
		"bar":  "asdf",
		"baz":  stringArr,
		"kaz":  1,
		"taz":  true,
		"uzz":  NewJSON(strJson),
		"quux": simplePros,
	}

	ninchatProps := Props{
		m: val,
	}

	bt, err := ninchatProps.MarshalJSON()
	if err != nil {
		t.Error("should be able to marshal ninchat json", err)
	}
	if bt == nil {
		t.Error("byte representation of the json will not be empty")
	}
	err = ninchatProps.UnMarshalJSON(bt)
	if err != nil {
		t.Error("should be able to unmarshal ninchat string", err)
	}

	nextBt, _ := ninchatProps.MarshalJSON()
	t.Logf("%+v\n", string(bt))
	t.Logf("%+v\n", string(nextBt))
	t.Logf("%+v\n", ninchatProps)

}

func TestPropsUtil_JSONString(t *testing.T) {
	t.Run("Convert a props to string", func(t *testing.T) {
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
		val := map[string]interface{}{
			"foo":  3.14159,
			"bar":  "asdf",
			"baz":  stringArr,
			"kaz":  1,
			"taz":  true,
			"uzz":  NewJSON(strJson),
			"quux": simplePros,
		}

		ninchatProps := Props{
			m: val,
		}
		value := PropsUtil{}.JSONString(ninchatProps)
		if value == "" {
			t.Error("should be able to marshal props")
		}
		t.Logf("%+v\n", value)
	})
}

func TestPropsUtil_FromJsonString(t *testing.T) {
	t.Run("Convert a props string to props object", func(t *testing.T) {
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
		val := map[string]interface{}{
			"foo":  3.14159,
			"bar":  "asdf",
			"baz":  stringArr,
			"kaz":  1,
			"taz":  true,
			"uzz":  NewJSON(strJson),
			"quux": simplePros,
		}

		ninchatProps := Props{
			m: val,
		}
		value := PropsUtil{}.JSONString(ninchatProps)
		if value == "" {
			t.Error("should be able to marshal props")
		}
		propsObject := PropsUtil{}.FromJsonString(value)
		t.Logf("%+v\n", propsObject)
	})
}
