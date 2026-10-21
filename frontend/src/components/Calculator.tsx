import { useState, type ChangeEvent, type FormEvent } from 'react'
import { calculate, CalculatorApiError, type Operation } from '../api/calculator'
import './Calculator.css'

const OPERATIONS: { value: Operation; label: string; unary?: boolean }[] = [
  { value: 'add', label: 'Add (a + b)' },
  { value: 'subtract', label: 'Subtract (a − b)' },
  { value: 'multiply', label: 'Multiply (a × b)' },
  { value: 'divide', label: 'Divide (a ÷ b)' },
  { value: 'exponentiate', label: 'Exponentiate (a ^ b)' },
  { value: 'sqrt', label: 'Square root (√a)', unary: true },
  { value: 'percentage', label: 'Percentage (a% of b)' },
]

function isUnary(operation: Operation): boolean {
  return OPERATIONS.find((candidate) => candidate.value === operation)?.unary ?? false
}

export function Calculator() {
  const [operation, setOperation] = useState<Operation>('add')
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [aTouched, setATouched] = useState(false)
  const [bTouched, setBTouched] = useState(false)
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const unary = isUnary(operation)
  const aNumeric = a.trim() !== '' && !Number.isNaN(Number(a))
  const bNumeric = unary || (b.trim() !== '' && !Number.isNaN(Number(b)))
  const aMessage = aTouched && !aNumeric ? '"a" needs a number' : null
  const bMessage = !unary && bTouched && !bNumeric ? '"b" needs a number' : null
  const canSubmit = aNumeric && bNumeric && !loading

  function clearOutput() {
    setResult(null)
    setError(null)
  }

  function handleOperationChange(event: ChangeEvent<HTMLSelectElement>) {
    clearOutput()
    setOperation(event.target.value as Operation)
  }

  function handleAChange(event: ChangeEvent<HTMLInputElement>) {
    clearOutput()
    setA(event.target.value)
  }

  function handleBChange(event: ChangeEvent<HTMLInputElement>) {
    clearOutput()
    setB(event.target.value)
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setATouched(true)
    if (!unary) setBTouched(true)
    if (!canSubmit) return

    setLoading(true)
    clearOutput()

    try {
      const value = await calculate({
        operation,
        a: Number(a),
        b: unary ? undefined : Number(b),
      })
      setResult(value)
    } catch (err) {
      setError(
        err instanceof CalculatorApiError
          ? err.message
          : 'Something went wrong. Please try again.',
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="calculator" onSubmit={handleSubmit} aria-busy={loading}>
      <h1>Calculator</h1>

      <div className="field">
        <label htmlFor="operation">Operation</label>
        <select
          id="operation"
          value={operation}
          onChange={handleOperationChange}
          disabled={loading}
        >
          {OPERATIONS.map((op) => (
            <option key={op.value} value={op.value}>
              {op.label}
            </option>
          ))}
        </select>
      </div>

      <div className="field">
        <label htmlFor="operand-a">a</label>
        <input
          id="operand-a"
          type="text"
          inputMode="decimal"
          value={a}
          onChange={handleAChange}
          onBlur={() => setATouched(true)}
          placeholder="e.g. 4"
          disabled={loading}
          aria-invalid={aMessage ? true : undefined}
          aria-describedby={aMessage ? 'operand-a-message' : undefined}
        />
        {aMessage && (
          <p className="field-message" id="operand-a-message">
            {aMessage}
          </p>
        )}
      </div>

      {!unary && (
        <div className="field">
          <label htmlFor="operand-b">b</label>
          <input
            id="operand-b"
            type="text"
            inputMode="decimal"
            value={b}
            onChange={handleBChange}
            onBlur={() => setBTouched(true)}
            placeholder="e.g. 2"
            disabled={loading}
            aria-invalid={bMessage ? true : undefined}
            aria-describedby={bMessage ? 'operand-b-message' : undefined}
          />
          {bMessage && (
            <p className="field-message" id="operand-b-message">
              {bMessage}
            </p>
          )}
        </div>
      )}

      <button type="submit" disabled={!canSubmit}>
        {loading ? (
          <>
            <span className="spinner" aria-hidden="true" />
            Calculating…
          </>
        ) : (
          '='
        )}
      </button>

      <div className={`output${error ? ' output--error' : ''}`} aria-live="polite">
        {error !== null ? (
          <p className="output-error">{error}</p>
        ) : result !== null ? (
          <>
            <span className="output-label">Result</span>
            <p className="output-result">{result}</p>
          </>
        ) : (
          <p className="output-placeholder">Enter values and press =</p>
        )}
      </div>
    </form>
  )
}
