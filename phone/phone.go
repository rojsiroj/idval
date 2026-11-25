// Package phone parses and validates Indonesian phone numbers.
package phone

import (
	"strings"

	"github.com/rojsiroj/idval"
)

// Kind is the type of phone number.
type Kind string

// Kind values.
const (
	Mobile   Kind = "mobile"
	Landline Kind = "landline"
)

// Operator is the mobile brand that owns a prefix. Indonesia has no mobile
// number portability, so the prefix identifies the brand. XL, Axis and
// Smartfren are brands of XLSmart. Indosat and Tri are brands of Indosat
// Ooredoo Hutchison.
type Operator string

// Operator values.
const (
	Telkomsel Operator = "telkomsel"
	Indosat   Operator = "indosat"
	XL        Operator = "xl"
	Axis      Operator = "axis"
	Tri       Operator = "tri"
	Smartfren Operator = "smartfren"
)

// prefixes maps the first 3 digits of a mobile number (after the trunk 0)
// to the operator.
var prefixes = map[string]Operator{
	"811": Telkomsel, "812": Telkomsel, "813": Telkomsel,
	"821": Telkomsel, "822": Telkomsel, "823": Telkomsel,
	"851": Telkomsel, "852": Telkomsel, "853": Telkomsel,
	"814": Indosat, "815": Indosat, "816": Indosat,
	"855": Indosat, "856": Indosat, "857": Indosat, "858": Indosat,
	"817": XL, "818": XL, "819": XL, "859": XL, "877": XL, "878": XL,
	"831": Axis, "832": Axis, "833": Axis, "838": Axis,
	"895": Tri, "896": Tri, "897": Tri, "898": Tri, "899": Tri,
	"881": Smartfren, "882": Smartfren, "883": Smartfren, "884": Smartfren,
	"885": Smartfren, "886": Smartfren, "887": Smartfren, "888": Smartfren,
	"889": Smartfren,
}

// Number is a parsed phone number.
type Number struct {
	E164     string   `json:"e164"`               // For example "+6281234567890".
	National string   `json:"national"`           // For example "081234567890".
	Kind     Kind     `json:"kind"`               //
	Operator Operator `json:"operator,omitempty"` // Empty for a landline.
}

// Parse parses s. Spaces, dots, dashes and parentheses in s are ignored.
// The number must start with "0", "62" or "+62". Parse accepts a trunk "0"
// after the country code, for example "+62 0812...".
//
// The error is an *idval.Error.
func Parse(s string) (Number, error) {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == ' ' || r == '.' || r == '-' || r == '(' || r == ')':
			continue
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '+' && b.Len() == 0:
			// A plus sign is permitted only before the first digit.
			b.WriteRune(r)
		default:
			return Number{}, idval.ErrInvalidCharacters
		}
	}
	d := b.String()

	var nsn string
	switch {
	case strings.HasPrefix(d, "+62"):
		nsn = d[3:]
	case strings.HasPrefix(d, "62"):
		nsn = d[2:]
	case strings.HasPrefix(d, "0"):
		nsn = d[1:]
	default:
		return Number{}, idval.ErrInvalidPrefix
	}
	if d[0] != '0' && strings.HasPrefix(nsn, "0") {
		nsn = nsn[1:]
	}
	if nsn == "" {
		return Number{}, idval.ErrInvalidLength
	}

	n := Number{E164: "+62" + nsn, National: "0" + nsn}
	switch c := nsn[0]; {
	case c == '8':
		if len(nsn) < 9 || len(nsn) > 12 {
			return Number{}, idval.ErrInvalidLength
		}
		op, ok := prefixes[nsn[:3]]
		if !ok {
			return Number{}, idval.ErrUnknownOperator
		}
		n.Kind, n.Operator = Mobile, op
	case c >= '2' && c <= '9':
		if len(nsn) < 8 || len(nsn) > 11 {
			return Number{}, idval.ErrInvalidLength
		}
		n.Kind = Landline
	default:
		return Number{}, idval.ErrInvalidPrefix
	}
	return n, nil
}

// Valid reports whether s is a valid Indonesian phone number.
func Valid(s string) bool {
	_, err := Parse(s)
	return err == nil
}
