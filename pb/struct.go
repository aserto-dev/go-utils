package pb

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/structpb"
)

// NewStruct, returns *structpb.Struct instance with initialized Fields map.
func NewStruct() *structpb.Struct {
	return &structpb.Struct{Fields: make(map[string]*structpb.Value)}
}

// Attempts to construct a structpb.Struct from an interface{}, using mapstructure
// to map from 'from' to map[string]interface{}
func ToStructValue(from interface{}) (*structpb.Value, error) {
	result := map[string]interface{}{}
	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  &result,
	})
	if err != nil {
		return nil, err
	}

	err = dec.Decode(from)
	if err != nil {
		return nil, err
	}

	val, err := structpb.NewStruct(result)
	if err != nil {
		return nil, err
	}

	return structpb.NewStructValue(val), err
}

// Attempts to map a *structpb.Struct to an interface{}, using mapstructure
// to map from the struct's fields to 'to'
func FromStructValue(from *structpb.Value, to interface{}) error {
	val := from.GetStructValue()
	if val == nil {
		return errors.New("no struct value present")
	}

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		TagName: "json",
		Result:  to,
	})
	if err != nil {
		return err
	}

	err = dec.Decode(val.AsMap())
	if err != nil {
		return err
	}

	return nil
}
