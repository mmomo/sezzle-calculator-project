## Inception

- Let's start with the Inception phase. Don't write any code in this step. Read CLAUDE.md and docs/assignment.md. List any ambiguities or decisions i need to make, each with 2 options and a recommendation. Stop there and wait for my answers.

- Here are my answers: 1: B, 2: A, 3: A, 4: A, 5: A, 6: A, 7: A, 8: A, 9: A, 10: A

- Proceed


## Construction

- Create first unit of work for backend
- Continue with next unit of work for frontend

- Improve the frontend UI. Goals, in priority order:
1. Clarity: obvious inputs, operation selection, and a prominent result area.
2. Feedback: loading state during API calls, clear inline validation errors,
   and distinct styling for server errors vs. validation errors.
3. Responsive: works well at 375px mobile width and on desktop.
4. Accessibility: proper labels, keyboard navigation, visible focus states,
   aria-live on the result so screen readers announce it.
5. Visual: clean and minimal, one accent color, CSS variables for theming.
Constraints: plain CSS, no UI libraries, keep components and tests passing.
Propose the changes first.


-Review test coverage. Run the backend coverage and report the numbers. Then propose the most valuable missing tests
(edge cases, error paths, UI error states). Don't chase 100%; focus on key behavior. Plan first.

## Operations

- Audit the current code against docs/requirements.md and the deliverables in docs/assignment.md. Produce a checklist: each requirement  done / partial / missing, with file references. Don't change any code yet.