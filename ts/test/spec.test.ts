// Runs the shared test vectors in ../spec. The Go tests run the same vectors.
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { describe, expect, it } from "vitest";
import { parseNIK, parseNITKU, parseNPWP, parsePhone, type Result } from "../src/index.js";

interface Case {
  name: string;
  input: string;
  result?: Record<string, unknown>;
  error?: string;
}

function load<T>(file: string): T {
  const path = fileURLToPath(new URL(`../../spec/${file}`, import.meta.url));
  return JSON.parse(readFileSync(path, "utf8")) as T;
}

/** Parses "YYYY-MM-DD" as a local calendar date, to match the Go reference date. */
function localDate(s: string): Date {
  const [y, m, d] = s.split("-").map(Number) as [number, number, number];
  return new Date(y, m - 1, d);
}

function check(c: Case, got: Result<unknown>): void {
  if (c.error !== undefined) {
    expect(got.ok ? got.value : got.error.code).toBe(c.error);
  } else {
    expect(got.ok ? got.value : got.error.code).toStrictEqual(c.result);
  }
}

describe("nik", () => {
  const spec = load<{ now: string; cases: Case[] }>("nik.json");
  const now = localDate(spec.now);
  it.each(spec.cases)("$name", (c) => check(c, parseNIK(c.input, { now })));
});

describe("npwp", () => {
  const spec = load<{ now: string; npwp: Case[]; nitku: Case[] }>("npwp.json");
  const now = localDate(spec.now);
  it.each(spec.npwp)("npwp: $name", (c) => check(c, parseNPWP(c.input, { now })));
  it.each(spec.nitku)("nitku: $name", (c) => check(c, parseNITKU(c.input, { now })));
});

describe("phone", () => {
  const spec = load<{ cases: Case[] }>("phone.json");
  it.each(spec.cases)("$name", (c) => check(c, parsePhone(c.input)));
});
