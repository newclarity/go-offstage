package util

import (
	"fmt"
	"github.com/mikeschinkel/go-only"
	"reflect"
	"regexp"
)

type NamedVarsRegexp struct {
	*regexp.Regexp
}

func NewNamedVarsRegexp(s string) *NamedVarsRegexp {
	return &NamedVarsRegexp{
		Regexp: regexp.MustCompile(s),
	}
}

type StringStringMap = map[string]string

var nameCaptureRegex = regexp.MustCompile("\\?P<([^>]+)>")

func (cr *NamedVarsRegexp) StructCapture(obj interface{}, s string) (err error) {
	var fieldname string
	var value reflect.Value
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case *reflect.ValueError:
				err = fmt.Errorf("unable to capture property '%s' into struct '%s' via regexp '%s'; %w",
					fieldname,
					value.Type().String(),
					cr.String(),
					r.(*reflect.ValueError))
			default:
				err = fmt.Errorf("unable to capture struct '%s'; %#v", s, r)
			}
			return
		}
	}()
	for range only.Once {
		props := nameCaptureRegex.FindAllStringSubmatch(cr.String(), -1)
		m := cr.MapCapture(s)
		if len(m) != len(props) {
			err = ErrorMatchingAndCapturing.Errorf(s, cr.String())
			break
		}
		value = reflect.ValueOf(obj).Elem()
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		for _, prop := range props {
			fieldname = prop[1]
			f := value.FieldByName(fieldname)
			f.SetString(m[fieldname])
		}
	}
	return err
}
func (cr *NamedVarsRegexp) MapCapture(str string) StringStringMap {
	match := cr.FindStringSubmatch(str)
	m := make(StringStringMap, 0)
	for i, name := range cr.SubexpNames() {
		if i == 0 {
			continue
		}
		m[name] = ""
		if match == nil {
			break
		}
		if i >= len(match) {
			continue
		}
		m[name] = match[i]
	}
	return m
}
