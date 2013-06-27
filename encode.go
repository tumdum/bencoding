package bencoding

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"strconv"
)

// Marshal returns bencode encoding of v.
//
// Marshal traverses the value v recursively.
// Each traversed value is encoded based on its type, using following rules:
//
// String values are encoded as bencode strings.
//
// Arrays of bytes are encoded as bencode strings.
//
// Arrays of elements which type is not byte are encoded as bencode lists.
//
// Slices of bytes are encoded as bencode strings.
//
// Slices of elements which type is not byte are encoded as bencode lists.
//
// Integer types are encoded as bencode integer.
//
// Maps from string to interface{} are encoded as bencode dictionaries.
//
// Structs are encoded as dictionaries. Each exported field becomes
// a member of dictionary unless
//  - the field's bencoding tag is "" or "-"
//
// Pointers are encoded as values to which they point.
//
// Any other type is not supported and trying to encode it will result in an error.
func Marshal(v interface{}) ([]byte, error) {
	var e encodeState
	if err := e.Marshal(v); err != nil {
		return nil, err
	}
	return e.Bytes(), nil
}

type Encoder struct {
	w io.Writer
	e encodeState
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (enc *Encoder) Encode(v interface{}) error {
	enc.e.Reset()
	if err := enc.e.Marshal(v); err != nil {
		return err
	}
	_, err := enc.w.Write(enc.e.Bytes())
	return err
}

type encodeState struct {
	bytes.Buffer
}

func (e *encodeState) Marshal(v interface{}) error {
	val := reflect.ValueOf(v)
	return e.marshal(val)
}

func (e *encodeState) marshal(val reflect.Value) error {
	switch val.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return e.marshalInt(val)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return e.marshalUnsignedInt(val)
	case reflect.String:
		return e.marshalString(val)
	case reflect.Slice:
		return e.marshalSlice(val)
	case reflect.Array:
		return e.marshalArray(val)
	case reflect.Map:
		return e.marshalMap(val)
	case reflect.Struct:
		return e.marshalStruct(val)
	case reflect.Ptr:
		return e.marshalPtr(val)
	case reflect.Interface:
		return e.marshalInterface(val)
	default:
		return errors.New("Unknown kind: " + val.Kind().String())
	}
}

func (e *encodeState) marshalInt(val reflect.Value) error {
	if err := e.WriteByte('i'); err != nil {
		return nil
	}
	if _, err := e.Write(strconv.AppendInt([]byte{}, val.Int(), 10)); err != nil {
		return err
	}
	if err := e.WriteByte('e'); err != nil {
		return err
	}
	return nil
}

func (e *encodeState) marshalUnsignedInt(val reflect.Value) error {
	if err := e.WriteByte('i'); err != nil {
		return err
	}
	if _, err := e.Write(strconv.AppendUint([]byte{}, val.Uint(), 10)); err != nil {
		return err
	}
	if err := e.WriteByte('e'); err != nil {
		return err
	}
	return nil
}

func (e *encodeState) marshalString(val reflect.Value) error {
	if _, err := e.Write(strconv.AppendInt([]byte{}, int64(len(val.String())), 10)); err != nil {
		return err
	}
	if err := e.WriteByte(':'); err != nil {
		return err
	}
	if _, err := e.Write([]byte(val.String())); err != nil {
		return err
	}
	return nil
}

func (e *encodeState) marshalSlice(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		return e.marshalCollectionAsList(val)
	}
	// treat slice like string
	valBytes := val.Bytes()
	if _, err := e.Write(strconv.AppendInt([]byte{}, int64(len(valBytes)), 10)); err != nil {
		return err
	}
	if err := e.WriteByte(':'); err != nil {
		return err
	}
	if _, err := e.Write(valBytes); err != nil {
		return err
	}
	return nil
}

func (e *encodeState) marshalCollectionAsList(val reflect.Value) error {
	if err := e.WriteByte('l'); err != nil {
		return err
	}
	for i := 0; i < val.Len(); i++ {
		// array of interface{} values, need to extract unterling type of element
		element := reflect.ValueOf(val.Index(i).Interface())
		if err := e.marshal(element); err != nil {
			return err
		}
	}
	return e.WriteByte('e')
}

func (e *encodeState) marshalArray(val reflect.Value) error {
	elemKind := val.Type().Elem().Kind()
	if elemKind != reflect.Uint8 {
		return e.marshalCollectionAsList(val)
	}
	if _, err := e.Write(strconv.AppendInt([]byte{}, int64(val.Len()), 10)); err != nil {
		return err
	}
	if err := e.WriteByte(':'); err != nil {
		return err
	}
	for i := 0; i < val.Len(); i++ {
		if err := e.WriteByte(byte(val.Index(i).Uint())); err != nil {
			return err
		}
	}
	return nil
}

func (e *encodeState) marshalMap(val reflect.Value) error {
	keys := val.MapKeys()
	if err := e.WriteByte('d'); err != nil {
		return err
	}
	for _, key := range keys {
		if key.Kind() != reflect.String {
			return errors.New("Map can be marshaled only if keys are of type 'string'")
		}
		if err := e.marshal(key); err != nil {
			return err
		}
		value := val.MapIndex(key)
		if err := e.marshal(value); err != nil {
			return err
		}
	}
	return e.WriteByte('e')
}

func (e *encodeState) marshalStruct(val reflect.Value) error {
	if err := e.WriteByte('d'); err != nil {
		return err
	}
	valType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldValue := val.Field(i)
		fieldOpt := extractFieldOptions(val, valType.Field(i).Name)
		if fieldOpt == "" {
			continue
		}
		fieldName := reflect.ValueOf(fieldOpt)
		if err := e.marshal(fieldName); err != nil {
			return err
		}
		if err := e.marshal(fieldValue); err != nil {
			return err
		}
	}
	return e.WriteByte('e')
}

func (e *encodeState) marshalPtr(val reflect.Value) error {
	return e.marshal(val.Elem())
}

func (e *encodeState) marshalInterface(val reflect.Value) error {
	return e.marshal(val.Elem())
}
