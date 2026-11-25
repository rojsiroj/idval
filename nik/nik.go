// Package nik parses and validates the Indonesian population identity
// number (Nomor Induk Kependudukan, NIK).
//
// The package checks the format only. It does not check Dukcapil records.
package nik

import (
	"encoding/json"
	"time"

	"github.com/rojsiroj/idval"
	"github.com/rojsiroj/idval/internal/normalize"
	"github.com/rojsiroj/idval/internal/region"
)

// Gender is the gender that a NIK encodes.
type Gender string

// Gender values.
const (
	Male   Gender = "male"
	Female Gender = "female"
)

// NIK is a parsed NIK.
type NIK struct {
	NIK          string    // 16 digits, without separators.
	ProvinceCode string    // 2 digits.
	ProvinceName string    // Kemendagri province name.
	RegencyCode  string    // 4 digits (province + regency).
	DistrictCode string    // 6 digits (province + regency + district).
	BirthDate    time.Time // Midnight UTC.
	Gender       Gender
	Serial       string // 4 digits.
}

// MarshalJSON writes BirthDate as "YYYY-MM-DD". The field names match the
// TypeScript package.
func (n NIK) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NIK          string `json:"nik"`
		ProvinceCode string `json:"provinceCode"`
		ProvinceName string `json:"provinceName"`
		RegencyCode  string `json:"regencyCode"`
		DistrictCode string `json:"districtCode"`
		BirthDate    string `json:"birthDate"`
		Gender       Gender `json:"gender"`
		Serial       string `json:"serial"`
	}{n.NIK, n.ProvinceCode, n.ProvinceName, n.RegencyCode, n.DistrictCode,
		n.BirthDate.Format(time.DateOnly), n.Gender, n.Serial})
}

// Parse parses s with the current time as the reference date.
// See ParseAt.
func Parse(s string) (NIK, error) {
	return ParseAt(s, time.Now())
}

// ParseAt parses s. Spaces, dots and dashes in s are ignored.
//
// A NIK stores the birth year as 2 digits. ParseAt puts the year in the
// 2000s. If that date is after the date of now (in the location of now),
// ParseAt puts the year in the 1900s.
//
// The error is an *idval.Error.
func ParseAt(s string, now time.Time) (NIK, error) {
	d, ok := normalize.Digits(s)
	if !ok {
		return NIK{}, idval.ErrInvalidCharacters
	}
	if len(d) != 16 {
		return NIK{}, idval.ErrInvalidLength
	}

	prov := d[0:2]
	name, ok := region.Provinces[prov]
	if !ok || d[2:4] == "00" || d[4:6] == "00" {
		return NIK{}, idval.ErrUnknownRegion
	}

	day, month, yy := atoi(d[6:8]), atoi(d[8:10]), atoi(d[10:12])
	gender := Male
	if day > 40 {
		gender = Female
		day -= 40
	}
	if day < 1 || day > 31 || month < 1 || month > 12 {
		return NIK{}, idval.ErrInvalidBirthDate
	}

	year := 2000 + yy
	ny, nm, nd := now.Date()
	if year > ny || (year == ny && (month > int(nm) || (month == int(nm) && day > nd))) {
		year -= 100
	}
	birth := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if birth.Day() != day {
		// time.Date normalized an impossible date, for example 31 February.
		return NIK{}, idval.ErrInvalidBirthDate
	}

	if d[12:16] == "0000" {
		return NIK{}, idval.ErrInvalidSerial
	}

	return NIK{
		NIK:          d,
		ProvinceCode: prov,
		ProvinceName: name,
		RegencyCode:  d[0:4],
		DistrictCode: d[0:6],
		BirthDate:    birth,
		Gender:       gender,
		Serial:       d[12:16],
	}, nil
}

// Valid reports whether s is a valid NIK, with the current time as the
// reference date.
func Valid(s string) bool {
	_, err := Parse(s)
	return err == nil
}

// atoi converts a string of ASCII digits. The caller must check the digits.
func atoi(s string) int {
	n := 0
	for i := 0; i < len(s); i++ {
		n = n*10 + int(s[i]-'0')
	}
	return n
}
