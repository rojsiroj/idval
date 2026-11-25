/**
 * Parses and validates the Indonesian taxpayer number (Nomor Pokok Wajib Pajak, NPWP)
 * and the business-location number (Nomor Identitas Tempat Kegiatan Usaha, NITKU).
 *
 * The module checks the format only. It does not check DJP records, and it does not
 * check the legacy check digit, because DJP does not publish the algorithm.
 * @module
 */
import { fail, ok, type Result } from "./errors.js";
import { digits } from "./internal/normalize.js";
import { parseNIK, type ParseOptions } from "./nik.js";

export type { ParseOptions } from "./nik.js";

/** "legacy" is "0" followed by the 15-digit NPWP. "nik" is the NIK of a resident individual. */
export type NPWPKind = "legacy" | "nik";

export interface NPWP {
  /** 16 digits. */
  npwp16: string;
  /** 15 digits. Present only when `kind` is "legacy". */
  npwp15?: string;
  kind: NPWPKind;
  /** XX.XXX.XXX.X-XXX.XXX for "legacy", 16 digits for "nik". */
  formatted: string;
}

export interface NITKU {
  /** 22 digits. */
  nitku: string;
  /** The first 16 digits. */
  npwp16: string;
  /** The last 6 digits. "000000" is the head office. */
  suffix: string;
}

/**
 * Parses a 15-digit or 16-digit NPWP. Spaces, dots and dashes in `input` are ignored.
 * A 16-digit NPWP that does not start with "0" must be a valid NIK.
 * `options.now` is the reference date for that NIK (see `parseNIK`).
 */
export function parseNPWP(input: string, options: ParseOptions = {}): Result<NPWP> {
  const d = digits(input);
  if (d === undefined) return fail("INVALID_CHARACTERS");
  return parseDigits(d, options);
}

function parseDigits(d: string, options: ParseOptions): Result<NPWP> {
  if (d.length === 15) return ok(legacy(d));
  if (d.length !== 16) return fail("INVALID_LENGTH");
  if (d.startsWith("0")) return ok(legacy(d.slice(1)));
  const n = parseNIK(d, options);
  if (!n.ok) return n;
  return ok({ npwp16: d, kind: "nik", formatted: d });
}

function legacy(d15: string): NPWP {
  const formatted = `${d15.slice(0, 2)}.${d15.slice(2, 5)}.${d15.slice(5, 8)}.${d15.slice(8, 9)}-${d15.slice(9, 12)}.${d15.slice(12, 15)}`;
  return { npwp16: `0${d15}`, npwp15: d15, kind: "legacy", formatted };
}

/**
 * Parses a 22-digit NITKU. Spaces, dots and dashes in `input` are ignored.
 * The first 16 digits must be a valid 16-digit NPWP.
 */
export function parseNITKU(input: string, options: ParseOptions = {}): Result<NITKU> {
  const d = digits(input);
  if (d === undefined) return fail("INVALID_CHARACTERS");
  if (d.length !== 22) return fail("INVALID_LENGTH");
  const n = parseDigits(d.slice(0, 16), options);
  if (!n.ok) return n;
  return ok({ nitku: d, npwp16: d.slice(0, 16), suffix: d.slice(16) });
}

/** Returns true if `input` is a valid 15-digit or 16-digit NPWP. */
export function isValidNPWP(input: string): boolean {
  return parseNPWP(input).ok;
}
