/** Identifies the cause of a validation error. The Go and TypeScript packages use the same values. */
export type ErrorCode =
  | "INVALID_CHARACTERS"
  | "INVALID_LENGTH"
  | "UNKNOWN_REGION"
  | "INVALID_BIRTH_DATE"
  | "INVALID_SERIAL"
  | "INVALID_PREFIX"
  | "UNKNOWN_OPERATOR";

export type Lang = "en" | "id";

/** Default messages for each error code, in English and Indonesian. */
export const messages: Readonly<Record<Lang, Readonly<Record<ErrorCode, string>>>> = {
  en: {
    INVALID_CHARACTERS: "The value contains characters that are not permitted.",
    INVALID_LENGTH: "The value has an incorrect number of digits.",
    UNKNOWN_REGION: "The region code is not known.",
    INVALID_BIRTH_DATE: "The birth date is not a valid date.",
    INVALID_SERIAL: "The serial number must not be 0000.",
    INVALID_PREFIX: "The number does not start with 0, 62 or +62, or the area code is not valid.",
    UNKNOWN_OPERATOR: "The mobile prefix does not belong to a known operator.",
  },
  id: {
    INVALID_CHARACTERS: "Nilai berisi karakter yang tidak diizinkan.",
    INVALID_LENGTH: "Jumlah digit tidak sesuai.",
    UNKNOWN_REGION: "Kode wilayah tidak dikenal.",
    INVALID_BIRTH_DATE: "Tanggal lahir tidak valid.",
    INVALID_SERIAL: "Nomor urut tidak boleh 0000.",
    INVALID_PREFIX: "Nomor tidak diawali 0, 62 atau +62, atau kode area tidak valid.",
    UNKNOWN_OPERATOR: "Prefiks nomor seluler tidak dikenal.",
  },
};

export interface ValidationError {
  code: ErrorCode;
  /** The English message. Use `messages.id[code]` for Indonesian. */
  message: string;
}

/** The result of a parse function. Check `ok` before you read `value` or `error`. */
export type Result<T> = { ok: true; value: T } | { ok: false; error: ValidationError };

export function ok<T>(value: T): Result<T> {
  return { ok: true, value };
}

export function fail<T>(code: ErrorCode): Result<T> {
  return { ok: false, error: { code, message: messages.en[code] } };
}
