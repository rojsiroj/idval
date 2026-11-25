package npwp_test

import (
	"fmt"
	"testing"

	"github.com/rojsiroj/idval/internal/spectest"
	"github.com/rojsiroj/idval/npwp"
)

func TestSpec(t *testing.T) {
	var spec struct {
		Now   string          `json:"now"`
		NPWP  []spectest.Case `json:"npwp"`
		NITKU []spectest.Case `json:"nitku"`
	}
	spectest.Load(t, "npwp.json", &spec)
	now := spectest.Date(t, spec.Now)
	for _, c := range spec.NPWP {
		t.Run("npwp/"+c.Name, func(t *testing.T) {
			got, err := npwp.ParseAt(c.Input, now)
			spectest.Check(t, c, got, err)
		})
	}
	for _, c := range spec.NITKU {
		t.Run("nitku/"+c.Name, func(t *testing.T) {
			got, err := npwp.ParseNITKUAt(c.Input, now)
			spectest.Check(t, c, got, err)
		})
	}
}

func ExampleParse() {
	n, err := npwp.Parse("012345678901000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n.Formatted, n.NPWP16)
	// Output: 01.234.567.8-901.000 0012345678901000
}
