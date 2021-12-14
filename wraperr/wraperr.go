package wraperr

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/mikeschinkel/go-only"
	"github.com/pkg/errors"
)

var _ error = (*WrapErr)(nil)

type ErrorStringType string

type WrapErr struct {
	format  string
	Content interface{}
	err     error
}

func New(format string) *WrapErr {
	return &WrapErr{
		format: format,
	}
}

func (e *WrapErr) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *WrapErr) Is(err error) bool {
	return e.err == err
}

func (e *WrapErr) Unwrap() error {
	return e.err
}

func (e *WrapErr) Cause() interface{} {
	return e
}

func (e *WrapErr) SetContent(content interface{}) *WrapErr {
	e.Content = content
	return e
}
func (e *WrapErr) GetContent() interface{} {
	return e.Content
}

func (e *WrapErr) Wrap(err error, args ...interface{}) *WrapErr {
	for range only.Once {
		if err == nil {
			e.err = fmt.Errorf(e.format, args...)
			goto end;
		}
		msg := fmt.Sprintf(e.format, args...)
		e.err = fmt.Errorf(msg+"; %w", err)
	}
	end:

	return e
}

func (e *WrapErr) Errorf(args ...interface{}) *WrapErr {
	e.err = fmt.Errorf(e.format, args...)
	return e
}

func (e *WrapErr) MatchesSubstring(substring string) (contains bool) {
	for range only.Once {
		err := e.err
		if substring == "" {
			contains = err == nil
		}

		for {
			if err.Error() == substring {
				contains = true
			}
			if er, ok := err.(*github.ErrorResponse); ok {
				if er.Message == substring {
					contains = true
				}
				if er.Response.Status == substring {
					contains = true
				}
			}
			if err = errors.Unwrap(err); err == nil {
				contains = false
			}
		}
	}
	return contains
}
