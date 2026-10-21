# Calculator

A full-stack calculator: a Go REST backend (stdlib only) and a React +
TypeScript frontend (Vite), built as a take-home assignment
(`docs/assignment.md`). Scope decisions and the detailed design are in
[`docs/requirements.md`](docs/requirements.md) and
[`docs/design.md`](docs/design.md) — this README covers setup, running, and
a summary of the decisions that matter to a reviewer.

## Development process

Built using **AI-DLC** (AI-Driven Development Life Cycle), working through
three phases with Claude Code:

- **Inception** — read the assignment, surfaced ambiguities as explicit
  options + a recommendation, and got a decision on each before any code was
  written. Locked into [`docs/requirements.md`](docs/requirements.md).
- **Construction** — designed the package/component layout
  ([`docs/design.md`](docs/design.md)), then implemented backend and
  frontend as separate reviewed units of work, tests alongside the code.
- **Operations** — ran the app end-to-end (backend + frontend dev server),
  manually verified golden-path and error cases in the browser, and audited
  the result against the requirements and assignment deliverables.

The full prompt history for each phase is in [`PROMPTS.md`](PROMPTS.md).

## Prerequisites

- Go 1.27+ (developed against `go1.27.2`)
- Node 20+ and npm (developed against Node `v26.12.0`, npm `11.20.0`)

## Setup

```bash
# backend — no dependencies to install, stdlib only
cd backend && go build ./...

# frontend
cd frontend && npm install
```

## Running the app

Two terminals, both from the repo root:

```bash
# terminal 1 — backend, listens on :8080 (override with PORT=...)
cd backend && go run ./cmd/server

# terminal 2 — frontend dev server
cd frontend && npm run dev
```

Open the URL Vite prints (typically `http://localhost:5173`). The dev
server proxies `/api/*` to the backend on `:8080` (`frontend/vite.config.ts`),
so there's no CORS setup needed and no separate base URL to configure.

To build the frontend for production: `cd frontend && npm run build` (output
in `frontend/dist`, not served by anything in this repo — it's a static
build you'd put behind your own web server).

## Running tests

```bash
cd backend && go test ./... -cover
```

```
ok  	.../backend/internal/api		coverage: 100.0% of statements
ok  	.../backend/internal/calculator	coverage: 100.0% of statements
```

(`backend/cmd/server` shows 0% — it's just the `main` wiring, not tested
directly.)

**Frontend has no automated tests.** This was a deliberate scope cut for
the assignment's timebox, made explicitly against `CLAUDE.md`'s stated
Vitest/React Testing Library plan — not an oversight. Manual verification
was done instead: golden path (e.g. `4 + 2 = 6`), the `sqrt` unary case
(operand `b` hides correctly), and a backend error (`divide` by zero)
rendering the server's message in the UI. See "Known limitations" below.

## API

### `POST /api/calculate`

Request:

```json
{
  "operation": "add" | "subtract" | "multiply" | "divide" | "exponentiate" | "sqrt" | "percentage",
  "a": <number>,
  "b": <number>   // omit for "sqrt"
}
```

`percentage` is defined as **"a percent of b"**: `(a / 100) * b` — e.g.
`a=50, b=200` → `100`. This is an assumption (the term is ambiguous); see
`docs/requirements.md` for the alternative considered.

### Examples

Success:

```bash
curl -s -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","a":2,"b":3}'
# {"result":5}
```

```bash
curl -s -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","a":16}'
# {"result":4}
```

```bash
curl -s -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"percentage","a":50,"b":200}'
# {"result":100}
```

Malformed input → `400`:

```bash
curl -s -o /dev/null -w '%{http_code}\n' -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"modulo","a":1,"b":2}'
# 400, body: {"error":"unknown operation: \"modulo\""}
```

Valid shape, invalid math → `422`:

```bash
curl -s -o /dev/null -w '%{http_code}\n' -X POST localhost:8080/api/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","a":1,"b":0}'
# 422, body: {"error":"division by zero"}
```

Full status-code table and validation rules are in
`docs/requirements.md` → "API contract".

## Design decisions (summary)

Full rationale in [`docs/requirements.md`](docs/requirements.md) and
[`docs/design.md`](docs/design.md). Highlights:

- **Single generic endpoint** (`/api/calculate`) rather than one route per
  operation — less duplication, same validation path for every operation.
- **Two-operand model, no expression parsing** — matches a physical
  calculator; an expression parser wasn't worth the timebox.
- **`float64` throughout**, no arbitrary-precision decimal — sufficient for
  a calculator demo.
- **Status codes carry meaning**: `400` for malformed requests, `422` for
  well-formed-but-invalid math (divide by zero, negative sqrt, a result
  that overflows to `±Inf`).
- **Backend is the sole source of truth for validation**; the frontend only
  does UX-level checks (disabling submit) and renders backend errors as-is.
- **Stateless** — no calculation history.
- **No Dockerfile** — marked optional in the assignment, skipped given the
  timebox.

## Project structure

```
backend/
  cmd/server/            entrypoint
  internal/calculator/   pure arithmetic + domain errors
  internal/api/          HTTP handler + request/response types
frontend/
  src/api/               typed fetch client
  src/components/        Calculator form/UI
docs/
  assignment.md          original brief (input only)
  requirements.md        locked scope decisions + API contract (source of truth)
  design.md              package/component layout and rationale
PROMPTS.md               prompts used during development
```

## Known limitations

- No frontend automated tests (see "Running tests" above).
- Responsive/mobile layout has a CSS breakpoint but hasn't been verified on
  an actual narrow viewport, only visually inspected at desktop width.
- No persistence/history, by design — each calculation is a single
  stateless request.
