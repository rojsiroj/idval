/**
 * Parses and validates Indonesian phone numbers.
 * @module
 */
import { fail, ok, type Result } from "./errors.js";

export type PhoneKind = "mobile" | "landline";

/**
 * The mobile brand that owns a prefix. Indonesia has no mobile number portability,
 * so the prefix identifies the brand. XL, Axis and Smartfren are brands of XLSmart.
 * Indosat and Tri are brands of Indosat Ooredoo Hutchison.
 */
export type Operator = "telkomsel" | "indosat" | "xl" | "axis" | "tri" | "smartfren";

export interface PhoneNumber {
  /** For example "+6281234567890". */
  e164: string;
  /** For example "081234567890". */
  national: string;
  kind: PhoneKind;
  /** Present only when `kind` is "mobile". */
  operator?: Operator;
}

const ranges: Record<Operator, string[]> = {
  telkomsel: ["811", "812", "813", "821", "822", "823", "851", "852", "853"],
  indosat: ["814", "815", "816", "855", "856", "857", "858"],
  xl: ["817", "818", "819", "859", "877", "878"],
  axis: ["831", "832", "833", "838"],
  tri: ["895", "896", "897", "898", "899"],
  smartfren: ["881", "882", "883", "884", "885", "886", "887", "888", "889"],
};

/** Maps the first 3 digits of a mobile number (after the trunk 0) to the operator. */
const prefixes = new Map<string, Operator>(
  (Object.keys(ranges) as Operator[]).flatMap((op) => ranges[op].map((p) => [p, op] as const)),
);

/**
 * Parses an Indonesian phone number. Spaces, dots, dashes and parentheses in `input`
 * are ignored. The number must start with "0", "62" or "+62". The function accepts a
 * trunk "0" after the country code, for example "+62 0812...".
 */
export function parsePhone(input: string): Result<PhoneNumber> {
  const d = input.replace(/[ .\-()]/g, "");
  // A plus sign is permitted only before the first digit.
  if (!/^\+?[0-9]*$/.test(d)) return fail("INVALID_CHARACTERS");

  let nsn: string;
  if (d.startsWith("+62")) nsn = d.slice(3);
  else if (d.startsWith("62")) nsn = d.slice(2);
  else if (d.startsWith("0")) nsn = d.slice(1);
  else return fail("INVALID_PREFIX");
  if (!d.startsWith("0") && nsn.startsWith("0")) nsn = nsn.slice(1);
  if (nsn === "") return fail("INVALID_LENGTH");

  const e164 = `+62${nsn}`;
  const national = `0${nsn}`;
  const c = nsn[0] as string;
  if (c === "8") {
    if (nsn.length < 9 || nsn.length > 12) return fail("INVALID_LENGTH");
    const operator = prefixes.get(nsn.slice(0, 3));
    if (operator === undefined) return fail("UNKNOWN_OPERATOR");
    return ok({ e164, national, kind: "mobile", operator });
  }
  if (c >= "2" && c <= "9") {
    if (nsn.length < 8 || nsn.length > 11) return fail("INVALID_LENGTH");
    return ok({ e164, national, kind: "landline" });
  }
  return fail("INVALID_PREFIX");
}

/** Returns true if `input` is a valid Indonesian phone number. */
export function isValidPhone(input: string): boolean {
  return parsePhone(input).ok;
}
