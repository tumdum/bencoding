package bencoding

import (
	"bytes"
	"reflect"
	"testing"
)

func TestIntegerMarshaling(t *testing.T) {
	data := []struct {
		in  interface{}
		out string
	}{
		{1, "i1e"},
		{-1, "i-1e"},
		{0, "i0e"},
		{43210, "i43210e"},
		{int(1), "i1e"},
		{int8(2), "i2e"},
		{int16(3), "i3e"},
		{int32(4), "i4e"},
		{int64(5), "i5e"},
		{uint(6), "i6e"},
		{uint8(7), "i7e"},
		{uint16(8), "i8e"},
		{uint32(9), "i9e"},
		{uint64(10), "i10e"},
	}

	for _, test := range data {
		output, err := Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if test.out != string(output) {
			t.Fatalf("got different output then expected: %v != %v", test.out, string(output))
		}
	}
}

func TestStringMarshaling(t *testing.T) {
	data := []struct {
		in  string
		out string
	}{
		{"", "0:"},
		{"x", "1:x"},
		{"foobarbaz", "9:foobarbaz"},
	}
	for _, test := range data {
		output, err := Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if test.out != string(output) {
			t.Fatalf("got different output then expected: %v != %v", test.out, string(output))
		}
	}
}

func TestSliceAsStringMarshaling(t *testing.T) {
	data := []struct {
		in  []byte
		out []byte
	}{
		{nil, []byte{'0', ':'}},
		{[]byte{}, []byte{'0', ':'}},
		{[]byte("test"), []byte("4:test")},
		{[]byte{0, 1, 2, 3}, []byte{'4', ':', 0, 1, 2, 3}},
	}
	for _, test := range data {
		output, err := Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(test.out, output) != 0 {
			t.Fatalf("got different output then expected: %v != %v", test.out, output)
		}
	}
}

func TestArrayAsStringMarshaling(t *testing.T) {
	data := []struct {
		in  interface{}
		out []byte
	}{
		{[...]byte{}, []byte("0:")},
		{[...]byte{'x', 'y'}, []byte("2:xy")},
		{[...]byte{0, 1, 2}, []byte{'3', ':', 0, 1, 2}},
	}
	for _, test := range data {
		output, err := Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(test.out, output) != 0 {
			t.Fatalf("got different output then expected: %v != %v", string(test.out), string(output))
		}
	}
}

func TestSliceAsListMarshaling(t *testing.T) {
	data := []struct {
		in  []interface{}
		out []byte
	}{
		{[]interface{}{}, []byte("le")},
		{[]interface{}{1, 2, 3}, []byte("li1ei2ei3ee")},
		{[]interface{}{"foo", "bar", "baz"}, []byte("l3:foo3:bar3:baze")},
		{[]interface{}{1, "foo", 2, "bar"}, []byte("li1e3:fooi2e3:bare")},
		{[]interface{}{1, []interface{}{"bar"}}, []byte("li1el3:baree")},
	}
	for _, test := range data {
		output, err := Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(test.out, output) != 0 {
			t.Fatalf("got different output then expected: %v != %v", string(test.out), string(output))
		}
	}
}

func TestArrayAsListMarshaling(t *testing.T) {
	data := []struct {
		in  interface{}
		out []byte
	}{
		{[...]interface{}{}, []byte("le")},
		{[...]interface{}{1, 2, 3}, []byte("li1ei2ei3ee")},
		{[...]interface{}{"foo", "bar", "baz"}, []byte("l3:foo3:bar3:baze")},
		{[...]interface{}{1, "foo", 2, "bar"}, []byte("li1e3:fooi2e3:bare")},
		{[...]interface{}{1, [...]interface{}{"bar"}}, []byte("li1el3:baree")},
	}
	for _, test := range data {
		output, err := Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(test.out, output) != 0 {
			t.Fatalf("got different output then expected: %v != %v", string(test.out), string(output))
		}
	}
}

func TestMapAsDictMarshaling(t *testing.T) {
	data := []struct {
		in  interface{}
		out []byte
	}{
		{map[string]string{}, []byte{'d', 'e'}},
		{map[string]int{"1": 1, "3": 3, "123": 123}, []byte("d1:1i1e3:123i123e1:3i3ee")},
		{map[string]string{"publisher": "bob", "publisher-webpage": "www.example.com", "publisher.location": "home"}, []byte("d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee")},
		{map[string]interface{}{"1": "one"}, []byte("d1:13:onee")},
	}
	for _, test := range data {
		output, err := Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(test.out, output) != 0 {
			t.Fatalf("got different output then expected: %v != %v", string(test.out), string(output))
		}
	}
}

func TestStructMarshaling(t *testing.T) {
	type Test struct {
		Foo string
		Bar int
		Baz []string
	}
	test := Test{"FooValue", 42, []string{"A", "BB", "CCC"}}
	expected := "d3:Bari42e3:Bazl1:A2:BB3:CCCe3:Foo8:FooValuee"
	output, err := Marshal(test)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare([]byte(expected), output) != 0 {
		t.Fatalf("got different output then expected: %v != %v", string(expected), string(output))
	}
}

func TestPointerMarshaling(t *testing.T) {
	s := ""
	i := 1
	data := []struct {
		in  interface{}
		out []byte
	}{
		{&map[string]string{}, []byte{'d', 'e'}},
		{&[]string{}, []byte("le")},
		{&s, []byte("0:")},
		{&i, []byte("i1e")},
		{&[]*[]string{&[]string{}, &[]string{}}, []byte("llelee")},
	}
	for _, test := range data {
		output, err := Marshal(test.in)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Compare(test.out, output) != 0 {
			t.Fatalf("got different output then expected: %v != %v", string(test.out), string(output))
		}
	}
}

func TestEncoderCorrectness(t *testing.T) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	e.Encode(map[string]string{"publisher": "bob", "publisher-webpage": "www.example.com", "publisher.location": "home"})
	output := string(b.Bytes())
	if output != "d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee" {
		t.Fatalf("Encoder output incorrect: %v", output)
	}
}

func TestTagParsing(t *testing.T) {
	type T struct {
		Field1 int    `bencoding:"port"`
		Field2 int    `json:"foo"`
		Field3 string `bencoding:"-"`
	}
	testVar := T{42, 52, "test"}
	test := reflect.ValueOf(testVar)
	if tag := tagForFieldNamed(test, "Field1"); tag != `bencoding:"port"` {
		t.Fatalf("Expected 'bencoding:\"port\"' got '%v'", tag)
	}
	if tag := tagForFieldNamed(test, "Field2"); tag != "json:\"foo\"" {
		t.Fatalf("Expected empty tag got '%v'", tag)
	}
	if tag := tagForFieldNamed(test, "Field3"); tag != `bencoding:"-"` {
		t.Fatalf("Expected 'bencoding:\"-\"' got '%v'", tag)
	}
}

func TestParsingTags(t *testing.T) {
	data := []struct {
		in  string
		out string
	}{
		{`bencoding:"foo"`, "foo"},
		{`bencoding:""`, ""},
		{`bencoding:"-"`, "-"},
	}

	for _, test := range data {
		out := parseTag(test.in)
		if out == nil {
			t.Fatalf("Expected '%v' got nil", test.out)
		} else if *out != test.out {
			t.Fatalf("Expected '%v' got '%v'", test.out, *out)
		}
	}
	fails := []string{
		`bencoding:`,
		`bencoding:"foo`,
		`bencoding"foo"`,
		`foo:"bar"`,
	}
	for _, test := range fails {
		out := parseTag(test)
		if out != nil {
			t.Fatalf("Expected to get nil but got '%v'", out)
		}
	}
}

func TestParsingOfFieldOptions(t *testing.T) {
	type T struct {
		Field1 int `bencoding:"foo"`
		Field2 int `other:"bar"`
		Field3 int `bencoding:"-"`
		Field4 int
	}
	var testVar T
	test := reflect.ValueOf(testVar)
	if opt := extractFieldOptions(test, "Field1"); opt != "foo" {
		t.Fatalf("Expected to have Field1 renamed to 'foo' instead it is renamed to '%v'", opt)
	}
	if opt := extractFieldOptions(test, "Field2"); opt != "Field2" {
		t.Fatalf("Field2 should not be renamed, yet it is to '%v'", opt)
	}
	if opt := extractFieldOptions(test, "Field3"); opt != "" {
		t.Fatalf("Field3 should be ignored, yet it is renamed to '%v'", opt)
	}
	if opt := extractFieldOptions(test, "Field4"); opt != "Field4" {
		t.Fatalf("Field4 should not be renamed, yet it is to '%v'", opt)
	}
}

func TestMarshalingOfTaggetStruct(t *testing.T) {
	type T struct {
		VeryLongName  int `bencoding:"s"`
		NotNeeded     int `bencoding:""`
		AlsoNotNeeded int `bencoding:"-"`
		i             int
	}
	test := T{42, 0, 1, 1337}
	expected := "d1:ii1337e1:si42ee"
	if s, e := Marshal(test); e != nil {
		t.Fatal(e)
	} else if string(s) != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, string(s))
	}
}

func TestMarshalingTaggedStructShouldOrderFieldsBasedOnTaggedNames(t *testing.T) {
	type T struct {
		A int `bencoding:"Z"`
		B int `bencoding:"Y"`
		C int `bencoding:"X"`
	}
	test := T{3, 2, 1}
	expected := "d1:Xi1e1:Yi2e1:Zi3ee"
	if s, e := Marshal(test); e != nil {
		t.Fatal(e)
	} else if string(s) != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, string(s))
	}
}
