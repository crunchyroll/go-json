// Package converter converts a JSON message into a set of Go structs that describe the minimal
// schema into which it can be unmarshalled.
package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

const (
	// This is probably better done with a text template. Oops.
	getterFmt = `func (x *%s) Get%s() *%s {
	if x.%[2]s == nil {
		return &%s{}
	}

	return x.%[2]s
}
`
)

var (
	collisions = map[string]*conversionType{}
)

type conversionType struct {
	Parent         *conversionType
	IsPrimitive    bool
	arrayLevels    int
	emptyArray     bool
	Type           string
	jsonTag        string
	Fields         map[string]*conversionType
	typeCollisions map[string]bool
}

func convertFieldName(fn string) string {
	words := strings.Split(fn, "_")
	for i, w := range words {
		words[i] = strings.Title(w)
	}
	return strings.Join(words, "")
}

// Returns the qualified type of this conversion type. This return value is meaningless if c
// is a primitive type (e.g., bool, string, float64)
func (c *conversionType) QualifiedType() string {
	var qt string

	if c.emptyArray {
		return "interface{}"
	}
	if c.IsPrimitive || c.Parent == nil {
		return c.Type
	}

	for {
		if len(c.typeCollisions) > 0 {
			qt = fmt.Sprintf("%s_%s", c.Parent.QualifiedType(), c.Type)
		} else {
			if c.typeCollisions == nil {
				c.typeCollisions = map[string]bool{}
			}
			qt = c.Type
		}

		col, ok := collisions[qt]
		if !ok || col == c {
			break
		}

		c.typeCollisions[qt] = true
		collisions[qt] = c

		col.QualifiedType()

		continue
	}

	collisions[qt] = c

	return qt
}

func (c *conversionType) FromJSONObject(jsonObject map[string]interface{}) {
	c.Fields = map[string]*conversionType{}
	for field, value := range jsonObject {
		ct := &conversionType{Parent: c, jsonTag: field}
		field = convertFieldName(field)
		c.Fields[field] = ct
		ct.Convert(field, value)
	}
}

func (c *conversionType) FromArray(name string, object interface{}) {
	c.arrayLevels += 1
	v := reflect.ValueOf(object)
	if v.Len() == 0 {
		c.emptyArray = true
		return
	}

nestingLoop:
	for {
		v = v.Index(0)

		switch v.Kind() {
		case reflect.Slice, reflect.Array:
			c.arrayLevels += 1
			if v.Len() == 0 {
				c.emptyArray = true
				return
			}
		default:
			break nestingLoop
		}
	}

	c.Convert(name, v.Interface())
}

func (c *conversionType) Convert(name string, object interface{}) {
	c.Type = name

	t := reflect.TypeOf(object)

	switch t.Kind() {
	case reflect.Map:
		// Field type is an object
		jsonObject := object.(map[string]interface{})
		c.FromJSONObject(jsonObject)
	case reflect.Slice, reflect.Array:
		c.FromArray(name, object)
	default:
		c.IsPrimitive = true
		c.Type = t.Name()
	}

	c.QualifiedType()
}

func (c *conversionType) WriteGetters(w io.Writer) {
	for fn, fv := range c.Fields {
		// Don't generate getters for primitive or array types.
		if fv.IsPrimitive || fv.arrayLevels != 0 {
			continue
		}

		fmt.Fprintln(w, "")

		qt := c.QualifiedType()
		fieldQT := fv.QualifiedType()
		fmt.Fprintf(w, getterFmt, qt, fn, fieldQT)
	}
}

func (c *conversionType) WriteStructs(w io.Writer) {
	if c.IsPrimitive || c.emptyArray {
		return
	}

	// First write the sub-types out.
	for _, fv := range c.Fields {
		fv.WriteStructs(w)
	}

	fmt.Fprintf(w, "\ntype %s struct {\n", c.QualifiedType())

	for fn, fv := range c.Fields {
		arrayPrefix := strings.Repeat("[]", fv.arrayLevels)

		var ptrPrefix string
		if !fv.IsPrimitive && !fv.emptyArray {
			ptrPrefix = "*"
		}

		fmt.Fprintf(w, "\t%s %s%s%s `json:%q`\n", fn, arrayPrefix, ptrPrefix, fv.QualifiedType(), fv.jsonTag)
	}

	fmt.Fprintln(w, "}")

	c.WriteGetters(w)
}

func Convert(name string, raw []byte) (*conversionType, error) {
	jsonObject := map[string]interface{}{}
	err := json.Unmarshal(raw, &jsonObject)
	if err != nil {
		return nil, err
	}

	c := &conversionType{}
	c.Convert(name, jsonObject)

	return c, nil
}
