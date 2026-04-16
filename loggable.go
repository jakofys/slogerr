package slogerr

import (
	"errors"
	"log/slog"

	"github.com/jakofys/xerrors"
)

type Loggable interface {
	LogAttr() []slog.Attr
}

func AttrFromError(err error) []slog.Attr {
	if err == nil {
		return nil
	}
	loggable := xerrors.AsInterface[Loggable](err)
	if loggable == nil {
		return nil
	}
	attrs := loggable.LogAttr()
	if childAttr := AttrFromError(errors.Unwrap(loggable.(error))); len(childAttr) > 0 {
		attrs = append(attrs, childAttr...)
	}
	return attrs
}
