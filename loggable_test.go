package slogerr_test

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/jakofys/slogerr"
	"github.com/stretchr/testify/assert"
)

type (
	errL2 struct {
		err error
	}
	errL1 struct{}
)

func (errL1) Error() string   { return "errl1" }
func (errL2) Error() string   { return "errl2" }
func (e errL2) Unwrap() error { return e.err }
func (errL1) LogAttr() []slog.Attr {
	return []slog.Attr{
		slog.String("key1", "value1"),
	}
}

func (errL2) LogAttr() []slog.Attr {
	return []slog.Attr{
		slog.String("key2", "value2"),
	}
}

func TestLogAttr(t *testing.T) {
	cases := []struct {
		name, msg string
		expected  []slog.Attr
		input     error
	}{
		{
			name:  "direct implementation",
			msg:   "should return key1 string attribute",
			input: errL1{},
			expected: []slog.Attr{
				slog.String("key1", "value1"),
			},
		},
		{
			name:     "no implementation",
			msg:      "should return nil",
			input:    errors.New("no implementation"),
			expected: nil,
		},
		{
			name:  "indirect implementation",
			msg:   "should return key1 string attribute",
			input: fmt.Errorf("errl1: %w", errL1{}),
			expected: []slog.Attr{
				slog.String("key1", "value1"),
			},
		},
		{
			name:  "aggregate implementation",
			msg:   "should return key1 and key2 string attribute",
			input: errL2{err: errL1{}},
			expected: []slog.Attr{
				slog.String("key1", "value1"),
				slog.String("key2", "value2"),
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := slogerr.AttrFromError(c.input)
			assert.ElementsMatch(t, result, c.expected)
		})
	}
}
