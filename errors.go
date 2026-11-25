// Package idval holds the error type and error codes that the idval
// validators share. The validators are in the nik, npwp and phone packages.
package idval

// Code identifies the cause of a validation error. The Go and TypeScript
// packages use the same code values.
type Code string

// Error codes.
const (
	CodeInvalidCharacters Code = "INVALID_CHARACTERS"
	CodeInvalidLength     Code = "INVALID_LENGTH"
	CodeUnknownRegion     Code = "UNKNOWN_REGION"
	CodeInvalidBirthDate  Code = "INVALID_BIRTH_DATE"
	CodeInvalidSerial     Code = "INVALID_SERIAL"
	CodeInvalidPrefix     Code = "INVALID_PREFIX"
	CodeUnknownOperator   Code = "UNKNOWN_OPERATOR"
)

// Sentinel errors. Use them with errors.Is:
//
//	if errors.Is(err, idval.ErrInvalidLength) { ... }
var (
	ErrInvalidCharacters = &Error{Code: CodeInvalidCharacters}
	ErrInvalidLength     = &Error{Code: CodeInvalidLength}
	ErrUnknownRegion     = &Error{Code: CodeUnknownRegion}
	ErrInvalidBirthDate  = &Error{Code: CodeInvalidBirthDate}
	ErrInvalidSerial     = &Error{Code: CodeInvalidSerial}
	ErrInvalidPrefix     = &Error{Code: CodeInvalidPrefix}
	ErrUnknownOperator   = &Error{Code: CodeUnknownOperator}
)

var messages = map[string]map[Code]string{
	"en": {
		CodeInvalidCharacters: "The value contains characters that are not permitted.",
		CodeInvalidLength:     "The value has an incorrect number of digits.",
		CodeUnknownRegion:     "The region code is not known.",
		CodeInvalidBirthDate:  "The birth date is not a valid date.",
		CodeInvalidSerial:     "The serial number must not be 0000.",
		CodeInvalidPrefix:     "The number does not start with 0, 62 or +62, or the area code is not valid.",
		CodeUnknownOperator:   "The mobile prefix does not belong to a known operator.",
	},
	"id": {
		CodeInvalidCharacters: "Nilai berisi karakter yang tidak diizinkan.",
		CodeInvalidLength:     "Jumlah digit tidak sesuai.",
		CodeUnknownRegion:     "Kode wilayah tidak dikenal.",
		CodeInvalidBirthDate:  "Tanggal lahir tidak valid.",
		CodeInvalidSerial:     "Nomor urut tidak boleh 0000.",
		CodeInvalidPrefix:     "Nomor tidak diawali 0, 62 atau +62, atau kode area tidak valid.",
		CodeUnknownOperator:   "Prefiks nomor seluler tidak dikenal.",
	},
}

// Error is a validation error. Code identifies the cause.
type Error struct {
	Code Code
}

// Error returns the English message.
func (e *Error) Error() string {
	return e.Message("en")
}

// Message returns the message in the given language ("en" or "id").
// For an unknown language, Message returns the English message.
func (e *Error) Message(lang string) string {
	m, ok := messages[lang]
	if !ok {
		m = messages["en"]
	}
	if s, ok := m[e.Code]; ok {
		return s
	}
	return string(e.Code)
}

// Is reports whether target is an *Error with the same code.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t.Code == e.Code
}
