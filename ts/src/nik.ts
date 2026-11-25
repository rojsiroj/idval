/**
 * Parses and validates the Indonesian population identity number (Nomor Induk Kependudukan, NIK).
 * The module checks the format only. It does not check Dukcapil records.
 * @module
 */
import { fail, ok, type Result } from "./errors.js";
import { digits } from "./internal/normalize.js";
import { provinces } from "./internal/provinces.gen.js";

export type Gender = "male" | "female";

export interface NIK {
  /** 16 digits, without separators. */
  nik: string;
  /** 2 digits. */
  provinceCode: string;
  /** Kemendagri province name. */
  provinceName: string;
  /** 4 digits (province + regency). */
  regencyCode: string;
  /** 6 digits (province + regency + district). */
  districtCode: string;
  /** "YYYY-MM-DD". */
  birthDate: string;
  gender: Gender;
  /** 4 digits. */
  serial: string;
}

export interface ParseOptions {
  /**
   * The reference date for the 2-digit birth year. The default is the current date.
   * The function uses the local calendar date of this value.
   */
  now?: Date;
}

/**
 * Parses a NIK. Spaces, dots and dashes in `input` are ignored.
 *
 * A NIK stores the birth year as 2 digits. The function puts the year in the 2000s.
 * If that date is after the date of `now`, the function puts the year in the 1900s.
 */
export function parseNIK(input: string, options: ParseOptions = {}): Result<NIK> {
  const d = digits(input);
  if (d === undefined) return fail("INVALID_CHARACTERS");
  if (d.length !== 16) return fail("INVALID_LENGTH");

  const provinceCode = d.slice(0, 2);
  const provinceName = provinces[provinceCode];
  if (provinceName === undefined || d.slice(2, 4) === "00" || d.slice(4, 6) === "00") {
    return fail("UNKNOWN_REGION");
  }

  let day = Number(d.slice(6, 8));
  const month = Number(d.slice(8, 10));
  const yy = Number(d.slice(10, 12));
  let gender: Gender = "male";
  if (day > 40) {
    gender = "female";
    day -= 40;
  }
  if (day < 1 || day > 31 || month < 1 || month > 12) return fail("INVALID_BIRTH_DATE");

  const now = options.now ?? new Date();
  const ny = now.getFullYear();
  const nm = now.getMonth() + 1;
  const nd = now.getDate();
  let year = 2000 + yy;
  if (year > ny || (year === ny && (month > nm || (month === nm && day > nd)))) year -= 100;
  if (day > daysInMonth(year, month)) return fail("INVALID_BIRTH_DATE");

  const serial = d.slice(12, 16);
  if (serial === "0000") return fail("INVALID_SERIAL");

  return ok({
    nik: d,
    provinceCode,
    provinceName,
    regencyCode: d.slice(0, 4),
    districtCode: d.slice(0, 6),
    birthDate: `${year}-${pad(month)}-${pad(day)}`,
    gender,
    serial,
  });
}

/** Returns true if `input` is a valid NIK, with the current date as the reference date. */
export function isValidNIK(input: string): boolean {
  return parseNIK(input).ok;
}

function daysInMonth(year: number, month: number): number {
  if (month === 2) {
    const leap = (year % 4 === 0 && year % 100 !== 0) || year % 400 === 0;
    return leap ? 29 : 28;
  }
  return [4, 6, 9, 11].includes(month) ? 30 : 31;
}

function pad(n: number): string {
  return String(n).padStart(2, "0");
}
