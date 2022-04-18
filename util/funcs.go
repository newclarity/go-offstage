package util

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"
)
import (
	"bytes"
	"fmt"
	"github.com/mikeschinkel/go-only"
	"github.com/newclarity/go-offstage/wraperr"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

// ApplyTemplate compile template and write it to output argument
func ApplyTemplate(filename string, data interface{}) (text string, err error) {
	for range only.Once {

		t := template.New(filepath.Base(filename))
		tmpl, err2 := t.ParseFiles(filename)
		if err2 != nil {
			err = err2
			break
		}

		var output bytes.Buffer
		err = tmpl.Execute(&output, data)
		if err != nil {
			break
		}

		text = output.String()

	}
	if err != nil {
		err = ErrorApplyingTemplate.Wrap(err,
			filename)
	}
	return text, err
}

var ErrorAccessingCurrentDirectory = wraperr.New("failed to access current directory")

func GetCurrentDir() (dir string, err error) {
	dir, err = os.Getwd()
	if err != nil {
		err = ErrorAccessingCurrentDirectory.Wrap(err)
	}
	return dir, err
}

func Fail(msg string) {
	_, _ = fmt.Fprintf(os.Stderr, msg)
	os.Exit(1)

}

func ReadResponseBody(r *http.Response) string {
	var body []byte
	body, _ = ioutil.ReadAll(r.Body)
	_ = r.Body.Close()
	var v interface{}
	_ = json.Unmarshal(body, &v)
	body, _ = json.MarshalIndent(v, "", "   ")
	return string(body)
}

// flattenString returns as string cleaned of newlines and tabs.
// Replaces line feeds, carriage returns and tabs with a "\n", "\r" and "\t," respectively.
// TODO Maybe rename with a more canonical name than this?
func flattenString(s string) string {
	var m = map[string]string{
		"\n": "\\n",
		"\r": "\\r",
		"\t": "\\t",
	}
	for f, r := range m {
		s = strings.ReplaceAll(s, f, r)
	}
	return s

}

func collapseWhitespace(s string) string {
	spaces := regexp.MustCompile(`\s+`)
	return spaces.ReplaceAllString(strings.TrimSpace(s), " ")
}

var indent = strings.Repeat(" ", 4)

func toIndentedJson(i interface{}) string {
	return _toJson(i, func(i interface{}) ([]byte, error) {
		return json.MarshalIndent(i, indent, "   ")
	})
}

func toJson(i interface{}) string {
	return _toJson(i, func(i interface{}) ([]byte, error) {
		return json.Marshal(i)
	})
}

type marshaler func(interface{}) ([]byte, error)

func _toJson(i interface{}, fn marshaler) string {
	j, err := fn(i)
	if err != nil {
		println(err.Error())
	}
	return string(j)
}

// filenameWithoutExtension strips path and extension and return bare filename as a string
// e.g. filenameWithoutExtension("/foo/bar/baz.ext") => "bar"
func filenameWithoutExtension(fp string) string {
	fp = filepath.Base(fp)
	return fp[:len(fp)-len(filepath.Ext(fp))]
}
