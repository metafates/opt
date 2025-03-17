package opt

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
)

var _ interface {
	sql.Scanner
	driver.Valuer
} = (*Opt[any])(nil)

// Scan implements the [sql.Scanner] interface.
func (o *Opt[T]) Scan(src any) error {
	if src == nil {
		*o = None[T]()

		return nil
	}

	// is is only possible to assert interfaces, so convert first
	// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#why-not-permit-type-assertions-on-values-whose-type-is-a-type-parameter
	var value T

	if scanner, ok := any(&value).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return fmt.Errorf("failed to scan: %w", err)
		}

		*o = Some(value)

		return nil
	}

	if converted, err := driver.DefaultParameterConverter.ConvertValue(src); err == nil {
		if v, ok := converted.(T); ok {
			*o = Some(v)

			return nil
		}
	}

	return o.scanConvertValue(src)
}

// Value implements the [driver.Valuer] interface.
//
// Use unwrap methods (e.g. [Opt.TryGet]) instead for getting the go value
func (o Opt[T]) Value() (driver.Value, error) {
	if !o.hasValue {
		return nil, nil
	}

	return driver.DefaultParameterConverter.ConvertValue(o.value)
}

func (o *Opt[T]) scanConvertValue(src any) error {
	var nullable sql.Null[T]

	if err := nullable.Scan(src); err != nil {
		return err
	}

	if nullable.Valid {
		*o = Some(nullable.V)
	} else {
		*o = None[T]()
	}

	return nil
}
