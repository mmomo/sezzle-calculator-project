# Requirements

Source of truth for this project. Supersedes `docs/assignment.md` wherever the
two differ — this document captures the decisions made during the Inception
phase to resolve ambiguity in the original assignment.

## Scope decisions

| # | Decision | Choice |
|---|----------|--------|
| 1 | Operation scope | Basic four (add, subtract, multiply, divide) + exponentiation, square root, percentage |
| 2 | API shape | Single generic `POST /api/calculate` endpoint |
| 3 | Input model | Two-operand only (one calculation per request), no expression parsing |
| 4 | Precision | `float64` throughout (frontend and backend); no arbitrary-precision decimal |
| 5 | Error handling | HTTP status codes carry meaning: `400` malformed input, `422` valid-shape-but-invalid-math |
| 6 | Validation | Backend is the single source of truth; frontend validation is UX-only (e.g. disable submit on empty field) |
| 7 | State | Stateless — no calculation history persisted |
| 8 | Dev connectivity | Vite dev-server proxies `/api/*` to the Go backend; same-origin in dev, no CORS needed |
| 9 | Dockerfile | Out of scope |
| 10 | Coverage reporting | Plain `go test -cover` / `vitest --coverage` terminal summary pasted into README |

## Operations

| Operation | Arity | Formula | Notes |
|---|---|---|---|
| `add` | binary | `a + b` | |
| `subtract` | binary | `a - b` | |
| `multiply` | binary | `a * b` | |
| `divide` | binary | `a / b` | `b == 0` is a domain error |
| `exponentiate` | binary | `a ^ b` | |
| `sqrt` | unary | `sqrt(a)` | `b` not required/used; `a < 0` is a domain error |
| `percentage` | binary | `(a / 100) * b` | reads as "a percent of b" (e.g. `a=50, b=200` → `100`); assumption, see below |

**Assumption:** "percentage" has no single universally-agreed calculator
semantics. We define it as "a percent of b" rather than "a as a percentage of
b" (`a / b * 100`) because it mirrors the common physical-calculator `%` key
behavior and avoids a division (so no extra divide-by-zero case). This is
called out explicitly in the README.

## API contract

### Request

```
POST /api/calculate
Content-Type: application/json

{
  "operation": "add" | "subtract" | "multiply" | "divide" | "exponentiate" | "sqrt" | "percentage",
  "a": number,
  "b": number   // omit for "sqrt"; required for all other operations
}
```

### Success response

```
200 OK

{ "result": number }
```

### Error response

```
400 Bad Request | 422 Unprocessable Entity

{ "error": string }
```

### Status code usage

| Status | Meaning | Examples |
|---|---|---|
| `400` | Malformed request — not a valid calculation to even attempt | invalid/missing JSON body, unknown `operation` value, `a` missing or not a number, `b` missing when required by the operation |
| `422` | Well-formed request, but the math is undefined/invalid | `divide` with `b == 0`, `sqrt` with `a < 0`, or any computed result that is `NaN`/`±Inf` (generic overflow guard applied after computing any operation) |
| `200` | Success | finite numeric result |

## Frontend requirements

- Inputs for operand `a`, operand `b` (hidden/disabled for `sqrt`), and an
  operation selector.
- Submit triggers a call to `POST /api/calculate`; result or error is
  displayed inline.
- Input validation (UX-only, see decision #6): prevent submitting
  non-numeric input or a missing required operand; backend errors are
  surfaced as-is.
- Responsive layout: usable on a mobile-width viewport (basic support, not a
  full design system).

## Non-functional

- Backend: Go standard library only (`net/http`), table-driven unit tests for
  the calculator logic and HTTP handlers.
- Frontend: React + TypeScript via Vite, plain CSS, Vitest + React Testing
  Library for the API client, hooks, and components.
- README documents setup, run instructions, API examples (including error
  cases), design decisions/assumptions, and pasted test coverage summaries.
