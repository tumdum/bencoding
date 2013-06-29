package bencoding

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"net/url"
	"reflect"
	"strconv"
)

// Unmarshal parses the bencoded data and stores the result
// in the value pointed to by v.
//
// Unmarshal uses the inverse of the encodings that Marshal
// uses, allocating maps, slices, pointers and strings as
// necessary.
func Unmarshal(data []byte, v interface{}) error {
	d := NewBytesDecoder(data)
	return d.Decode(v)
}

type Decoder struct {
	torrentDecoder *TorrentDecoder
}

func NewDecoder(r io.Reader) *Decoder {
	td := NewTorrentDecoder(r)
	return &Decoder{td}
}

func NewBytesDecoder(b []byte) *Decoder {
	td := NewBytesTorrentDecoder(b)
	return &Decoder{td}
}

func NewStringDecoder(s string) *Decoder {
	td := NewStringTorrentDecoder(s)
	return &Decoder{td}
}

func (d *Decoder) Decode(v interface{}) error {
	return d.torrentDecoder.unmarshal(v)
}

type TorrentDecoder struct {
	b *hashingRreader
}

func NewTorrentDecoder(r io.Reader) *TorrentDecoder {
	hr := hashingRreader{bufio.NewReader(r), nil, false}
	d := TorrentDecoder{&hr}
	return &d
}

func NewBytesTorrentDecoder(b []byte) *TorrentDecoder {
	return NewTorrentDecoder(bytes.NewBuffer(b))
}

func NewStringTorrentDecoder(s string) *TorrentDecoder {
	return NewBytesTorrentDecoder([]byte(s))
}

type InfoHash []byte

func (i InfoHash) String() string {
	return url.QueryEscape(string(i))
}

func UnmarshalTorrent(data []byte, v interface{}) (InfoHash, error) {
	d := NewBytesTorrentDecoder(data)
	return d.Decode(v)
}

func (d *TorrentDecoder) Decode(v interface{}) (InfoHash, error) {
	if e := d.unmarshal(v); e != nil {
		return nil, e
	} else if d.b.hash != nil {
		return InfoHash(d.b.hash.Sum(nil)), nil
	} else {
		return nil, errors.New("missing info key")
	}
}

func (d *TorrentDecoder) peek() (byte, error) {
	if b, e := d.b.Peek(1); e != nil {
		return ' ', e
	} else {
		return b[0], nil
	}
}

func (d *TorrentDecoder) unmarshal(v interface{}) error {
	val := reflect.ValueOf(v)
	return d.unmarshalToVal(val)
}

func (d *TorrentDecoder) unmarshalToVal(val reflect.Value) error {
	if val.Kind() != reflect.Ptr {
		return errors.New("Can only unmarshal pointers")
	}
	val = val.Elem()
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return d.unmarshalInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return d.unmarshalUint(val)
	case reflect.String:
		return d.unmarshalString(val)
	case reflect.Slice:
		return d.unmarshalSlice(val)
	case reflect.Map:
		return d.unmarshalMap(val)
	case reflect.Struct:
		return d.unmarshalStruct(val)
	}
	return errors.New("Unsupported type encountered")
}

func (d *TorrentDecoder) unmarshalSlice(v reflect.Value) error {
	if v.Type().Elem().Kind() == reflect.Uint8 {
		var s string
		sv := reflect.ValueOf(&s).Elem()
		if e := d.unmarshalString(sv); e != nil {
			return e
		}
		v.SetBytes([]byte(s))
		return nil
	}

	if b, e := d.peek(); e != nil {
		return e
	} else if b != 'l' {
		return errors.New("malformed list beggining (missing 'l')")
	}
	d.b.ReadByte()
	for {
		if b, e := d.peek(); e != nil {
			return e
		} else if b == 'e' {
			break
		}
		if newValue, e := d.unmarshalUnknownItem(); e != nil {
			return nil
		} else {
			v.Set(reflect.Append(v, newValue))
		}
	}

	if b, e := d.peek(); e != nil {
		return e
	} else if b != 'e' {
		return errors.New("malformed list end (missing 'e')")
	}
	d.b.ReadByte()
	return nil
}

func (d *TorrentDecoder) unmarshalMap(v reflect.Value) error {
	if b, e := d.peek(); e != nil {
		return e
	} else if b != 'd' {
		return errors.New("Malformed dict")
	}
	d.b.ReadByte()
	for {
		if b, e := d.peek(); e != nil {
			return e
		} else if b == 'e' {
			break
		}
		var key string
		keyv := reflect.ValueOf(&key).Elem()
		if e := d.unmarshalString(keyv); e != nil {
			return e
		}

		var infoEncountered bool
		if key == "info" {
			infoEncountered = true
			d.b.StartHasing()
		}

		if newValue, e := d.unmarshalUnknownItem(); e != nil {
			return e
		} else {
			v.SetMapIndex(keyv, newValue)
		}

		if infoEncountered {
			d.b.StopHashing()
		}

	}
	if b, e := d.peek(); e != nil {
		return e
	} else if b != 'e' {
		return errors.New("Malformed dict")
	}
	d.b.ReadByte()
	return nil
}

func (d *TorrentDecoder) unmarshalStruct(v reflect.Value) error {
	di := make(map[string]interface{})
	dv := reflect.ValueOf(&di).Elem()
	if e := d.unmarshalMap(dv); e != nil {
		return e
	}
	return bind(di, v)
}

func (d *TorrentDecoder) unmarshalUnknownItem() (reflect.Value, error) {
	b, e := d.peek()
	if e != nil {
		return reflect.Value{}, e
	}
	switch b {
	case 'i':
		var i int64
		iv := reflect.ValueOf(&i).Elem()
		if e := d.unmarshalInt(iv); e != nil {
			return reflect.Value{}, e
		} else {
			return iv, nil
		}
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		var s string
		sv := reflect.ValueOf(&s).Elem()
		if e := d.unmarshalString(sv); e != nil {
			return reflect.Value{}, e
		} else {
			return sv, nil
		}
	case 'l':
		var l []interface{}
		lv := reflect.ValueOf(&l).Elem()
		if e := d.unmarshalSlice(lv); e != nil {
			return reflect.Value{}, e
		} else {
			return lv, nil
		}
	case 'd':
		di := make(map[string]interface{})
		dv := reflect.ValueOf(&di).Elem()
		if e := d.unmarshalMap(dv); e != nil {
			return reflect.Value{}, e
		} else {
			return dv, nil
		}
	default:
		return reflect.Value{}, errors.New("Unknonw item")
	}
}

func (d *TorrentDecoder) unmarshalInt(v reflect.Value) error {
	if b, e := d.peek(); e != nil {
		return e
	} else if b != 'i' {
		return errors.New("Malformed integer input")
	}
	d.b.ReadByte()
	if data, e := d.b.ReadBytes('e'); e != nil {
		return e
	} else {
		if val, e := strconv.ParseInt(string(data[:len(data)-1]), 10, 64); e != nil {
			return e
		} else {
			v.SetInt(val)
		}
	}
	return nil
}

func (d *TorrentDecoder) unmarshalUint(v reflect.Value) error {
	if b, e := d.peek(); e != nil {
		return e
	} else if b != 'i' {
		return errors.New("Malformed integer input")
	}
	d.b.ReadByte()
	if data, e := d.b.ReadBytes('e'); e != nil {
		return e
	} else {
		if val, e := strconv.ParseUint(string(data[:len(data)-1]), 10, 64); e != nil {
			return e
		} else {
			v.SetUint(val)
		}
	}
	return nil
}

// TODO: maybe this should be implemented in a more
// efficent manner...
func readExactly(b *hashingRreader, out []byte) error {
	pos := 0
	for pos < len(out) {
		if b, e := b.ReadByte(); e != nil {
			return e
		} else {
			out[pos] = b
			pos++
		}
	}
	return nil
}

func (d *TorrentDecoder) unmarshalString(v reflect.Value) error {
	lStr, e := d.b.ReadBytes(':')
	if e != nil {
		return e
	}
	length, e := strconv.ParseInt(string(lStr[:len(lStr)-1]), 10, 64)
	if e != nil {
		return e
	}
	content := make([]byte, length)
	e = readExactly(d.b, content)
	if e != nil {
		return e
	}
	v.SetString(string(content))
	return nil
}
