# Design

Implementation-level design following from `docs/requirements.md`. Describes
package/component layout and the responsibilities of each piece — not the
code itself.

## Backend layout

```
backend/
  go.mod
  cmd/server/main.go           entrypoint: build mux, start http.Server
  internal/calculator/
    calculator.go               pure arithmetic, no HTTP concerns
    calculator_test.go
  internal/api/
    types.go                    request/response structs
    handler.go                  decode, structurally validate, dispatch, encode
    handler_test.go
```

`internal/...` because this is a standalone binary, not a library other
modules import.

### Layering and where each status code is decided

Two layers, deliberately kept separate so each layer owns exactly one class
of error (and thus one status code):

- **`internal/api` (handler)** — HTTP concerns and *structural* validation:
  is the JSON well-formed, is `operation` one of the known values, is `a`
  present, is `b` present when the operation needs it. Any failure here is a
  **400**. The handler never computes anything itself.
- **`internal/calculator`** — pure functions, no knowledge of HTTP. Takes
  already-structurally-valid input and returns `(float64, error)`. Any
  failure here is a **domain** error (divide by zero, negative sqrt, a
  non-finite result) and maps to **422**.

This means `calculator.Calculate` is only ever called with a known
operation — the handler's whitelist check guarantees that — so the package
never has to reason about "unknown operation" as a case of its own.

### `internal/calculator`

```go
type Operation string

const (
    Add          Operation = "add"
    Subtract     Operation = "subtract"
    Multiply     Operation = "multiply"
    Divide       Operation = "divide"
    Exponentiate Operation = "exponentiate"
    Sqrt         Operation = "sqrt"
    Percentage   Operation = "percentage"
)

func Calculate(op Operation, a, b float64) (float64, error)
```

- `Sqrt` is unary; the handler passes `b = 0` for it and `Calculate` simply
  ignores `b` in that branch. `b`'s "optionality" is a JSON/request concern,
  not a calculator concern — `Calculate` always takes two plain `float64`s.
- Sentinel errors for domain failures: `ErrDivideByZero`, `ErrNegativeSqrt`,
  `ErrNotFinite`. `Divide` and `Sqrt` check their own precondition explicitly
  (gives a precise message); additionally, `Calculate` runs one generic
  `math.IsNaN`/`math.IsInf` check on the final result as a catch-all safety
  net (e.g. covers `exponentiate` overflow) before returning, so no operation
  can silently return `NaN`/`±Inf` to the client.
- `Percentage(a, b) = (a / 100) * b` — see `docs/requirements.md` for the
  semantics assumption. No division involved, so no extra zero-check needed.

### `internal/api`

```go
type CalculateRequest struct {
    Operation string   `json:"operation"`
    A         *float64 `json:"a"`
    B         *float64 `json:"b"`
}

type CalculateResponse struct {
    Result float64 `json:"result"`
}

type ErrorResponse struct {
    Error string `json:"error"`
}
```

`A`/`B` are pointers so a missing JSON field is distinguishable from an
explicit `0` — required for the "`a` is required" / "`b` is required for
this operation" validation messages.

Handler flow for `POST /api/calculate`:

1. Decode body → malformed JSON is `400`.
2. Look up `Operation` in a whitelist (`map[string]arity`, arity = unary or
   binary) → unknown value is `400`.
3. `A == nil` → `400`. If arity is binary and `B == nil` → `400`.
   (If arity is unary and `B` is present anyway, it's ignored rather than
   rejected — being lenient about an unused extra field costs nothing and
   keeps the client contract simple.)
4. Call `calculator.Calculate(op, *A, b)` (`b` defaulted to `0` when `nil`).
5. Error from step 4 → `422` with `{"error": err.Error()}`.
6. Success → `200` with `{"result": ...}`.

### Routing

Go's stdlib `net/http.ServeMux` (Go 1.22+ method+pattern syntax, e.g.
`"POST /api/calculate"`) is sufficient — no third-party router needed, which
keeps us on stdlib-only per `CLAUDE.md`.

### Testing

- `calculator_test.go`: table-driven, one table entry per operation covering
  the happy path plus its specific edge cases (divide by zero, negative
  sqrt, an exponent large enough to overflow to `+Inf`).
- `handler_test.go`: table-driven over raw request bodies using
  `httptest.NewRecorder`/`httptest.NewRequest`, asserting status code and
  JSON body shape for each of the 400/422/200 cases.

## Frontend layout

```
frontend/src/
  api/
    calculator.ts                typed fetch client + CalculatorApiError
    calculator.test.ts
  components/
    Calculator.tsx                form: operands, operation select, submit, result/error display
    Calculator.test.tsx
  App.tsx                         renders <Calculator />
```

- **`api/calculator.ts`**: one function, `calculate(request): Promise<number>`,
  typed against the same request/response shapes as the backend contract.
  On a non-2xx response it reads `{error}` from the body and throws a
  `CalculatorApiError` carrying the message and status; on network failure
  it lets the fetch rejection propagate. This is the single place that knows
  about the HTTP contract — the component never touches `fetch` directly,
  which is what makes it mockable in component tests without a network
  layer.
- **`components/Calculator.tsx`**: local `useState` for operand `a`,
  operand `b`, selected operation, result, error, and loading — no separate
  state-management hook. The logic here is small enough that a hook would be
  an abstraction without a second consumer to justify it (per `CLAUDE.md`:
  no abstractions beyond what's needed). Deriving unary-vs-binary from the
  selected operation hides/disables the `b` input for `sqrt`.
  UX-only validation (decision #6 in requirements): submit is disabled when
  a required operand is empty or non-numeric; the backend's `400`/`422`
  messages are rendered as-is on failure, no re-interpretation client-side.
- Responsive layout: a single CSS flex/grid layout for the form, one
  `max-width` + media-query breakpoint to stack inputs vertically on narrow
  viewports. No UI library, per constraints.

### Testing

- `api/calculator.test.ts`: mocks `global.fetch`, asserts the request body
  sent matches the contract and that both success and error (400/422/network
  failure) responses are translated correctly.
- `components/Calculator.test.tsx`: React Testing Library, mocks the
  `api/calculator` module (not `fetch`) so component tests stay focused on
  rendering/interaction — fill in operands, pick an operation, submit, and
  assert the rendered result or error message, including the `sqrt` unary
  case (no `b` input shown) and a disabled-submit case for invalid input.

## Dev connectivity

Per requirements decision #8: `vite.config.ts` proxies `/api/*` to the Go
backend's port in dev, so the frontend calls same-origin relative paths
(`/api/calculate`) and no CORS handling is needed on the backend.
