# Project: Full-stack calculator (Sezzle take-home)

## Context
- A calculator app with a React + TypeScript frontend and a Go REST backend.
The frontend consumes the backend API for all arithmetic operations.
- Timebox: small scope. Prioritize correctness, clarity and tests over features.

## Workflow
1. Before each unit of work, propose a short plan and wait for approval.
2. Write tests before or alongside the implementation.
3. Keep each change small enough for one commit and suggest a commit message.
4. Briefly explain any non-obvious language or framework idioms you use.
5. If a requirement is ambiguous, ask instead of assuming.

## Stack and constraints
- Backend: Go, standard library only (net/http). No frameworks.
- Frontend: React + TypeScript, built with Vite. Plain CSS, no UI libraries.
- Testing: Go `testing` package (table-driven tests), Vitest + React Testing Library.
- Do not add dependencies without asking first.

## Repository structure
- backend/    Go module (calculator logic, HTTP handlers, entrypoint)
- frontend/   React + TS app (API client, hooks, components)
- docs/       Requirements and design documents

## Key docs
- Original assignment: docs/assignment.md (input only)
- Approved requirements: docs/requirements.md (source of truth once written)
- Design: docs/design.md

