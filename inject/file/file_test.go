package file

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type JsonType struct {
	Test string `json:"test"`
}

type XmlType struct {
	XMLName xml.Name `xml:"xml"`
	Test    string   `xml:"test"`
}

type GobType struct {
	Test string
}

func TestInjector(t *testing.T) {

	for _, testCase := range []struct {
		value    reflect.Value
		tagValue string
		expected reflect.Value
	}{
		{
			value:    reflect.ValueOf(new(string)).Elem(),
			tagValue: "test.txt",
			expected: reflect.ValueOf("test"),
		},
		{
			value:    reflect.ValueOf(&JsonType{}),
			tagValue: "test.json",
			expected: reflect.ValueOf(&JsonType{"test"}),
		},
		{
			value:    reflect.ValueOf(&XmlType{}),
			tagValue: "test.xml",
			expected: reflect.ValueOf(&XmlType{
				XMLName: xml.Name{Local: "xml"},
				Test:    "test",
			}),
		},
		{
			value:    reflect.ValueOf(&GobType{}),
			tagValue: "test.gob",
			expected: reflect.ValueOf(&GobType{"test"}),
		},
		{
			value:    reflect.ValueOf(new(string)).Elem(),
			tagValue: "test_text,txt",
			expected: reflect.ValueOf("test"),
		},
	} {
		if ok, err := Inject(testCase.value, testCase.tagValue); err != nil {
			t.Error(err)
		} else if !ok {
			t.Errorf("%s: expected value to be set", testCase.tagValue)
		} else if !reflect.DeepEqual(testCase.value.Interface(), testCase.expected.Interface()) {
			t.Errorf("%s: expected %q but got %q", testCase.tagValue, testCase.expected, testCase.value)
		}
	}

	var f os.File
	if ok, err := Inject(reflect.ValueOf(&f), "test.txt"); err != nil {
		t.Error(err)
	} else if !ok {
		t.Error("expected file value to be set")
	} else if bytes, err := ioutil.ReadAll(&f); err != nil {
		t.Error(err)
	} else if string(bytes) != "test" {
		t.Errorf(`%s: expected "test" but got %q`, string(bytes))
	}
}
