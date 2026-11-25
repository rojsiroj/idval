/**
 * Removes spaces, dots and dashes from `s`.
 * Returns undefined if the result contains a character that is not an ASCII digit.
 */
export function digits(s: string): string | undefined {
  const d = s.replace(/[ .-]/g, "");
  return /^[0-9]*$/.test(d) ? d : undefined;
}
