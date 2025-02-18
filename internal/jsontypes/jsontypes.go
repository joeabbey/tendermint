// Package jsontypes supports decoding for interface types whose concrete
// implementations need to be stored as JSON. To do this, concrete values are
// packaged in wrapper objects having the form:
//
//   {
//     "type": "<type-tag>",
//     "value": <json-encoding-of-value>
//   }
//
// This package provides a registry for type tag strings and functions to
// encode and decode wrapper objects.
package jsontypes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

// The Tagged interface must be implemented by a type in order to register it
// with the jsontypes package. The TypeTag method returns a string label that
// is used to distinguish objects of that type.
type Tagged interface {
	TypeTag() string
}

// registry records the mapping from type tags to value types.  Values in this
// map must be normalized to non-pointer types.
var registry = struct {
	types map[string]reflect.Type
}{types: make(map[string]reflect.Type)}

// register adds v to the type registry. It reports an error if the tag
// returned by v is already registered.
func register(v Tagged) error {
	tag := v.TypeTag()
	if t, ok := registry.types[tag]; ok {
		return fmt.Errorf("type tag %q already registered to %v", tag, t)
	}
	typ := reflect.TypeOf(v)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	registry.types[tag] = typ
	return nil
}

// MustRegister adds v to the type registry. It will panic if the tag returned
// by v is already registered. This function is meant for use during program
// initialization.
func MustRegister(v Tagged) {
	if err := register(v); err != nil {
		panic(err)
	}
}

type wrapper struct {
	Type  string          `json:"type"`
	Value json.RawMessage `json:"value"`
}

// Marshal marshals a JSON wrapper object containing v. If v == nil, Marshal
// returns the JSON "null" value without error.
func Marshal(v Tagged) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return json.Marshal(wrapper{
		Type:  v.TypeTag(),
		Value: data,
	})
}

// Unmarshal unmarshals a JSON wrapper object into v. It reports an error if
// the data do not encode a valid wrapper object, if the wrapper's type tag is
// not registered with jsontypes, or if the resulting value is not compatible
// with the type of v.
func Unmarshal(data []byte, v interface{}) error {
	// Verify that the target is some kind of pointer.
	target := reflect.ValueOf(v)
	if target.Kind() != reflect.Ptr {
		return fmt.Errorf("target %T is not a pointer", v)
	}

	var w wrapper
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&w); err != nil {
		return fmt.Errorf("invalid type wrapper: %w", err)
	}
	typ, ok := registry.types[w.Type]
	if !ok {
		return fmt.Errorf("unknown type tag: %q", w.Type)
	} else if !typ.AssignableTo(target.Elem().Type()) {
		return fmt.Errorf("type %v not assignable to %T", typ, v)
	}

	obj := reflect.New(typ)
	if err := json.Unmarshal(w.Value, obj.Interface()); err != nil {
		return fmt.Errorf("decoding wrapped value: %w", err)
	}
	target.Elem().Set(obj.Elem())
	return nil
}
