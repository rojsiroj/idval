package phone_test

import (
	"fmt"
	"testing"

	"github.com/rojsiroj/idval/internal/spectest"
	"github.com/rojsiroj/idval/phone"
)

func TestSpec(t *testing.T) {
	var spec struct {
		Cases []spectest.Case `json:"cases"`
	}
	spectest.Load(t, "phone.json", &spec)
	for _, c := range spec.Cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := phone.Parse(c.Input)
			spectest.Check(t, c, got, err)
		})
	}
}

func ExampleParse() {
	n, err := phone.Parse("+62 857-1234-5678")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n.E164, n.National, n.Operator)
	// Output: +6285712345678 085712345678 indosat
}
