package testclient

import (
	"fmt"
	"github.com/mikeschinkel/go-only"
	"os"
	"strings"
	"testing"
)

type Fixture struct {
	T        *testing.T
	Filename string
}

func NewFixture(t *testing.T, fn string) *Fixture {
	return &Fixture{
		T:        t,
		Filename: fn,
	}
}

// Filepath gets the filepath for a fixture from its base filename and the current working directory.
func (f *Fixture) Filepath() string {
	var fp string
	for range only.Once {
		wd, err := os.Getwd()
		if err != nil {
			f.T.Errorf("%s:%s",
				fmt.Sprintf("unable to get working directory for filename='%s'", f.Filename),
				err)
			break
		}
		fp = fmt.Sprintf("%s%cfixtures%c%s",
			wd,
			os.PathSeparator,
			os.PathSeparator,
			strings.TrimLeft(f.Filename, string(os.PathSeparator)))
	}
	return fp
}
