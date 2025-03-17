package opt

import "encoding/json"

var _ interface {
	json.Marshaler
	json.Unmarshaler
} = (*Opt[any])(nil)

// MarshalJSON implemenets [json.Marshaler] interface
func (o Opt[T]) MarshalJSON() ([]byte, error) {
	if o.hasValue {
		return json.Marshal(o.value)
	}

	return []byte("null"), nil
}

// UnmarshalJSON implemenets [json.Unmarshaler] interface
func (o *Opt[T]) UnmarshalJSON(b []byte) error {
	var value *T

	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	if value == nil {
		*o = None[T]()
	} else {
		*o = Some(*value)
	}

	return nil
}
