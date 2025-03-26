// Package opt provides safe abstractions over optional values.
//
// Inspired by the [Option type in Rust] and follows the same ideas and function signatures.
//
// [Option type in Rust]: https://doc.rust-lang.org/std/option/enum.Option.html
package opt

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// Opt (option) represents an optional value.
// Every option is either [Some] and contains a value, or [None], and does not.
//
// Opt also separates explicitly set values, see [Opt.IsExplicit].
// Implicit [Some] is unreachable state.
//
// Use cases:
//   - Initial values
//   - Return values for functions that are not defined over their entire input range (partial functions)
//   - Optional struct fields
//   - Optional function arguments
type Opt[T any] struct {
	value    T
	hasValue bool
	explicit bool
}

// Some returns option with some value
func Some[T any](value T) Opt[T] {
	return Opt[T]{value: value, hasValue: true, explicit: true}
}

// None returns an option with no value.
// None option could also be defined like that:
//
//	var none Opt[any]
func None[T any]() Opt[T] {
	return Opt[T]{explicit: true}
}

// FromPtr returns [Some] with the underlying pointer value if it's not nil or [None] otherwise.
//
// NOTE: Do not use this for converting [proto.Message]s, as they can not be dereferenced without breaking internal reflection mechanism.
// Use [FromProto] for that purpose
func FromPtr[T any](ptr *T) Opt[T] {
	if ptr == nil {
		return None[T]()
	}

	return Some(*ptr)
}

// FromZero returns [Some] with the given value if it's not zero value or [None] otherwise.
//
// This function requires value type to be comparable so that it can be checked for zero value without using reflection
func FromZero[T comparable](value T) Opt[T] {
	var zero T

	if value == zero {
		return None[T]()
	}

	return Some(value)
}

// FromProto converts [proto.Message] to either [Some] value, if the message is valid, or [None].
//
// An invalid message is an empty, read-only value.
// An invalid message often corresponds to a nil pointer of the concrete message type, but the details are implementation dependent.
//
// See [protoreflect.Message.IsValid]
func FromProto[T proto.Message](msg T) Opt[T] {
	if msg.ProtoReflect().IsValid() {
		return Some(msg)
	}

	return None[T]()
}

// FromTuple returns [Some] with the given value if ok is true, [None] otherwise
func FromTuple[T any](value T, ok bool) Opt[T] {
	if ok {
		return Some(value)
	}

	return None[T]()
}

// IsExplicit reports whether this option was explicitly specified as either [None] or [Some].
// This property is also applicable for decoded values, such as ones from [json.Unmarshal].
//
// If [Opt.IsSome] is true, it is guaranteed that this function will also return true
//
// This propery allows to represent all three possible states:
//   - The value is not set
//   - The value is explicitly set to [None]
//   - The value is explicitly set to a given [Some] value
func (o Opt[T]) IsExplicit() bool {
	return o.explicit
}

// IsSome returns true if the option is a [Some] value.
//
// If this option is [Some] it is guaranteed to be explicit. See [Opt.IsExplicit]
func (o Opt[T]) IsSome() bool {
	return o.hasValue
}

// IsSomeAnd returns true if the option is a [Some] and the value inside of it matches a predicate.
func (o Opt[T]) IsSomeAnd(and func(T) bool) bool {
	if o.hasValue {
		return and(o.value)
	}

	return false
}

// IsNone returns true if the option is a [None] value.
func (o Opt[T]) IsNone() bool {
	return !o.hasValue
}

// IsNoneOr returns true if the option is a [None] or the value inside of it matches a predicate.
func (o Opt[T]) IsNoneOr(orElse func(T) bool) bool {
	if !o.hasValue {
		return true
	}

	return orElse(o.value)
}

// GetOrEmpty returns the contained [Some] value or an empty value for this type.
func (o Opt[T]) GetOrEmpty() T {
	// we could just return o.value ignoring the o.hasValue, but see [Opt.TryGet] explanation
	if o.hasValue {
		return o.value
	}

	var empty T
	return empty
}

// TryGet returns the contained [Some] value or an empty value for this type
// and boolean stating if this option is [Some]
//
// If you only need a contained value or an empty one use [Opt.GetOrEmpty]
func (o Opt[T]) TryGet() (T, bool) {
	// we could just return o.value, o.hasValue
	// but if T is a pointer-value it makes it possible to modify underlying empty value for all the future calls.
	// the risk is still there for non-empty values (unless we deep-clone), but it is usually expected behaviour
	if o.hasValue {
		return o.value, true
	}

	var empty T
	return empty, false
}

// MustGet returns the contained [Some] value.
//
// Panics if the self value equals [None].
func (o Opt[T]) MustGet() T {
	if o.hasValue {
		return o.value
	}

	panic("called MustGet on empty option")
}

// GetOr returns the contained [Some] value or a provided default.
func (o Opt[T]) GetOr(or T) T {
	if o.hasValue {
		return o.value
	}

	return or
}

// GetOrElse returns the contained [Some] value or computes it from a function.
func (o Opt[T]) GetOrElse(orElse func() T) T {
	if o.hasValue {
		return o.value
	}

	return orElse()
}

// Inspect calls a function with a contained value if [Some].
//
// Returns the original option.
func (o Opt[T]) Inspect(f func(T)) Opt[T] {
	if o.hasValue {
		f(o.value)
	}

	return o
}

// Map maps a value by applying a function to a contained value (if [Some]) or returns [None] (if [None]).
//
// See [Map] if you need to return a different type.
func (o Opt[T]) Map(f func(T) T) Opt[T] {
	if o.hasValue {
		return Some(f(o.value))
	}

	return o
}

// And returns [None] if the option is [None], otherwise returns `and`.
func (o Opt[T]) And(and Opt[T]) Opt[T] {
	if o.hasValue {
		return and
	}

	return o
}

// AndThen returns [None] if the option is [None], otherwise calls `andThen` with
// the wrapped value and returns the result.
//
// See [AndThen] if you need to return a different type
func (o Opt[T]) AndThen(andThen func(T) Opt[T]) Opt[T] {
	if o.hasValue {
		return andThen(o.value)
	}

	return o
}

// Or returns itself if it contains a value, otherwise returns `or`.
func (o Opt[T]) Or(or Opt[T]) Opt[T] {
	if o.hasValue {
		return o
	}

	return or
}

// OrElse returns itself if it contains a value, otherwise calls `orElse` and returns the result.
func (o Opt[T]) OrElse(orElse func() Opt[T]) Opt[T] {
	if o.hasValue {
		return o
	}

	return orElse()
}

// Filter returns [None] if the option is [None], otherwise calls predicate with the wrapped value and returns:
//   - [Some] if predicate returns true.
//   - [None] if predicate returns false.
func (o Opt[T]) Filter(predicate func(T) bool) Opt[T] {
	if !o.hasValue {
		return o
	}

	if !predicate(o.value) {
		return None[T]()
	}

	return o
}

// ToPtr returns pointer to the value if the option is [Some] or nil otherwise.
//
// The underlying value of the pointer is safe to modify, as it is copied before return
// to avoid changes to the original value.
func (o Opt[T]) ToPtr() *T {
	if o.hasValue {
		value := o.value
		return &value
	}

	return nil
}

func (o Opt[T]) String() string {
	if o.hasValue {
		return fmt.Sprintf("Some(%v)", o.value)
	}

	return "None"
}

// ToSlice returns singleton slice if option is [Some] or nil otherwise.
func (o Opt[T]) ToSlice() []T {
	if o.hasValue {
		return []T{o.value}
	}

	return nil
}

// IndexSlice returns closure that accepts index and returns [Some] slice value at the given index
// if the index exists or [None] otherwise
func IndexSlice[S ~[]T, T any](slice S) func(index int) Opt[T] {
	return func(index int) Opt[T] {
		if index >= len(slice) || index < 0 {
			return None[T]()
		}

		return Some(slice[index])
	}
}

// IndexMap returns closure that accepts key and returns [Some] map value at the given key
// if the key exists or [None] otherwise
func IndexMap[M ~map[K]V, K comparable, V any](m M) func(key K) Opt[V] {
	return func(key K) Opt[V] {
		value, ok := m[key]

		if ok {
			return Some(value)
		}

		return None[V]()
	}
}

// Map maps a value by applying a function to a contained value (if [Some]) or returns [None] (if [None]).
//
// This function allows `f` to return a different type in contrast to the [Opt.Map] which is limited
// by the lack of method type parameters in Go.
func Map[T, U any](option Opt[T], f func(T) U) Opt[U] {
	if option.hasValue {
		return Some(f(option.value))
	}

	return None[U]()
}

// AndThen returns [None] if the option is [None], otherwise calls `f` with
// the wrapped value and returns the result.
//
// This function allows `f` to return a different type in contrast to the [Opt.AndThen] which is limited
// by the lack of method type parameters in Go.
func AndThen[T, U any](option Opt[T], f func(T) Opt[U]) Opt[U] {
	if option.hasValue {
		return f(option.value)
	}

	return None[U]()
}
