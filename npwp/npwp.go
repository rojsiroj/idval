// Package npwp parses and validates the Indonesian taxpayer number
// (Nomor Pokok Wajib Pajak, NPWP) and the business-location number
// (Nomor Identitas Tempat Kegiatan Usaha, NITKU).
//
// The package checks the format only. It does not check DJP records, and it
// does not check the legacy check digit, because DJP does not publish the
// algorithm.
package npwp

import (
	"time"

	"github.com/rojsiroj/idval"
	"github.com/rojsiroj/idval/internal/normalize"
	"github.com/rojsiroj/idval/nik"
)

// Kind is the form of a 16-digit NPWP.
type Kind string

// Kind values.
const (
	// Legacy is "0" followed by the 15-digit NPWP.
	Legacy Kind = "legacy"
	// NIK is the NIK of a resident individual.
	NIK Kind = "nik"
)

// NPWP is a parsed NPWP.
type NPWP struct {
	NPWP16    string `json:"npwp16"`           // 16 digits.
	NPWP15    string `json:"npwp15,omitempty"` // 15 digits. Empty when Kind is NIK.
	Kind      Kind   `json:"kind"`
	Formatted string `json:"formatted"` // XX.XXX.XXX.X-XXX.XXX for Legacy, 16 digits for NIK.
}

// NITKU is a parsed NITKU.
type NITKU struct {
	NITKU  string `json:"nitku"`  // 22 digits.
	NPWP16 string `json:"npwp16"` // The first 16 digits.
	Suffix string `json:"suffix"` // The last 6 digits. "000000" is the head office.
}

// Parse parses s with the current time as the reference date. See ParseAt.
func Parse(s string) (NPWP, error) {
	return ParseAt(s, time.Now())
}

// ParseAt parses a 15-digit or 16-digit NPWP. Spaces, dots and dashes in s
// are ignored. A 16-digit NPWP that does not start with "0" must be a valid
// NIK. ParseAt uses now as the reference date for that NIK (see nik.ParseAt).
//
// The error is an *idval.Error.
func ParseAt(s string, now time.Time) (NPWP, error) {
	d, ok := normalize.Digits(s)
	if !ok {
		return NPWP{}, idval.ErrInvalidCharacters
	}
	return parseDigits(d, now)
}

func parseDigits(d string, now time.Time) (NPWP, error) {
	switch {
	case len(d) == 15:
		return legacy(d), nil
	case len(d) == 16 && d[0] == '0':
		return legacy(d[1:]), nil
	case len(d) == 16:
		if _, err := nik.ParseAt(d, now); err != nil {
			return NPWP{}, err
		}
		return NPWP{NPWP16: d, Kind: NIK, Formatted: d}, nil
	default:
		return NPWP{}, idval.ErrInvalidLength
	}
}

func legacy(d15 string) NPWP {
	f := d15[0:2] + "." + d15[2:5] + "." + d15[5:8] + "." + d15[8:9] + "-" + d15[9:12] + "." + d15[12:15]
	return NPWP{NPWP16: "0" + d15, NPWP15: d15, Kind: Legacy, Formatted: f}
}

// ParseNITKU parses s with the current time as the reference date.
// See ParseNITKUAt.
func ParseNITKU(s string) (NITKU, error) {
	return ParseNITKUAt(s, time.Now())
}

// ParseNITKUAt parses a 22-digit NITKU. Spaces, dots and dashes in s are
// ignored. The first 16 digits must be a valid 16-digit NPWP.
//
// The error is an *idval.Error.
func ParseNITKUAt(s string, now time.Time) (NITKU, error) {
	d, ok := normalize.Digits(s)
	if !ok {
		return NITKU{}, idval.ErrInvalidCharacters
	}
	if len(d) != 22 {
		return NITKU{}, idval.ErrInvalidLength
	}
	if _, err := parseDigits(d[:16], now); err != nil {
		return NITKU{}, err
	}
	return NITKU{NITKU: d, NPWP16: d[:16], Suffix: d[16:]}, nil
}

// Valid reports whether s is a valid 15-digit or 16-digit NPWP.
func Valid(s string) bool {
	_, err := Parse(s)
	return err == nil
}
