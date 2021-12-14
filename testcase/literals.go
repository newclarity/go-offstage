package testcase

import (
	"fmt"
	"github.com/mikeschinkel/go-only"
	"github.com/onsi/ginkgo"
	"strings"
)

var literals map[string]string
var used map[string]bool

func init() {
	literals = make(map[string]string, 0)
	used = make(map[string]bool)
	Literal("", "Empty String")
}

// Literal is used to register a string literal for later use
// in a Gingko test wrapped with the L() function.
func Literal(literal, description string) {
	literals[literal] = description
}

// recordUsage records that a literal was used
func recordUsage(s string) {
	used[s] = true
}

// L checks to see if the string passed has been registered as a string literal.
func L(s string, args ...interface{}) (_s string) {
	var ok bool
	for range only.Once {
		_s = s
		_, ok = literals[s]
		if !ok {
			break
		}
		recordUsage(s)
		if len(args) == 0 {
			break
		}
		for _, arg := range args {
			__s, _ok := arg.(string)
			if !_ok {
				// TODO Handle any data type
				continue
			}
			L(__s)
		}
		if !ok {
			break
		}
		_s = fmt.Sprintf(s, args...)
		_, ok = literals[_s]
		if ok {
			recordUsage(_s)
		}
	}
	if !ok {
		ginkgo.Fail(fmt.Sprintf(`Unregistered string literal: "%s"`,
			_s))
	}
	return _s
}

// CheckLiteralUsage should be run in Gingko's AfterSuite() to ensure all
// literals registered with Literal() were actually referenced with L().
func CheckLiteralUsage() {
	for range only.Once {
		unused := make([]string, 0)
		for l := range literals {
			if _, ok := used[l]; ok {
				continue
			}
			unused = append(unused, l)
		}
		if len(unused) == 0 {
			break
		}
		ginkgo.Fail(fmt.Sprintf(`Registered literals never used: "%s"`,
			strings.Join(unused, "', '")))
	}
}
