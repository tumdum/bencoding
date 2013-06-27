package bencoding

import (
	"errors"
	"reflect"
)

func bind(d map[string]interface{}, o reflect.Value) error {
	for k, v := range d {
		if e := bindField(k, v, o); e != nil {
			return e
		}
	}
	return nil
}

func findCorrectlyTaggedField(name string, o reflect.Value) reflect.Value {
	ot := o.Type()
	for i := 0; i < o.NumField(); i++ {
		fname := ot.Field(i).Name
		fieldOpts := extractFieldOptions(o, fname)
		if string(fieldOpts) == name {
			return o.Field(i)
		}
	}
	return reflect.Value{}
}

func bindField(k string, value interface{}, o reflect.Value) error {
	if field := findCorrectlyTaggedField(k, o); !field.IsValid() {
		return nil // missing struct fields are not errors
	} else {
		vvalue := reflect.ValueOf(value)
		if vvalue.Type().AssignableTo(field.Type()) {
			field.Set(vvalue)
		} else if isBindableStructAndDict(field, value) {
			return bind(value.(map[string]interface{}), field)
		} else if field.Kind() == reflect.Ptr {
			return bindPtr(value, field)
		} else {
			return errors.New("field '" + k + "' failed to bind")
		}
	}
	return nil
}

type withKind interface {
	Kind() reflect.Kind
}

func isBindableStructAndDict(structure withKind, dictionary interface{}) bool {
	if _, isDict := dictionary.(map[string]interface{}); structure.Kind() == reflect.Struct && isDict {
		return true
	}
	return false
}

func bindPtr(value interface{}, field reflect.Value) error {
	vvalue := reflect.ValueOf(value)
	if vvalue.Type().AssignableTo(field.Type().Elem()) {
		prepare(field)
		field.Elem().Set(vvalue)
	} else if isBindableStructAndDict(field.Type().Elem(), value) {
		prepare(field)
		bind(value.(map[string]interface{}), field.Elem())
	} else {
		return errors.New("unable to bind")
	}
	return nil
}

func prepare(v reflect.Value) {
	if v.IsNil() {
		v.Set(reflect.New(v.Type().Elem()))
	}
}
