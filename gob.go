package opt

import (
	"encoding/gob"
)

var _ interface {
	gob.GobEncoder
	gob.GobDecoder
} = (*Opt[any])(nil)

// GobEncode implemenets [gob.GobEncoder] interface
func (o Opt[T]) GobEncode() ([]byte, error) {
	return o.MarshalBinary()
}

// GobDecode implemenets [gob.GobDecoder] interface
func (o *Opt[T]) GobDecode(data []byte) error {
	return o.UnmarshalBinary(data)
}
