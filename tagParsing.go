package bencoding

import (
	"reflect"
	"regexp"
)

func tagForFieldNamed(value reflect.Value, name string) string {
	field, hasField := value.Type().FieldByName(name)
	if !hasField {
		return ""
	}
	return string(field.Tag)
}

func parseTag(tag string) *string {
	const fieldRegexp = `bencoding:"([\w- ]*)"`
	reg := regexp.MustCompile(fieldRegexp)
	if matches := reg.FindStringSubmatch(tag); len(matches) > 2 {
		panic("regexp for parsing fields seems to be wrong -- more then two groups returned")
	} else if len(matches) == 2 {
		return &matches[1]
	} else {
		return nil
	}
}

func extractFieldOptions(v reflect.Value, name string) string {
	tag := tagForFieldNamed(v, name)
	bencodingTag := parseTag(tag)
	if bencodingTag == nil {
		return name
	} else if *bencodingTag == "" || *bencodingTag == "-" {
		return ""
	} else {
		return *bencodingTag
	}
}
