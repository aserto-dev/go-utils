package protoutil

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/pkg/errors"
	"google.golang.org/protobuf/runtime/protoimpl"
	"google.golang.org/protobuf/types/known/structpb"
)

// NewValue extends `structpb.NewValue` so it understands `json.Number`
func NewValue(v interface{}) (*structpb.Value, error) {
	switch v := v.(type) {
	case json.Number:
		floatValue, err := v.Float64()
		if err != nil {
			return nil, errors.Wrapf(err, "error converting json.Number '%s' to float64", v.String())
		}
		return structpb.NewNumberValue(floatValue), nil
	case map[string]interface{}:
		v2, err := NewStruct(v)
		if err != nil {
			return nil, err
		}
		return structpb.NewStructValue(v2), nil
	case []interface{}:
		v2, err := NewList(v)
		if err != nil {
			return nil, err
		}
		return structpb.NewListValue(v2), nil
	default:
		return structpb.NewValue(v)
	}
}

// NewStruct constructs a Struct from a general-purpose Go map.
// The map keys must be valid UTF-8.
// The map values are converted using NewValue.
func NewStruct(v map[string]interface{}) (*structpb.Struct, error) {
	x := &structpb.Struct{Fields: make(map[string]*structpb.Value, len(v))}
	for k, v := range v {
		if !utf8.ValidString(k) {
			return nil, protoimpl.X.NewError("invalid UTF-8 in string: %q", k)
		}
		var err error
		x.Fields[k], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}

// NewList constructs a ListValue from a general-purpose Go slice.
// The slice elements are converted using NewValue.
func NewList(v []interface{}) (*structpb.ListValue, error) {
	x := &structpb.ListValue{Values: make([]*structpb.Value, len(v))}
	for i, v := range v {
		var err error
		x.Values[i], err = NewValue(v)
		if err != nil {
			return nil, err
		}
	}
	return x, nil
}
