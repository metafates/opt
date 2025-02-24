package opt

import (
	"bytes"
	"encoding"
	"encoding/gob"
	"errors"
)

var _ interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
} = (*Opt[any])(nil)

// MarshalBinary implemenets [encoding.BinaryMarshaler] interface
func (o Opt[T]) MarshalBinary() ([]byte, error) {
	if !o.ok {
		return []byte{0}, nil
	}

	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(o.value); err != nil {
		return []byte{}, err
	}

	return append([]byte{1}, buf.Bytes()...), nil
}

// UnmarshalBinary implemenets [encoding.BinaryUnmarshaler] interface
func (o *Opt[T]) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return errors.New("Opt[T].UnmarshalBinary: no data")
	}

	if data[0] == 0 {
		*o = None[T]()
		return nil
	}

	buf := bytes.NewBuffer(data[1:])
	dec := gob.NewDecoder(buf)

	var value T

	if err := dec.Decode(&value); err != nil {
		return err
	}

	*o = Some(value)

	return nil
}
