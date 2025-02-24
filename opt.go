// Package opt provides safe abstractions over optional values.
//
// Inspired by the [Option type in Rust] and follows the same ideas and function signatures.
//
// [Option type in Rust]: https://doc.rust-lang.org/std/option/enum.Option.html
package opt

import (
	"fmt"
)

// Opt (option) represents an optional value.
// Every option is either [Some] and contains a value, or [None], and does not.
//
// Use cases:
//   - Initial values
//   - Return values for functions that are not defined over their entire input range (partial functions)
//   - Optional struct fields
//   - Optional function arguments
type Opt[T any] struct {
	value T
	ok    bool
}

// Some returns option with some value
func Some[T any](value T) Opt[T] {
	return Opt[T]{value: value, ok: true}
}

// None returns an option with no value.
// None option could also be defined like that:
//
//	var none Opt[any]
func None[T any]() Opt[T] {
	return Opt[T]{}
}

// FromPtr returns [Some] with the underlying pointer value if it's not nil or [None] otherwise
func FromPtr[T any](ptr *T) Opt[T] {
	if ptr == nil {
		return None[T]()
	}

	return Some(*ptr)
}

// IsSome returns true if the option is a [Some] value.
func (o Opt[T]) IsSome() bool {
	return o.ok
}

// IsSomeAnd returns true if the option is a [Some] and the value inside of it matches a predicate.
func (o Opt[T]) IsSomeAnd(and func(T) bool) bool {
	if o.ok {
		return and(o.value)
	}

	return false
}

// IsNone returns true if the option is a [None] value.
func (o Opt[T]) IsNone() bool {
	return !o.ok
}

// IsNoneOr returns true if the option is a [None] or the value inside of it matches a predicate.
func (o Opt[T]) IsNoneOr(orElse func(T) bool) bool {
	if !o.ok {
		return true
	}

	return orElse(o.value)
}

// UnwrapOrEmpty returns the contained [Some] value or an empty value for this type.
func (o Opt[T]) UnwrapOrEmpty() T {
	// since value is not a pointer it contains empty value even if the option is none
	return o.value
}

// TryUnwrap returns the contained [Some] value or an empty value for this type
// and boolean stating if this option is [Some]
//
// If you only need a contained value or an empty one use [Opt.UnwrapOrEmpty]
func (o Opt[T]) TryUnwrap() (T, bool) {
	return o.value, o.ok
}

// MustUnwrap returns the contained [Some] value.
//
// Panics if the self value equals [None].
func (o Opt[T]) MustUnwrap() T {
	if o.ok {
		return o.value
	}

	panic("called MustUnwrap on empty option")
}

// UnwrapOr returns the contained [Some] value or a provided default.
func (o Opt[T]) UnwrapOr(or T) T {
	if o.ok {
		return o.value
	}

	return or
}

// UnwrapOrElse returns the contained [Some] value or computes it from a function.
func (o Opt[T]) UnwrapOrElse(orElse func() T) T {
	if o.ok {
		return o.value
	}

	return orElse()
}

// Inspect calls a function with a contained value if [Some].
//
// Returns the original option.
func (o Opt[T]) Inspect(f func(T)) Opt[T] {
	if o.ok {
		f(o.value)
	}

	return o
}

// Map maps a value by applying a function to a contained value (if [Some]) or returns [None] (if [None]).
//
// See [Map] if you need to return a different type.
func (o Opt[T]) Map(f func(T) T) Opt[T] {
	if o.ok {
		return Some(f(o.value))
	}

	return o
}

// And returns [None] if the option is [None], otherwise returns `and`.
func (o Opt[T]) And(and Opt[T]) Opt[T] {
	if o.ok {
		return and
	}

	return o
}

// AndThen returns [None] if the option is [None], otherwise calls `andThen` with
// the wrapped value and returns the result.
//
// See [AndThen] if you need to return a different type
func (o Opt[T]) AndThen(andThen func(T) Opt[T]) Opt[T] {
	if o.ok {
		return andThen(o.value)
	}

	return o
}

// Or returns itself if it contains a value, otherwise returns `or`.
func (o Opt[T]) Or(or Opt[T]) Opt[T] {
	if o.ok {
		return o
	}

	return or
}

// OrElse returns itself if it contains a value, otherwise calls `orElse` and returns the result.
func (o Opt[T]) OrElse(orElse func() Opt[T]) Opt[T] {
	if o.ok {
		return o
	}

	return orElse()
}

// Filter returns [None] if the option is [None], otherwise calls predicate with the wrapped value and returns:
//   - [Some] if predicate returns true.
//   - [None] if predicate returns false.
func (o Opt[T]) Filter(predicate func(T) bool) Opt[T] {
	if !o.ok {
		return o
	}

	if !predicate(o.value) {
		return None[T]()
	}

	return o
}

// ToPtr returns pointer to the underlying value if the option is [Some] or nil otherwise
func (o Opt[T]) ToPtr() *T {
	if o.ok {
		return &o.value
	}

	return nil
}

func (o Opt[T]) String() string {
	if o.ok {
		return fmt.Sprintf("Some(%v)", o.value)
	}

	return "None"
}

// IndexSlice returns [Some] slice value at the given index if the index exists or [None] otherwise
func IndexSlice[S ~[]T, T any](slice S, index int) Opt[T] {
	if index >= len(slice) {
		return None[T]()
	}

	return Some(slice[index])
}

// IndexMap returns [Some] map value at the given key if the key exists or [None] otherwise
func IndexMap[M ~map[K]V, K comparable, V any](m M, key K) Opt[V] {
	value, ok := m[key]

	return Opt[V]{value: value, ok: ok}
}

// Map maps a value by applying a function to a contained value (if [Some]) or returns [None] (if [None]).
//
// This function allows `f` to return a different type in contrast to the [Opt.Map] which is limited
// by the lack of method type parameters in Go.
func Map[T, U any](option Opt[T], f func(T) U) Opt[U] {
	if option.ok {
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
	if option.ok {
		return f(option.value)
	}

	return None[U]()
}
