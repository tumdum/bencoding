package bencoding

import (
	"testing"
)

func TestUnmarshalAnInt(t *testing.T) {
	var i int
	b := []byte("i5e")
	e := Unmarshal(b, &i)
	if e != nil {
		t.Fatal(e)
	}
	if i != 5 {
		t.Fatalf("Expected 5 got %v", i)
	}
	var u uint
	b = []byte("i5e")
	e = Unmarshal(b, &u)
	if e != nil {
		t.Fatal(e)
	}
	if u != 5 {
		t.Fatalf("Expected 5 got %v", u)
	}
}

func TestUnmarshalTwoInts(t *testing.T) {
	var i int
	var j uint32
	var k int64
	var l uint
	b := []byte("i5ei4ei-12345ei77e")
	d := NewBytesDecoder(b)
	if e := d.Decode(&i); e != nil {
		t.Fatal(e)
	} else if i != 5 {
		t.Fatalf("Expected 5 got '%v'", i)
	}
	if e := d.Decode(&j); e != nil {
		t.Fatal(e)
	} else if j != 4 {
		t.Fatalf("Expected 4 got '%v'", j)
	}
	if e := d.Decode(&k); e != nil {
		t.Fatal(e)
	} else if k != -12345 {
		t.Fatalf("Expected -12345 got '%v'", k)
	}
	if e := d.Decode(&l); e != nil {
		t.Fatal(e)
	} else if l != 77 {
		t.Fatalf("Expected 77 got '%v'", l)
	}
}

func TestUnmarshalString(t *testing.T) {
	var s1, s2, s3 string
	var i int
	b := []byte("3:foo4:baari42e5:baaaz")
	d := NewBytesDecoder(b)
	if e := d.Decode(&s1); e != nil {
		t.Fatal(e)
	} else if s1 != "foo" {
		t.Fatalf("Expected 'foo' got '%v'", s1)
	}
	if e := d.Decode(&s2); e != nil {
		t.Fatal(e)
	} else if s2 != "baar" {
		t.Fatalf("Expected 'baar' got '%v'", s2)
	}
	if e := d.Decode(&i); e != nil {
		t.Fatal(e)
	} else if i != 42 {
		t.Fatalf("Expected 42 got '%v'", i)
	}
	if e := d.Decode(&s3); e != nil {
		t.Fatal(e)
	} else if s3 != "baaaz" {
		t.Fatalf("Expected 'baaaz' got '%v'", s3)
	}
}

func TestUnmarshalEmptyList(t *testing.T) {
	var l []interface{}
	b := []byte("le")
	d := NewBytesDecoder(b)
	if e := d.Decode(&l); e != nil {
		t.Fatal(e)
	} else if len(l) != 0 {
		t.Fatalf("Expected empty slice got '%v'", l)
	}
}

func TestUnmarshalListWithOneElement(t *testing.T) {
	var l []interface{}
	b := []byte("li3ee")
	d := NewBytesDecoder(b)
	if e := d.Decode(&l); e != nil {
		t.Fatal(e)
	} else if len(l) != 1 || l[0].(int64) != 3 {
		t.Fatalf("Expected [3] got '%v'", l)
	}
}

func TestUnmarshalListWithOneString(t *testing.T) {
	var l []interface{}
	b := []byte("l3:fooe")
	d := NewBytesDecoder(b)
	if e := d.Decode(&l); e != nil {
		t.Fatal(e)
	} else if len(l) != 1 || l[0].(string) != "foo" {
		t.Fatalf("Expected ['foo'] got '%v'", l)
	}
}

func TestUnmarshalListWithOneList(t *testing.T) {
	var l []interface{}
	b := []byte("llei3eex")
	d := NewBytesDecoder(b)
	if e := d.Decode(&l); e != nil {
		t.Fatal(e)
	} else if len(l) != 2 || len(l[0].([]interface{})) != 0 || l[1].(int64) != 3 {
		t.Fatalf("Expected [[],3] got '%v'", l)
	}
}

func TestUnmarshalListWithMultipleElements(t *testing.T) {
	var l []interface{}
	b := []byte("li3e6:foobarli42eee")
	d := NewBytesDecoder(b)
	if e := d.Decode(&l); e != nil {
		t.Fatal(e)
	} else if len(l) != 3 || l[0].(int64) != 3 || l[1].(string) != "foobar" || len(l[2].([]interface{})) != 1 || l[2].([]interface{})[0].(int64) != 42 {
		t.Fatalf("Expected [3,'foobar',[42]] got '%v'", l)
	}
}

func TestUnmarshalEmptyDict(t *testing.T) {
	dict := make(map[string]interface{})
	b := []byte("de")
	d := NewBytesDecoder(b)
	if e := d.Decode(&dict); e != nil {
		t.Fatal(e)
	} else if len(dict) != 0 {
		t.Fatalf("Expected {} got '%v'", d)
	}
}

func TestUnmarshalListWithDict(t *testing.T) {
	var l []interface{}
	b := []byte("ldei3ee")
	d := NewBytesDecoder(b)
	if e := d.Decode(&l); e != nil {
		t.Fatal(e)
	} else if len(l) != 2 || len(l[0].(map[string]interface{})) != 0 {
		t.Fatalf("Expected [{}] got '%v'", l)
	}
}

func TestUnmarshalDictWithOneEntry(t *testing.T) {
	d := make(map[string]interface{})
	b := []byte("d3:fooi4ee")
	dec := NewBytesDecoder(b)
	if e := dec.Decode(&d); e != nil {
		t.Fatal(e)
	} else if d["foo"].(int64) != 4 {
		t.Fatalf("Expected {foo:4} got '%v'", d)
	}
}

func TestUnmarshalTwoEmptyDicts(t *testing.T) {
	d1 := make(map[string]interface{})
	d2 := make(map[string]interface{})
	b := []byte("dede")
	dec := NewBytesDecoder(b)
	if e := dec.Decode(&d1); e != nil {
		t.Fatal(e)
	} else if len(d1) != 0 {
		t.Fatalf("Expected {} got '%v'", d1)
	}
	if e := dec.Decode(&d2); e != nil {
		t.Fatal(e)
	} else if len(d2) != 0 {
		t.Fatalf("Expected {} got '%v'", d2)
	}
}

func TestUnmarshalDictWithMultipleKeys(t *testing.T) {
	d := make(map[string]interface{})
	b := []byte("d9:publisher3:bob17:publisher-webpage15:www.example.com18:publisher.location4:homee")
	dec := NewBytesDecoder(b)
	if e := dec.Decode(&d); e != nil {
		t.Fatal(e)
	} else if d["publisher"] != "bob" || d["publisher-webpage"] != "www.example.com" || d["publisher.location"] != "home" || len(d) != 3 {
		t.Fatalf("Expected something different then '%v'", d)
	}
}

func TestUnmarshalingHandlesMissingInput(t *testing.T) {
	var i int
	d := NewStringDecoder("i3")
	if e := d.Decode(&i); e == nil {
		t.Fatalf("Expected to fail while parsing 'i3'")
	}
	var s string
	d = NewStringDecoder("4:foo")
	if e := d.Decode(&s); e == nil {
		t.Fatalf("Expected to fail while parsing '4:foo'")
	}
	var l []interface{}
	d = NewStringDecoder("l")
	if e := d.Decode(&l); e == nil {
		t.Fatalf("Expected to fail while parsing 'l'")
	}
	di := make(map[string]interface{})
	d = NewStringDecoder("d")
	if e := d.Decode(&di); e == nil {
		t.Fatalf("Expected to fail while parsing 'd'")
	}
}

func TestUnmarshalingSliceOfBytesShouldProduceString(t *testing.T) {
	var s []byte
	d := NewStringDecoder("3:foo")
	if e := d.Decode(&s); e != nil {
		t.Fatal(e)
	} else if string(s) != "foo" {
		t.Fatalf("Expected 'foo' got '%v'", string(s))
	}
}

func TestUnmarshalingOfStruct(t *testing.T) {
	type T struct {
		Number int64
		Text   string
		List   []interface{}
		Dict   map[string]interface{}
	}

	v := T{}
	di := make(map[string]interface{})
	s := "d6:Numberi42e4:Text3:foo4:Listli1ei2ei3ee4:Dictd3:fooi100e3:bari200e3:bazi300eee"
	d := NewStringDecoder(s)
	if e := d.Decode(&di); e != nil {
		t.Fatal(e)
	} else if len(di) != 4 {
		t.Fatalf("Not expected: '%v'", di)
	}

	d = NewStringDecoder(s)
	if e := d.Decode(&v); e != nil {
		t.Fatal(e)
	} else if v.Number != 42 || v.Text != "foo" ||
		v.List[0].(int64) != 1 || v.List[1].(int64) != 2 || v.List[2].(int64) != 3 ||
		v.Dict["foo"].(int64) != 100 || v.Dict["bar"].(int64) != 200 || v.Dict["baz"].(int64) != 300 {
		t.Fatalf("Expected T{Number:42 Text:'foo' List:[1 2 3] Dict:{foo:100 bar:200 baz:300}} got '%v'", v)
	}
}

func TestUnmarshalingStructWithTags(t *testing.T) {
	type T struct {
		LongName   int64 `bencoding:"n"`
		Ignore     int   `bencoding:""`
		AlsoIgnore int   `bencoding:"-"`
		F          int64
		G          int64 `bencoding:"a b"`
	}
	v := T{}
	s := "d1:Fi1337e3:foo3:bar1:ni42e3:a bi55ee"
	d := NewStringDecoder(s)
	if e := d.Decode(&v); e != nil {
		t.Fatal(e)
	} else if v.LongName != 42 || v.F != 1337 || v.G != 55 {
		t.Fatalf("Expected T{42,0,0,1337,55} got '%v'", v)
	}
}

func TestUnmarshalStructInAStruct(t *testing.T) {
	type T1 struct {
		Foo int64
		Bar string
	}
	type T2 struct {
		Baz  int64
		Test T1
		XYZ  string
	}
	v := T2{}
	s := "d3:Bazi10e4:Testd3:Fooi20e3:Bar3:aaae3:XYZ4:teste"
	d := NewStringDecoder(s)
	if e := d.Decode(&v); e != nil {
		t.Fatal(e)
	} else if v.Baz != 10 || v.Test.Foo != 20 || v.Test.Bar != "aaa" || v.XYZ != "test" {
		t.Fatalf("Did not expect to get '%v'", v)
	}
}

func TestUnmarshalingStructWithAPointer(t *testing.T) {
	type T2 struct {
		K int64
	}
	type T struct {
		I *int64
		S *string
		K *T2
	}
	var ts T
	s := "d1:Ii42e1:S3:foo1:Kd1:Ki10eee"
	d := NewStringDecoder(s)
	if e := d.Decode(&ts); e != nil {
		t.Fatal(e)
	}
	if *ts.I != 42 || *ts.S != "foo" || ts.K.K != 10 {
		t.Fatalf("Expected that {42 \"foo\"} '%v'", ts)
	}
}

func TestUnmarshalTorrentFromStringWitoutTorrentData(t *testing.T) {
	var i int
	b := []byte("i5e")
	h, e := UnmarshalTorrent(b, &i)
	if e == nil || h != nil {
		t.Fatalf("Expected to receive nil hash and error, got '%v' and '%v'", h, e)
	}
}

func TestUnmarshalStructWithStringSlice(t *testing.T) {
	type T struct {
		X []string
	}
	var v T
	s := "d1:Xl3:foo3:bar3:bazee"
	d := NewStringDecoder(s)
	if e := d.Decode(&v); e != nil {
		t.Fatal(e)
	}
}
