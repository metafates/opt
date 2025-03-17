package opt

import (
	"bytes"
	"database/sql/driver"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpt_ToPtr(t *testing.T) {
	x := Some("a")

	*x.ToPtr() = "b"

	require.Equal(t, Some("a"), x)
}

func TestOpt_FromZero(t *testing.T) {
	require.Equal(t, None[string](), FromZero(""))
	require.Equal(t, Some("foo"), FromZero("foo"))

	require.Equal(t, None[int](), FromZero(0))
	require.Equal(t, Some(1), FromZero(1))

	require.Equal(t, None[bool](), FromZero(false))
	require.Equal(t, Some(true), FromZero(true))
}

func TestOpt_Scan(t *testing.T) {
	t.Run("nil scan", func(t *testing.T) {
		var option Opt[string]

		err := option.Scan(nil)

		require.NoError(t, err)
		require.Equal(t, None[string](), option)
	})

	t.Run("scan scanner", func(t *testing.T) {
		t.Run("null", func(t *testing.T) {
			var option Opt[string]

			err := option.Scan(driver.Value(nil))

			require.NoError(t, err)
			require.Equal(t, None[string](), option)
		})

		t.Run("not null", func(t *testing.T) {
			var option Opt[string]

			err := option.Scan(driver.Value("go"))

			require.NoError(t, err)
			require.Equal(t, Some("go"), option)
		})
	})

	t.Run("scan regular value", func(t *testing.T) {
		var option Opt[string]

		err := option.Scan("go")

		require.NoError(t, err)
		require.Equal(t, Some("go"), option)
	})
}

func TestOpt_Value(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		option := None[string]()

		value, err := option.Value()
		require.NoError(t, err)
		require.Equal(t, nil, value)
	})

	t.Run("some", func(t *testing.T) {
		option := Some("apple")

		value, err := option.Value()
		require.NoError(t, err)
		require.Equal(t, "apple", value)
	})
}

func TestOpt_IsExplicit(t *testing.T) {
	var foo, bar struct {
		Name string
		Age  Opt[int]
	}

	JSONEncoder{}.Decode(t, []byte(`{"name":"bar","age":null}`), &foo)
	require.True(t, foo.Age.IsExplicit())

	JSONEncoder{}.Decode(t, []byte(`{"name":"bar"}`), &bar)
	require.False(t, bar.Age.IsExplicit())
}

func TestEncode(t *testing.T) {
	testCases := []struct {
		name      string
		wantOpt   Opt[string]
		wantBytes []byte
		encoder   Encoder
	}{
		{
			name:      "json some",
			wantOpt:   Some("apple"),
			wantBytes: []byte(`"apple"`),
			encoder:   JSONEncoder{},
		},
		{
			name:      "json none",
			wantOpt:   None[string](),
			wantBytes: []byte(`null`),
			encoder:   JSONEncoder{},
		},
		{
			name:      "text some",
			wantOpt:   Some("apple"),
			wantBytes: []byte(`"apple"`),
			encoder:   TextEncoder{},
		},
		{
			name:      "text none",
			wantOpt:   None[string](),
			wantBytes: []byte(`null`),
			encoder:   TextEncoder{},
		},
		{
			name:      "binary some",
			wantOpt:   Some("apple"),
			wantBytes: []byte{1, 0x8, 0xC, 0x0, 0x5, 0x61, 0x70, 0x70, 0x6C, 0x65},
			encoder:   BinaryEncoder{},
		},
		{
			name:      "binary none",
			wantOpt:   None[string](),
			wantBytes: []byte{0},
			encoder:   BinaryEncoder{},
		},
		{
			name:      "gob some",
			wantOpt:   Some("apple"),
			wantBytes: []byte{0x16, 0x7F, 0x5, 0x1, 0x1, 0xB, 0x4F, 0x70, 0x74, 0x5B, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0x5D, 0x1, 0xFF, 0x80, 0x0, 0x0, 0x0, 0xE, 0xFF, 0x80, 0x0, 0xA, 0x1, 0x8, 0xC, 0x0, 0x5, 0x61, 0x70, 0x70, 0x6C, 0x65},
			encoder:   GobEncoder{},
		},
		{
			name:      "gob none",
			wantOpt:   None[string](),
			wantBytes: []byte{0x16, 0x7F, 0x5, 0x1, 0x1, 0xB, 0x4F, 0x70, 0x74, 0x5B, 0x73, 0x74, 0x72, 0x69, 0x6E, 0x67, 0x5D, 0x1, 0xFF, 0x80, 0x0, 0x0, 0x0, 0x5, 0xFF, 0x80, 0x0, 0x1, 0x0},
			encoder:   GobEncoder{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("unmarshal", func(t *testing.T) {
				var opt Opt[string]

				tc.encoder.Decode(t, tc.wantBytes, &opt)

				require.Equal(t, tc.wantOpt, opt)
			})

			t.Run("marshal", func(t *testing.T) {
				bytes := tc.encoder.Encode(t, tc.wantOpt)

				require.Equal(t, tc.wantBytes, bytes)
			})
		})
	}
}

type Encoder interface {
	Encode(t *testing.T, v any) []byte
	Decode(t *testing.T, data []byte, v any)
}

type JSONEncoder struct{}

func (JSONEncoder) Encode(t *testing.T, v any) []byte {
	t.Helper()

	b, err := json.Marshal(v)
	require.NoError(t, err)

	return b
}

func (JSONEncoder) Decode(t *testing.T, data []byte, v any) {
	t.Helper()

	err := json.Unmarshal(data, v)
	require.NoError(t, err)
}

type TextEncoder struct{}

func (TextEncoder) Encode(t *testing.T, v any) []byte {
	t.Helper()

	b, err := v.(encoding.TextMarshaler).MarshalText()
	require.NoError(t, err)

	return b
}

func (TextEncoder) Decode(t *testing.T, data []byte, v any) {
	t.Helper()

	err := v.(encoding.TextUnmarshaler).UnmarshalText(data)
	require.NoError(t, err)
}

type BinaryEncoder struct{}

func (BinaryEncoder) Encode(t *testing.T, v any) []byte {
	t.Helper()

	b, err := v.(encoding.BinaryMarshaler).MarshalBinary()
	require.NoError(t, err)

	return b
}

func (BinaryEncoder) Decode(t *testing.T, data []byte, v any) {
	t.Helper()

	err := v.(encoding.BinaryUnmarshaler).UnmarshalBinary(data)
	require.NoError(t, err)
}

type GobEncoder struct{}

func (GobEncoder) Encode(t *testing.T, v any) []byte {
	t.Helper()

	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(v)
	require.NoError(t, err)

	return buf.Bytes()
}

func (GobEncoder) Decode(t *testing.T, data []byte, v any) {
	t.Helper()

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(v)
	require.NoError(t, err)
}
