package converters

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNullStringToString(t *testing.T) {
	var tests = []struct {
		name  string
		input sql.NullString
		want  string
	}{
		{name: "valid non-empty nullString", input: sql.NullString{String: "valid", Valid: true}, want: "valid"},
		{name: "invalid empty nullString", input: sql.NullString{String: "", Valid: false}, want: ""},
		{name: "valid empty nullString", input: sql.NullString{String: "", Valid: true}, want: ""},
		{name: "invalid non-empty nullString", input: sql.NullString{String: "valid", Valid: false}, want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NullStringToString(tt.input))
		})
	}
}

func TestStringToNullString(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  sql.NullString
	}{
		{name: "non-empty string", input: "non-empty", want: sql.NullString{String: "non-empty", Valid: true}},
		{name: "empty string", input: "", want: sql.NullString{String: "", Valid: false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, StringToNullString(tt.input))
		})
	}
}

func TestGenderStringToUint8(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  uint8
		err   error
	}{
		{name: "female gender", input: "female", want: 0, err: nil},
		{name: "male gender", input: "male", want: 1, err: nil},
		{name: "unknown gender", input: "unknown", want: 2, err: fmt.Errorf("unsupported gender type: unknown")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := GenderStringToUint8(tt.input)
			assert.Equal(t, tt.want, res)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestGenderUint8ToString(t *testing.T) {
	var tests = []struct {
		name  string
		input uint8
		want  string
	}{
		{name: "female gender", input: 0, want: "female"},
		{name: "male gender", input: 1, want: "male"},
		{name: "unknown gender", input: 2, want: "unknown"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GenderUint8ToString(tt.input))
		})
	}
}

func TestStringToStringPointer(t *testing.T) {
	var tests = []struct {
		name  string
		input string
		want  *string
	}{
		{name: "empty string", input: "", want: nil},
		{name: "non-empty string", input: "non-empty", want: ptr("non-empty")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, StringToStringPointer(tt.input))
		})
	}
}

func ptr(s string) *string { return &s }

/*
func Test(t *testing.T) {
	var tests = []struct {
		name string
		input
		want
	}{
		{name: , input: , want: }
		{name: , input: , want: }
	}
	for _, tt := range tests {
		t.Run(tt.name, func (t *testing.T) {
			assert.Equal(t, tt.want, )
		})
	}
}
*/
