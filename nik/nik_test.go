package nik_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rojsiroj/idval"
	"github.com/rojsiroj/idval/internal/spectest"
	"github.com/rojsiroj/idval/nik"
)

func TestSpec(t *testing.T) {
	var spec struct {
		Now   string          `json:"now"`
		Cases []spectest.Case `json:"cases"`
	}
	spectest.Load(t, "nik.json", &spec)
	now := spectest.Date(t, spec.Now)
	for _, c := range spec.Cases {
		t.Run(c.Name, func(t *testing.T) {
			got, err := nik.ParseAt(c.Input, now)
			spectest.Check(t, c, got, err)
		})
	}
}

func TestErrorsIs(t *testing.T) {
	_, err := nik.Parse("123")
	if !errors.Is(err, idval.ErrInvalidLength) {
		t.Fatalf("got %v, want ErrInvalidLength", err)
	}
}

func ExampleParseAt() {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	n, err := nik.ParseAt("3174054512050003", now)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(n.ProvinceName, n.BirthDate.Format(time.DateOnly), n.Gender)
	// Output: DKI Jakarta 2005-12-05 female
}
