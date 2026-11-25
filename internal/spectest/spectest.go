// Package spectest loads the shared test vectors in spec/ and compares
// results with them. Only tests use this package.
package spectest

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/rojsiroj/idval"
)

// Case is one test vector. A case has a Result or an Error, not both.
type Case struct {
	Name   string         `json:"name"`
	Input  string         `json:"input"`
	Result map[string]any `json:"result"`
	Error  idval.Code     `json:"error"`
}

// Load reads spec/<file> into v.
func Load(t *testing.T, file string, v any) {
	t.Helper()
	_, here, _, _ := runtime.Caller(0)
	path := filepath.Join(filepath.Dir(here), "..", "..", "spec", file)
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
}

// Date parses a "YYYY-MM-DD" reference date from a spec file.
func Date(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse(time.DateOnly, s)
	if err != nil {
		t.Fatalf("parse date %q: %v", s, err)
	}
	return d
}

// Check compares got and err with the expected values in c.
func Check(t *testing.T, c Case, got any, err error) {
	t.Helper()
	if c.Error != "" {
		if !errors.Is(err, &idval.Error{Code: c.Error}) {
			t.Fatalf("input %q: got error %v, want code %s", c.Input, err, c.Error)
		}
		return
	}
	if err != nil {
		t.Fatalf("input %q: got error %v, want a result", c.Input, err)
	}
	b, merr := json.Marshal(got)
	if merr != nil {
		t.Fatalf("marshal result: %v", merr)
	}
	var m map[string]any
	if uerr := json.Unmarshal(b, &m); uerr != nil {
		t.Fatalf("unmarshal result: %v", uerr)
	}
	if !reflect.DeepEqual(m, c.Result) {
		t.Fatalf("input %q:\n got  %v\n want %v", c.Input, m, c.Result)
	}
}
