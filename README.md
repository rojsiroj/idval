# idval

Parse and validate Indonesian data in Go and TypeScript:

- **NIK**: region, birth date, gender and serial number.
- **NPWP**: the 15-digit legacy form and the 16-digit form, and **NITKU** (22 digits).
- **Phone numbers**: convert to E.164, detect mobile or landline, and find the mobile operator.

The Go module and the npm package run the same test vectors in [`spec/`](spec). A Go backend and a TypeScript frontend give the same result for the same input.

Both packages have zero runtime dependencies and make no network calls.

> **Format validation only.** idval does not query Dukcapil or DJP records. A valid result means that the value has a correct structure. It does not mean that the person or taxpayer exists.

## Go

```bash
go get github.com/rojsiroj/idval
```

```go
import (
	"errors"
	"fmt"
	"time"

	"github.com/rojsiroj/idval"
	"github.com/rojsiroj/idval/nik"
	"github.com/rojsiroj/idval/npwp"
	"github.com/rojsiroj/idval/phone"
)

n, err := nik.Parse("3174054512050003")
if errors.Is(err, idval.ErrInvalidBirthDate) {
	// Show your own message, or use err.(*idval.Error).Message("id").
}
fmt.Println(n.ProvinceName, n.BirthDate.Format(time.DateOnly), n.Gender) // DKI Jakarta 2005-12-05 female

p, _ := npwp.Parse("012345678901000")
fmt.Println(p.Formatted, p.NPWP16) // 01.234.567.8-901.000 0012345678901000

ph, _ := phone.Parse("+62 857-1234-5678")
fmt.Println(ph.E164, ph.Operator) // +6285712345678 indosat
```

## TypeScript

```bash
npm install idval
```

```ts
import { parseNIK, parseNPWP, parsePhone, messages } from "idval";
// Or import one module only: import { parseNIK } from "idval/nik";

const r = parseNIK("3174054512050003");
if (r.ok) {
  r.value.birthDate; // "2005-12-05"
} else {
  r.error.code; // for example "INVALID_BIRTH_DATE"
  messages.id[r.error.code]; // Indonesian message
}
```

## Rules

### NIK

| Digits | Meaning |
|---|---|
| 1–2 | Province code (Kemendagri). Must be a known province. |
| 3–4 | Regency or city code. Must not be `00`. |
| 5–6 | District code. Must not be `00`. |
| 7–12 | Birth date as `DDMMYY`. For women, Dukcapil adds 40 to the day. |
| 13–16 | Serial number. Must not be `0000`. |

The birth year has 2 digits. idval puts the year in the 2000s. If that date is after the reference date, idval puts the year in the 1900s. The reference date is the current date by default. Use `nik.ParseAt` in Go or the `now` option in TypeScript to set a different date.

### NPWP and NITKU

- 15 digits: the legacy form. idval returns the 16-digit form (`"0"` + 15 digits) and the formatted value `XX.XXX.XXX.X-XXX.XXX`.
- 16 digits that start with `0`: the legacy form.
- 16 digits that do not start with `0`: the NIK of a resident individual. idval validates it as a NIK.
- NITKU: 22 digits. The first 16 digits are a 16-digit NPWP. The last 6 digits identify the business location. `000000` is the head office.

idval ignores the legacy check digit, because DJP does not publish the algorithm.

### Phone numbers

- The number must start with `0`, `62` or `+62`. idval also accepts `+62 0…`.
- Mobile numbers start with `8` after the trunk `0` and have 9–12 digits after it. The prefix must belong to a known operator.
- Landline numbers start with `2`–`7` or `9` after the trunk `0` and have 8–11 digits after it.

## Error codes

| Code | Cause |
|---|---|
| `INVALID_CHARACTERS` | The input has a character that is not a digit or a permitted separator. |
| `INVALID_LENGTH` | The input has an incorrect number of digits. |
| `UNKNOWN_REGION` | The province code is not known, or the regency or district code is `00`. |
| `INVALID_BIRTH_DATE` | The birth date in the NIK is not a real date. |
| `INVALID_SERIAL` | The NIK serial number is `0000`. |
| `INVALID_PREFIX` | The phone number does not start with `0`, `62` or `+62`, or the next digit is not valid. |
| `UNKNOWN_OPERATOR` | The mobile prefix does not belong to a known operator. |

Both packages include English and Indonesian messages for each code.

## Repository layout

```
spec/       Shared test vectors. Go and TypeScript tests read these files.
data/       Source data, for example province codes.
scripts/    Code generators. They write the tables for Go and TypeScript from data/.
nik/ npwp/ phone/   Go packages.
ts/         npm package.
```

## Contributing

1. Add or change a test vector in `spec/` first.
2. Make the Go tests pass: `go test ./...`
3. Make the TypeScript tests pass: `cd ts && npm test`
4. If you change `data/`, run `node scripts/gen-regions.mjs` and commit the generated files.

CI fails if the Go and TypeScript results are different, or if the generated files are not current.

## Roadmap

- Regency and district codes, so that idval validates the full region of a NIK.
- NIK generator for test data.
- Postal codes and vehicle plates (TNKB).
- Zod and Valibot adapters for TypeScript.

## License

MIT
