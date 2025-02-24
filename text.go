package opt

import (
	"encoding"
	"encoding/json"
)

var _ interface {
	encoding.TextMarshaler
	encoding.TextUnmarshaler
} = (*Opt[any])(nil)

// MarshalText implemenets [encoding.TextMarshaler] interface
func (o Opt[T]) MarshalText() ([]byte, error) {
	return json.Marshal(o)
}

// UnmarshalText implemenets [encoding.TextUnmarshaler] interface
func (o *Opt[T]) UnmarshalText(data []byte) error {
	return json.Unmarshal(data, o)
}
